package kimlikruntime

import (
	"context"
	"github.com/devingen/api-core/dvnruntime"
	core_model "github.com/devingen/api-core/model"
	"github.com/devingen/kimlik-api/token-service"
	"net/http"
)

func AssertAuthenticationAndBody(ctx context.Context, req dvnruntime.Request, bodyValue interface{}) (string, *token_service.TokenPayload, error) {
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

	if req.Body == "" {
		return "", tokenPayload, core_model.NewError(http.StatusBadRequest, "body-missing")
	}

	err := dvnruntime.ParseBody(req.Body, &bodyValue)
	if err != nil {
		return base, tokenPayload, err
	}
	return base, tokenPayload, nil
}
