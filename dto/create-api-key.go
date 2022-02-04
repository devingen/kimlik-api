package dto

type CreateApiKeyRequest struct {
	Name      string   `json:"name" validate:"required"`
	ProductId string   `json:"productId" validate:"required"`
	Scopes    []string `json:"scopes" validate:"required"`
}

type CreateApiKeyResponse struct {
	Key string `json:"key"`
}
