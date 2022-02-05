package kimlik

import (
	"encoding/base64"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/model"
	"net/http"
	"strings"
)

func VerifyApiKey(apiKey string) (*model.ApiKeyPayload, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(apiKey)
	if err != nil {
		return nil, core.NewError(http.StatusBadRequest, "api-key-cannot-be-decoded")
	}

	// TODO check
	//   1 - api key exists in database
	//   2 - key value matches the hash saved in the database
	keyParts := strings.Split(string(decodedKey), ":")
	if len(keyParts) != 2 {
		return nil, core.NewError(http.StatusBadRequest, "api-key-is-malformed")
	}
	return &model.ApiKeyPayload{
		Name:  keyParts[0],
		Value: keyParts[1],
	}, nil
}
