REPO_NAME := devingen/kimlik-api
IMAGE_TAG := 0.1.0

.PHONY: build-docker
build-docker:
	@echo "Building Docker image"
	export GO111MODULE=on
	docker buildx build --platform linux/amd64 -t $(REPO_NAME):$(IMAGE_TAG) --build-arg GIT_TOKEN=$(GIT_TOKEN) .

.PHONY: push-docker
push-docker:
	@echo "Pushing Docker image"
	docker push $(REPO_NAME):$(IMAGE_TAG)

.PHONY: release-docker
release-docker: build-docker push-docker



.PHONY: build clean deploy

build:
	export GO111MODULE=on
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/oauth-token/bootstrap aws/oauth-token/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/anonymize-user/bootstrap aws/anonymize-user/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/activate-user/bootstrap aws/activate-user/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/create-session/bootstrap aws/create-session/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/get-session/bootstrap aws/get-session/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/build-saml-auth-url/bootstrap aws/build-saml-auth-url/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/change-password/bootstrap aws/change-password/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/consume-saml-auth-response/bootstrap aws/consume-saml-auth-response/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/create-api-key/bootstrap aws/create-api-key/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/create-saml-config/bootstrap aws/create-saml-config/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/delete-api-key/bootstrap aws/delete-api-key/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/delete-saml-config/bootstrap aws/delete-saml-config/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/find-api-keys/bootstrap aws/find-api-keys/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/find-saml-configs/bootstrap aws/find-saml-configs/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/find-users/bootstrap aws/find-users/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/login-with-email/bootstrap aws/login-with-email/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/login-with-saml/bootstrap aws/login-with-saml/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/register-with-email/bootstrap aws/register-with-email/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/update-api-key/bootstrap aws/update-api-key/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/update-saml-config/bootstrap aws/update-saml-config/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/verify-api-key/bootstrap aws/verify-api-key/main.go

zip:
	zip -j bin/oauth-token.zip bin/oauth-token/bootstrap
	zip -j bin/anonymize-user.zip bin/anonymize-user/bootstrap
	zip -j bin/activate-user.zip bin/activate-user/bootstrap
	zip -j bin/create-session.zip bin/create-session/bootstrap
	zip -j bin/get-session.zip bin/get-session/bootstrap
	zip -j bin/build-saml-auth-url.zip bin/build-saml-auth-url/bootstrap
	zip -j bin/change-password.zip bin/change-password/bootstrap
	zip -j bin/consume-saml-auth-response.zip bin/consume-saml-auth-response/bootstrap
	zip -j bin/create-api-key.zip bin/create-api-key/bootstrap
	zip -j bin/create-saml-config.zip bin/create-saml-config/bootstrap
	zip -j bin/delete-api-key.zip bin/delete-api-key/bootstrap
	zip -j bin/delete-saml-config.zip bin/delete-saml-config/bootstrap
	zip -j bin/find-api-keys.zip bin/find-api-keys/bootstrap
	zip -j bin/find-saml-configs.zip bin/find-saml-configs/bootstrap
	zip -j bin/find-users.zip bin/find-users/bootstrap
	zip -j bin/login-with-email.zip bin/login-with-email/bootstrap
	zip -j bin/login-with-saml.zip bin/login-with-saml/bootstrap
	zip -j bin/register-with-email.zip bin/register-with-email/bootstrap
	zip -j bin/update-api-key.zip bin/update-api-key/bootstrap
	zip -j bin/update-saml-config.zip bin/update-saml-config/bootstrap
	zip -j bin/verify-api-key.zip bin/verify-api-key/bootstrap

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy-devingen: clean build zip
	serverless deploy --stage prod --region eu-central-1 --verbose

teardown-devingen: clean
	serverless remove --stage prod --region eu-central-1 --verbose

deploy-devingen-staging: clean build zip
	serverless deploy --stage staging --region ca-central-1 --verbose

teardown-devingen-staging: clean
	serverless remove --stage staging --region ca-central-1 --verbose