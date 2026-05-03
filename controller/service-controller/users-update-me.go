package service_controller

import (
	"context"
	"net/http"

	core "github.com/devingen/api-core"
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

	fullName := body.FirstName + " " + body.LastName
	_, _, err = c.DataService.UpdateUser(ctx, base, &model.User{
		ID:        user.ID,
		FirstName: &body.FirstName,
		LastName:  &body.LastName,
		Email:     &body.Email,
		Name:      &fullName,
	})
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body: dto.GetUserInfoResponse{
			Sub:        user.ID.Hex(),
			Name:       fullName,
			GivenName:  body.FirstName,
			FamilyName: body.LastName,
			Email:      body.Email,
		},
	}, nil
}
