package kimlik

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/dto"
	"github.com/go-resty/resty/v2"
)

type KimlikAPIClient struct {
	Address string
	Client  *resty.Client
}

// New generates new KimlikAPIClient
// address: The complete URL of the webhook api.
// headersValue: Key=value list of headers joined by ",". E.g.("X-Api-Key=abc,X-Client=web-hook")
func New(address string, headersValue string) KimlikAPIClient {
	httpClient := resty.New().
		SetBaseURL(address).
		SetHeader("Content-Type", "application/json")

	if headersValue != "" {
		for _, keyAndValue := range strings.Split(headersValue, ",") {
			headerParts := strings.SplitN(keyAndValue, "=", 2)
			httpClient.SetHeader(headerParts[0], headerParts[1])
		}
	}

	return KimlikAPIClient{
		Address: address,
		Client:  httpClient,
	}
}

func (client KimlikAPIClient) OAuthToken(ctx context.Context, data dto.OAuthTokenRequest) (*dto.OAuthTokenResponse, error) {

	resp, err := client.Client.R().EnableTrace().
		SetBody(data).
		SetResult(&dto.OAuthTokenResponse{}).
		SetError(&map[string]interface{}{}).
		Post("/oauth/token")

	if err != nil {
		switch err.(type) {
		case *url.Error:
			return nil, core.NewError(http.StatusInternalServerError, "kimlik-api-is-unreachable:"+err.Error())
		}
		return nil, err
	}
	if resp.IsError() {
		body := map[string]interface{}{}
		unmErr := json.Unmarshal(resp.Body(), &body)
		if unmErr == nil {
			errorMessage, ok := body["error"].(string)
			if ok {
				return nil, core.NewError(resp.StatusCode(), errorMessage)
			}
		}
		return nil, core.NewError(resp.StatusCode(), "kimlik-api-returned-error: "+string(resp.Body()))
	}

	return resp.Result().(*dto.OAuthTokenResponse), nil
}

//func (client KimlikAPIClient) CreateSession(ctx context.Context, data dto.CreateSession) (*dto.LoginResponse, error) {
//
//	resp, err := client.Client.R().EnableTrace().
//		SetBody(data).
//		SetResult(&dto.LoginResponse{}).
//		SetError(&map[string]interface{}{}).
//		Post("/sessions")
//
//	if err != nil {
//		switch err.(type) {
//		case *url.Error:
//			return nil, core.NewError(http.StatusInternalServerError, "auth-api-is-unreachable")
//		}
//		return nil, err
//	}
//	if resp.IsError() {
//		body := map[string]interface{}{}
//		unmErr := json.Unmarshal(resp.Body(), &body)
//		if unmErr == nil {
//			errorMessage, ok := body["error"].(string)
//			if ok {
//				return nil, core.NewError(resp.StatusCode(), errorMessage)
//			}
//		}
//		return nil, core.NewError(resp.StatusCode(), "auth-api-returned-error: "+string(resp.Body()))
//	}
//
//	return resp.Result().(*dto.LoginResponse), nil
//}

func (client KimlikAPIClient) GetSession(ctx context.Context, headers map[string]string) (*dto.GetSessionResponse, error) {

	requestHeaders := map[string]string{
		"authorization": headers["authorization"],
	}
	resp, err := client.Client.R().EnableTrace().
		SetHeaders(requestHeaders).
		SetResult(&dto.GetSessionResponse{}).
		SetError(&map[string]interface{}{}).
		Get("/session")

	if err != nil {
		switch err.(type) {
		case *url.Error:
			return nil, core.NewError(http.StatusInternalServerError, "auth-api-is-unreachable")
		}
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, core.NewError(resp.StatusCode(), "auth-api-returned-error:"+resp.String())
	}

	return resp.Result().(*dto.GetSessionResponse), nil
}

func (client KimlikAPIClient) AnonymizeUser(ctx context.Context, headers map[string]string, id string) error {

	requestHeaders := map[string]string{
		"authorization": headers["authorization"],
	}
	resp, err := client.Client.R().EnableTrace().
		SetHeaders(requestHeaders).
		SetError(&map[string]interface{}{}).
		Post("/users/" + id + "/anonymize")

	if err != nil {
		switch err.(type) {
		case *url.Error:
			return core.NewError(http.StatusInternalServerError, "auth-api-is-unreachable")
		}
		return err
	}
	if resp.IsError() {
		return core.NewError(resp.StatusCode(), "auth-api-returned-error:"+resp.String())
	}

	return nil
}
