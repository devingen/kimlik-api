.PHONY: build clean deploy

build:
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/build-saml-auth-url aws/build-saml-auth-url/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/change-password aws/change-password/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/consume-saml-auth-response aws/consume-saml-auth-response/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/create-api-key aws/create-api-key/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/create-saml-config aws/create-saml-config/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/delete-api-key aws/delete-api-key/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/delete-saml-config aws/delete-saml-config/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/find-api-keys aws/find-api-keys/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/find-saml-configs aws/find-saml-configs/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/find-users aws/find-users/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/get-session aws/get-session/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/login-with-email aws/login-with-email/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/register-with-email aws/register-with-email/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/update-api-key aws/update-api-key/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/update-saml-config aws/update-saml-config/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/verify-api-key aws/verify-api-key/main.go

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