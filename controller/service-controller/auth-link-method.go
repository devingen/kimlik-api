package service_controller

import (
	"context"
	"net/http"

	core "github.com/devingen/api-core"
	kimlik "github.com/devingen/kimlik-api"
	"github.com/devingen/kimlik-api/dto"
)

type LinkAuthenticationMethodRequest struct {
	IDToken string `json:"id_token" validate:"required"`
}

// LinkAuthenticationMethod links an OIDC identity to the authenticated user's account.
func (c ServiceController) LinkAuthenticationMethod(ctx context.Context, req core.Request) (*core.Response, error) {
	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	tokenPayload, err := kimlik.AssertAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	var body LinkAuthenticationMethodRequest
	err = req.AssertBody(&body)
	if err != nil {
		return nil, err
	}

	idToken := IDToken{RawIDToken: body.IDToken}
	if err = idToken.Parse(); err != nil {
		return nil, core.NewError(http.StatusBadRequest, "could-not-parse-token:"+err.Error())
	}
	if err = idToken.Verify(ctx); err != nil {
		return nil, core.NewError(http.StatusBadRequest, "could-not-verify-token:"+err.Error())
	}

	// Reject if another user already owns this OIDC identity
	existing, err := c.DataService.FindOIDCAuthByIssuerAndSubject(ctx, base, idToken.Claims.Issuer, idToken.Claims.Subject)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, core.NewError(http.StatusConflict, "oidc-identity-already-linked-to-another-user")
	}

	// Reject if the current user already has an OIDC auth with the same issuer
	existingOfUser, err := c.DataService.FindOIDCAuthOfUser(ctx, base, tokenPayload.UserID, idToken.Claims.Issuer)
	if err != nil {
		return nil, err
	}
	if existingOfUser != nil {
		return nil, core.NewError(http.StatusConflict, "oidc-identity-already-linked-to-this-user")
	}

	user, err := c.DataService.FindUserWithId(ctx, base, tokenPayload.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, core.NewError(http.StatusNotFound, "user-not-found")
	}

	_, err = c.DataService.CreateAuthWithIDToken(ctx, base, idToken.Claims.ToMap(), user)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body: dto.GetUserInfoResponse{
			Sub:        user.ID.Hex(),
			Name:       user.FullName(),
			GivenName:  *user.FirstName,
			FamilyName: *user.LastName,
			Email:      *user.Email,
		},
	}, nil
}
