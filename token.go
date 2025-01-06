package kimlik

import (
	"net/http"
	"strings"

	core "github.com/devingen/api-core"

	token_service "github.com/devingen/kimlik-api/token-service"
)

const JWTPrefix = "Bearer"

func ExtractToken(jwtService token_service.ITokenService, req core.Request) (*token_service.TokenPayload, error) {
	authHeader, hasAuthHeader := req.Headers["authorization"]
	if !hasAuthHeader || authHeader == "" {
		// check cookie if there is no header
		return ExtractTokenFromCookie(jwtService, req)
	}

	// check if the header starts with 'Bearer ' prefix
	if strings.Index(authHeader, JWTPrefix) != 0 {
		// skip if auth header is not a valid JWT
		return nil, nil
	}

	tokenPayload, err := jwtService.ParseAccessToken(authHeader[len(JWTPrefix)+1:])
	if err != nil {
		// return error if the JWT is not a valid or expired
		return nil, err
	}
	return tokenPayload, nil
}

func ExtractTokenFromCookie(jwtService token_service.ITokenService, req core.Request) (*token_service.TokenPayload, error) {
	cookieHeader, hasCookieHeader := req.GetHeader("cookie")
	if hasCookieHeader {
		cookies, err := http.ParseCookie(cookieHeader)
		if err != nil {
			panic(err)
		}

		accessToken := ""
		for _, cookie := range cookies {
			if cookie.Name == "kimlik_token" {
				accessToken = cookie.Value
			}
		}

		if accessToken != "" {
			tokenPayload, err := jwtService.ParseAccessToken(accessToken)
			if err != nil {
				// return error if the JWT is not a valid or expired
				return nil, err
			}
			return tokenPayload, nil
		}
	}
	return nil, nil
}
