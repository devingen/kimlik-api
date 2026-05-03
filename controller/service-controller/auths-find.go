package service_controller

import (
	"context"
	"net/http"

	core "github.com/devingen/api-core"
	kimlik "github.com/devingen/kimlik-api"
	"github.com/devingen/kimlik-api/dto"
)

// FindAuths returns all authentication methods of the authenticated user.
func (c ServiceController) FindAuths(ctx context.Context, req core.Request) (*core.Response, error) {
	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	tokenPayload, err := kimlik.AssertAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	auths, err := c.DataService.FindAuthsOfUser(ctx, base, tokenPayload.UserID)
	if err != nil {
		return nil, err
	}

	response := make([]dto.AuthResponse, len(auths))
	for i, auth := range auths {
		response[i] = dto.AuthResponse{
			ID:   auth.ID.Hex(),
			Type: string(auth.Type),
		}
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body:       response,
	}, nil
}
