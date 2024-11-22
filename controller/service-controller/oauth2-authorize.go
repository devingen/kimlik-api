package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	"net/http"
)

// OAuth2Authorize handles the OAuth2 authorization process. Creates a
//
//	Redirects user to login if not logged in. Redirects user to authorize page.
func (c ServiceController) OAuth2Authorize(ctx context.Context, req core.Request) (*core.Response, error) {

	return &core.Response{
		StatusCode: http.StatusNotImplemented,
	}, nil

	//base, hasBase := req.PathParameters["base"]
	//if !hasBase {
	//	return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	//}
	//
	//var body dto.OAuth2AuthorizeRequest
	//err := req.AssertBody(&body)
	//if err != nil {
	//	return nil, err
	//}
	//
	//token, err := kimlik.AssertAuthentication(ctx)
	//if err != nil {
	//	return nil, err
	//}
	//
	//// TODO find the OAuth2Client from database with body.ClientID
	//// TODO check if the redirect uri matches the ones in OAuth2Client from database
	//// TODO check if the scopes matches the ones in OAuth2Client from database
	//
	//user, err := c.DataService.FindUserWithId(ctx, base, token.UserID)
	//if err != nil {
	//	return nil, err
	//}
	//
	//code, err := c.TokenService.GenerateAuthorizationCode()
	//if err != nil {
	//	return nil, err
	//}
	//
	//oac := &model.OAuthAccessCode{
	//	CreatedBy:   user.DBRef(base),
	//	Code:        code,
	//	ClientID:    core.String("base.devingen.io"),
	//	RedirectURI: core.String("https://base.devingen.io/oauth/redirect/das"),
	//	Scope:       core.String("openid"),
	//}
	//
	//_, err = c.DataService.CreateOAuthAccessCode(
	//	ctx,
	//	base,
	//	oac,
	//)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return &core.Response{
	//	StatusCode: http.StatusFound,
	//	Headers: map[string]string{
	//		"Location": core.StringValue(body.RedirectURI) + "?code=" + *code + "&state=" + core.StringValue(body.State),
	//	},
	//}, nil
}
