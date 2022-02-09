.PHONY: build clean deploy

build:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/create-saml-config aws/create-saml-config/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/build-saml-auth-url aws/build-saml-auth-url/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/consume-saml-auth-response aws/consume-saml-auth-response/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/register-with-email aws/register-with-email/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/login-with-email aws/login-with-email/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/get-session aws/get-session/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy-devingen: clean build
	serverless deploy --stage prod --region eu-central-1 --verbose

teardown-devingen: clean
	serverless remove --stage prod --region eu-central-1 --verbose

deploy-devingen-dev: clean build
	serverless deploy --stage dev --region ca-central-1 --verbose

teardown-devingen-dev: clean
	serverless remove --stage dev --region ca-central-1 --verbose