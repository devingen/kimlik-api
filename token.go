package kimlik

import (
	token_service "github.com/devingen/kimlik-api/token-service"
	json_web_token_service "github.com/devingen/kimlik-api/token-service/json-web-token-service"
	"strings"
)

func RetrieveToken(authHeader, signKey string) (*token_service.TokenPayload, error) {
	if authHeader == "" {
		// skip if there is no header
		return nil, nil
	}

	startsWithBearer := strings.Index(authHeader, JWTPrefix) == 0
	if !startsWithBearer {
		// skip if auth header is not a valid JWT
		return nil, nil
	}

	tokenPayload, err := json_web_token_service.New(signKey).ParseToken(authHeader[len(JWTPrefix)+1:])
	if err != nil {
		// return error if the JWT is not a valid or expired
		return nil, err
	}
	return tokenPayload, nil
}
