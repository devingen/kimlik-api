package kimlikruntime

import (
	"context"
	"github.com/devingen/api-core/dvnruntime"
	core_model "github.com/devingen/api-core/model"
	"github.com/devingen/kimlik-api/token-service"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func AssertAuthenticationAndBody(ctx context.Context, req dvnruntime.Request, bodyValue interface{}) (string, *token_service.TokenPayload, error) {
	// assert token is present
	if ctx.Value(token_service.ContextKeyTokenPayload) == nil {
		return "", nil, core_model.NewStatusError(http.StatusUnauthorized)
	}

	tokenPayload := ctx.Value(token_service.ContextKeyTokenPayload).(*token_service.TokenPayload)
	if tokenPayload == nil {
		return "", nil, core_model.NewStatusError(http.StatusUnauthorized)
	}

	// retrieve base from path
	base := req.PathParameters["base"]
	if base == "" {
		return "", tokenPayload, core_model.NewError(http.StatusBadRequest, "base-missing")
	}

	// assert body is present
	if req.Body == "" {
		return "", tokenPayload, core_model.NewError(http.StatusBadRequest, "body-missing")
	}

	// parse body
	err := dvnruntime.ParseBody(req.Body, &bodyValue)
	if err != nil {
		return base, tokenPayload, err
	}

	// validate body
	if validate == nil {
		validate = validator.New()
	}
	err = validate.Struct(bodyValue)

	// return proper validation error
	if err != nil {
		switch castedError := err.(type) {
		case validator.ValidationErrors:
			return base, tokenPayload, core_model.NewError(http.StatusBadRequest, castedError.Error())
		default:
			return base, tokenPayload, err
		}
	}
	return base, tokenPayload, nil
}
