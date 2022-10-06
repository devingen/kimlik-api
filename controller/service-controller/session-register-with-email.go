package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/dto"
	token_service "github.com/devingen/kimlik-api/token-service"
	"net/http"
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
	)
	if err != nil {
		return nil, err
	}

	_, err = c.DataService.CreateAuthWithPassword(ctx, base, body.Password, user)
	if err != nil {
		return nil, err
	}

	userAgent := req.Headers["user-agent"]
	client := req.Headers["client"]
	ip := req.IP
	session, err := c.DataService.CreateSession(ctx, base, client, userAgent, ip, user)
	if err != nil {
		return nil, err
	}

	jwt, err := c.TokenService.GenerateToken(
		user.ID.Hex(),
		session.ID.Hex(),
		[]token_service.Scope{ScopeAll},
		240,
	)
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
