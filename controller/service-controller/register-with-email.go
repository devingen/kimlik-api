package service_controller

import (
	"context"
	"github.com/devingen/api-core/dvnruntime"
	coremodel "github.com/devingen/api-core/model"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/kimlikruntime"
	token_service "github.com/devingen/kimlik-api/token-service"
	"net/http"
)

func (controller ServiceController) RegisterWithEmail(ctx context.Context, req dvnruntime.Request) (interface{}, int, error) {

	var body dto.RegisterWithEmailRequest
	base, err := kimlikruntime.AssertBaseAndBody(ctx, req, &body)
	if err != nil {
		return nil, 0, err
	}

	userWithSameEmail, err := controller.Service.FindUserUserWithEmail(base, body.Email)
	if err != nil {
		return nil, 0, err
	}

	if userWithSameEmail != nil {
		return nil, 0, coremodel.NewStatusError(http.StatusConflict)
	}

	user, err := controller.Service.CreateUser(
		base,
		body.FirstName,
		body.LastName,
		body.Email,
	)
	if err != nil {
		return nil, 0, err
	}

	_, err = controller.Service.CreateAuthWithPassword(base, body.Password, user)
	if err != nil {
		return nil, 0, err
	}

	userAgent := req.Headers["User-Agent"]
	client := req.Headers["Client"]
	ip := req.IP
	session, err := controller.Service.CreateSession(base, client, userAgent, ip, user)
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
