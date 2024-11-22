package service_controller

import (
	"context"
	"github.com/devingen/kimlik-api"
	"github.com/devingen/kimlik-api/dto"
	"net/http"

	core "github.com/devingen/api-core"
)

func (c ServiceController) GetUserInfo(ctx context.Context, req core.Request) (*core.Response, error) {

	base, err := req.AssertPathParameter("base")
	if err != nil {
		return nil, err
	}

	tokenPayload, err := kimlik.AssertAuthentication(ctx)
	if err != nil {
		return nil, err
	}
	userID := tokenPayload.UserID

	user, err := c.DataService.FindUserWithId(ctx, base, userID)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body: dto.GetUserInfoResponse{
			Sub:               user.ID.Hex(),
			Name:              user.FullName(),
			GivenName:         *user.FirstName,
			FamilyName:        *user.LastName,
			PreferredUsername: "",
			Email:             *user.Email,
			Picture:           "",
		},
	}, nil
}
