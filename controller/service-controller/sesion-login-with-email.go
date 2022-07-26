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

	user, err := c.DataService.FindUserWithEmail(ctx, base, body.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, core.NewStatusError(http.StatusNotFound)
	}

	auth, err := c.DataService.FindAuthOfUser(ctx, base, user.ID.Hex(), model.AuthTypePassword)
	if err != nil {
		return nil, err
	}

	if auth == nil {
		return nil, core.NewError(http.StatusInternalServerError, "auth-missing")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(body.Password)); err != nil {
		return nil, core.NewError(http.StatusUnauthorized, "password-mismatch")
	}

	userAgent := req.Headers["User-Agent"]
	client := req.Headers["Client"]
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
		StatusCode: http.StatusOK,
		Body: dto.LoginWithEmailResponse{
			UserID: user.ID.Hex(),
			JWT:    jwt,
		},
	}, nil
}
