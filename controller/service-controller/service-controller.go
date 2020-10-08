package service_controller

import (
	coremodel "github.com/devingen/api-core/model"
	"github.com/devingen/kimlik-api/service"
	"github.com/devingen/kimlik-api/util"
	"github.com/dgrijalva/jwt-go"
	"net/http"
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

type TokenData struct {
	Version   string `bson:"ver,omitempty"`
	Expires   string `bson:"exp,omitempty"`
	UserId    string `bson:"userId,omitempty"`
	SessionId string `bson:"sessionId,omitempty"`
	Scopes    string `bson:"scopes,omitempty"`
}

func (controller *ServiceController) GenerateToken(userId, sessionId string, scopes []Scope, duration int32) (string, error) {

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

func (controller *ServiceController) ParseToken(accessToken string) (*TokenData, error) {
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

	var data TokenData
	convertErr := util.ConvertMapToStruct(token.Claims.(jwt.MapClaims), &data)
	if convertErr != nil {
		coremodel.NewError(http.StatusInternalServerError, convertErr.Error())
	}

	return &data, nil
}
