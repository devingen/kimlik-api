package service_controller

import (
	"context"
	"net/http"

	core "github.com/devingen/api-core"
	core_dto "github.com/devingen/api-core/dto"
)

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=6,max=24"`
}

func (c ServiceController) ResetPassword(ctx context.Context, req core.Request) (*core.Response, error) {

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	var body ResetPasswordRequest
	err := req.AssertBody(&body)
	if err != nil {
		return nil, err
	}

	// Validate token
	token, err := c.TokenService.ParseAccessToken(body.Token)
	if err != nil {
		return nil, core.NewError(http.StatusUnauthorized, "invalid-token")
	}

	// Verify token has "set-password" scope
	if !token.ContainsScope("set-password") {
		return nil, core.NewError(http.StatusForbidden, "insufficient-scope")
	}

	// Find user's password auth
	auth, err := c.DataService.FindPasswordAuthOfUser(ctx, base, token.UserID)
	if err != nil {
		return nil, err
	}

	if auth == nil {
		return nil, core.NewError(http.StatusInternalServerError, "auth-missing")
	}

	// Update password
	auth.Password = body.NewPassword

	updatedAt, revision, err := c.DataService.UpdateAuth(ctx, base, auth)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body: core_dto.UpdateEntryResponse{
			ID:        auth.ID.Hex(),
			UpdatedAt: *updatedAt,
			Revision:  revision,
		},
	}, nil
}
