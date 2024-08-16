package service_controller

import (
	"context"
	"net/http"

	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/model"
	token_service "github.com/devingen/kimlik-api/token-service"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
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
		auth, user, isNewUser, err = c.validateSessionWithIDToken(ctx, base, *body.IDToken)
	} else {
		return nil, core.NewError(http.StatusBadRequest, "no-authentication-method-found-for-provided-parameters")
	}
	if err != nil {
		if user != nil {
			c.createFailedSession(ctx, req, base, auth, user, err)
		}
		return nil, err
	}

	jwt, err := c.createSuccessfulSessionAndGenerateToken(ctx, req, base, auth, user)
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

func (c ServiceController) validateSessionWithIDToken(ctx context.Context, base, idToken string) (*model.Auth, *model.User, bool, error) {

	// TODO find the open id config from database with the audience and issuer in the token
	// TODO this validate is for Google. Check others later
	payload, err := idtoken.Validate(ctx, idToken, "")
	if err != nil {
		return nil, nil, false, core.NewError(http.StatusBadRequest, "could-not-validate-token")
	}

	email, ok := payload.Claims["email"].(string)
	if !ok {
		return nil, nil, false, core.NewError(http.StatusBadRequest, "email-missing-in-token-claims")
	}

	user, err := c.DataService.FindUserWithEmail(ctx, base, email)
	if err != nil {
		return nil, nil, false, err
	}

	if user != nil {
		auth, err := c.DataService.FindAuthOfUser(ctx, base, user.ID.Hex(), model.AuthTypeOpenID)
		if err != nil {
			return nil, nil, false, err
		}
		if auth == nil {
			// User exists but has no Open ID authentication method meaning that the account is not created via Open ID
			auth, err = c.DataService.CreateAuthWithIDToken(ctx, base, payload.Claims, user)
			if err != nil {
				return nil, nil, false, err
			}
			return auth, user, false, nil
		}
		return auth, user, false, nil
	}

	firstName, ok := payload.Claims["given_name"].(string)
	if !ok {
		return nil, nil, false, core.NewError(http.StatusBadRequest, "first-name-missing-in-token-claims")
	}

	lastName, ok := payload.Claims["family_name"].(string)
	if !ok {
		return nil, nil, false, core.NewError(http.StatusBadRequest, "last-name-missing-in-token-claims")
	}

	var isEmailVerified bool
	isEmailVerifiedBoolValue, hasEmailVerifiedBool := payload.Claims["email_verified"].(bool)
	if hasEmailVerifiedBool {
		isEmailVerified = isEmailVerifiedBoolValue
	} else {
		isEmailVerifiedStringValue, hasEmailVerifiedString := payload.Claims["email_verified"].(string)
		if hasEmailVerifiedString {
			isEmailVerified = isEmailVerifiedStringValue == "true"
		} else {
			isEmailVerified = false
		}
	}

	user, err = c.DataService.CreateUser(ctx, base, firstName, lastName, email, model.UserStatusActive, isEmailVerified)
	if err != nil {
		return nil, nil, false, err
	}

	auth, err := c.DataService.CreateAuthWithIDToken(ctx, base, payload.Claims, user)
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
		return nil, nil, core.NewStatusError(http.StatusNotFound)
	}

	auth, err := c.DataService.FindAuthOfUser(ctx, base, user.ID.Hex(), model.AuthTypePassword)
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

func (c ServiceController) createSuccessfulSessionAndGenerateToken(ctx context.Context, req core.Request, base string, auth *model.Auth, user *model.User) (string, error) {
	userAgent := req.Headers["user-agent"]
	client := req.Headers["client"]
	ip := req.IP
	session, err := c.DataService.CreateSession(ctx, base, client, userAgent, ip, "", auth, user)
	if err != nil {
		return "", err
	}

	jwt, err := c.TokenService.GenerateToken(
		user.ID.Hex(),
		session.ID.Hex(),
		[]token_service.Scope{ScopeAll},
		240,
	)
	if err != nil {
		return "", err
	}
	return jwt, nil
}

func (c ServiceController) createFailedSession(ctx context.Context, req core.Request, base string, auth *model.Auth, user *model.User, err error) {
	userAgent := req.Headers["user-agent"]
	client := req.Headers["client"]
	ip := req.IP
	c.DataService.CreateSession(ctx, base, client, userAgent, ip, err.Error(), auth, user)
	// TODO print error properly
}
