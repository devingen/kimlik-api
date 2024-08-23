package service_controller

import (
	"context"
	"net/http"
	"time"

	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson"
)

const AccessTokenExpirationTime = 240 * time.Hour

// OAuthToken authenticates user and returns access token with given grant type.
func (c ServiceController) OAuthToken(ctx context.Context, req core.Request) (*core.Response, error) {
	// TODO check client ID and secrets
	// TODO get these values as "Content-Type: application/x-www-form-urlencoded" ?????

	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	var params dto.OAuthTokenRequest
	err := req.AssertBody(&params)
	if err != nil {
		return nil, err
	}

	switch params.GrantType {
	case dto.GrantTypeOIDC:
		return c.handleGrantTypeKimlikOIDC(ctx, req, base, params)
	case dto.GrantTypePassword:
		return c.handleGrantTypePassword(ctx, req, base, params)
	case dto.GrantTypeRefreshToken:
		return c.handleGrantTypeRefreshToken(ctx, base, params)
	}

	return nil, core.NewError(http.StatusBadRequest, "invalid-grant-type")
}

func (c ServiceController) handleGrantTypePassword(ctx context.Context, req core.Request, base string, params dto.OAuthTokenRequest) (*core.Response, error) {
	if params.Username == nil {
		return nil, core.NewError(http.StatusBadRequest, "username-missing")
	}
	auth, user, err := c.validateSessionWithPassword(ctx, base, *params.Username, *params.Password)

	accessToken, refreshToken, err := c.createSuccessfulSessionAndGenerateToken(ctx, req, base, auth, user)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body: dto.OAuthTokenResponse{
			AccessToken:  accessToken,
			TokenType:    "Bearer",
			ExpiresIn:    AccessTokenExpirationTime.Seconds(),
			RefreshToken: refreshToken,
		},
	}, nil
}

func (c ServiceController) handleGrantTypeKimlikOIDC(ctx context.Context, req core.Request, base string, params dto.OAuthTokenRequest) (*core.Response, error) {
	auth, user, isNewUser, err := c.validateSessionWithIDToken(ctx, base, *params.IDToken)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := c.createSuccessfulSessionAndGenerateToken(ctx, req, base, auth, user)
	if err != nil {
		return nil, err
	}

	responseStatusCode := http.StatusOK
	if isNewUser {
		responseStatusCode = http.StatusCreated
	}

	return &core.Response{
		StatusCode: responseStatusCode,
		Body: dto.OAuthTokenResponse{
			AccessToken:  accessToken,
			TokenType:    "Bearer",
			ExpiresIn:    AccessTokenExpirationTime.Seconds(),
			RefreshToken: refreshToken,
		},
	}, nil

}

func (c ServiceController) handleGrantTypeRefreshToken(ctx context.Context, base string, params dto.OAuthTokenRequest) (*core.Response, error) {
	sessions, err := c.DataService.FindSessions(ctx, base, bson.M{
		"refreshToken": c.TokenService.HashRefreshToken(*params.RefreshToken),
	})
	if err != nil {
		return nil, err
	}
	if len(sessions) == 0 {
		return nil, core.NewStatusError(http.StatusUnauthorized)
	}
	session := sessions[0]

	refreshToken, err := c.generateRefreshToken()
	if err != nil {
		return nil, err
	}

	accessToken, err := c.generateAccessToken(session.User.ID.Hex(), session.ID.Hex())
	if err != nil {
		return nil, err
	}

	_, _, err = c.DataService.UpdateSession(ctx, base, &model.Session{
		ID:           session.ID,
		RefreshToken: &refreshToken.HashedToken,
	})
	if err != nil {
		return nil, err
	}

	return &core.Response{
		StatusCode: http.StatusOK,
		Body: dto.OAuthTokenResponse{
			AccessToken:  accessToken,
			TokenType:    "Bearer",
			ExpiresIn:    AccessTokenExpirationTime.Seconds(),
			RefreshToken: refreshToken.RawToken,
		},
	}, nil
}
