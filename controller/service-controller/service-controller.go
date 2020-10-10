package service_controller

import (
	"github.com/devingen/kimlik-api/service"
	token_service "github.com/devingen/kimlik-api/token-service"
)

const (
	ScopeAll token_service.Scope = "all"
)

// ServiceController implements IServiceController interface by using KimlikService
type ServiceController struct {
	Service      service.KimlikService
	TokenService token_service.ITokenService
}

type TokenData struct {
	Version   string `bson:"ver,omitempty"`
	Expires   string `bson:"exp,omitempty"`
	UserId    string `bson:"userId,omitempty"`
	SessionId string `bson:"sessionId,omitempty"`
	Scopes    string `bson:"scopes,omitempty"`
}
