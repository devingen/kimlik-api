package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api"
	"github.com/devingen/kimlik-api/dto"
	"net/http"
)

func (controller ServiceController) GetSession(ctx context.Context, req core.Request) (interface{}, int, error) {

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, 0, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	token, err := kimlik.AssertAuthentication(ctx)
	if err != nil {
		return nil, 0, err
	}

	user, err := controller.DataService.FindUserUserWithId(base, token.UserId)
	if err != nil {
		return nil, 0, err
	}

	if user == nil {
		return nil, 0, core.NewStatusError(http.StatusNotFound)
	}

	return &dto.GetSessionResponse{User: user}, http.StatusOK, err
}
