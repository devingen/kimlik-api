package dto

import "github.com/devingen/kimlik-api/model"

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

type LoginWithEmailRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
}

type LoginWithEmailResponse struct {
	UserID string `json:"userId"`
	JWT    string `json:"jwt"`
}

type ChangePasswordRequest struct {
	Password string `json:"password" validate:"min=6,max=32"`
}

type ChangePasswordResponse struct {
}

type GetSessionResponse struct {
	User    *model.User    `json:"user"`
	Session *model.Session `json:"session"`
}
