package server

import (
	"context"
	"net/http"

	core "github.com/devingen/api-core"
	"github.com/devingen/api-core/database"
	"github.com/devingen/api-core/server"
	"github.com/devingen/api-core/wrapper"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"github.com/devingen/kimlik-api/config"
	service_controller "github.com/devingen/kimlik-api/controller/service-controller"
	ds "github.com/devingen/kimlik-api/data-service"
	mongods "github.com/devingen/kimlik-api/data-service/mongo-data-service"
	webhookis "github.com/devingen/kimlik-api/interceptor-service/webhook-interceptor-service"
	token_service "github.com/devingen/kimlik-api/token-service"
	json_web_token_service "github.com/devingen/kimlik-api/token-service/json-web-token-service"
	kimlikwrapper "github.com/devingen/kimlik-api/wrapper"
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
	router.HandleFunc("/{base}/oauth2/authorize", wrap(serviceController.OAuth2Authorize)).Methods(http.MethodGet)
	router.HandleFunc("/{base}/oauth2/token", wrap(serviceController.OAuth2Token)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/oauth2/certs", wrap(serviceController.OAuth2GetJWKS)).Methods(http.MethodGet)
	router.HandleFunc("/{base}/.well-known/openid-configuration", wrap(serviceController.OAuth2GetOIDCConfiguration)).Methods(http.MethodGet)

	router.HandleFunc("/{base}/authenticate", wrap(serviceController.Authenticate)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/authorization-url", wrap(serviceController.GetAuthorizationURL)).Methods(http.MethodGet)
	router.HandleFunc("/{base}/auth-methods", wrap(serviceController.GetAuthMethods)).Methods(http.MethodGet)
	router.HandleFunc("/{base}/register", wrap(serviceController.RegisterWithEmail)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/userinfo", wrap(serviceController.GetUserInfo)).Methods(http.MethodGet)
	router.HandleFunc("/{base}/session", wrap(serviceController.GetSession)).Methods(http.MethodGet)
	router.HandleFunc("/{base}/activate", wrap(serviceController.ActivateUser)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/auth/password", wrap(serviceController.ChangePassword)).Methods(http.MethodPut)
	router.HandleFunc("/{base}/login", wrap(serviceController.LoginWithEmail)).Methods(http.MethodPost)   // TODO delete
	router.HandleFunc("/{base}/sessions", wrap(serviceController.CreateSession)).Methods(http.MethodPost) // TODO delete

	router.HandleFunc("/{base}/users", wrap(serviceController.FindUsers)).Methods(http.MethodGet)
	router.HandleFunc("/{base}/users/{id}/anonymize", wrap(serviceController.AnonymizeUser)).Methods(http.MethodPost)

	router.HandleFunc("/{base}/api-keys", wrap(serviceController.CreateAPIKey)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/api-keys", wrap(serviceController.FindAPIKeys)).Methods(http.MethodGet)
	router.HandleFunc("/{base}/api-keys/{id}", wrap(serviceController.UpdateAPIKey)).Methods(http.MethodPut)
	router.HandleFunc("/{base}/api-keys/{id}", wrap(serviceController.DeleteAPIKey)).Methods(http.MethodDelete)
	router.HandleFunc("/{base}/api-keys/verify", wrap(serviceController.VerifyAPIKey)).Methods(http.MethodGet)

	router.HandleFunc("/{base}/app-integrations", wrap(serviceController.CreateAppIntegration)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/app-integrations", wrap(serviceController.FindAppIntegrations)).Methods(http.MethodGet)
	router.HandleFunc("/{base}/app-integrations/{id}", wrap(serviceController.UpdateAppIntegration)).Methods(http.MethodPut)
	router.HandleFunc("/{base}/app-integrations/{id}", wrap(serviceController.DeleteAppIntegration)).Methods(http.MethodDelete)

	router.HandleFunc("/{base}/oauth2-configs", wrap(serviceController.CreateOAuth2Config)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/oauth2-configs", wrap(serviceController.FindOAuth2Configs)).Methods(http.MethodGet)
	router.HandleFunc("/{base}/oauth2-configs/{id}", wrap(serviceController.UpdateOAuth2Config)).Methods(http.MethodPut)
	router.HandleFunc("/{base}/oauth2-configs/{id}", wrap(serviceController.DeleteOAuth2Config)).Methods(http.MethodDelete)

	router.HandleFunc("/{base}/saml-configs", wrap(serviceController.CreateSAMLConfig)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/saml-configs", wrap(serviceController.FindSAMLConfigs)).Methods(http.MethodGet)
	router.HandleFunc("/{base}/saml-configs/{id}", wrap(serviceController.UpdateSAMLConfig)).Methods(http.MethodPut)
	router.HandleFunc("/{base}/saml-configs/{id}", wrap(serviceController.DeleteSAMLConfig)).Methods(http.MethodDelete)
	router.HandleFunc("/{base}/saml-configs/{id}/build", wrap(serviceController.BuildSAMLAuthURL)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/saml-configs/{id}/login", wrap(serviceController.LoginWithSAML)).Methods(http.MethodGet)
	router.HandleFunc("/{base}/saml-configs/{id}/consume", wrap(serviceController.ConsumeSAMLAuthResponse)).Methods(http.MethodPost)

	router.HandleFunc("/{base}/setup", wrap(serviceController.Setup)).Methods(http.MethodPost)
	router.HandleFunc("/{base}/tenant-info", wrap(serviceController.GetTenantInfo)).Methods(http.MethodGet)
	router.HandleFunc("/{base}/tenant-info", wrap(serviceController.UpdateTenantInfo)).Methods(http.MethodPut)

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
