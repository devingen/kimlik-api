package wrapper

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api"
)

// WithAuth wraps the controller func by adding the authentication to the context.
func WithAuth(f core.Controller, signKey string) core.Controller {
	return func(ctx context.Context, req core.Request) (interface{}, int, error) {

		// add auth to the context
		ctxWithAuth, err := kimlik.WithJWTAuth(ctx, req, signKey)
		if err != nil {
			return nil, 0, err
		}

		ctxWithAuth, err = kimlik.WithAPIKeyAuth(ctxWithAuth, req)
		if err != nil {
			return nil, 0, err
		}

		// execute function
		result, status, err := f(ctxWithAuth, req)

		return result, status, err
	}
}
