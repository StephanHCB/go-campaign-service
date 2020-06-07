package downstreamcall

import (
	"context"
	"encoding/json"
	"github.com/StephanHCB/go-campaign-service/web/util/media"
	"github.com/go-chi/chi/middleware"
	"github.com/go-http-utils/headers"
	"io/ioutil"
	"net/http"
	"strings"
)

// helper functions for dealing with json

func RenderJson(dto interface{}) (string, error) {
	representationBytes, err := json.Marshal(dto)
	if err != nil {
		return "", err
	}
	return string(representationBytes), nil
}

// tip: dto := &whatever.WhateverDto{}
func ParseJson(body string, dto interface{}) error {
	err := json.Unmarshal([]byte(body), dto)
	return err
}

// performs a http POST, returning the response body and status and passing on the request id if present in the context
func PerformPOST(ctx context.Context, httpClient *http.Client, url string, requestBody string) (string, int, error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(requestBody))
	if err != nil {
		return "", 0, err
	}

	req.Header.Set(headers.ContentType, media.ContentTypeApplicationJson)

	requestId := middleware.GetReqID(ctx)
	if requestId != "" {
		req.Header.Set(middleware.RequestIDHeader, requestId)
	}

	response, err := httpClient.Do(req)
	if err != nil {
		return "", 0, err
	}

	responseBody, err := responseBodyString(response)
	return responseBody, response.StatusCode, err
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
