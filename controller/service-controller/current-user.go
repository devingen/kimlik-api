package service_controller

import (
	"context"
	"net/http"

	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api"
)

func (c ServiceController) GetCurrentUser(ctx context.Context, req core.Request) (*core.Response, error) {

	base, err := req.AssertPathParameter("base")
	if err != nil {
		return nil, err
	}

	tokenPayload, err := kimlik.AssertAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	user, err := c.DataService.FindUserWithId(ctx, base, tokenPayload.UserID)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body:       user,
	}, nil
}
