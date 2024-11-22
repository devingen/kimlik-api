package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/model"
	"net/http"
)

func (c ServiceController) Setup(ctx context.Context, req core.Request) (*core.Response, error) {

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	var body dto.SetupRequest
	err := req.AssertBody(&body)
	if err != nil {
		return nil, err
	}

	item, err := c.DataService.CreateTenantInfo(ctx, base, &model.TenantInfo{
		Name: &body.TenantName,
	})
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body:       item,
	}, nil
}
