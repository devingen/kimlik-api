package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/dto"
	"go.mongodb.org/mongo-driver/bson"
)

func PreGetQueryEnhance(controller ServiceController, ctx context.Context, req core.Request) (bson.M, *dto.WebhookPreResponse, int, interface{}) {

	interceptorResponse, interceptorStatusCode, interceptorError := controller.InterceptorService.Pre(ctx, req)
	if interceptorError != nil {
		return nil, interceptorResponse, interceptorStatusCode, interceptorError
	}

	if interceptorResponse == nil {
		return bson.M{}, nil, 0, nil
	}

	query := bson.M{}
	if interceptorResponse.QueryEnhance != nil {
		if interceptorResponse.QueryEnhance.IDsIn != nil {
			query["_id"] = bson.M{"$in": interceptorResponse.QueryEnhance.IDsIn}
		}
	}

	return query, interceptorResponse, interceptorStatusCode, interceptorError
}
