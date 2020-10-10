package kimlikruntime

import (
	"context"
	"github.com/devingen/api-core/dvnruntime"
	core_model "github.com/devingen/api-core/model"
	"net/http"
)

func AssertBaseAndBody(ctx context.Context, req dvnruntime.Request, bodyValue interface{}) (string, error) {
	base := req.PathParameters["base"]
	if base == "" {
		return "", core_model.NewError(http.StatusBadRequest, "base-missing")
	}

	if req.Body == "" {
		return "", core_model.NewError(http.StatusBadRequest, "body-missing")
	}

	err := dvnruntime.ParseBody(req.Body, &bodyValue)
	if err != nil {
		return base, err
	}
	return base, nil
}
