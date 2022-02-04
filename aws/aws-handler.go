package aws

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/api-core/database"
	"github.com/devingen/kimlik-api/config"
	"github.com/devingen/kimlik-api/controller"
	service_controller "github.com/devingen/kimlik-api/controller/service-controller"
	mongods "github.com/devingen/kimlik-api/data-service/mongo-data-service"
	"github.com/kelseyhightower/envconfig"
	"log"

	"github.com/devingen/api-core/wrapper"
)

var db *database.Database

// InitDeps creates the dependencies for the AWS Lambda functions.
func InitDeps() (controller.IServiceController, func(f core.Controller) wrapper.AWSLambdaHandler) {
	var appConfig config.App
	err := envconfig.Process("sepet", &appConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	dataService := mongods.New(appConfig.Mongo.Database, getDatabase(appConfig))
	serviceController := service_controller.New(dataService)

	wrap := generateWrapper(appConfig)
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

func generateWrapper(appConfig config.App) func(f core.Controller) wrapper.AWSLambdaHandler {
	return func(f core.Controller) wrapper.AWSLambdaHandler {
		ctx := context.Background()

		// add logger
		withLogger := wrapper.WithLogger(appConfig.LogLevel, f)

		// convert to HTTP handler
		handler := wrapper.WithLambdaHandler(ctx, withLogger)
		return handler
	}
}
