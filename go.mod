module github.com/devingen/kimlik-api

go 1.12

//replace github.com/devingen/api-core => ../api-core

require (
	github.com/aws/aws-lambda-go v1.16.0
	github.com/devingen/api-core v0.0.18
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-playground/validator/v10 v10.4.0
	github.com/gorilla/mux v1.7.4
	github.com/stretchr/testify v1.4.0
	go.mongodb.org/mongo-driver v1.3.2
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
)
