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

func (client KimlikAPIClient) OAuth2Token(ctx context.Context, data dto.OAuth2TokenRequest) (*dto.OAuth2TokenResponse, int, error) {

	resp, err := client.Client.R().EnableTrace().
		SetBody(data).
		SetResult(&dto.OAuth2TokenResponse{}).
		SetError(&map[string]interface{}{}).
		Post("/oauth2/token")

	if err != nil {
		switch err.(type) {
		case *url.Error:
			return nil, http.StatusInternalServerError, core.NewError(http.StatusInternalServerError, "kimlik-api-is-unreachable:"+err.Error())
		}
		return nil, http.StatusInternalServerError, err
	}
	if resp.IsError() {
		body := map[string]interface{}{}
		unmErr := json.Unmarshal(resp.Body(), &body)
		if unmErr == nil {
			errorMessage, ok := body["error"].(string)
			if ok {
				return nil, resp.StatusCode(), core.NewError(resp.StatusCode(), errorMessage)
			}
		}
		return nil, resp.StatusCode(), core.NewError(resp.StatusCode(), "kimlik-api-returned-error: "+string(resp.Body()))
	}

	return resp.Result().(*dto.OAuth2TokenResponse), resp.StatusCode(), nil
}

func (client KimlikAPIClient) Authenticate(ctx context.Context, data dto.AuthorizeRequest) (*dto.OAuth2TokenResponse, int, error) {

	resp, err := client.Client.R().EnableTrace().
		SetBody(data).
		SetResult(&dto.OAuth2TokenResponse{}).
		SetError(&map[string]interface{}{}).
		Post("/authenticate")

	if err != nil {
		switch err.(type) {
		case *url.Error:
			return nil, http.StatusInternalServerError, core.NewError(http.StatusInternalServerError, "kimlik-api-is-unreachable:"+err.Error())
		}
		return nil, http.StatusInternalServerError, err
	}
	if resp.IsError() {
		body := map[string]interface{}{}
		unmErr := json.Unmarshal(resp.Body(), &body)
		if unmErr == nil {
			errorMessage, ok := body["error"].(string)
			if ok {
				return nil, resp.StatusCode(), core.NewError(resp.StatusCode(), errorMessage)
			}
		}
		return nil, resp.StatusCode(), core.NewError(resp.StatusCode(), "kimlik-api-returned-error: "+string(resp.Body()))
	}

	return resp.Result().(*dto.OAuth2TokenResponse), resp.StatusCode(), nil
}

func (client KimlikAPIClient) LinkAuthMethod(ctx context.Context, headers map[string]string, data dto.LinkAuthMethodRequest) (*dto.LinkAuthMethodResponse, int, error) {

	resp, err := client.Client.R().EnableTrace().
		SetHeaders(client.mergeHeaders(headers)).
		SetBody(data).
		SetResult(&dto.LinkAuthMethodResponse{}).
		SetError(&map[string]interface{}{}).
		Post("/link-authentication")

	if err != nil {
		switch err.(type) {
		case *url.Error:
			return nil, http.StatusInternalServerError, core.NewError(http.StatusInternalServerError, "kimlik-api-is-unreachable:"+err.Error())
		}
		return nil, http.StatusInternalServerError, err
	}
	if resp.IsError() {
		body := map[string]interface{}{}
		unmErr := json.Unmarshal(resp.Body(), &body)
		if unmErr == nil {
			errorMessage, ok := body["error"].(string)
			if ok {
				return nil, resp.StatusCode(), core.NewError(resp.StatusCode(), errorMessage)
			}
		}
		return nil, resp.StatusCode(), core.NewError(resp.StatusCode(), "kimlik-api-returned-error: "+string(resp.Body()))
	}

	return resp.Result().(*dto.LinkAuthMethodResponse), resp.StatusCode(), nil
}

func (client KimlikAPIClient) FindAuths(ctx context.Context, headers map[string]string) ([]dto.AuthResponse, int, error) {

	resp, err := client.Client.R().EnableTrace().
		SetHeaders(headers).
		SetResult(&[]dto.AuthResponse{}).
		SetError(&map[string]interface{}{}).
		Get("/auths")

	if err != nil {
		switch err.(type) {
		case *url.Error:
			return nil, http.StatusInternalServerError, core.NewError(http.StatusInternalServerError, "kimlik-api-is-unreachable:"+err.Error())
		}
		return nil, http.StatusInternalServerError, err
	}
	if resp.IsError() {
		body := map[string]interface{}{}
		unmErr := json.Unmarshal(resp.Body(), &body)
		if unmErr == nil {
			errorMessage, ok := body["error"].(string)
			if ok {
				return nil, resp.StatusCode(), core.NewError(resp.StatusCode(), errorMessage)
			}
		}
		return nil, resp.StatusCode(), core.NewError(resp.StatusCode(), "kimlik-api-returned-error: "+string(resp.Body()))
	}

	return *resp.Result().(*[]dto.AuthResponse), resp.StatusCode(), nil
}

func (client KimlikAPIClient) DeleteAuthMethod(ctx context.Context, headers map[string]string, id string) error {

	resp, err := client.Client.R().EnableTrace().
		SetHeaders(headers).
		SetError(&map[string]interface{}{}).
		Delete("/auths/" + id)

	if err != nil {
		switch err.(type) {
		case *url.Error:
			return core.NewError(http.StatusInternalServerError, "kimlik-api-is-unreachable:"+err.Error())
		}
		return err
	}
	if resp.IsError() {
		return core.NewError(resp.StatusCode(), "kimlik-api-returned-error: "+string(resp.Body()))
	}
	return nil
}

func (client KimlikAPIClient) UpdateUser(ctx context.Context, headers map[string]string, id string, data dto.UpdateUserRequest) (*dto.GetUserInfoResponse, int, error) {

	resp, err := client.Client.R().EnableTrace().
		SetHeaders(client.mergeHeaders(headers)).
		SetBody(data).
		SetResult(&dto.GetUserInfoResponse{}).
		SetError(&map[string]interface{}{}).
		Put("/users/" + id)

	if err != nil {
		switch err.(type) {
		case *url.Error:
			return nil, http.StatusInternalServerError, core.NewError(http.StatusInternalServerError, "kimlik-api-is-unreachable:"+err.Error())
		}
		return nil, http.StatusInternalServerError, err
	}
	if resp.IsError() {
		body := map[string]interface{}{}
		unmErr := json.Unmarshal(resp.Body(), &body)
		if unmErr == nil {
			errorMessage, ok := body["error"].(string)
			if ok {
				return nil, resp.StatusCode(), core.NewError(resp.StatusCode(), errorMessage)
			}
		}
		return nil, resp.StatusCode(), core.NewError(resp.StatusCode(), "kimlik-api-returned-error: "+string(resp.Body()))
	}

	return resp.Result().(*dto.GetUserInfoResponse), resp.StatusCode(), nil
}

func (client KimlikAPIClient) RegisterWithEmail(ctx context.Context, data dto.RegisterWithEmailRequest) (*dto.RegisterWithEmailResponse, int, error) {

	resp, err := client.Client.R().EnableTrace().
		SetBody(data).
		SetResult(&dto.RegisterWithEmailResponse{}).
		SetError(&map[string]interface{}{}).
		Post("/register")

	if err != nil {
		switch err.(type) {
		case *url.Error:
			return nil, http.StatusInternalServerError, core.NewError(http.StatusInternalServerError, "kimlik-api-is-unreachable:"+err.Error())
		}
		return nil, http.StatusInternalServerError, err
	}
	if resp.IsError() {
		body := map[string]interface{}{}
		unmErr := json.Unmarshal(resp.Body(), &body)
		if unmErr == nil {
			errorMessage, ok := body["error"].(string)
			if ok {
				return nil, resp.StatusCode(), core.NewError(resp.StatusCode(), errorMessage)
			}
		}
		return nil, resp.StatusCode(), core.NewError(resp.StatusCode(), "kimlik-api-returned-error: "+string(resp.Body()))
	}

	return resp.Result().(*dto.RegisterWithEmailResponse), resp.StatusCode(), nil
}

func (client KimlikAPIClient) GetUserInfo(ctx context.Context, headers map[string]string) (*dto.GetUserInfoResponse, error) {

	requestHeaders := map[string]string{
		"authorization": headers["authorization"],
	}
	resp, err := client.Client.R().EnableTrace().
		SetHeaders(requestHeaders).
		SetResult(&dto.GetUserInfoResponse{}).
		SetError(&map[string]interface{}{}).
		Get("/userinfo")

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

	return resp.Result().(*dto.GetUserInfoResponse), nil
}

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

func (client KimlikAPIClient) mergeHeaders(headers map[string]string) map[string]string {

	merged := make(map[string]string, len(headers)+len(client.Client.Header))
	for k, vals := range client.Client.Header {
		if len(vals) > 0 {
			merged[k] = vals[0]
		}
	}
	for k, v := range headers {
		merged[k] = v
	}
	return merged
}
