package service

import (
	"context"
	"github.com/devingen/kimlik-api/model"
	"time"
)

type IKimlikDataService interface {
	CreateAPIKey(base, name, productId string, scopes []string, keyPrefix, hash string, user *model.User) (*model.ApiKey, error)
	CreateAuthWithPassword(base, password string, user *model.User) (*model.Auth, error)
	CreateSession(base, client, userAgent, ip string, user *model.User) (*model.Session, error)
	CreateUser(base, firstName, lastName, email string) (*model.User, error)
	FindAuthOfUser(base, userId string, authType model.AuthType) (*model.Auth, error)
	FindUserUserWithEmail(base, email string) (*model.User, error)
	FindUserUserWithId(base, id string) (*model.User, error)
	UpdateAuth(base string, auth *model.Auth) (*time.Time, int, error)

	CreateSAMLConfig(ctx context.Context, base string, item *model.SAMLConfig) (*model.SAMLConfig, error)
	GetSAMLConfig(ctx context.Context, base, id string) (*model.SAMLConfig, error)
}
