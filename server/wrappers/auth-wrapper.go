package wrappers

import (
	"context"
	"github.com/devingen/api-core/dvnruntime"
	"github.com/devingen/api-core/server"
	"github.com/devingen/api-core/util"
	"github.com/devingen/kimlik-api/kimlikruntime"
	"net/http"
)

func WithAuth(f dvnruntime.ControllerFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, httpReq *http.Request) {

		// convert http request to our custom request
		req := server.AdaptRequest(httpReq)

		// parse jwt and build context
		requestContext, err := kimlikruntime.BuildContext(context.Background(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// execute function
		result, status, err := f(requestContext, req)

		// convert response to our custom response
		response, err := util.BuildResponse(status, result, err)

		// write response data
		server.ReturnResponse(w, response, err)
	}
}
