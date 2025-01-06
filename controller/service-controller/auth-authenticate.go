package service_controller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	core "github.com/devingen/api-core"
	"github.com/go-resty/resty/v2"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/devingen/kimlik-api/dto"
)

// Authenticate authenticates user and returns access token with given auth type.
func (c ServiceController) Authenticate(ctx context.Context, req core.Request) (*core.Response, error) {
	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	var body dto.AuthorizeRequest
	err := req.AssertBody(&body)
	if err != nil {
		return nil, err
	}

	switch body.AuthType {
	case dto.AuthorizationTypePassword:
		return nil, core.NewError(http.StatusNotImplemented, "not-implemented")
	case dto.AuthorizationTypeCode:
		return c.handleAuthTypeCode(ctx, req, base, body)
	case dto.AuthorizationTypeOIDC:
		return nil, core.NewError(http.StatusNotImplemented, "not-implemented")
	}

	return nil, core.NewError(http.StatusBadRequest, "invalid-grant-type")
}

// Exchanges the oauth2 authorization code with access token from the IdP for OAuth2 authorization process.
func (c ServiceController) handleAuthTypeCode(ctx context.Context, req core.Request, base string, body dto.AuthorizeRequest) (*core.Response, error) {

	// Find authentication request that's created for authentication URL generation that forwards the user
	// to IdP. The state value is unique per auth request. It's added to the authentication URL, it's sent back
	// by IdP and used here to find the details of the related authentication request.
	authenticationRequests, err := c.DataService.FindOAuth2AuthenticationRequests(ctx, base, bson.M{"state": body.State})
	if err != nil {
		return nil, err
	}

	if len(authenticationRequests) == 0 {
		return nil, core.NewError(http.StatusNotFound, "request-not-found")
	}
	authenticationRequest := authenticationRequests[0]

	// Find the OAuth2 config used for initiating authorization process.
	oAuth2Configs, err := c.DataService.FindOAuth2Configs(ctx, base, bson.M{"clientId": authenticationRequest.ClientID})
	if err != nil {
		return nil, err
	}
	if len(oAuth2Configs) == 0 {
		return nil, core.NewError(http.StatusNotFound, "oauth2-config-not-found")
	}
	oAuth2Config := oAuth2Configs[0]

	// Get the tenant info to send the original redirection URL to IdP.
	tenantInfo, err := c.DataService.GetTenantInfo(ctx, base)
	if err != nil {
		return nil, err
	}

	if tenantInfo == nil {
		return nil, core.NewError(http.StatusNotFound, "tenant-info-not-initialised")
	}

	if tenantInfo.OAuth2RedirectionURL == nil {
		return nil, core.NewError(http.StatusNotAcceptable, "tenant-info-oauth2-redirection-url-missing")
	}

	tokenResponse, err := exchangeToken(*oAuth2Config.TokenEndpoint, dto.OAuth2TokenRequest{
		GrantType:    dto.OAuth2GrantTypeAuthorizationCode,
		ClientID:     oAuth2Config.ClientID,
		RedirectURI:  tenantInfo.OAuth2RedirectionURL,
		Code:         body.Code,
		CodeVerifier: authenticationRequest.CodeVerifier,
	})
	if err != nil {
		return nil, err
	}

	auth, user, isNewUser, err := c.validateSessionWithIDToken(ctx, base, tokenResponse.IDToken, "", "")
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := c.createSuccessfulSessionAndGenerateToken(ctx, req, base, auth, user)
	if err != nil {
		return nil, err
	}

	responseStatusCode := http.StatusOK
	if isNewUser {
		responseStatusCode = http.StatusCreated
	}

	return &core.Response{
		StatusCode: responseStatusCode,
		Body: dto.OAuth2TokenResponse{
			AccessToken:  accessToken,
			TokenType:    "Bearer",
			ExpiresIn:    AccessTokenExpirationTime.Seconds(),
			RefreshToken: refreshToken,
		},
	}, nil
}

func exchangeToken(tokenEndpoint string, requestBody dto.OAuth2TokenRequest) (*dto.OAuth2TokenResponse, error) {

	resp, err := resty.New().
		SetBaseURL(tokenEndpoint).
		SetHeader("Content-Type", "application/json").R().EnableTrace().
		SetBody(requestBody).
		SetResult(&dto.OAuth2TokenResponse{}).
		SetError(&map[string]interface{}{}).
		Post("")

	if err != nil {
		switch err.(type) {
		case *url.Error:
			return nil, core.NewError(http.StatusInternalServerError, "token-api-is-unreachable:"+err.Error())
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

	return resp.Result().(*dto.OAuth2TokenResponse), nil
}
