package dto

type CreateApiKeyRequest struct {
	Name   *string  `json:"name" validate:"required"`
	Scopes []string `json:"scopes" validate:"required"`
}

type CreateApiKeyResponse struct {
	Name  string `json:"name"`
	KeyID string `json:"keyId"`
	Key   string `json:"key"`
}

type UpdateApiKeyRequest struct {
	Name   *string  `json:"name" validate:"required"`
	Scopes []string `json:"scopes" validate:"required"`
}
