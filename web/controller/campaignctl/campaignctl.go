package campaignctl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/StephanHCB/go-campaign-service/api/v1/apierrors"
	"github.com/StephanHCB/go-campaign-service/api/v1/campaign"
	"github.com/StephanHCB/go-campaign-service/internal/service/campaignsrv"
	"github.com/StephanHCB/go-campaign-service/web/middleware/authentication"
	"github.com/StephanHCB/go-campaign-service/web/util/media"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-http-utils/headers"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"time"
)

type CampaignControllerImpl struct {
	s campaignsrv.CampaignService
}

func Create(server chi.Router, campaignService campaignsrv.CampaignService) campaign.CampaignApi {
	controller := &CampaignControllerImpl{s: campaignService}
	controller.SetupRoutes(server)
	return controller
}

func (c *CampaignControllerImpl) SetupRoutes(server chi.Router) {
	server.Put("/api/rest/v1/campaigns", c.CreateCampaign)
	server.Post("/api/rest/v1/campaigns/{id:[1-9][0-9]*}", c.UpdateCampaign)
	server.Get("/api/rest/v1/campaigns/{id:[1-9][0-9]*}", c.GetCampaign)
	server.Post("/api/rest/v1/campaigns/{id:[1-9][0-9]*}/execute", c.ExecuteCampaign)
}

func (c *CampaignControllerImpl) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// check authentication and role - in a real world scenario we would place these checks in a middleware
	// with a list of allowed url patterns
	err := authentication.CheckUserIsLoggedIn(ctx)
	if err != nil {
		campaignUnauthenticatedErrorHandler(ctx, w, r, err)
		return
	}
	err = authentication.CheckUserHasRole(ctx, "admin")
	if err != nil {
		campaignUnauthorizedErrorHandler(ctx, w, r, err)
		return
	}

	dto, err := parseBodyToCampaignDto(ctx, w, r)
	if err != nil {
		campaignParseErrorHandler(ctx, w, r, err)
		return
	}

	validationErrs := validate(ctx, dto)
	if len(validationErrs) != 0 {
		campaignValidationErrorHandler(ctx, w, r, validationErrs)
		return
	}

	newCampaign := c.s.NewCampaign(ctx)
	err = mapDtoToCampaign(dto, newCampaign)
	if err != nil {
		campaignParseErrorHandler(ctx, w, r, err)
		return
	}

	id, err := c.s.CreateCampaign(ctx, newCampaign)
	if err != nil {
		campaignWriteErrorHandler(ctx, w, r, err)
		return
	}

	location := fmt.Sprintf("%s/%d", r.RequestURI, id)
	log.Ctx(ctx).Info().Msg("sending new Location " + location)
	w.Header().Set(headers.Location, location)
	w.WriteHeader(http.StatusCreated)
}

func (c *CampaignControllerImpl) UpdateCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := idFromVars(ctx, w, r)
	if err != nil {
		return
	}
	dto, err := parseBodyToCampaignDto(ctx, w, r)
	if err != nil {
		campaignParseErrorHandler(ctx, w, r, err)
		return
	}
	campaignRetrieved, err := c.s.GetCampaign(ctx, id)
	if err != nil {
		campaignNotFoundErrorHandler(ctx, w, r, id)
		return
	}
	validationErrs := validate(ctx, dto)
	if len(validationErrs) != 0 {
		campaignValidationErrorHandler(ctx, w, r, validationErrs)
		return
	}
	err = mapDtoToCampaign(dto, campaignRetrieved)
	if err != nil {
		campaignParseErrorHandler(ctx, w, r, err)
		return
	}
	err = c.s.UpdateCampaign(ctx, campaignRetrieved)
	if err != nil {
		campaignWriteErrorHandler(ctx, w, r, err)
		return
	}
	w.Header().Add(headers.Location, r.RequestURI)
}

func (c *CampaignControllerImpl) GetCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := idFromVars(ctx, w, r)
	if err != nil {
		return
	}
	existingCampaign, err := c.s.GetCampaign(ctx, id)
	if err != nil {
		campaignNotFoundErrorHandler(ctx, w, r, id)
		return
	}
	dto := campaign.CampaignDto{}
	mapCampaignToDto(existingCampaign, &dto)

	w.Header().Set(headers.Location, r.RequestURI)
	w.Header().Add(headers.ContentType, media.ContentTypeApplicationJson)
	writeJson(ctx, w, dto)
}

func (c *CampaignControllerImpl) ExecuteCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := idFromVars(ctx, w, r)
	if err != nil {
		return
	}
	existingCampaign, err := c.s.GetCampaign(ctx, id)
	if err != nil {
		campaignNotFoundErrorHandler(ctx, w, r, id)
		return
	}

	resultMap, err := c.s.ExecuteCampaign(ctx, existingCampaign)
	if err != nil {
		campaignExecutionGlobalErrorHandler(ctx, w, r, err, id)
		return
	}

	dto := campaign.ExecutionResultDto{}
	mapResultMapToExecutionResultDto(resultMap, &dto)

	w.Header().Set(headers.Location, fmt.Sprintf("/api/rest/v1/campaigns/%d", id))
	w.Header().Add(headers.ContentType, media.ContentTypeApplicationJson)
	writeJson(ctx, w, dto)
}

func idFromVars(ctx context.Context, w http.ResponseWriter, r *http.Request) (uint, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		invalidIdErrorHandler(ctx, w, r, err, idStr)
	}
	return uint(id), err
}

func parseBodyToCampaignDto(ctx context.Context, w http.ResponseWriter, r *http.Request) (*campaign.CampaignDto, error) {
	decoder := json.NewDecoder(r.Body)
	dto := &campaign.CampaignDto{}
	err := decoder.Decode(dto)
	if err != nil {
		dto = &campaign.CampaignDto{}
	}
	return dto, err
}

func campaignUnauthenticatedErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	log.Ctx(ctx).Warn().Msgf("SECURITY: Unauthenticated call: %v", err)
	errorHandler(ctx, w, r, "campaign.security.error", http.StatusUnauthorized, []string{})
}

func campaignUnauthorizedErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	log.Ctx(ctx).Warn().Msgf("SECURITY: Unauthorized call: %v", err)
	errorHandler(ctx, w, r, "campaign.security.error", http.StatusForbidden, []string{})
}

func campaignParseErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	log.Ctx(ctx).Warn().Err(err).Msgf("campaign body could not be parsed: %v", err)
	errorHandler(ctx, w, r, "campaign.parse.error", http.StatusBadRequest, []string{})
}

func campaignValidationErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, validationErrors []string) {
	log.Ctx(ctx).Warn().Msgf("received campaign data with validation errors: %v", validationErrors)
	errorHandler(ctx, w, r, "campaign.data.invalid", http.StatusBadRequest, validationErrors)
}

func invalidIdErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error, id string) {
	log.Ctx(ctx).Warn().Err(err).Msgf("received invalid attendee id '%s'", id)
	errorHandler(ctx, w, r, "campaign.id.invalid", http.StatusBadRequest, []string{})
}

func campaignNotFoundErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, id uint) {
	log.Ctx(ctx).Warn().Msgf("campaign id %v not found", id)
	errorHandler(ctx, w, r, "campaign.id.notfound", http.StatusNotFound, []string{})
}

func campaignExecutionGlobalErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error, id uint) {
	log.Ctx(ctx).Warn().Err(err).Msgf("campaign id %v execution error - mails were not sent: %v", id, err)
	errorHandler(ctx, w, r, "campaign.execution.failed", http.StatusBadGateway, []string{err.Error()})
}

func campaignWriteErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	log.Ctx(ctx).Warn().Err(err).Msgf("campaign could not be written: %v", err)
	if err.Error() == "duplicate campaign subject" {
		errorHandler(ctx, w, r, "campaign.data.duplicate", http.StatusBadRequest, []string{"there is already a campaign with this subject"})
	} else {
		errorHandler(ctx, w, r, "campaign.write.error", http.StatusInternalServerError, []string{})
	}
}

func errorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, msg string, status int, details []string) {
	timestamp := time.Now().Format(time.RFC3339)
	response := apierrors.ErrorDto{Message: msg, Timestamp: timestamp, Details: details, RequestId: middleware.GetReqID(ctx)}
	w.Header().Set(headers.ContentType, media.ContentTypeApplicationJson)
	w.WriteHeader(status)
	writeJson(ctx, w, response)
}

func writeJson(ctx context.Context, w http.ResponseWriter, v interface{}) {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msgf("error while encoding json response: %v", err)
	}
}
