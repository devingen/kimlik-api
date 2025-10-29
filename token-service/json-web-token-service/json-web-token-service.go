package json_web_token_service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	core "github.com/devingen/api-core"
	token_service "github.com/devingen/kimlik-api/token-service"
	"github.com/devingen/kimlik-api/util"
	"github.com/dgrijalva/jwt-go"
)

const (
	JWTVersion   = "ver"
	JWTExpires   = "exp"
	JWTUserID    = "userId"
	JWTSessionID = "sessionId"
	JWTScopes    = "scopes"
)

const envJWTSignKey = "KIMLIK_JWT_SIGN_KEY"

type JWTService struct {
	signKey string
}

// New generates new JWTService
func New(signKey string) *JWTService {
	return &JWTService{
		signKey: signKey,
	}
}

func (jwtService *JWTService) Init() {
	signKey, hasVar := os.LookupEnv(envJWTSignKey)
	if !hasVar {
		log.Fatalf("Missing environment variable %s", envJWTSignKey)
	}
	jwtService.signKey = signKey
}

func (jwtService *JWTService) GenerateAccessToken(userId, sessionId string, scopes []string, exp int64) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	mapClaims := token.Claims.(jwt.MapClaims)
	mapClaims[JWTVersion] = "1.0"
	mapClaims[JWTExpires] = exp
	mapClaims[JWTUserID] = userId
	mapClaims[JWTSessionID] = sessionId
	mapClaims[JWTScopes] = scopes

	tokenString, err := token.SignedString([]byte(jwtService.signKey))
	if err != nil {
		return tokenString, err
	}
	return tokenString, nil
}

func (jwtService *JWTService) ParseAccessToken(accessToken string) (*token_service.TokenPayload, error) {
	token, tokenErr := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtService.signKey), nil
	})

	if tokenErr != nil {
		if tokenErr.Error() == "Token is expired" {
			return nil, core.NewError(http.StatusUnauthorized, "token-expired")
		}
		return nil, core.NewError(http.StatusUnauthorized, tokenErr.Error())
	}

	if !token.Valid {
		return nil, core.NewError(http.StatusUnauthorized, "invalid-token")
	}

	var data token_service.TokenPayload
	convertErr := util.ConvertMapToStruct(token.Claims.(jwt.MapClaims), &data)
	if convertErr != nil {
		core.NewError(http.StatusInternalServerError, convertErr.Error())
	}

	return &data, nil
}

func (jwtService *JWTService) GenerateRefreshToken() (*token_service.RefreshToken, error) {
	// Generate a raw refresh token
	rawToken, err := GenerateSecureToken(32) // 32 bytes = 64 hex characters
	if err != nil {
		return nil, errors.New("Error generating refresh token:" + err.Error())
	}

	// Hash the token before storing it
	hashedToken := jwtService.HashRefreshToken(rawToken)

	// Store `hashedToken` in database, and send `rawToken` to the client
	return &token_service.RefreshToken{
		HashedToken: hashedToken,
		RawToken:    rawToken,
	}, nil
}

// HashRefreshToken creates a HMAC SHA-256 hash of the token using a secret key
func (jwtService *JWTService) HashRefreshToken(token string) string {
	h := hmac.New(sha256.New, []byte(jwtService.signKey))
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

func (jwtService *JWTService) GenerateAuthorizationCode() (*string, error) {
	// Generate random code
	code, err := GenerateSecureToken(32) // 32 bytes = 64 hex characters
	if err != nil {
		return nil, errors.New("Error generating code:" + err.Error())
	}

	return &code, nil
}

// GenerateSecureToken generates a cryptographically secure random token
func GenerateSecureToken(length int) (string, error) {
	token := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, token); err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}
