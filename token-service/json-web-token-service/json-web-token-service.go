package json_web_token_service

import (
	coremodel "github.com/devingen/api-core/model"
	token_service "github.com/devingen/kimlik-api/token-service"
	"github.com/devingen/kimlik-api/util"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

const (
	JWTVersion   = "ver"
	JWTExpires   = "exp"
	JWTUserID    = "userId"
	JWTSessionID = "sessionId"
	JWTScopes    = "scopes"
)

type JWTService struct {
}

func (jwtService *JWTService) GenerateToken(userId, sessionId string, scopes []token_service.Scope, duration int32) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	mapClaims := token.Claims.(jwt.MapClaims)
	mapClaims[JWTVersion] = "1.0"
	mapClaims[JWTExpires] = time.Now().Add(time.Hour * time.Duration(duration)).Unix()
	mapClaims[JWTUserID] = userId
	mapClaims[JWTSessionID] = sessionId
	mapClaims[JWTScopes] = scopes

	tokenString, err := token.SignedString([]byte("a.SignKey"))
	if err != nil {
		return tokenString, err
	}
	return tokenString, nil
}

func (jwtService *JWTService) ParseToken(accessToken string) (*token_service.TokenPayload, error) {
	token, tokenErr := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		return []byte("a.SignKey"), nil
	})

	if tokenErr != nil {
		if tokenErr.Error() == "Token is expired" {
			return nil, coremodel.NewError(http.StatusUnauthorized, "token-expired")
		}
		return nil, coremodel.NewError(http.StatusUnauthorized, tokenErr.Error())
	}

	if !token.Valid {
		return nil, coremodel.NewError(http.StatusUnauthorized, "invalid-token")
	}

	var data token_service.TokenPayload
	convertErr := util.ConvertMapToStruct(token.Claims.(jwt.MapClaims), &data)
	if convertErr != nil {
		coremodel.NewError(http.StatusInternalServerError, convertErr.Error())
	}

	return &data, nil
}
