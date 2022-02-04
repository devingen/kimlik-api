package server

//import (
//	"github.com/devingen/api-core/database"
//	"github.com/devingen/api-core/server"
//	"github.com/devingen/kimlik-api/controller"
//	"github.com/devingen/kimlik-api/server/wrappers"
//	"github.com/devingen/kimlik-api/service"
//	"github.com/devingen/kimlik-api/token-service/json-web-token-service"
//	"github.com/gorilla/mux"
//	"log"
//	"net/http"
//)
//
//// Runs the server that contains all the services
//func main() {
//
//	db, err := database.NewDatabase()
//	if err != nil {
//		log.Fatalf("Database connection failed %s", err.Error())
//	}
//
//	databaseService := service.NewDatabaseService(db)
//	tokenService := json_web_token_service.NewTokenService()
//	serviceController := controller.NewServiceController(databaseService, tokenService)
//
//	router := mux.NewRouter()
//	router.HandleFunc("/{base}/register", wrappers.WithLog(wrappers.WithAuth(serviceController.RegisterWithEmail))).Methods(http.MethodPost)
//	router.HandleFunc("/{base}/login", wrappers.WithLog(wrappers.WithAuth(serviceController.LoginWithEmail))).Methods(http.MethodPost)
//	router.HandleFunc("/{base}/session", wrappers.WithLog(wrappers.WithAuth(serviceController.GetSession))).Methods(http.MethodGet)
//	router.HandleFunc("/{base}/auth/password", wrappers.WithLog(wrappers.WithAuth(serviceController.ChangePassword))).Methods(http.MethodPut)
//	router.HandleFunc("/{base}/api-keys", wrappers.WithLog(wrappers.WithAuth(serviceController.CreateAPIKey))).Methods(http.MethodPost)
//
//	http.Handle("/", &server.CORSRouterDecorator{R: router})
//	err = http.ListenAndServe(":1001", &server.CORSRouterDecorator{R: router})
//	if err != nil {
//		log.Fatalf("Listen and serve failed %s", err.Error())
//	}
//}
