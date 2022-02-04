package kimlik

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/model"
	token_service "github.com/devingen/kimlik-api/token-service"
	"net/http"
)

// tokenKey type is an opaque type for the token lookup in a given context.
type tokenKey struct{}

// apiKeyKey type is an opaque type for the api key lookup in a given context.
type apiKeyKey struct{}

const JWTPrefix = "Bearer"

func WithJWTAuth(ctx context.Context, req core.Request, signKey string) (context.Context, error) {

	tokenPayload, err := RetrieveToken(req.Headers["Authorization"], signKey)
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

func WithAPIKeyAuth(ctx context.Context, req core.Request) (context.Context, error) {

	apiKey, hasApiKey := req.Headers["Api-Key"]
	if !hasApiKey || apiKey == "" {
		return ctx, nil
	}

	apiKeyPayload, err := VerifyApiKey(apiKey)
	if err != nil {
		return nil, core.NewError(http.StatusBadRequest, "invalid-api-key")
	}

	return context.WithValue(ctx, apiKeyKey{}, apiKeyPayload), nil
}

// OfApiKey function extracts a valid ApiKeyPayload object from a given context.
func OfApiKey(ctx context.Context) *model.ApiKeyPayload {
	token, ok := ctx.Value(apiKeyKey{}).(*model.ApiKeyPayload)
	if !ok {
		return nil
	}
	return token
}
