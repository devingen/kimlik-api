package service_controller

import (
	"context"
	"net/http"

	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/model"
)

func (c ServiceController) RegisterWithEmail(ctx context.Context, req core.Request) (*core.Response, error) {

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	var body dto.RegisterWithEmailRequest
	err := req.AssertBody(&body)
	if err != nil {
		return nil, err
	}

	userWithSameEmail, err := c.DataService.FindUserWithEmail(ctx, base, body.Email)
	if err != nil {
		return nil, err
	}

	if userWithSameEmail != nil {
		return nil, core.NewStatusError(http.StatusConflict)
	}

	user, err := c.DataService.CreateUser(
		ctx,
		base,
		body.FirstName,
		body.LastName,
		body.Email,
		model.UserStatusActive,
		false,
	)
	if err != nil {
		return nil, err
	}

	auth, err := c.DataService.CreateAuthWithPassword(ctx, base, body.Password, user)
	if err != nil {
		return nil, err
	}

	jwt, _, err := c.createSuccessfulSessionAndGenerateToken(ctx, req, base, auth, user)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusCreated,
		Body: dto.RegisterWithEmailResponse{
			UserID: user.ID.Hex(),
			JWT:    jwt,
		},
	}, nil
}
