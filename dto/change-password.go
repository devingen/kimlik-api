package dto

type ChangePasswordRequest struct {
	Password string `json:"password" validate:"min=6,max=32"`
}

type ChangePasswordResponse struct {
}
