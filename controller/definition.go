package controller

import (
	"context"
	"github.com/devingen/api-core/dvnruntime"
)

type ControllerFunc func(ctx context.Context, req dvnruntime.Request) (interface{}, int, error)

type IServiceController interface {
	RegisterWithEmail(ctx context.Context, req dvnruntime.Request) (interface{}, int, error)
	LoginWithEmail(ctx context.Context, req dvnruntime.Request) (interface{}, int, error)
	GetSession(ctx context.Context, req dvnruntime.Request) (interface{}, int, error)
}
