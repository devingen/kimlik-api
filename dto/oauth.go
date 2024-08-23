package dto

type GrantType string

const (
	// GrantTypeOIDC Authenticates user via ID Token taken from provider SDK. Requires 'id_token' in the request params.
	GrantTypeOIDC GrantType = "https://kimlik.devingen.io/oauth/grant-type/oidc"

	// GrantTypePassword uses username and password to authenticate user.
	GrantTypePassword GrantType = "password"

	// GrantTypeRefreshToken uses an existing refresh token to generate new access and refresh tokens.
	// Client needs to continue using the new refresh token as the old refresh token expires.
	GrantTypeRefreshToken GrantType = "refresh_token"
)

type OAuthTokenRequest struct {
	//ClientID     *string `json:"client_id" validate:"required,oneof=refresh_token"`
	//ClientSecret *string `json:"client_secret" validate:"required,oneof=refresh_token"`

	// GrantType defines the way to authenticate.
	// See https://auth0.com/docs/get-started/applications/application-grant-types
	GrantType    GrantType `json:"grant_type" validate:"required,oneof=password refresh_token https://kimlik.devingen.io/oauth/grant-type/oidc"`
	RefreshToken *string   `json:"refresh_token" validate:"required_if=GrantType refresh_token"`

	// IDToken is required if grant type is GrantTypeOIDC
	IDToken *string `json:"id_token"`

	// Username is required if grant type is GrantTypePassword
	Username *string `json:"username" validate:"required_if=GrantType password"`

	// Password is required if grant type is GrantTypePassword
	Password *string `json:"password" validate:"required_if=GrantType password"`
}

type OAuthTokenResponse struct {
	AccessToken  string  `json:"access_token"`
	TokenType    string  `json:"token_type"`
	ExpiresIn    float64 `json:"expires_in"`
	RefreshToken string  `json:"refresh_token"`
}
