package kimlik

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/model"
	token_service "github.com/devingen/kimlik-api/token-service"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func SetValidator(v *validator.Validate) {
	core.SetValidator(v)
}

func AssertAuthentication(ctx context.Context) (*token_service.TokenPayload, error) {
	tokenPayload := Of(ctx)
	if tokenPayload == nil {
		return nil, core.NewStatusError(http.StatusUnauthorized)
	}
	return tokenPayload, nil
}

func AssertApiKey(ctx context.Context) (*model.ApiKeyPayload, error) {
	apiKeyPayload := OfApiKey(ctx)
	if apiKeyPayload == nil {
		return nil, core.NewStatusError(http.StatusUnauthorized)
	}
	return apiKeyPayload, nil
}

func AssertAuthenticationAndBody(ctx context.Context, req core.Request, bodyValue interface{}) (*token_service.TokenPayload, error) {
	tokenPayload := Of(ctx)
	if tokenPayload == nil {
		return nil, core.NewStatusError(http.StatusUnauthorized)
	}

	// assert body is present
	if req.Body == "" {
		return tokenPayload, core.NewError(http.StatusBadRequest, "body-missing")
	}

	// parse body
	err := req.ParseBody(&bodyValue)
	if err != nil {
		return tokenPayload, err
	}

	// validate body
	validate := core.GetValidator()
	if validate == nil {
		validate = validator.New()
	}
	err = validate.Struct(bodyValue)

	// return proper validation error
	if err != nil {
		switch castedError := err.(type) {
		case validator.ValidationErrors:
			return tokenPayload, core.NewError(http.StatusBadRequest, castedError.Error())
		default:
			return tokenPayload, err
		}
	}
	return tokenPayload, nil
}
