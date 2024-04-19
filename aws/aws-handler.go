package aws

import (
	"context"
	"log"

	core "github.com/devingen/api-core"
	"github.com/devingen/api-core/database"
	"github.com/devingen/kimlik-api/config"
	"github.com/devingen/kimlik-api/controller"
	service_controller "github.com/devingen/kimlik-api/controller/service-controller"
	ds "github.com/devingen/kimlik-api/data-service"
	mongods "github.com/devingen/kimlik-api/data-service/mongo-data-service"
	webhookis "github.com/devingen/kimlik-api/interceptor-service/webhook-interceptor-service"
	token_service "github.com/devingen/kimlik-api/token-service"
	json_web_token_service "github.com/devingen/kimlik-api/token-service/json-web-token-service"
	kimlikwrapper "github.com/devingen/kimlik-api/wrapper"
	"github.com/kelseyhightower/envconfig"

	"github.com/devingen/api-core/wrapper"
)

var db *database.Database

// InitDeps creates the dependencies for the AWS Lambda functions.
func InitDeps() (controller.IServiceController, func(f core.Controller) wrapper.AWSLambdaHandler) {
	var appConfig config.App
	err := envconfig.Process("kimlik_api", &appConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	dataService := mongods.New(getDatabase(appConfig))
	jwtService := json_web_token_service.New(appConfig.JWTSignKey)
	interceptorService := webhookis.New(appConfig.Webhook.URL, appConfig.Webhook.Headers)
	serviceController := service_controller.New(dataService, jwtService, interceptorService)

	wrap := generateWrapper(appConfig, jwtService, dataService)
	return serviceController, wrap
}

func getDatabase(appConfig config.App) *database.Database {
	if db == nil {
		var err error
		db, err = database.New(appConfig.Mongo.URI)
		if err != nil {
			log.Fatalf("Database connection failed when creating a new database %s", err.Error())
		}
	} else if !db.IsConnected() {
		err := db.ConnectWithURI(appConfig.Mongo.URI)
		if err != nil {
			log.Fatalf("Database connection failed when connecting to an existing database %s", err.Error())
		}
	}
	return db
}

func generateWrapper(appConfig config.App, jwtService token_service.ITokenService, dataService ds.IKimlikDataService) func(f core.Controller) wrapper.AWSLambdaHandler {
	return func(f core.Controller) wrapper.AWSLambdaHandler {
		ctx := context.Background()

		// add logger and auth handler
		withLogger := wrapper.WithLogger(
			appConfig.LogLevel,
			kimlikwrapper.WithAuth(f, jwtService, dataService),
		)

		// convert to HTTP handler
		handler := wrapper.WithLambdaHandler(ctx, withLogger)
		return handler
	}
}
