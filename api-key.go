package kimlik

import (
	"encoding/base64"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/model"
	"net/http"
	"strings"
)

// ExtractApiKey reads the api key from header. The value in the header is expected to be
// a base64 encoded version of the api key ID and api key concatenated with a colon.
// In order to pass a valid api key header in the requests;
//   1 - Concatenate the key ID and key like apiKeyID:apiKey
//   2 - Base64 encode the concatenated string
//   3 - Put the generated base64 string into request headers with 'Api-Key' or 'api-key' header name
func ExtractApiKey(req core.Request) (*model.ApiKeyPayload, error) {
	if req.Headers == nil {
		return nil, nil
	}

	apiKey, hasApiKey := req.Headers["Api-Key"]
	if !hasApiKey || apiKey == "" {
		// aws lambda converts the custom headers to lowercase
		apiKey, hasApiKey = req.Headers["api-key"]
		if !hasApiKey || apiKey == "" {
			return nil, nil
		}
	}

	decodedKey, err := base64.StdEncoding.DecodeString(apiKey)
	if err != nil {
		return nil, core.NewError(http.StatusBadRequest, "api-key-cannot-be-decoded")
	}

	keyParts := strings.Split(string(decodedKey), ":")
	if len(keyParts) != 2 {
		return nil, core.NewError(http.StatusBadRequest, "api-key-is-malformed")
	}

	return &model.ApiKeyPayload{
		KeyID: keyParts[0],
		Key:   keyParts[1],
	}, nil
}
