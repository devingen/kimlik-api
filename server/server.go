package server

import (
	"context"
	"net/http"

	core "github.com/devingen/api-core"
	"github.com/devingen/api-core/database"
	"github.com/devingen/api-core/server"
	"github.com/devingen/api-core/wrapper"
	"github.com/devingen/kimlik-api/config"
	service_controller "github.com/devingen/kimlik-api/controller/service-controller"
	ds "github.com/devingen/kimlik-api/data-service"
	mongods "github.com/devingen/kimlik-api/data-service/mongo-data-service"
	webhookis "github.com/devingen/kimlik-api/interceptor-service/webhook-interceptor-service"
	token_service "github.com/devingen/kimlik-api/token-service"
	json_web_token_service "github.com/devingen/kimlik-api/token-service/json-web-token-service"
	kimlikwrapper "github.com/devingen/kimlik-api/wrapper"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
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

	wrap := generateWrapper(appConfig, jwtService, dataService)

	router := mux.NewRouter()
	router.HandleFunc("/{base}/oauth/token", wrap(serviceController.OAuthToken)).Methods(http.MethodPost)

	router.HandleFunc("/{base}/session", wrap(serviceController.GetSession)).Methods(http.MethodGet)
	router.HandleFunc("/{base}/sessions", wrap(serviceController.CreateSession)).Methods(http.MethodPost)

	router.HandleFunc("/{base}/register", wrap(serviceController.RegisterWithEmail)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/login", wrap(serviceController.LoginWithEmail)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/activate", wrap(serviceController.ActivateUser)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/auth/password", wrap(serviceController.ChangePassword)).Methods(http.MethodPut)

	router.HandleFunc("/{base}/users", wrap(serviceController.FindUsers)).Methods(http.MethodGet)
	router.HandleFunc("/{base}/users/{id}/anonymize", wrap(serviceController.AnonymizeUser)).Methods(http.MethodPost)

	router.HandleFunc("/{base}/api-keys", wrap(serviceController.CreateAPIKey)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/api-keys", wrap(serviceController.FindAPIKeys)).Methods(http.MethodGet)
	router.HandleFunc("/{base}/api-keys/{id}", wrap(serviceController.UpdateAPIKey)).Methods(http.MethodPut)
	router.HandleFunc("/{base}/api-keys/{id}", wrap(serviceController.DeleteAPIKey)).Methods(http.MethodDelete)
	router.HandleFunc("/{base}/api-keys/verify", wrap(serviceController.VerifyAPIKey)).Methods(http.MethodGet)

	router.HandleFunc("/{base}/saml-configs", wrap(serviceController.CreateSAMLConfig)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/saml-configs", wrap(serviceController.FindSAMLConfigs)).Methods(http.MethodGet)
	router.HandleFunc("/{base}/saml-configs/{id}", wrap(serviceController.UpdateSAMLConfig)).Methods(http.MethodPut)
	router.HandleFunc("/{base}/saml-configs/{id}", wrap(serviceController.DeleteSAMLConfig)).Methods(http.MethodDelete)
	router.HandleFunc("/{base}/saml-configs/{id}/build", wrap(serviceController.BuildSAMLAuthURL)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/saml-configs/{id}/login", wrap(serviceController.LoginWithSAML)).Methods(http.MethodGet)
	router.HandleFunc("/{base}/saml-configs/{id}/consume", wrap(serviceController.ConsumeSAMLAuthResponse)).Methods(http.MethodPost)

	http.Handle("/", &server.CORSRouterDecorator{
		R: router,
		Headers: map[string]string{
			server.CORSAccessControlAllowHeaders: server.CORSAccessControlAllowHeadersDefaultValue + ",devingen-product-id",
			server.CORSAccessControlAllowMethods: server.CORSAccessControlAllowMethodsDefaultValue,
		},
		AllowSenderOrigin: true,
	})
	return srv
}

func generateWrapper(appConfig config.App, jwtService token_service.ITokenService, dataService ds.IKimlikDataService) func(f core.Controller) func(http.ResponseWriter, *http.Request) {
	return func(f core.Controller) func(http.ResponseWriter, *http.Request) {
		ctx := context.Background()

		// add logger and auth handler
		withLogger := wrapper.WithLogger(
			appConfig.LogLevel,
			kimlikwrapper.WithAuth(f, jwtService, dataService),
		)

		// convert to HTTP handler
		handler := wrapper.WithHTTPHandler(ctx, withLogger)
		return handler
	}
}
