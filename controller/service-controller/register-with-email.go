package service_controller

import (
	coremodel "github.com/devingen/api-core/model"
	"github.com/devingen/kimlik-api/dto"
	"net/http"
	"strings"
)

func (controller ServiceController) RegisterWithEmail(base, client, userAgent, ip string, request *dto.RegisterWithEmailRequest) (*dto.RegisterWithEmailResponse, error) {

	email := strings.TrimSpace(request.Email)

	if !IsValidEmail(email) {
		return nil, coremodel.NewError(http.StatusBadRequest, "invalid-email")
	}

	usersWithSameEmail, err := controller.Service.FindUserUserWithEmail(base, email)
	if err != nil {
		return nil, err
	}

	if usersWithSameEmail != nil && len(usersWithSameEmail) > 0 {
		return nil, coremodel.NewStatusError(http.StatusConflict)
	}

	user, err := controller.Service.CreateUser(
		base,
		request.FirstName,
		request.LastName,
		email,
	)
	if err != nil {
		return nil, err
	}

	_, err = controller.Service.CreateAuth(base, request.Password, user)
	if err != nil {
		return nil, err
	}

	session, err := controller.Service.CreateSession(base, client, userAgent, ip, user)
	if err != nil {
		return nil, err
	}

	jwt, err := controller.GenerateToken(user.ID.String(), session.ID.String(), []Scope{ScopeAll}, 240)
	if err != nil {
		return nil, err
	}

	return &dto.RegisterWithEmailResponse{
		UserID: user.ID.Hex(),
		JWT:    jwt,
	}, err
}
