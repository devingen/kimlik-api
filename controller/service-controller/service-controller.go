package service_controller

import (
	"github.com/devingen/kimlik-api/controller"
	ds "github.com/devingen/kimlik-api/data-service"
	is "github.com/devingen/kimlik-api/interceptor-service"
	token_service "github.com/devingen/kimlik-api/token-service"
)

const (
	ScopeAll token_service.Scope = "all"
)

// ServiceController implements IServiceController interface by using KimlikService
type ServiceController struct {
	DataService        ds.IKimlikDataService
	TokenService       token_service.ITokenService
	InterceptorService is.IKimlikInterceptorService
}

// New generates new ServiceController
func New(dataService ds.IKimlikDataService, tokenService token_service.ITokenService, interceptorService is.IKimlikInterceptorService) controller.IServiceController {
	return ServiceController{
		DataService:        dataService,
		InterceptorService: interceptorService,
		TokenService:       tokenService,
	}
}

type TokenData struct {
	Version   string `bson:"ver,omitempty"`
	Expires   string `bson:"exp,omitempty"`
	UserId    string `bson:"userId,omitempty"`
	SessionId string `bson:"sessionId,omitempty"`
	Scopes    string `bson:"scopes,omitempty"`
}
