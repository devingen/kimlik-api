package service

import (
	"context"
	"time"

	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IKimlikDataService interface {
	// region user

	AnonymizeUser(ctx context.Context, base string, id primitive.ObjectID) error
	CreateUser(ctx context.Context, base, firstName, lastName, email string, status model.UserStatus, isEmailVerified bool) (*model.User, error)
	FindUserWithEmail(ctx context.Context, base, email string) (*model.User, error)
	FindUserWithId(ctx context.Context, base, id string) (*model.User, error)
	FindUsers(ctx context.Context, base string, query bson.M) ([]*model.User, error)
	UpdateUser(ctx context.Context, base string, user *model.User) (*time.Time, int, error)

	// endregion

	// region auth

	CreateAuthWithPassword(ctx context.Context, base, password string, user *model.User) (*model.Auth, error)
	CreateAuthWithIDToken(ctx context.Context, base string, claims map[string]interface{}, user *model.User) (*model.Auth, error)
	FindPasswordAuthOfUser(ctx context.Context, base, userId string) (*model.Auth, error)
	FindOIDCAuthOfUser(ctx context.Context, base, userId string, issuer string) (*model.Auth, error)
	FindAuthsOfUser(ctx context.Context, base, userId string) ([]*model.Auth, error)
	UpdateAuth(ctx context.Context, base string, auth *model.Auth) (*time.Time, int, error)
	DeleteAuth(ctx context.Context, base string, id primitive.ObjectID) error

	// endregion

	// region sessions
	CreateSession(ctx context.Context, base, client, userAgent, ip, refreshToken, error string, auth *model.Auth, user *model.User) (*model.Session, error)
	FindSessionWithId(ctx context.Context, base, id string) (*model.Session, error)
	FindSessions(ctx context.Context, base string, query bson.M) ([]*model.Session, error)
	UpdateSession(ctx context.Context, base string, session *model.Session) (*time.Time, int, error)
	// endregion

	// region api key
	CreateAPIKey(ctx context.Context, base, name string, scopes []string, keyID, hash string) (*model.APIKey, error)
	GetAPIKey(ctx context.Context, base, id string) (*model.APIKey, error)
	FindAPIKeys(ctx context.Context, base string, query bson.M) ([]*model.APIKey, error)
	UpdateAPIKey(ctx context.Context, base string, apiKey *model.APIKey) (*time.Time, int, error)
	DeleteAPIKey(ctx context.Context, base string, id primitive.ObjectID) error
	// endregion

	// region saml

	CreateSAMLConfig(ctx context.Context, base string, item *model.SAMLConfig) (*model.SAMLConfig, error)
	GetSAMLConfig(ctx context.Context, base, id string) (*model.SAMLConfig, error)
	FindSAMLConfigs(ctx context.Context, base string, query bson.M) ([]*model.SAMLConfig, error)
	UpdateSAMLConfig(ctx context.Context, base string, samlConfig *model.SAMLConfig) (*time.Time, int, error)
	DeleteSAMLConfig(ctx context.Context, base string, id primitive.ObjectID) error

	// endregion

	// region tenant info
	CreateTenantInfo(ctx context.Context, base string, item *model.TenantInfo) (*model.TenantInfo, error)
	GetTenantInfo(ctx context.Context, base string) (*model.TenantInfo, error)
	UpdateTenantInfo(ctx context.Context, base string, item *model.TenantInfo) (*time.Time, int, error)
	//endregion

	// region oauth
	CreateOAuthAccessCode(ctx context.Context, base string, item *model.OAuthAccessCode) (*model.OAuthAccessCode, error)
	// endregion
}
