package kimlikruntime

import (
	"context"
	"github.com/devingen/api-core/dvnruntime"
	core_model "github.com/devingen/api-core/model"
	"github.com/devingen/kimlik-api/token-service"
	"net/http"
)

func AssertAuthentication(ctx context.Context, req dvnruntime.Request) (string, *token_service.TokenPayload, error) {
	if ctx.Value(token_service.ContextKeyTokenPayload) == nil {
		return "", nil, core_model.NewStatusError(http.StatusUnauthorized)
	}

	tokenPayload := ctx.Value(token_service.ContextKeyTokenPayload).(*token_service.TokenPayload)
	if tokenPayload == nil {
		return "", nil, core_model.NewStatusError(http.StatusUnauthorized)
	}

	base := req.PathParameters["base"]
	if base == "" {
		return "", tokenPayload, core_model.NewError(http.StatusBadRequest, "base-missing")
	}
	return base, tokenPayload, nil
}
