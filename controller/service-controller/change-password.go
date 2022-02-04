package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	core_dto "github.com/devingen/api-core/dto"
	"github.com/devingen/kimlik-api"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/model"
	"net/http"
)

func (controller ServiceController) ChangePassword(ctx context.Context, req core.Request) (interface{}, int, error) {

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, 0, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	var body dto.ChangePasswordRequest
	token, err := kimlik.AssertAuthenticationAndBody(ctx, req, &body)
	if err != nil {
		return nil, 0, err
	}

	auth, err := controller.DataService.FindAuthOfUser(base, token.UserId, model.AuthTypePassword)
	if err != nil {
		return nil, 0, err
	}

	if auth == nil {
		return nil, 0, core.NewError(http.StatusInternalServerError, "auth-missing")
	}
	auth.Password = body.Password

	updatedAt, revision, err := controller.DataService.UpdateAuth(base, auth)
	if err != nil {
		return nil, 0, err
	}

	return core_dto.UpdateEntryResponse{
		ID:        auth.ID.Hex(),
		UpdatedAt: *updatedAt,
		Revision:  revision,
	}, http.StatusOK, nil
}
