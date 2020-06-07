package mailserviceclient

import (
	"context"
	"fmt"
	"github.com/StephanHCB/go-campaign-service/api/v1/apierrors"
	"github.com/StephanHCB/go-campaign-service/internal/repository/configuration"
	"github.com/StephanHCB/go-campaign-service/internal/repository/mailservice"
	"github.com/StephanHCB/go-campaign-service/internal/repository/util/downstreamcall"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type MailSenderRepositoryImpl struct {
	netClient *http.Client
}

const HystrixCommandName = "mailservice_send"

const sendmailEndpoint = "api/rest/v1/sendmail"

// --- instance creation ---

func Create() mailservice.MailSenderRepository {
	timeoutMs := configuration.MailerServiceTimeoutMs()

	downstreamcall.ConfigureHystrixCommand(HystrixCommandName, int(timeoutMs))

	return &MailSenderRepositoryImpl{
		netClient: &http.Client{
			// theoretically, this is no longer necessary with hystrix
			Timeout: time.Millisecond * time.Duration(timeoutMs) * 2,
		},
	}
}

// --- implementation of repository interface ---

type EmailDto struct {
	ToAddress string `json:"to_address"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

func (r *MailSenderRepositoryImpl) SendEmail(ctx context.Context, address string, subject string, body string) error {
	requestDto := EmailDto{
		ToAddress: address,
		Subject:   subject,
		Body:      body,
	}
	requestBody, err := downstreamcall.RenderJson(requestDto)
	if err != nil {
		return err
	}

	responseBody, httpstatus, err := downstreamcall.HystrixPerformPOST(ctx, HystrixCommandName, r.netClient,
		configuration.MailerServiceUrl()+sendmailEndpoint, requestBody)

	if err != nil || httpstatus != http.StatusOK {
		if err == nil {
			err = fmt.Errorf("unexpected http status %d, was expecting %d", httpstatus, http.StatusOK)
		}

		errorResponseDto := &apierrors.ErrorDto{}
		err2 := downstreamcall.ParseJson(responseBody, errorResponseDto)
		if err2 == nil {
			log.Ctx(ctx).Error().Err(err).Msgf("Error sending mail to '%s' via mailer-service: error from response is %s, local error is %s", address, errorResponseDto.Message, err.Error())
		} else {
			log.Ctx(ctx).Error().Err(err).Msgf("Error sending mail to '%s' via mailer-service with no structured response available: local error is %s", address, err.Error())
		}

		return err
	}

	return nil
}
