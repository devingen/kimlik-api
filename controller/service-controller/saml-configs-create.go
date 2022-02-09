package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/dto"
	"net/http"
)

func (controller ServiceController) CreateSAMLConfig(ctx context.Context, req core.Request) (interface{}, int, error) {

	//_, interceptorStatusCode, interceptorError := controller.InterceptorService.Pre(ctx, req)
	//if interceptorError != nil {
	//	return interceptorError, interceptorStatusCode, nil
	//}

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, 0, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	var body dto.CreateSAMLConfigRequest
	err := req.AssertBody(&body)
	//_, err := kimlik.AssertAuthenticationAndBody(ctx, req, &body)
	if err != nil {
		return nil, 0, err
	}

	domain, err := controller.DataService.CreateSAMLConfig(ctx, base, &body)
	if err != nil {
		return nil, 0, err
	}

	//controller.InterceptorService.Final(ctx, req, domain)

	return &domain, http.StatusCreated, err
}
