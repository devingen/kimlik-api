package kimlik

import (
	"context"
	core "github.com/devingen/api-core"
	ds "github.com/devingen/kimlik-api/data-service"
	"github.com/devingen/kimlik-api/model"
	token_service "github.com/devingen/kimlik-api/token-service"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
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

func AssertApiKey(ctx context.Context, base string, dataService ds.IKimlikDataService) (*model.APIKey, error) {
	// retrieve the api key extracted from the 'api-key' header
	apiKeyPayload := OfApiKey(ctx)
	if apiKeyPayload == nil {
		return nil, core.NewStatusError(http.StatusUnauthorized)
	}

	return AssertApiKeyPayload(ctx, base, dataService, apiKeyPayload)
}

func AssertApiKeyPayload(ctx context.Context, base string, dataService ds.IKimlikDataService, apiKeyPayload *model.ApiKeyPayload) (*model.APIKey, error) {
	// retrieve the corresponding api key entry from the database
	apiKeys, err := dataService.FindAPIKeys(ctx, base, bson.M{
		"keyId": apiKeyPayload.KeyID,
	})
	if err != nil {
		return nil, err
	}

	if apiKeys == nil || len(apiKeys) == 0 {
		return nil, core.NewError(http.StatusUnauthorized, "key-not-found")
	}
	apiKey := apiKeys[0]

	// check if the api key sent by the client matches the hash saved in the database
	if err := bcrypt.CompareHashAndPassword([]byte(*apiKey.Hash), []byte(apiKeyPayload.Key)); err != nil {
		return nil, core.NewError(http.StatusUnauthorized, "key-mismatch")
	}
	return apiKey, nil
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
