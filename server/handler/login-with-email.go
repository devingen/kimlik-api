package handler

import (
	"encoding/json"
	"github.com/devingen/api-core/server"
	"github.com/devingen/api-core/util"
	"github.com/devingen/kimlik-api/dto"
	"github.com/gorilla/mux"
	"net/http"
)

func (handler ServerHandler) loginWithEmail(w http.ResponseWriter, r *http.Request) {
	pathVariables := mux.Vars(r)
	base := pathVariables["base"]

	var body dto.LoginWithEmailRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ip, _ := util.GetClientIPHelper(r)
	userAgent := r.Header.Get("User-Agent")
	client := r.Header.Get("Client")

	result, err := handler.Controller.LoginWithEmail(base, client, userAgent, ip, &body)
	response, err := util.BuildResponse(http.StatusOK, result, err)
	server.ReturnResponse(w, response, err)
}
