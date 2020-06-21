package handler

import (
	"github.com/devingen/kimlik-api/controller"
	"github.com/gorilla/mux"
	"net/http"
)

type ServerHandler struct {
	Controller controller.ServiceController
	Router     *mux.Router
}

func NewHttpServiceHandler(controller controller.ServiceController) ServerHandler {
	handler := ServerHandler{Controller: controller}

	handler.Router = mux.NewRouter()
	handler.Router.HandleFunc("/{base}/register/email", handler.RegisterWithEmail).Methods(http.MethodPost)
	handler.Router.HandleFunc("/{base}/login/email", handler.loginWithEmail).Methods(http.MethodPost)

	return handler
}
