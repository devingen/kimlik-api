package controller

import (
	"context"

	core "github.com/devingen/api-core"
)

type IServiceController interface {
	// GetSession returns the current session details within the authorizarion header.
	GetSession(ctx context.Context, req core.Request) (*core.Response, error)

	// CreateSession logs in existing user (password) and registers new user (openid connect) and creates/returns new session.
	CreateSession(ctx context.Context, req core.Request) (*core.Response, error)

	RegisterWithEmail(ctx context.Context, req core.Request) (*core.Response, error)
	LoginWithEmail(ctx context.Context, req core.Request) (*core.Response, error)
	ChangePassword(ctx context.Context, req core.Request) (*core.Response, error)

	FindUsers(ctx context.Context, req core.Request) (*core.Response, error)

	CreateAPIKey(ctx context.Context, req core.Request) (*core.Response, error)
	FindAPIKeys(ctx context.Context, req core.Request) (*core.Response, error)
	UpdateAPIKey(ctx context.Context, req core.Request) (*core.Response, error)
	DeleteAPIKey(ctx context.Context, req core.Request) (*core.Response, error)
	VerifyAPIKey(ctx context.Context, req core.Request) (*core.Response, error)

	CreateSAMLConfig(ctx context.Context, req core.Request) (*core.Response, error)
	FindSAMLConfigs(ctx context.Context, req core.Request) (*core.Response, error)
	UpdateSAMLConfig(ctx context.Context, req core.Request) (*core.Response, error)
	DeleteSAMLConfig(ctx context.Context, req core.Request) (*core.Response, error)
	BuildSAMLAuthURL(ctx context.Context, req core.Request) (*core.Response, error)
	ConsumeSAMLAuthResponse(ctx context.Context, req core.Request) (*core.Response, error)
	LoginWithSAML(ctx context.Context, req core.Request) (*core.Response, error)
}
