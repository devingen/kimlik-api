module github.com/devingen/kimlik-api

go 1.12

//replace github.com/devingen/api-core => ../api-core

require (
	github.com/aws/aws-lambda-go v1.40.0
	github.com/devingen/api-core v0.0.27
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-playground/validator/v10 v10.4.1
	github.com/go-resty/resty/v2 v2.7.0
	github.com/gorilla/mux v1.7.4
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/russellhaering/gosaml2 v0.6.0
	github.com/russellhaering/goxmldsig v1.1.1
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.9.0
	go.mongodb.org/mongo-driver v1.16.0
	golang.org/x/crypto v0.25.0
	google.golang.org/api v0.188.0 // indirect
)
