package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (c ServiceController) DeleteAPIKey(ctx context.Context, req core.Request) (*core.Response, error) {

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

	id, err := primitive.ObjectIDFromHex(req.PathParameters["id"])
	if err != nil {
		return nil, err
	}

	err = c.DataService.DeleteAPIKey(ctx, base, id)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusNoContent,
	}, nil
}
