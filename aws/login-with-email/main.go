package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/devingen/kimlik-api/aws"
)

func main() {
	serviceController := aws.GenerateController()
	lambda.Start(aws.GenerateHandler(serviceController.LoginWithEmail))
}
