package acceptance

import (
	"encoding/json"
	"errors"
	"github.com/StephanHCB/go-campaign-service/api/v1/campaign"
	"github.com/StephanHCB/go-campaign-service/web/util/media"
	"github.com/go-http-utils/headers"
	"io/ioutil"
	"net/http"
	"strings"
)

const tstEmailAddress = "test@mailinator.com"

// placing these here because they are package global

type tstWebResponse struct {
	status      int
	body        string
	contentType string
	location    string
}

func tstWebResponseFromResponse(response *http.Response) (tstWebResponse, error) {
	status := response.StatusCode
	ct := ""
	if val, ok := response.Header[headers.ContentType]; ok {
		ct = val[0]
	}
	loc := ""
	if val, ok := response.Header[headers.Location]; ok {
		loc = val[0]
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return tstWebResponse{}, err
	}
	err = response.Body.Close()
	if err != nil {
		return tstWebResponse{}, err
	}
	return tstWebResponse{
		status:      status,
		body:        string(body),
		contentType: ct,
		location:    loc,
	}, nil
}

func tstPerformGet(relativeUrlWithLeadingSlash string, bearerToken string) (tstWebResponse, error) {
	if ts == nil {
		return tstWebResponse{}, errors.New("test web server was not initialized")
	}
	request, err := http.NewRequest(http.MethodGet, ts.URL+relativeUrlWithLeadingSlash, nil)
	if err != nil {
		return tstWebResponse{}, err
	}
	if bearerToken != "" {
		request.Header.Set(headers.Authorization, "Bearer "+bearerToken)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return tstWebResponse{}, err
	}
	return tstWebResponseFromResponse(response)
}

func tstPerformPut(relativeUrlWithLeadingSlash string, requestBody string, bearerToken string) (tstWebResponse, error) {
	request, err := http.NewRequest(http.MethodPut, ts.URL+relativeUrlWithLeadingSlash, strings.NewReader(requestBody))
	if err != nil {
		return tstWebResponse{}, err
	}
	if bearerToken != "" {
		request.Header.Set(headers.Authorization, "Bearer "+bearerToken)
	}
	request.Header.Set(headers.ContentType, media.ContentTypeApplicationJson)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return tstWebResponse{}, err
	}
	return tstWebResponseFromResponse(response)
}

func tstPerformPost(relativeUrlWithLeadingSlash string, requestBody string, bearerToken string) (tstWebResponse, error) {
	request, err := http.NewRequest(http.MethodPost, ts.URL+relativeUrlWithLeadingSlash, strings.NewReader(requestBody))
	if err != nil {
		return tstWebResponse{}, err
	}
	if bearerToken != "" {
		request.Header.Set(headers.Authorization, "Bearer "+bearerToken)
	}
	request.Header.Set(headers.ContentType, media.ContentTypeApplicationJson)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return tstWebResponse{}, err
	}
	return tstWebResponseFromResponse(response)
}

func tstReadCampaign(location string) (campaign.CampaignDto, error) {
	readResponse, err := tstPerformGet(location, tstAuthAdmin())
	if err != nil {
		return campaign.CampaignDto{}, err
	}
	campaignRead := campaign.CampaignDto{}
	err = tstParseJson(readResponse.body, &campaignRead)
	return campaignRead, err
}

func tstBuildValidCampaign(testcase string) campaign.CampaignDto {
	return campaign.CampaignDto{
		Subject: "Subject for " + testcase,
		Body: "Body for " + testcase + "\n\nThis is some more text in the body\n",
		Recipients: []campaign.RecipientDto{
			{ToAddress: tstEmailAddress},
		},
	}
}

func tstRenderJson(v interface{}) string {
	representationBytes, err := json.Marshal(v)
	if err != nil {
		return "ERROR - returning invalid json by intention {"
	}
	return string(representationBytes)
}

// tip: dto := &whatever.WhateverDto{}
func tstParseJson(body string, dto interface{}) error {
	err := json.Unmarshal([]byte(body), dto)
	return err
}
