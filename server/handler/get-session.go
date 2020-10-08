package handler

import (
	"github.com/devingen/api-core/server"
	"github.com/devingen/api-core/util"
	"github.com/gorilla/mux"
	"net/http"
)

func (handler ServerHandler) getSession(w http.ResponseWriter, r *http.Request) {
	pathVariables := mux.Vars(r)
	base := pathVariables["base"]

	jwt := r.Header.Get("Authorization")[len("Bearer "):]

	result, err := handler.Controller.GetSession(base, jwt)
	response, err := util.BuildResponse(http.StatusOK, result, err)
	server.ReturnResponse(w, response, err)
}
