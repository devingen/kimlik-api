package kimlik

import (
	core "github.com/devingen/api-core"
	token_service "github.com/devingen/kimlik-api/token-service"
	"strings"
)

const JWTPrefix = "Bearer"

func ExtractToken(jwtService token_service.ITokenService, req core.Request) (*token_service.TokenPayload, error) {
	authHeader, hasAuthHeader := req.Headers["authorization"]
	if !hasAuthHeader || authHeader == "" {
		// skip if there is no header
		return nil, nil
	}

	// check if the header starts with 'Bearer ' prefix
	if strings.Index(authHeader, JWTPrefix) != 0 {
		// skip if auth header is not a valid JWT
		return nil, nil
	}

	tokenPayload, err := jwtService.ParseToken(authHeader[len(JWTPrefix)+1:])
	if err != nil {
		// return error if the JWT is not a valid or expired
		return nil, err
	}
	return tokenPayload, nil
}
