package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/devingen/kimlik-api/aws"
)

func main() {
	serviceController, wrap := aws.InitDeps()
	lambda.Start(wrap(serviceController.OAuthToken))
}
