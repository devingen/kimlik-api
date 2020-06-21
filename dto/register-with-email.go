package dto

type RegisterWithEmailRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type RegisterWithEmailResponse struct {
	UserID string `json:"userId"`
	JWT    string `json:"jwt"`
}
