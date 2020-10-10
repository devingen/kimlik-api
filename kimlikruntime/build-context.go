package kimlikruntime

import (
	"context"
	"github.com/devingen/api-core/dvnruntime"
	token_service "github.com/devingen/kimlik-api/token-service"
	json_web_token_service "github.com/devingen/kimlik-api/token-service/json-web-token-service"
	"strings"
)

const JWTPrefix = "Bearer"

func BuildContext(ctx context.Context, req dvnruntime.Request) (context.Context, error) {

	authHeader, hasAuthHeader := req.Headers["Authorization"]
	if !hasAuthHeader {
		// skip if there is no header
		return ctx, nil
	}

	startsWithBearer := strings.Index(authHeader, JWTPrefix) == 0
	if !startsWithBearer {
		// skip if auth header is not a valid JWT
		return ctx, nil
	}

	tokenPayload, err := json_web_token_service.NewTokenService().ParseToken(authHeader[len(JWTPrefix)+1:])
	if err != nil {
		// return error if the JWT is not a valid or expired
		return nil, err
	}

	requestContext := context.WithValue(ctx, token_service.ContextKeyTokenPayload, tokenPayload)
	return requestContext, nil
}
