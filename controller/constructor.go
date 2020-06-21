package controller

import (
	"github.com/devingen/kimlik-api/controller/service-controller"
	"github.com/devingen/kimlik-api/service"
)

// NewServiceController generates new ServiceController
func NewServiceController(service service.KimlikService) *service_controller.ServiceController {
	return &service_controller.ServiceController{
		Service: service,
	}
}
