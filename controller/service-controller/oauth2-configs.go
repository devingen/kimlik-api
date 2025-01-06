package service_controller

import (
	"context"
	"net/http"

	core "github.com/devingen/api-core"
	coredto "github.com/devingen/api-core/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/devingen/kimlik-api"
	"github.com/devingen/kimlik-api/model"
)

func (c ServiceController) CreateOAuth2Config(ctx context.Context, req core.Request) (*core.Response, error) {

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

	var body model.OAuth2Config
	err := req.AssertBody(&body)
	if err != nil {
		return nil, err
	}

	item, err := c.DataService.CreateOAuth2Config(ctx, base, &body)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusCreated,
		Body:       item,
	}, nil
}

func (c ServiceController) FindOAuth2Configs(ctx context.Context, req core.Request) (*core.Response, error) {

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

	items, err := c.DataService.FindOAuth2Configs(ctx, base, nil)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body:       coredto.GetListResponse{Results: items},
	}, nil
}

func (c ServiceController) UpdateOAuth2Config(ctx context.Context, req core.Request) (*core.Response, error) {

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

	var body model.OAuth2Config
	_, err = kimlik.AssertAuthenticationAndBody(ctx, req, &body)
	if err != nil {
		return nil, err
	}

	body.ID = id
	updatedAt, revision, err := c.DataService.UpdateOAuth2Config(ctx, base, &body)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body: coredto.UpdateEntryResponse{
			ID:        body.ID.Hex(),
			UpdatedAt: *updatedAt,
			Revision:  revision,
		},
	}, nil
}

func (c ServiceController) DeleteOAuth2Config(ctx context.Context, req core.Request) (*core.Response, error) {

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

	err = c.DataService.DeleteOAuth2Config(ctx, base, id)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusNoContent,
	}, nil
}
