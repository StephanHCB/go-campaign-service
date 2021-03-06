package downstreamcall

import (
	"context"
	"github.com/StephanHCB/go-campaign-service/web/middleware/authentication"
	"github.com/StephanHCB/go-campaign-service/web/util/media"
	"github.com/go-chi/chi/middleware"
	"github.com/go-http-utils/headers"
	"io/ioutil"
	"net/http"
	"strings"
)

// performs a http POST, returning the response body and status and passing on the request id if present in the context
func PerformPOST(ctx context.Context, httpClient *http.Client, url string, requestBody string) (string, int, error) {
	return performWithBody(ctx, http.MethodPost, httpClient, url, requestBody)
}

// performs a http PUT, returning the response body and status and passing on the request id if present in the context
func PerformPUT(ctx context.Context, httpClient *http.Client, url string, requestBody string) (string, int, error) {
	return performWithBody(ctx, http.MethodPut, httpClient, url, requestBody)
}

// performs a http GET, returning the response body and status and passing on the request id if present in the context
func PerformGET(ctx context.Context, httpClient *http.Client, url string, requestBody string) (string, int, error) {
	return performNoBody(ctx, http.MethodGet, httpClient, url)
}

// --- internal helper functions ---

func performNoBody(ctx context.Context, method string, httpClient *http.Client, url string) (string, int, error) {
	return performWithBody(ctx, method, httpClient, url, "")
}

func performWithBody(ctx context.Context, method string, httpClient *http.Client, url string, requestBody string) (string, int, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(requestBody))
	if err != nil {
		return "", 0, err
	}

	if requestBody != "" {
		req.Header.Set(headers.ContentType, media.ContentTypeApplicationJson)
	}

	requestId := middleware.GetReqID(ctx)
	if requestId != "" {
		req.Header.Set(middleware.RequestIDHeader, requestId)
	}

	if rawtoken, err := authentication.ExtractRawTokenFromContext(ctx); err == nil {
		req.Header.Set(headers.Authorization, "Bearer " + rawtoken)
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
