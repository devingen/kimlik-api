package dto

type LoginWithEmailRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginWithEmailResponse struct {
	UserID string `json:"userId"`
	JWT    string `json:"jwt"`
}
