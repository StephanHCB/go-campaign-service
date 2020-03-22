package campaign

import "net/http"

// --- models ---

// Model for CampaignDto.
//
// swagger:model campaignDto
type CampaignDto struct {
	// The email subject
	Subject string `json:"subject"`
	// The email body
	Body string `json:"body"`
	// The list of recipients
	Recipients []RecipientDto `json:"recipients"`
}

// Model for RecipientDto.
//
// swagger:model recipientDto
type RecipientDto struct {
	// The email address to send to
	ToAddress string `json:"to_address"`
}

// --- parameters and responses --- needed to use models

// Parameters for creating a Campaign
//
// swagger:parameters createCampaignParams
type CreateCampaignParams struct {
	// in:body
	Body CampaignDto
}

// Parameters for updating a Campaign
//
// swagger:parameters updateCampaignParams
type UpdateCampaignParams struct {
	// The id of the campaign
	//
	// in:path
	Id string

	// The changed data of the campaign to be set
	//
	// in:body
	Body CampaignDto
}

// Parameters for fetching a Campaign
//
// swagger:parameters getCampaignParams
type GetCampaignParams struct {
	// The id of the campaign
	//
	// in:path
	Id string
}

// The campaign location response with just a Location header
//
// swagger:response campaignLocationResponse
type CampaignLocationResponse struct {
	// Location header
	Location string `json:"Location"`
}

// The campaign data response including a Location header
//
// swagger:response campaignDataResponse
type CampaignDataResponse struct {
	// Location header
	Location string `json:"Location"`

	// The data of the campaign
	//
	// in:body
	Body CampaignDto
}

// --- routes ---

type CampaignApi interface {
	// swagger:route PUT /api/rest/v1/campaigns campaign-tag createCampaignParams
	// This will create a new campaign and return its location
	//
	// responses:
	//   201: campaignLocationResponse
	//   400: errorResponse
	//   401: errorResponse
	//   403: errorResponse
	CreateCampaign(w http.ResponseWriter, r *http.Request)

	// swagger:route POST /api/rest/v1/campaigns/{Id} campaign-tag updateCampaignParams
	// This will update an existing campaign
	//
	// responses:
	//   200: campaignLocationResponse
	//   400: errorResponse
	//   401: errorResponse
	//   403: errorResponse
	//   404: errorResponse
	UpdateCampaign(w http.ResponseWriter, r *http.Request)

	// swagger:route GET /api/rest/v1/campaigns/{Id} campaign-tag getCampaignParams
	// This will return an existing campaign
	//
	// responses:
	//   200: campaignDataResponse
	//   401: errorResponse
	//   403: errorResponse
	//   404: errorResponse
	GetCampaign(w http.ResponseWriter, r *http.Request)
}
