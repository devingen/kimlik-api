package aws

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	coreaws "github.com/devingen/api-core/aws"
	"github.com/devingen/api-core/dvnruntime"
	"github.com/devingen/api-core/util"
	"github.com/devingen/kimlik-api/controller"
	service_controller "github.com/devingen/kimlik-api/controller/service-controller"
	"github.com/devingen/kimlik-api/kimlikruntime"
	"github.com/devingen/kimlik-api/service"
	json_web_token_service "github.com/devingen/kimlik-api/token-service/json-web-token-service"
)

func GenerateController() *service_controller.ServiceController {
	db := GetDatabase()
	databaseService := service.NewDatabaseService(db)
	tokenService := json_web_token_service.NewTokenService()
	return controller.NewServiceController(databaseService, tokenService)
}

func GenerateHandler(f dvnruntime.ControllerFunc) func(ctx context.Context, awsReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	return func(ctx context.Context, awsReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		req := coreaws.AdaptRequest(awsReq)
		requestContext, err := kimlikruntime.BuildContext(ctx, req)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		result, status, err := f(requestContext, req)
		response, err := util.BuildResponse(status, result, err)
		return coreaws.AdaptResponse(response, err)
	}
}
