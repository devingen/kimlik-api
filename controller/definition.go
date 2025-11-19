package controller

import (
	"context"

	core "github.com/devingen/api-core"
)

type IServiceController interface {
	// Setup creates initial configuration.
	Setup(ctx context.Context, req core.Request) (*core.Response, error)

	// OAuth2GetOIDCConfiguration returns the .well-known/openid-configuration.
	OAuth2GetOIDCConfiguration(ctx context.Context, req core.Request) (*core.Response, error)

	// OAuth2GetJWKS returns the .well-known/jwks.json.
	OAuth2GetJWKS(ctx context.Context, req core.Request) (*core.Response, error)

	// OAuth2Token is used to get an access token or a refresh token.
	// https://datatracker.ietf.org/doc/html/rfc6749#section-2.3
	OAuth2Token(ctx context.Context, req core.Request) (*core.Response, error)

	// OAuth2Authorize is used to generate authorization code and redirect users to client redirect uri.
	// https://datatracker.ietf.org/doc/html/rfc6749#section-4.1.1
	OAuth2Authorize(ctx context.Context, req core.Request) (*core.Response, error)

	// GetUserInfo is used as userinfo endpoint for OpenID protocol and returns user's details.
	// https://openid.net/specs/openid-connect-core-1_0.html#UserInfo
	GetUserInfo(ctx context.Context, req core.Request) (*core.Response, error)

	// GetSession returns the current session details within the authorization header.
	GetSession(ctx context.Context, req core.Request) (*core.Response, error)

	// CreateSession logs in existing user (password), registers new user (openid connect) and creates & returns new session.
	CreateSession(ctx context.Context, req core.Request) (*core.Response, error)

	// AnonymizeUser removes all the personal details from the User and deletes all authentication methods of the user
	AnonymizeUser(ctx context.Context, req core.Request) (*core.Response, error)

	// GetAuthorizationURL generates and returns authentication URL for SSO login.
	GetAuthorizationURL(ctx context.Context, req core.Request) (*core.Response, error)

	// GetAuthMethods returns available authentication methods for a given email address.
	GetAuthMethods(ctx context.Context, req core.Request) (*core.Response, error)

	Authenticate(ctx context.Context, req core.Request) (*core.Response, error)
	RegisterWithEmail(ctx context.Context, req core.Request) (*core.Response, error)
	LoginWithEmail(ctx context.Context, req core.Request) (*core.Response, error)
	ActivateUser(ctx context.Context, req core.Request) (*core.Response, error)
	ChangePassword(ctx context.Context, req core.Request) (*core.Response, error)

	FindUsers(ctx context.Context, req core.Request) (*core.Response, error)

	CreateAPIKey(ctx context.Context, req core.Request) (*core.Response, error)
	FindAPIKeys(ctx context.Context, req core.Request) (*core.Response, error)
	UpdateAPIKey(ctx context.Context, req core.Request) (*core.Response, error)
	DeleteAPIKey(ctx context.Context, req core.Request) (*core.Response, error)
	VerifyAPIKey(ctx context.Context, req core.Request) (*core.Response, error)

	CreateAppIntegration(ctx context.Context, req core.Request) (*core.Response, error)
	FindAppIntegrations(ctx context.Context, req core.Request) (*core.Response, error)
	UpdateAppIntegration(ctx context.Context, req core.Request) (*core.Response, error)
	DeleteAppIntegration(ctx context.Context, req core.Request) (*core.Response, error)

	CreateOAuth2Config(ctx context.Context, req core.Request) (*core.Response, error)
	FindOAuth2Configs(ctx context.Context, req core.Request) (*core.Response, error)
	UpdateOAuth2Config(ctx context.Context, req core.Request) (*core.Response, error)
	DeleteOAuth2Config(ctx context.Context, req core.Request) (*core.Response, error)

	CreateSAMLConfig(ctx context.Context, req core.Request) (*core.Response, error)
	FindSAMLConfigs(ctx context.Context, req core.Request) (*core.Response, error)
	UpdateSAMLConfig(ctx context.Context, req core.Request) (*core.Response, error)
	DeleteSAMLConfig(ctx context.Context, req core.Request) (*core.Response, error)
	BuildSAMLAuthURL(ctx context.Context, req core.Request) (*core.Response, error)
	ConsumeSAMLAuthResponse(ctx context.Context, req core.Request) (*core.Response, error)
	LoginWithSAML(ctx context.Context, req core.Request) (*core.Response, error)

	GetTenantInfo(ctx context.Context, req core.Request) (*core.Response, error)
	UpdateTenantInfo(ctx context.Context, req core.Request) (*core.Response, error)
}
