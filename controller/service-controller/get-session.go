package service_controller

import (
	coremodel "github.com/devingen/api-core/model"
	"github.com/devingen/kimlik-api/dto"
	"net/http"
)

func (controller ServiceController) GetSession(base, jwt string) (*dto.GetSessionResponse, error) {

	tokenData, tokenError := controller.ParseToken(jwt)
	if tokenError != nil {
		return nil, tokenError
	}

	user, err := controller.Service.FindUserUserWithId(base, tokenData.UserId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, coremodel.NewStatusError(http.StatusNotFound)
	}

	return &dto.GetSessionResponse{User: user}, err
}
