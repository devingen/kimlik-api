package dto

type SetupRequest struct {
	TenantName string `json:"tenantName" validate:"min=2,max=32"`
}

type SetupResponse struct {
}
