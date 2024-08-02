package service_controller

import (
	"context"
	"net/http"

	core "github.com/devingen/api-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c ServiceController) AnonymizeUser(ctx context.Context, req core.Request) (*core.Response, error) {

	base, err := req.AssertPathParameter("base")
	if err != nil {
		return nil, err
	}

	id, err := req.AssertPathParameter("id")
	if err != nil {
		return nil, err
	}

	_, interceptorStatusCode, interceptorError := c.InterceptorService.Pre(ctx, req)
	if interceptorError != nil {
		return &core.Response{
			StatusCode: interceptorStatusCode,
			Body:       interceptorError,
		}, nil
	}

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = c.DataService.AnonymizeUser(ctx, base, userID)
	if err != nil {
		return nil, err
	}

	auths, err := c.DataService.FindAuthsOfUser(ctx, base, id)
	if err != nil {
		return nil, err
	}

	for _, auth := range auths {
		err = c.DataService.DeleteAuth(ctx, base, auth.ID)
		if err != nil {
			return nil, err
		}
	}

	return &core.Response{
		StatusCode: http.StatusOK,
	}, nil
}
