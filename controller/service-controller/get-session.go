package service_controller

import (
	"context"
	"github.com/devingen/api-core/dvnruntime"
	coremodel "github.com/devingen/api-core/model"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/kimlikruntime"
	"net/http"
)

func (controller ServiceController) GetSession(ctx context.Context, req dvnruntime.Request) (interface{}, int, error) {

	base, token, err := kimlikruntime.AssertAuthentication(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	user, err := controller.Service.FindUserUserWithId(base, token.UserId)
	if err != nil {
		return nil, 0, err
	}

	if user == nil {
		return nil, 0, coremodel.NewStatusError(http.StatusNotFound)
	}

	return &dto.GetSessionResponse{User: user}, http.StatusOK, err
}
