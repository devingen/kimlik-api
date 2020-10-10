package controller

import (
	"context"
	"github.com/devingen/api-core/dvnruntime"
)

type IServiceController interface {
	RegisterWithEmail(ctx context.Context, req dvnruntime.Request) (interface{}, int, error)
	LoginWithEmail(ctx context.Context, req dvnruntime.Request) (interface{}, int, error)
	GetSession(ctx context.Context, req dvnruntime.Request) (interface{}, int, error)
	ChangePassword(ctx context.Context, req dvnruntime.Request) (interface{}, int, error)
}
