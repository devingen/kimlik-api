package service

import (
	"context"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IKimlikDataService interface {
	// region user

	CreateUser(ctx context.Context, base, firstName, lastName, email string) (*model.User, error)
	FindUserWithEmail(ctx context.Context, base, email string) (*model.User, error)
	FindUserWithId(ctx context.Context, base, id string) (*model.User, error)
	FindUsers(ctx context.Context, base string, query bson.M) ([]*model.User, error)

	// endregion

	// region auth

	CreateAuthWithPassword(ctx context.Context, base, password string, user *model.User) (*model.Auth, error)
	FindAuthOfUser(ctx context.Context, base, userId string, authType model.AuthType) (*model.Auth, error)
	UpdateAuth(ctx context.Context, base string, auth *model.Auth) (*time.Time, int, error)

	// endregion

	// region session

	CreateSession(ctx context.Context, base, client, userAgent, ip string, user *model.User) (*model.Session, error)
	FindSessionWithId(ctx context.Context, base, id string) (*model.Session, error)

	// endregion

	// region api key
	CreateAPIKey(ctx context.Context, base, name string, scopes []string, keyID, hash string, user *model.User) (*model.APIKey, error)
	GetAPIKey(ctx context.Context, base, id string) (*model.APIKey, error)
	FindAPIKeys(ctx context.Context, base string, query bson.M) ([]*model.APIKey, error)
	UpdateAPIKey(ctx context.Context, base string, apiKey *model.APIKey) (*time.Time, int, error)
	DeleteAPIKey(ctx context.Context, base string, id primitive.ObjectID) error
	// endregion

	// region saml

	CreateSAMLConfig(ctx context.Context, base string, item *model.SAMLConfig) (*model.SAMLConfig, error)
	GetSAMLConfig(ctx context.Context, base, id string) (*model.SAMLConfig, error)
	FindSAMLConfigs(ctx context.Context, base string, query bson.M) ([]*model.SAMLConfig, error)
	UpdateSAMLConfig(ctx context.Context, base string, apiKey *model.SAMLConfig) (*time.Time, int, error)
	DeleteSAMLConfig(ctx context.Context, base string, id primitive.ObjectID) error

	// endregion

}
