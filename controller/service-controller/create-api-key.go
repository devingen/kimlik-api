package service_controller

import (
	"context"
	"crypto/rand"
	"fmt"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api"
	"github.com/devingen/kimlik-api/dto"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"net/http"
)

func (controller ServiceController) CreateAPIKey(ctx context.Context, req core.Request) (interface{}, int, error) {

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, 0, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	var body dto.CreateApiKeyRequest
	token, err := kimlik.AssertAuthenticationAndBody(ctx, req, &body)
	if err != nil {
		return nil, 0, err
	}

	user, err := controller.DataService.FindUserUserWithId(base, token.UserId)
	if err != nil {
		return nil, 0, err
	}

	if user == nil {
		return nil, 0, core.NewStatusError(http.StatusNotFound)
	}

	key, hash, err := GenerateApiKey()
	if err != nil {
		return nil, 0, err
	}

	// TODO check product ownership
	_, err = controller.DataService.CreateAPIKey(base, body.Name, body.ProductId, body.Scopes, key[:7], hash, user)
	if err != nil {
		return nil, 0, err
	}

	return &dto.CreateApiKeyResponse{Key: key}, http.StatusOK, err
}

func GenerateApiKey() (string, string, error) {
	apiKey, err := GenerateRandomString(32)
	if err != nil {
		return "", "", err
	}

	hashedKey, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	fmt.Println("-----------")
	fmt.Println("ApiKey:   ", apiKey)
	fmt.Println("HashedKey:", string(hashedKey))
	fmt.Println("-----------")
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

func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}
