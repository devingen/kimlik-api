package service_controller

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"strings"

	core "github.com/devingen/api-core"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/devingen/kimlik-api/model"
)

func (c ServiceController) GetAuthorizationURL(ctx context.Context, req core.Request) (*core.Response, error) {
	protocol, err := req.AssertQueryStringParameter("protocol")
	if err != nil {
		return nil, err
	}

	authorizationURL := ""
	switch protocol {
	case "oauth2":
		authorizationURL, err = c.generateOAuth2AuthorizationURL(ctx, req)
	default:
		return nil, core.NewError(http.StatusBadRequest, "invalid-protocol")
	}

	if redirect, _ := req.GetQueryStringParameter("redirect"); redirect == "false" {
		// Return the authorization url if the redirection is disabled.
		return &core.Response{
			StatusCode: http.StatusOK,
			Body:       map[string]string{"authorization_url": authorizationURL},
		}, nil
	}

	// Redirect the request by default.
	return &core.Response{
		StatusCode: http.StatusFound,
		Headers: map[string]string{
			"Location": authorizationURL,
		},
	}, nil
}

func (c ServiceController) generateOAuth2AuthorizationURL(ctx context.Context, req core.Request) (string, error) {
	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return "", core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	clientID, err := req.AssertQueryStringParameter("clientId")
	if err != nil {
		return "", err
	}

	configs, err := c.DataService.FindOAuth2Configs(ctx, base, bson.M{"clientId": clientID})
	if err != nil {
		return "", err
	}

	if configs == nil || len(configs) == 0 {
		return "", core.NewError(http.StatusNotFound, "oauth2-config-not-found")
	}
	config := configs[0]

	tenantInfo, err := c.DataService.GetTenantInfo(ctx, base)
	if err != nil {
		return "", err
	}

	if tenantInfo == nil {
		return "", core.NewError(http.StatusNotFound, "tenant-info-not-initialised")
	}

	if tenantInfo.OAuth2RedirectionURL == nil {
		return "", core.NewError(http.StatusNotAcceptable, "tenant-info-oauth2-redirection-url-missing")
	}

	redirectURI, hasRedirectURI := req.GetQueryStringParameter("redirectURI")
	if !hasRedirectURI {
		redirectURI = *tenantInfo.OAuth2RedirectionURL
	}

	state, err := GenerateRandomString(16)
	if err != nil {
		return "", err
	}

	codeVerifier, codeChallenge, err := generateCodeVerifierAndChallenge()
	if err != nil {
		return "", err
	}

	authorizationURL, err := url.Parse(*config.AuthorizationEndpoint)
	if err != nil {
		return "", err
	}
	authorizationURLQuery := authorizationURL.Query()
	authorizationURLQuery.Add("scope", strings.Join(config.Scopes, " "))
	authorizationURLQuery.Add("response_type", "code")
	authorizationURLQuery.Add("client_id", clientID)
	authorizationURLQuery.Add("redirect_uri", redirectURI)
	authorizationURLQuery.Add("state", state)
	authorizationURLQuery.Add("code_challenge", codeChallenge)
	authorizationURLQuery.Add("code_challenge_method", "S256")
	authorizationURL.RawQuery = authorizationURLQuery.Encode()

	_, err = c.DataService.CreateOAuth2AuthenticationRequest(ctx, base, &model.OAuth2AuthenticationRequest{
		State:        core.String(state),
		ClientID:     core.String(clientID),
		CodeVerifier: core.String(codeVerifier),
	})
	if err != nil {
		return "", err
	}
	return authorizationURL.String(), nil
}

// Function to generate a random code verifier (string of 43 to 128 characters)
func generateCodeVerifierAndChallenge() (string, string, error) {
	codeVerifier, err := GenerateRandomString(8)
	if err != nil {
		return "", "", err
	}
	sha2 := sha256.New()
	_, err = io.WriteString(sha2, codeVerifier)
	if err != nil {
		return "", "", err
	}
	codeChallenge := base64.RawURLEncoding.EncodeToString(sha2.Sum(nil))
	return codeVerifier, codeChallenge, nil
}
