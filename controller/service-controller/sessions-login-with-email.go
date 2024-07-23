package service_controller

import (
	"context"
	"net/http"

	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/dto"
)

func (c ServiceController) LoginWithEmail(ctx context.Context, req core.Request) (*core.Response, error) {

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	var body dto.LoginWithEmailRequest
	err := req.AssertBody(&body)
	if err != nil {
		return nil, err
	}

	auth, user, err := c.validateSessionWithPassword(ctx, base, body.Email, body.Password)
	if err != nil {
		return nil, err
	}

	jwt, err := c.createSuccessfulSessionAndGenerateToken(ctx, req, base, auth, user)
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
