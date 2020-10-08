package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	coreaws "github.com/devingen/api-core/aws"
	"github.com/devingen/api-core/util"
	"github.com/devingen/kimlik-api/aws"
	"github.com/devingen/kimlik-api/controller"
	"github.com/devingen/kimlik-api/service"
	"net/http"
)

func main() {

	db := aws.GetDatabase()
	databaseService := service.NewDatabaseService(db)
	serviceController := controller.NewServiceController(databaseService)

	lambda.Start(func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		base := req.PathParameters["base"]
		jwt := req.Headers["Authorization"][len("Bearer "):]

		result, err := serviceController.GetSession(base, jwt)
		response, err := util.BuildResponse(http.StatusOK, result, err)
		return coreaws.AdaptResponse(response, err)
	})
}
