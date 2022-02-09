package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/devingen/kimlik-api/aws"
)

func main() {
	fmt.Println("build-saml-auth-url", 1)
	serviceController, wrap := aws.InitDeps()
	fmt.Println("build-saml-auth-url", 2)
	lambda.Start(wrap(serviceController.BuildSAMLAuthURL))
	fmt.Println("build-saml-auth-url", 3)
}
