package service_controller

import (
	"context"
	"net/http"
	"slices"
	"strings"

	core "github.com/devingen/api-core"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/devingen/kimlik-api"
	"github.com/devingen/kimlik-api/model"
)

// OAuth2Authorize handles the OAuth2 authorization process.
func (c ServiceController) OAuth2Authorize(ctx context.Context, req core.Request) (*core.Response, error) {
	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	token, err := kimlik.AssertAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	responseType, err := req.AssertQueryStringParameter("response_type")
	if err != nil {
		return nil, err
	}
	if responseType != "code" {
		return nil, core.NewError(http.StatusBadRequest, "invalid-response-type")
	}

	scope, err := req.AssertQueryStringParameter("scope")
	if err != nil {
		return nil, err
	}
	scopes := strings.Split(scope, " ")
	if len(scopes) == 0 {
		return nil, core.NewError(http.StatusBadRequest, "scopes-empty")
	}

	clientID, err := req.AssertQueryStringParameter("client_id")
	if err != nil {
		return nil, err
	}

	redirectURI, err := req.AssertQueryStringParameter("redirect_uri")
	if err != nil {
		return nil, err
	}

	state, err := req.AssertQueryStringParameter("state")
	if err != nil {
		return nil, err
	}

	codeChallenge, _ := req.GetQueryStringParameter("code_challenge")
	codeChallengeMethod, _ := req.GetQueryStringParameter("code_challenge_method")

	apps, err := c.DataService.FindAppIntegrations(ctx, base, bson.M{"clientId": clientID})
	if err != nil {
		return nil, err
	}
	if len(apps) == 0 {
		return nil, core.NewError(http.StatusNotFound, "client-not-found")
	}
	app := apps[0]

	for _, scopeInRequest := range scopes {
		if !slices.Contains(app.OAuth2Config.Scopes, scopeInRequest) {
			return nil, core.NewError(http.StatusBadRequest, "scope-not-allowed: "+scopeInRequest)
		}
	}

	if !slices.Contains(app.OAuth2Config.RedirectURLs, redirectURI) {
		return nil, core.NewError(http.StatusBadRequest, "redirect-url-not-allowed")
	}

	user, err := c.DataService.FindUserWithId(ctx, base, token.UserID)
	if err != nil {
		return nil, err
	}

	code, err := c.TokenService.GenerateAuthorizationCode()
	if err != nil {
		return nil, err
	}

	oac := &model.OAuth2AccessCode{
		CreatedBy:           user.DBRef(base),
		Code:                code,
		ClientID:            core.String(clientID),
		RedirectURI:         core.String(redirectURI),
		Scopes:              scopes,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}

	_, err = c.DataService.CreateOAuth2AccessCode(
		ctx,
		base,
		oac,
	)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusFound,
		Headers: map[string]string{
			"Location": redirectURI + "?code=" + *code + "&state=" + state,
		},
	}, nil
}
