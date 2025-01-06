package service_controller

import (
	"context"
	"net/http"

	core "github.com/devingen/api-core"
	coredto "github.com/devingen/api-core/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/devingen/kimlik-api"
	"github.com/devingen/kimlik-api/model"
)

func (c ServiceController) CreateAppIntegration(ctx context.Context, req core.Request) (*core.Response, error) {

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

	var body model.AppIntegration
	err := req.AssertBody(&body)
	if err != nil {
		return nil, err
	}

	item, err := c.DataService.CreateAppIntegration(ctx, base, &body)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusCreated,
		Body:       item,
	}, nil
}

func (c ServiceController) FindAppIntegrations(ctx context.Context, req core.Request) (*core.Response, error) {

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

	query := bson.M{}
	if clientID, hasClientID := req.GetQueryStringParameter("clientId"); hasClientID {
		query["clientId"] = clientID
	}

	items, err := c.DataService.FindAppIntegrations(ctx, base, query)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body:       coredto.GetListResponse{Results: items},
	}, nil
}

func (c ServiceController) UpdateAppIntegration(ctx context.Context, req core.Request) (*core.Response, error) {

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

	var body model.AppIntegration
	_, err = kimlik.AssertAuthenticationAndBody(ctx, req, &body)
	if err != nil {
		return nil, err
	}

	body.ID = id
	updatedAt, revision, err := c.DataService.UpdateAppIntegration(ctx, base, &body)
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

func (c ServiceController) DeleteAppIntegration(ctx context.Context, req core.Request) (*core.Response, error) {

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

	err = c.DataService.DeleteAppIntegration(ctx, base, id)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusNoContent,
	}, nil
}
