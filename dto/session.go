package dto

import "github.com/devingen/kimlik-api/model"

type AuthorizationType string

const (
	// AuthorizationTypePassword Authenticates user via email and password.
	AuthorizationTypePassword AuthorizationType = "password"

	// AuthorizationTypeCode Authenticates user via the code returned from an external OAuth2 Idp.
	AuthorizationTypeCode AuthorizationType = "authorization_code"

	// AuthorizationTypeOIDC Authenticates user via the id_token returned from an external OIDC Idp.
	AuthorizationTypeOIDC AuthorizationType = "oidc"
)

type AuthorizeRequest struct {
	// AuthType defines the way to authenticate.
	AuthType AuthorizationType `json:"auth_type" validate:"required,oneof=password authorization_code oidc"`

	// Username required if grant type is password
	Username *string `json:"username" validate:"required_if=AuthorizationType password"`

	// Password required if grant type is password
	Password *string `json:"password" validate:"required_if=AuthorizationType password"`

	// IDToken required if grant type is oidc
	IDToken *string `json:"id_token" validate:"required_if=AuthorizationType oidc"`

	// GivenName required if grant type is oidc, user is new and id token claims doesn't contain given_name (ex Apple)
	GivenName *string `json:"given_name"`

	// FamilyName required if grant type is oidc, user is new and id token claims doesn't contain family_name (ex Apple)
	FamilyName *string `json:"family_name"`

	// Code required if grant type is authorization_code
	Code *string `json:"code" validate:"required_if=AuthorizationType authorization_code"`

	// State required if grant type is authorization_code
	State *string `json:"state" validate:"required_if=AuthorizationType authorization_code"`
}

type ActivateUserRequest struct {
	UserActivationToken string `json:"userActivationToken" validate:"required"`
	FirstName           string `json:"firstName" validate:"min=2,max=32"`
	LastName            string `json:"lastName" validate:"min=2,max=32"`
	Password            string `json:"password" validate:"min=6,max=32"`
}

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
