package server

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/api-core/database"
	"github.com/devingen/api-core/server"
	"github.com/devingen/api-core/wrapper"
	"github.com/devingen/kimlik-api/config"
	service_controller "github.com/devingen/kimlik-api/controller/service-controller"
	mongods "github.com/devingen/kimlik-api/data-service/mongo-data-service"
	kimlikwrapper "github.com/devingen/kimlik-api/wrapper"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
)

// New creates a new HTTP server
func New(appConfig config.App, db *database.Database) *http.Server {
	validate := validator.New()
	core.SetValidator(validate)

	srv := &http.Server{Addr: ":" + appConfig.Port}

	dataService := mongods.New(appConfig.Mongo.Database, db)
	serviceController := service_controller.New(dataService)

	wrap := generateWrapper(appConfig)

	router := mux.NewRouter()
	router.HandleFunc("/{base}/register", wrap(serviceController.RegisterWithEmail)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/login", wrap(serviceController.LoginWithEmail)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/session", wrap(serviceController.GetSession)).Methods(http.MethodGet)
	router.HandleFunc("/{base}/auth/password", wrap(serviceController.ChangePassword)).Methods(http.MethodPut)
	router.HandleFunc("/{base}/api-keys", wrap(serviceController.CreateAPIKey)).Methods(http.MethodPost)

	http.Handle("/", &server.CORSRouterDecorator{R: router})
	return srv
}

func generateWrapper(appConfig config.App) func(f core.Controller) func(http.ResponseWriter, *http.Request) {
	return func(f core.Controller) func(http.ResponseWriter, *http.Request) {
		ctx := context.Background()

		// add logger and auth handler
		withLogger := wrapper.WithLogger(
			appConfig.LogLevel,
			kimlikwrapper.WithAuth(f),
		)

		// convert to HTTP handler
		handler := wrapper.WithHTTPHandler(ctx, withLogger)
		return handler
	}
}

// Runs the server that contains all the services
//func main() {
//
//	db, err := database.New("mongodb://devingen-services-dev:d3v8ng3ns34v8c3sd3v@mongodb0.mentornity.com/?authSource=admin")
//	if err != nil {
//		log.Fatalf("Database connection failed %s", err.Error())
//	}
//
//	databaseService := service.NewDatabaseService(db)
//	//kimlikService := kimlik_service.NewDatabaseService(db)
//	serviceController := controller.NewServiceController(*databaseService)
//
//	wrap := generateWrapper()
//
//	router := mux.NewRouter()
//	router.HandleFunc("/{base}/products", wrap(serviceController.CreateProduct)).Methods(http.MethodPost)
//	router.HandleFunc("/{base}/products", wrap(serviceController.GetProducts)).Methods(http.MethodGet)
//	router.HandleFunc("/{base}/products/{id}", wrap(serviceController.GetProductWithId)).Methods(http.MethodGet)
//	router.HandleFunc("/{base}/workspaces", wrap(serviceController.CreateWorkspace)).Methods(http.MethodPost)
//	router.HandleFunc("/{base}/workspace-ownerships", wrap(serviceController.GetWorkspaceOwnerships)).Methods(http.MethodGet)
//
//	router.HandleFunc("/webhooks/damga/pre", wrap(serviceController.HandleDamgaPreWebhook)).Methods(http.MethodPost)
//	router.HandleFunc("/webhooks/damga/final", wrap(serviceController.HandleDamgaFinalWebhook)).Methods(http.MethodPost)
//
//	http.Handle("/", &server.CORSRouterDecorator{R: router})
//	err = http.ListenAndServe(":1002", &server.CORSRouterDecorator{R: router})
//	if err != nil {
//		log.Fatalf("Listen and serve failed %s", err.Error())
//	}
//}
