package service_controller

import (
	"context"
	"crypto/rand"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api"
	"github.com/devingen/kimlik-api/dto"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"net/http"
)

func (c ServiceController) CreateAPIKey(ctx context.Context, req core.Request) (*core.Response, error) {

	_, interceptorStatusCode, interceptorError := c.InterceptorService.Pre(ctx, req)
	if interceptorError != nil {
		return &core.Response{
			StatusCode: interceptorStatusCode,
			Body:       interceptorError,
		}, nil
	}

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	var body dto.CreateApiKeyRequest
	token, err := kimlik.AssertAuthenticationAndBody(ctx, req, &body)
	if err != nil {
		return nil, err
	}

	user, err := c.DataService.FindUserWithId(ctx, base, token.UserId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, core.NewStatusError(http.StatusNotFound)
	}

	keyID, err := GenerateRandomString(20)
	if err != nil {
		return nil, err
	}

	key, hash, err := GenerateApiKey()
	if err != nil {
		return nil, err
	}

	apiKey, err := c.DataService.CreateAPIKey(ctx, base, *body.Name, body.Scopes, keyID, hash, user)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body: dto.CreateApiKeyResponse{
			Key:   key,
			KeyID: keyID,
			Name:  *apiKey.Name,
		},
	}, nil
}

func GenerateApiKey() (string, string, error) {
	apiKey, err := GenerateRandomString(40)
	if err != nil {
		return "", "", err
	}

	hashedKey, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	return apiKey, string(hashedKey), nil
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
