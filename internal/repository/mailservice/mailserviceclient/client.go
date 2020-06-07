package mailserviceclient

import (
	"context"
	"fmt"
	"github.com/StephanHCB/go-campaign-service/api/v1/apierrors"
	"github.com/StephanHCB/go-campaign-service/internal/repository/configuration"
	"github.com/StephanHCB/go-campaign-service/internal/repository/mailservice"
	"github.com/StephanHCB/go-campaign-service/internal/repository/util/downstreamcall"
	"github.com/afex/hystrix-go/hystrix"
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
	// configure circuit breaker
	hystrix.ConfigureCommand(HystrixCommandName, hystrix.CommandConfig{
		Timeout:               int(configuration.MailerServiceTimeoutMs()),
		MaxConcurrentRequests: 100,
		ErrorPercentThreshold: 25,
	})

	return &MailSenderRepositoryImpl{
		netClient: &http.Client{
			// theoretically, this is no longer necessary with hystrix
			Timeout: time.Millisecond * time.Duration(configuration.MailerServiceTimeoutMs()) * 2,
		},
	}
}

// --- circuit breaker layer ---

type responseInfo struct {
	body   string
	status int
}

func (r *MailSenderRepositoryImpl) hystrixPerformPOST(ctx context.Context, url string, requestBody string) (string, int, error) {
	output := make(chan responseInfo, 1)

	// hystrix.DoC blocks until either completed or error returned
	err := hystrix.DoC(ctx, HystrixCommandName, func(subctx context.Context) error {
		responseBody, httpstatus, innerErr := downstreamcall.PerformPOST(subctx, r.netClient, url, requestBody)
		output <- responseInfo{
			body:   responseBody,
			status: httpstatus,
		}

		// if we return an error at this point, it will count towards opening the circuit breaker
		if httpstatus >= 500 && innerErr == nil {
			// so let's make any http status in the 500 range causes us to return an error
			// in a real world situation this may need some more attention
			innerErr = fmt.Errorf("got unexpected http status %d", httpstatus)
		}
		return innerErr
	}, nil)

	responseData := responseInfo{}

	// non-blocking receive for optional output
	select {
	case out := <-output:
		responseData = out
	default:
		// presence of default branch means select will not block even if none of the channels are ready to read from
	}

	return responseData.body, responseData.status, err
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

	responseBody, httpstatus, err := r.hystrixPerformPOST(ctx, configuration.MailerServiceUrl()+sendmailEndpoint, requestBody)
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
