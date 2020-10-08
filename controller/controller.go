package controller

import (
	"github.com/devingen/kimlik-api/dto"
)

type ServiceController interface {
	RegisterWithEmail(base, client, userAgent, ip string, request *dto.RegisterWithEmailRequest) (*dto.RegisterWithEmailResponse, error)
	LoginWithEmail(base, client, userAgent, ip string, request *dto.LoginWithEmailRequest) (*dto.LoginWithEmailResponse, error)
	GetSession(base, userId string) (*dto.GetSessionResponse, error)
}
