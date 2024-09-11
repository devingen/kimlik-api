package service_controller

import (
	"context"
	"net/http"

	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/model"
)

func (c ServiceController) ActivateUser(ctx context.Context, req core.Request) (*core.Response, error) {

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	var body dto.ActivateUserRequest
	err := req.AssertBody(&body)
	if err != nil {
		return nil, err
	}

	payload, err := c.TokenService.ParseAccessToken(body.UserActivationToken)
	if err != nil {
		return nil, err
	}

	if payload.Scopes[0] != "activate" {
		return nil, core.NewError(http.StatusUnauthorized, "invalid-scope")
	}

	user, err := c.DataService.FindUserWithId(ctx, base, payload.UserId)
	if err != nil {
		return nil, err
	}

	if *user.Status != model.UserStatusNotActivated {
		return nil, core.NewError(http.StatusBadRequest, "user-already-activated")
	}

	auth, err := c.DataService.FindPasswordAuthOfUser(ctx, base, payload.UserId)
	if err != nil {
		return nil, err
	}

	if auth != nil {
		return nil, core.NewError(http.StatusConflict, "user-already-has-password-auth")
	}

	auth, err = c.DataService.CreateAuthWithPassword(ctx, base, body.Password, user)
	if err != nil {
		return nil, err
	}

	status := model.UserStatusActive
	_, _, err = c.DataService.UpdateUser(ctx, base, &model.User{
		ID:              user.ID,
		Status:          &status,
		IsEmailVerified: core.Bool(true), // email is verified because 'userActivationToken' is only sent by email
		FirstName:       core.String(body.FirstName),
		LastName:        core.String(body.LastName),
	})
	if err != nil {
		return nil, err
	}

	jwt, _, err := c.createSuccessfulSessionAndGenerateToken(ctx, req, base, auth, user)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body: dto.LoginResponse{
			UserID: user.ID.Hex(),
			JWT:    jwt,
		},
	}, nil
}
