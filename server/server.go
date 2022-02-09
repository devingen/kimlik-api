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
	webhookis "github.com/devingen/kimlik-api/interceptor-service/webhook-interceptor-service"
	json_web_token_service "github.com/devingen/kimlik-api/token-service/json-web-token-service"
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

	dataService := mongods.New(db)
	jwtService := json_web_token_service.New(appConfig.JWTSignKey)
	interceptorService := webhookis.New(appConfig.Webhook.URL, appConfig.Webhook.Headers)
	serviceController := service_controller.New(dataService, jwtService, interceptorService)

	wrap := generateWrapper(appConfig)

	router := mux.NewRouter()
	router.HandleFunc("/{base}/register", wrap(serviceController.RegisterWithEmail)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/login", wrap(serviceController.LoginWithEmail)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/session", wrap(serviceController.GetSession)).Methods(http.MethodGet)
	router.HandleFunc("/{base}/auth/password", wrap(serviceController.ChangePassword)).Methods(http.MethodPut)
	router.HandleFunc("/{base}/api-keys", wrap(serviceController.CreateAPIKey)).Methods(http.MethodPost)

	router.HandleFunc("/{base}/saml-configs", wrap(serviceController.CreateSAMLConfig)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/saml-configs/{id}/build", wrap(serviceController.BuildSAMLAuthURL)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/saml-configs/{id}/consume", wrap(serviceController.ConsumeSAMLAuthResponse)).Methods(http.MethodPost)

	http.Handle("/", &server.CORSRouterDecorator{R: router})
	return srv
}

func generateWrapper(appConfig config.App) func(f core.Controller) func(http.ResponseWriter, *http.Request) {
	return func(f core.Controller) func(http.ResponseWriter, *http.Request) {
		ctx := context.Background()

		// add logger and auth handler
		withLogger := wrapper.WithLogger(
			appConfig.LogLevel,
			kimlikwrapper.WithAuth(f, appConfig.JWTSignKey),
		)

		// convert to HTTP handler
		handler := wrapper.WithHTTPHandler(ctx, withLogger)
		return handler
	}
}
