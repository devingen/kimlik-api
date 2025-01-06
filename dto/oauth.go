package dto

type OAuth2GrantType string

const (
	// OAuth2GrantTypeAuthorizationCode Authenticates user via the code returned from the /authorize step.
	OAuth2GrantTypeAuthorizationCode OAuth2GrantType = "authorization_code"

	// OAuth2GrantTypeOIDC Authenticates user via ID Token taken from provider SDK. Requires 'id_token' in the request params.
	OAuth2GrantTypeOIDC OAuth2GrantType = "https://kimlik.devingen.io/oauth/grant-type/oidc"

	// OAuth2GrantTypePassword uses username and password to authenticate user.
	OAuth2GrantTypePassword OAuth2GrantType = "password"

	// OAuth2GrantTypeRefreshToken uses an existing refresh token to generate new access and refresh tokens.
	// Client needs to continue using the new refresh token as the old refresh token expires.
	OAuth2GrantTypeRefreshToken OAuth2GrantType = "refresh_token"
)

type OAuth2TokenRequest struct {
	//ClientID     *string `json:"client_id" validate:"required,oneof=refresh_token"`
	//ClientSecret *string `json:"client_secret" validate:"required,oneof=refresh_token"`

	// GrantType defines the way to authenticate.
	// See https://auth0.com/docs/get-started/applications/application-grant-types
	GrantType    OAuth2GrantType `json:"grant_type" validate:"required,oneof=password authorization_code refresh_token https://kimlik.devingen.io/oauth/grant-type/oidc"`
	RefreshToken *string         `json:"refresh_token" validate:"required_if=GrantType refresh_token"`

	// IDToken is required if grant type is OAuth2GrantTypeOIDC
	IDToken *string `json:"id_token"`

	// GivenName is required if grant type is OAuth2GrantTypeOIDC, user is new and id token claims doesn't contain given_name
	GivenName *string `json:"given_name"`

	// FamilyName is required if grant type is OAuth2GrantTypeOIDC, user is new and id token claims doesn't contain family_name
	FamilyName *string `json:"family_name"`

	// Username is required if grant type is OAuth2GrantTypePassword
	Username *string `json:"username" validate:"required_if=GrantType password"`

	// Password is required if grant type is OAuth2GrantTypePassword
	Password *string `json:"password" validate:"required_if=GrantType password"`

	// ClientID is required if grant type is OAuth2GrantTypeAuthorizationCode
	ClientID *string `json:"client_id" validate:"required_if=GrantType authorization_code"`

	// RedirectURI is required if grant type is OAuth2GrantTypeAuthorizationCode
	RedirectURI *string `json:"redirect_uri" validate:"required_if=GrantType authorization_code"`

	// Code is required if grant type is OAuth2GrantTypeAuthorizationCode
	Code *string `json:"code" validate:"required_if=GrantType authorization_code"`

	// CodeVerifier is required if grant type is OAuth2GrantTypeAuthorizationCode
	CodeVerifier *string `json:"code_verifier" validate:"required_if=GrantType authorization_code"`
}

type OAuth2TokenResponse struct {
	AccessToken  string  `json:"access_token,omitempty"`
	TokenType    string  `json:"token_type,omitempty"`
	ExpiresIn    float64 `json:"expires_in,omitempty"`
	RefreshToken string  `json:"refresh_token,omitempty"`
	IDToken      string  `json:"id_token,omitempty"`
}
