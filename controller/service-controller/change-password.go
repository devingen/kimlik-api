package service_controller

import (
	"context"
	core_dto "github.com/devingen/api-core/dto"
	"github.com/devingen/api-core/dvnruntime"
	coremodel "github.com/devingen/api-core/model"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/kimlikruntime"
	"github.com/devingen/kimlik-api/model"
	"net/http"
)

func (controller ServiceController) ChangePassword(ctx context.Context, req dvnruntime.Request) (interface{}, int, error) {
	var body dto.ChangePasswordRequest
	base, token, err := kimlikruntime.AssertAuthenticationAndBody(ctx, req, &body)
	if err != nil {
		return nil, 0, err
	}

	auth, err := controller.Service.FindAuthOfUser(base, token.UserId, model.AuthTypePassword)
	if err != nil {
		return nil, 0, err
	}

	if auth == nil {
		return nil, 0, coremodel.NewError(http.StatusInternalServerError, "auth-missing")
	}
	auth.Password = body.Password

	updatedAt, revision, err := controller.Service.UpdateAuth(base, auth)
	if err != nil {
		return nil, 0, err
	}

	return core_dto.UpdateEntryResponse{
		ID:        auth.ID.Hex(),
		UpdatedAt: *updatedAt,
		Revision:  revision,
	}, http.StatusOK, nil
}
