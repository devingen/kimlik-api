package controller

import (
	"github.com/devingen/kimlik-api/controller/service-controller"
	"github.com/devingen/kimlik-api/service"
	"github.com/devingen/kimlik-api/token-service"
)

// NewServiceController generates new ServiceController
func NewServiceController(service service.KimlikService, tokenService token_service.ITokenService) *service_controller.ServiceController {
	return &service_controller.ServiceController{
		Service:      service,
		TokenService: tokenService,
	}
}
