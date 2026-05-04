package service_controller

import (
	"context"
	"net/http"

	core "github.com/devingen/api-core"
	core_dto "github.com/devingen/api-core/dto"
	kimlik "github.com/devingen/kimlik-api"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/model"
)

// UpdateUser updates a user's name and email. Only the authenticated user can update their own record.
func (c ServiceController) UpdateUser(ctx context.Context, req core.Request) (*core.Response, error) {
	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	id, hasID := req.PathParameters["id"]
	if !hasID {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-id")
	}

	tokenPayload, err := kimlik.AssertAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	if tokenPayload.UserID != id {
		return nil, core.NewError(http.StatusForbidden, "forbidden")
	}

	var body dto.UpdateUserRequest
	err = req.AssertBody(&body)
	if err != nil {
		return nil, err
	}

	user, err := c.DataService.FindUserWithId(ctx, base, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, core.NewError(http.StatusNotFound, "user-not-found")
	}

	update := &model.User{
		ID:        user.ID,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
	}
	if body.FirstName != nil && body.LastName != nil {
		fullName := *body.FirstName + " " + *body.LastName
		update.Name = &fullName
	}

	updatedAt, revision, err := c.DataService.UpdateUser(ctx, base, update)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body: core_dto.UpdateEntryResponse{
			ID:        user.ID.Hex(),
			UpdatedAt: *updatedAt,
			Revision:  revision,
		},
	}, nil
}
