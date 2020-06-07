package mailserviceclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/StephanHCB/go-campaign-service/api/v1/apierrors"
	"github.com/StephanHCB/go-campaign-service/internal/repository/configuration"
	"github.com/StephanHCB/go-campaign-service/internal/repository/mailservice"
	"github.com/StephanHCB/go-campaign-service/web/util/media"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type MailSenderRepositoryImpl struct {
	netClient *http.Client
}

func Create() mailservice.MailSenderRepository {
	return &MailSenderRepositoryImpl{
		netClient: &http.Client{
			Timeout: time.Millisecond * time.Duration(configuration.MailerServiceTimeoutMs()),
		},
	}
}

const sendmailEndpoint = "api/rest/v1/sendmail"

type EmailDto struct {
	ToAddress string `json:"to_address"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

func (r *MailSenderRepositoryImpl) performPost(url string, requestBody string) (string, error) {
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

// TODO wrap low level call in circuit breaker

// implementation of repository interface

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

	responseBody, err := r.performPost(configuration.MailerServiceUrl() +sendmailEndpoint, requestBody)
	if err != nil {
		errorResponseDto := &apierrors.ErrorDto{}
		err2 := parseJson(responseBody, errorResponseDto)
		if err2 == nil {
			log.Ctx(ctx).Error().Err(err).Msgf("Error sending mail to '%s' via mailer-service: %s", address, errorResponseDto.Message)
		}
		return err
	}

	return nil
}
