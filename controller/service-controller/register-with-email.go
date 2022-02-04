package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/dto"
	token_service "github.com/devingen/kimlik-api/token-service"
	"net/http"
)

func (controller ServiceController) RegisterWithEmail(ctx context.Context, req core.Request) (interface{}, int, error) {

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, 0, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	var body dto.RegisterWithEmailRequest
	err := req.AssertBody(&body)
	if err != nil {
		return nil, 0, err
	}

	userWithSameEmail, err := controller.DataService.FindUserUserWithEmail(base, body.Email)
	if err != nil {
		return nil, 0, err
	}

	if userWithSameEmail != nil {
		return nil, 0, core.NewStatusError(http.StatusConflict)
	}

	user, err := controller.DataService.CreateUser(
		base,
		body.FirstName,
		body.LastName,
		body.Email,
	)
	if err != nil {
		return nil, 0, err
	}

	_, err = controller.DataService.CreateAuthWithPassword(base, body.Password, user)
	if err != nil {
		return nil, 0, err
	}

	userAgent := req.Headers["User-Agent"]
	client := req.Headers["Client"]
	ip := req.IP
	session, err := controller.DataService.CreateSession(base, client, userAgent, ip, user)
	if err != nil {
		return nil, 0, err
	}

	jwt, err := controller.TokenService.GenerateToken(
		user.ID.Hex(),
		session.ID.Hex(),
		[]token_service.Scope{ScopeAll},
		240,
	)
	if err != nil {
		return nil, 0, err
	}

	return &dto.RegisterWithEmailResponse{
		UserID: user.ID.Hex(),
		JWT:    jwt,
	}, http.StatusCreated, err
}
