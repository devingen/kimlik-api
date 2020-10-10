package kimlikruntime

import (
	"context"
	"github.com/devingen/api-core/dvnruntime"
	core_model "github.com/devingen/api-core/model"
	"net/http"
)

func AssertBase(ctx context.Context, req dvnruntime.Request) (string, error) {
	base := req.PathParameters["base"]
	if base == "" {
		return "", core_model.NewError(http.StatusBadRequest, "base-missing")
	}
	return base, nil
}
