package service_controller

import (
	"github.com/devingen/kimlik-api/service"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Scope string

const (
	ScopeAll Scope = "all"
)

const (
	JWTVersion   = "ver"
	JWTExpires   = "exp"
	JWTUserID    = "userId"
	JWTSessionID = "sessionId"
	JWTScopes    = "scopes"
)

// ServiceController implements ServiceController interface by using KimlikService
type ServiceController struct {
	Service service.KimlikService
}

func (controller *ServiceController) GenerateToken(userId, sessionId string, scope []Scope, duration int32) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	mapClaims := token.Claims.(jwt.MapClaims)
	mapClaims[JWTVersion] = "1.0"
	mapClaims[JWTExpires] = time.Now().Add(time.Hour * time.Duration(duration)).Unix()
	mapClaims[JWTUserID] = userId
	mapClaims[JWTSessionID] = sessionId
	mapClaims[JWTScopes] = scope

	tokenString, err := token.SignedString([]byte("a.SignKey"))
	if err != nil {
		return tokenString, err
	}
	return tokenString, nil
}
