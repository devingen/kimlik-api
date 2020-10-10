package service_controller

import (
	"context"
	"github.com/devingen/api-core/dvnruntime"
	coremodel "github.com/devingen/api-core/model"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/kimlikruntime"
	token_service "github.com/devingen/kimlik-api/token-service"
	"net/http"
	"strings"
)

func (controller ServiceController) RegisterWithEmail(ctx context.Context, req dvnruntime.Request) (interface{}, int, error) {

	var body dto.RegisterWithEmailRequest
	base, err := kimlikruntime.AssertBaseAndBody(ctx, req, &body)
	if err != nil {
		return nil, 0, err
	}

	email := strings.TrimSpace(body.Email)
	if !IsValidEmail(email) {
		return nil, 0, coremodel.NewError(http.StatusBadRequest, "invalid-email")
	}

	userWithSameEmail, err := controller.Service.FindUserUserWithEmail(base, email)
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
		email,
	)
	if err != nil {
		return nil, 0, err
	}

	_, err = controller.Service.CreateAuth(base, body.Password, user)
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
		user.ID.String(),
		session.ID.String(),
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
