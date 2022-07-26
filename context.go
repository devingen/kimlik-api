package kimlik

import (
	"context"
	core "github.com/devingen/api-core"
	ds "github.com/devingen/kimlik-api/data-service"
	"github.com/devingen/kimlik-api/model"
	token_service "github.com/devingen/kimlik-api/token-service"
	"net/http"
)

// tokenKey type is an opaque type for the token lookup in a given context.
type tokenKey struct{}

// apiKeyKey type is an opaque type for the api key lookup in a given context.
type apiKeyKey struct{}

func WithJWTAuth(jwtService token_service.ITokenService, ctx context.Context, req core.Request) (context.Context, error) {

	tokenPayload, err := ExtractToken(jwtService, req)
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, tokenKey{}, tokenPayload), nil
}

// Of function extracts a valid TokenPayload object from a given context.
func Of(ctx context.Context) *token_service.TokenPayload {
	token, ok := ctx.Value(tokenKey{}).(*token_service.TokenPayload)
	if !ok {
		return nil
	}
	return token
}

func WithAPIKeyAuth(ctx context.Context, req core.Request, dataService ds.IKimlikDataService) (context.Context, error) {

	apiKeyPayload, err := ExtractApiKey(req)
	if err != nil {
		return nil, core.NewError(http.StatusBadRequest, err.Error())
	}

	return context.WithValue(ctx, apiKeyKey{}, apiKeyPayload), nil
}

// OfApiKey function extracts a valid ApiKeyPayload object
// from a given context that is created by the WithAPIKeyAuth priorly.
func OfApiKey(ctx context.Context) *model.ApiKeyPayload {
	token, ok := ctx.Value(apiKeyKey{}).(*model.ApiKeyPayload)
	if !ok {
		return nil
	}
	return token
}
