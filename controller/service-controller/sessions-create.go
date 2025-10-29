package service_controller

import (
	"context"
	"net/http"
	"time"

	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/model"
	token_service "github.com/devingen/kimlik-api/token-service"
	"golang.org/x/crypto/bcrypt"
)

func (c ServiceController) CreateSession(ctx context.Context, req core.Request) (*core.Response, error) {
	base, hasBase := req.PathParameters["base"]
	if !hasBase {
		return nil, core.NewError(http.StatusInternalServerError, "missing-path-param-base")
	}

	var body dto.CreateSession
	err := req.AssertBody(&body)
	if err != nil {
		return nil, err
	}

	var isNewUser bool
	var auth *model.Auth
	var user *model.User
	if body.Email != nil && body.Password != nil {
		auth, user, err = c.validateSessionWithPassword(ctx, base, *body.Email, *body.Password)
	} else if body.IDToken != nil {
		auth, user, isNewUser, err = c.validateSessionWithIDToken(ctx, base, *body.IDToken, "", "")
	} else {
		return nil, core.NewError(http.StatusBadRequest, "no-authentication-method-found-for-provided-parameters")
	}
	if err != nil {
		if user != nil {
			c.createFailedSession(ctx, req, base, auth, user, err)
		}
		return nil, err
	}

	jwt, _, err := c.createSuccessfulSessionAndGenerateToken(ctx, req, base, auth, user)
	if err != nil {
		return nil, err
	}

	responseStatusCode := http.StatusOK
	if isNewUser {
		responseStatusCode = http.StatusCreated
	}

	return &core.Response{
		StatusCode: responseStatusCode,
		Body: dto.LoginResponse{
			UserID: user.ID.Hex(),
			JWT:    jwt,
		},
	}, nil
}

func (c ServiceController) validateSessionWithIDToken(ctx context.Context, base, rawIDToken, givenName, familyName string) (*model.Auth, *model.User, bool, error) {

	idToken := IDToken{RawIDToken: rawIDToken}
	if err := idToken.Parse(); err != nil {
		return nil, nil, false, core.NewError(http.StatusBadRequest, "could-not-parse-token:"+err.Error())
	}

	// TODO find the open id config from database with the audience and issuer in the token
	if err := idToken.Verify(ctx); err != nil {
		return nil, nil, false, core.NewError(http.StatusBadRequest, "could-not-verify-token:"+err.Error())
	}

	email := idToken.Claims.Email
	if email == "" {
		return nil, nil, false, core.NewError(http.StatusBadRequest, "email-missing-in-token-claims")
	}

	user, err := c.DataService.FindUserWithEmail(ctx, base, email)
	if err != nil {
		return nil, nil, false, err
	}

	if user != nil {
		auth, err := c.DataService.FindOIDCAuthOfUser(ctx, base, user.ID.Hex(), idToken.Claims.Issuer)
		if err != nil {
			return nil, nil, false, err
		}
		if auth == nil {
			// User exists but has no Open ID authentication method with this issuer
			// Meaning that the account is not created via Open ID with this issuer
			auth, err = c.DataService.CreateAuthWithIDToken(ctx, base, idToken.Claims.ToMap(), user)
			if err != nil {
				return nil, nil, false, err
			}
			return auth, user, false, nil
		}
		return auth, user, false, nil
	}

	firstName := idToken.Claims.GivenName
	if firstName == "" {
		if givenName == "" {
			return nil, nil, false, core.NewError(http.StatusBadRequest, "given-name-must-exist-in-claims-or-body-for-new-users")
		}
		firstName = givenName
	}

	lastName := idToken.Claims.FamilyName
	if lastName == "" {
		if familyName == "" {
			return nil, nil, false, core.NewError(http.StatusBadRequest, "family-name-must-exist-in-claims-or-body-for-new-users")
		}
		lastName = familyName
	}
	isEmailVerified := idToken.Claims.IsEmailVerified()

	user, err = c.DataService.CreateUser(ctx, base, firstName, lastName, email, model.UserStatusActive, isEmailVerified)
	if err != nil {
		return nil, nil, false, err
	}

	auth, err := c.DataService.CreateAuthWithIDToken(ctx, base, idToken.Claims.ToMap(), user)
	if err != nil {
		return nil, nil, false, err
	}
	return auth, user, true, nil
}

func (c ServiceController) validateSessionWithPassword(ctx context.Context, base, email, password string) (*model.Auth, *model.User, error) {
	user, err := c.DataService.FindUserWithEmail(ctx, base, email)
	if err != nil {
		return nil, nil, err
	}

	if user == nil {
		return nil, nil, core.NewError(http.StatusNotFound, "user-not-found")
	}

	auth, err := c.DataService.FindPasswordAuthOfUser(ctx, base, user.ID.Hex())
	if err != nil {
		return nil, user, err
	}

	if auth == nil {
		return nil, user, core.NewError(http.StatusBadRequest, "authentication-method-not-found-for-user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(password)); err != nil {
		return auth, user, core.NewError(http.StatusUnauthorized, "password-mismatch")
	}

	return auth, user, nil
}

func (c ServiceController) createSuccessfulSessionAndGenerateToken(ctx context.Context, req core.Request, base string, auth *model.Auth, user *model.User) (string, string, error) {
	userAgent := req.Headers["user-agent"]
	client := req.Headers["client"]
	ip := req.IP

	refreshToken, err := c.generateRefreshToken()
	if err != nil {
		return "", "", err
	}

	session, err := c.DataService.CreateSession(ctx, base, client, userAgent, ip, refreshToken.HashedToken, "", auth, user)
	if err != nil {
		return "", "", err
	}

	accessToken, err := c.generateAccessToken(user.ID.Hex(), session.ID.Hex())
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken.RawToken, nil
}

func (c ServiceController) generateRefreshToken() (*token_service.RefreshToken, error) {
	return c.TokenService.GenerateRefreshToken()
}

func (c ServiceController) generateAccessToken(userID, sessionID string) (string, error) {
	return c.TokenService.GenerateAccessToken(
		userID,
		sessionID,
		[]string{ScopeAll},
		time.Now().Add(AccessTokenExpirationTime).Unix(),
	)
}

func (c ServiceController) createFailedSession(ctx context.Context, req core.Request, base string, auth *model.Auth, user *model.User, err error) {
	userAgent := req.Headers["user-agent"]
	client := req.Headers["client"]
	ip := req.IP
	c.DataService.CreateSession(ctx, base, client, userAgent, ip, "", err.Error(), auth, user)
	// TODO print error properly
}
