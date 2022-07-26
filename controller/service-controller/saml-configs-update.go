package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	core_dto "github.com/devingen/api-core/dto"
	"github.com/devingen/kimlik-api"
	"github.com/devingen/kimlik-api/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (c ServiceController) UpdateSAMLConfig(ctx context.Context, req core.Request) (*core.Response, error) {

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

	var body dto.UpdateSAMLConfigRequest
	_, err = kimlik.AssertAuthenticationAndBody(ctx, req, &body)
	if err != nil {
		return nil, err
	}

	body.ID = id
	updatedAt, revision, err := c.DataService.UpdateSAMLConfig(ctx, base, &body)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body: core_dto.UpdateEntryResponse{
			ID:        body.ID.Hex(),
			UpdatedAt: *updatedAt,
			Revision:  revision,
		},
	}, nil
}
