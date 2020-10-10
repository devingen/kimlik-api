package dto

type RegisterWithEmailRequest struct {
	FirstName string `json:"firstName" validate:"min=2,max=32"`
	LastName  string `json:"lastName" validate:"min=2,max=32"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"min=6,max=32"`
}

type RegisterWithEmailResponse struct {
	UserID string `json:"userId"`
	JWT    string `json:"jwt"`
}
