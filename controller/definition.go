package controller

import (
	"context"
	core "github.com/devingen/api-core"
)

type IServiceController interface {
	RegisterWithEmail(ctx context.Context, req core.Request) (interface{}, int, error)
	LoginWithEmail(ctx context.Context, req core.Request) (interface{}, int, error)
	GetSession(ctx context.Context, req core.Request) (interface{}, int, error)
	ChangePassword(ctx context.Context, req core.Request) (interface{}, int, error)
	CreateAPIKey(ctx context.Context, req core.Request) (interface{}, int, error)

	CreateSAMLConfig(ctx context.Context, req core.Request) (interface{}, int, error)
	BuildSAMLAuthURL(ctx context.Context, req core.Request) (interface{}, int, error)
	ConsumeSAMLAuthResponse(ctx context.Context, req core.Request) (interface{}, int, error)
}
