.PHONY: build clean deploy

build:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/register-with-email aws/register-with-email/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy-prod: clean build
	serverless deploy --stage PROVIDE_STAGE_HERE --region eu-central-1 --verbose

teardown-prod: clean
	serverless remove --stage PROVIDE_STAGE_HERE --region eu-central-1 --verbose
