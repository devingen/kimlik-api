package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api"
	"github.com/devingen/kimlik-api/dto"
	"net/http"
)

func (c ServiceController) GetSession(ctx context.Context, req core.Request) (*core.Response, error) {

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	token, err := kimlik.AssertAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	session, err := c.DataService.FindSessionWithId(ctx, base, token.SessionID)
	if err != nil {
		return nil, err
	}

	user, err := c.DataService.FindUserWithId(ctx, base, token.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, core.NewStatusError(http.StatusNotFound)
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body:       dto.GetSessionResponse{Session: session, User: user},
	}, nil
}
