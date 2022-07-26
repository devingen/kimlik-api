package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api"
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
	token, err := kimlik.AssertAuthenticationAndBody(ctx, req, &body)
	if err != nil {
		return nil, err
	}

	user, err := c.DataService.FindUserWithId(ctx, base, token.UserId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, core.NewStatusError(http.StatusNotFound)
	}

	body.CreatedBy = user.DBRef(base)
	item, err := c.DataService.CreateSAMLConfig(ctx, base, &body)
	if err != nil {
		return nil, err
	}

	//controller.InterceptorService.Final(ctx, req, domain)

	return &core.Response{
		StatusCode: http.StatusCreated,
		Body:       item,
	}, nil
}
