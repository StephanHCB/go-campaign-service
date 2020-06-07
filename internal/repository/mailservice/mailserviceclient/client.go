package mailserviceclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/StephanHCB/go-campaign-service/api/v1/apierrors"
	"github.com/StephanHCB/go-campaign-service/internal/repository/configuration"
	"github.com/StephanHCB/go-campaign-service/internal/repository/mailservice"
	"github.com/StephanHCB/go-campaign-service/web/util/media"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type MailSenderRepositoryImpl struct {
	netClient *http.Client
}

const HystrixCommandName = "mailservice_send"

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

const sendmailEndpoint = "api/rest/v1/sendmail"

type EmailDto struct {
	ToAddress string `json:"to_address"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

func (r *MailSenderRepositoryImpl) performPost(ctx context.Context, url string, requestBody string) (string, error) {
	response, err := r.netClient.Post(url, media.ContentTypeApplicationJson, strings.NewReader(requestBody))
	if err != nil {
		return "", err
	}
	status := response.StatusCode
	if status != http.StatusOK {
		// still hand back the response body so an error message can potentially be extracted
		responseBody, _ := responseBodyString(response)
		return responseBody, fmt.Errorf("got unexpected http status %v", status)
	}
	return responseBodyString(response)
}

func responseBodyString(response *http.Response) (string, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	err = response.Body.Close()
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// wrap low level call in circuit breaker

func (r *MailSenderRepositoryImpl) HystrixPerformPost(ctx context.Context, url string, requestBody string) (string, error) {
	output := make(chan string, 1)
	// hystrix.DoC blocks until either completed or error returned
	err := hystrix.DoC(ctx, HystrixCommandName, func(subctx context.Context) error {
		responseBody, innerErr := r.performPost(subctx, url, requestBody)
		output <- responseBody
		// note: if we return an error at this point, it will count towards opening the circuit breaker
		return innerErr
	}, nil)

	// non-blocking receive for optional output
	responseBody := ""
	select {
	case out := <-output:
		responseBody = out
	default:
		// presence of default branch means select will not block even if none of the channels are ready to read from
	}

	return responseBody, err
}

// implementation of repository interface

func (r *MailSenderRepositoryImpl) SendEmail(ctx context.Context, address string, subject string, body string) error {
	requestDto := EmailDto{
		ToAddress: address,
		Subject:   subject,
		Body:      body,
	}
	requestBody, err := renderJson(requestDto)
	if err != nil {
		return err
	}

	responseBody, err := r.performPost(ctx, configuration.MailerServiceUrl() +sendmailEndpoint, requestBody)
	if err != nil {
		errorResponseDto := &apierrors.ErrorDto{}
		err2 := parseJson(responseBody, errorResponseDto)
		if err2 == nil {
			log.Ctx(ctx).Error().Err(err).Msgf("Error sending mail to '%s' via mailer-service: error from response is %s, local error is %s", address, errorResponseDto.Message, err.Error())
		} else {
			log.Ctx(ctx).Error().Err(err).Msgf("Error sending mail to '%s' via mailer-service with no response available: local error is %s", address, err.Error())
		}
		return err
	}

	return nil
}

func renderJson(v interface{}) (string, error) {
	representationBytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(representationBytes), nil
}

// tip: dto := &whatever.WhateverDto{}
func parseJson(body string, dto interface{}) error {
	err := json.Unmarshal([]byte(body), dto)
	return err
}
