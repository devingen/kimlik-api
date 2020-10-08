package service_controller

import (
	coremodel "github.com/devingen/api-core/model"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/model"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func (controller ServiceController) LoginWithEmail(base, client, userAgent, ip string, request *dto.LoginWithEmailRequest) (*dto.LoginWithEmailResponse, error) {

	email := strings.TrimSpace(request.Email)

	if !IsValidEmail(email) {
		return nil, coremodel.NewError(http.StatusBadRequest, "invalid-email")
	}

	user, err := controller.Service.FindUserUserWithEmail(base, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, coremodel.NewStatusError(http.StatusNotFound)
	}

	auth, err := controller.Service.FindAuthOfUser(base, user, model.AuthTypePassword)
	if err != nil {
		return nil, err
	}

	if auth == nil {
		return nil, coremodel.NewStatusError(http.StatusInternalServerError)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(request.Password)); err != nil {
		return nil, coremodel.NewStatusError(http.StatusUnauthorized)
	}

	session, err := controller.Service.CreateSession(base, client, userAgent, ip, user)
	if err != nil {
		return nil, err
	}

	jwt, err := controller.GenerateToken(user.ID.Hex(), session.ID.Hex(), []Scope{ScopeAll}, 240)
	if err != nil {
		return nil, err
	}

	return &dto.LoginWithEmailResponse{
		UserID: user.ID.Hex(),
		JWT:    jwt,
	}, err
}
