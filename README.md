# Kimlik API

## Structure
Project has different layers for the sake of modularity. Since that it supports both web server and AWS Lambda, 
they use common modules that contain the logic.

```
Request -> Server -> HTTP Handler  
                                  \
                                    -> Controller -> Service
                                  /
Request -> AWS Lambda Function --
```

### Web Server and AWS Lambda handlers
These are two interfaces that the API is exposed to. They are responsible of parsing the request,
building and passing the request parameters to the controllers and building the response.

### Controller
Gets the request data from Lambda or web server handlers, executes the logic through services and returns the result.

### Service
Services are responsible of data manipulation and it's usually used with databases.

## Integration Tests

Integration test suite connects to the Mongo Database running on `localhost` without any username and password.
Initially, it clears the database and inserts the sample data into the `dvn-atama-api-integration-test` 
database from the files in `integrationtests/data`.

```shell
go test ./integrationtests/...
```

## Deployment

### Docker

Docker build needs to access private repositories under github.com/devingen.

To let it pull the repos, you need to [generate a personal access token](https://github.com/settings/tokens)
and pass GIT_TOKEN to the command as follows.

```shell
GIT_TOKEN=GITHUB_TOKEN_GENERATED_ON_WEBSITE
make release-docker GIT_TOKEN=$GIT_TOKEN IMAGE_TAG=0.1.9
make release-docker GIT_TOKEN=$GIT_TOKEN IMAGE_TAG=latest
```

```shell
docker run -d \
  --restart always \
  --name devingen-api \
  -e KIMLIK_API_JWT_SIGN_KEY=... \
  -e KIMLIK_API_MONGO_URI=... \
  -e KIMLIK_API_MONGO_USER_DATABASE=... \
  -e KIMLIK_API_WEBHOOK_HEADERS=... \
  -e KIMLIK_API_WEBHOOK_URL=... \
  -p 1001:1001 \
  devingen/kimlik-api:latest
```

### AWS Lambda with Serverless Framework

This commands executes the command in `Makefile` which clears the previous builds,
generates executables and deploys the AWS Functions through Serverless Framework.

**dev**
```shell
make deploy-devingen-dev
```

**prod**
```shell
make deploy-devingen
```

## Development

### Adding local dependency

Add `replace` config to go.mod before require.

```
replace github.com/devingen/api-core => ../api-core
```

### Releasing a new version

Create a git tag with the desired version and push the tag.

```
# see tags
git tag --list

# create new tag
git tag -a v0.0.19 -m "Change user activation flow, change active user status enum"

# push new tag
git push origin v0.0.19
```