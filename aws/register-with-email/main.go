package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	coreaws "github.com/devingen/api-core/aws"
	"github.com/devingen/api-core/util"
	"github.com/devingen/kimlik-api/aws"
	"github.com/devingen/kimlik-api/controller"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/service"
	"net/http"
)

func main() {

	db := aws.GetDatabase()
	databaseService := service.NewDatabaseService(db)
	serviceController := controller.NewServiceController(databaseService)

	lambda.Start(func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

		var body dto.RegisterWithEmailRequest
		err := json.Unmarshal([]byte(req.Body), &body)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		ip := req.RequestContext.Identity.SourceIP
		userAgent := req.Headers["User-Agent"]
		client := req.Headers["Client"]
		base := req.PathParameters["base"]

		result, err := serviceController.RegisterWithEmail(base, client, userAgent, ip, &body)
		response, err := util.BuildResponse(http.StatusOK, result, err)
		return coreaws.AdaptResponse(response, err)
	})
}
