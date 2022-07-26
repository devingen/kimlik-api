package service_controller

import (
	"context"
	core "github.com/devingen/api-core"
	core_dto "github.com/devingen/api-core/dto"
	"github.com/devingen/kimlik-api"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/model"
	"net/http"
)

func (c ServiceController) ChangePassword(ctx context.Context, req core.Request) (*core.Response, error) {

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	var body dto.ChangePasswordRequest
	token, err := kimlik.AssertAuthenticationAndBody(ctx, req, &body)
	if err != nil {
		return nil, err
	}

	auth, err := c.DataService.FindAuthOfUser(ctx, base, token.UserId, model.AuthTypePassword)
	if err != nil {
		return nil, err
	}

	if auth == nil {
		return nil, core.NewError(http.StatusInternalServerError, "auth-missing")
	}
	auth.Password = body.Password

	updatedAt, revision, err := c.DataService.UpdateAuth(ctx, base, auth)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body: core_dto.UpdateEntryResponse{
			ID:        auth.ID.Hex(),
			UpdatedAt: *updatedAt,
			Revision:  revision,
		},
	}, nil
}
