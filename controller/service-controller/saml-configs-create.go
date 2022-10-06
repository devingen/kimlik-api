package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/dto"
	"net/http"
)

func (c ServiceController) CreateSAMLConfig(ctx context.Context, req core.Request) (*core.Response, error) {

	_, interceptorStatusCode, interceptorError := c.InterceptorService.Pre(ctx, req)
	if interceptorError != nil {
		return &core.Response{
			StatusCode: interceptorStatusCode,
			Body:       interceptorError,
		}, nil
	}

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	var body dto.CreateSAMLConfigRequest
	err := req.AssertBody(&body)
	if err != nil {
		return nil, err
	}

	item, err := c.DataService.CreateSAMLConfig(ctx, base, &body)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusCreated,
		Body:       item,
	}, nil
}
