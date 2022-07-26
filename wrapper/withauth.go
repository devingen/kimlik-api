package wrapper

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api"
	ds "github.com/devingen/kimlik-api/data-service"
	token_service "github.com/devingen/kimlik-api/token-service"
)

// WithAuth wraps the controller func by adding the authentication to the context.
func WithAuth(f core.Controller, jwtService token_service.ITokenService, dataService ds.IKimlikDataService) core.Controller {
	return func(ctx context.Context, req core.Request) (*core.Response, error) {

		// add auth to the context
		ctxWithAuth, err := kimlik.WithJWTAuth(jwtService, ctx, req)
		if err != nil {
			return nil, err
		}

		ctxWithAuth, err = kimlik.WithAPIKeyAuth(ctxWithAuth, req, dataService)
		if err != nil {
			return nil, err
		}

		// execute function
		response, err := f(ctxWithAuth, req)

		return response, err
	}
}
