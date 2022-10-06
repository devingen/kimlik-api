package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/api-core/dto"
	"net/http"
)

func (c ServiceController) FindAPIKeys(ctx context.Context, req core.Request) (*core.Response, error) {

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	_, interceptorStatusCode, interceptorError := c.InterceptorService.Pre(ctx, req)
	if interceptorError != nil {
		return &core.Response{
			StatusCode: interceptorStatusCode,
			Body:       interceptorError,
		}, nil
	}

	query, _, iStatusCode, iErr := PreGetQueryEnhance(c, ctx, req)
	if iErr != nil {
		return &core.Response{
			StatusCode: iStatusCode,
			Body:       iErr,
		}, nil
	}

	items, err := c.DataService.FindAPIKeys(ctx, base, query)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body:       dto.GetListResponse{Results: items},
	}, nil
}
