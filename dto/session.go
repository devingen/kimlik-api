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

type CreateSession struct {
	// exist if authentication type is 'password'
	Email    *string `json:"email" validate:"required_without=IDToken,omitempty,email"`
	Password *string `json:"password" validate:"required_without=IDToken,omitempty"`

	// exist if authentication type is 'openid'
	IDToken *string `json:"idToken" validate:"required_without_all=Email Password"`
}

type LoginResponse struct {
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
