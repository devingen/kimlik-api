package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/model"
	token_service "github.com/devingen/kimlik-api/token-service"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (controller ServiceController) LoginWithEmail(ctx context.Context, req core.Request) (interface{}, int, error) {

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, 0, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	var body dto.LoginWithEmailRequest
	err := req.AssertBody(&body)
	if err != nil {
		return nil, 0, err
	}

	user, err := controller.DataService.FindUserUserWithEmail(base, body.Email)
	if err != nil {
		return nil, 0, err
	}

	if user == nil {
		return nil, 0, core.NewStatusError(http.StatusNotFound)
	}

	auth, err := controller.DataService.FindAuthOfUser(base, user.ID.Hex(), model.AuthTypePassword)
	if err != nil {
		return nil, 0, err
	}

	if auth == nil {
		return nil, 0, core.NewError(http.StatusInternalServerError, "auth-missing")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(body.Password)); err != nil {
		return nil, 0, core.NewError(http.StatusUnauthorized, "password-mismatch")
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

	return &dto.LoginWithEmailResponse{
		UserID: user.ID.Hex(),
		JWT:    jwt,
	}, http.StatusOK, err
}
