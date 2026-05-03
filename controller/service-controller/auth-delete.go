package service_controller

import (
	"context"
	"net/http"

	core "github.com/devingen/api-core"
	kimlik "github.com/devingen/kimlik-api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeleteAuthMethod deletes an authentication method. Only the owner of the auth record can delete it.
func (c ServiceController) DeleteAuthMethod(ctx context.Context, req core.Request) (*core.Response, error) {
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

	authID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, core.NewError(http.StatusBadRequest, "invalid-auth-id")
	}

	auths, err := c.DataService.FindAuthsOfUser(ctx, base, tokenPayload.UserID)
	if err != nil {
		return nil, err
	}

	var found bool
	for _, auth := range auths {
		if auth.ID == authID {
			found = true
			break
		}
	}
	if !found {
		return nil, core.NewError(http.StatusForbidden, "forbidden")
	}

	err = c.DataService.DeleteAuthMethod(ctx, base, authID)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body:       map[string]interface{}{},
	}, nil
}
