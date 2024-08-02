package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthType string

const (
	AuthTypePassword AuthType = "password"
	AuthTypeOpenID            = "openid"
)

type Auth struct {
	// DBRef fields
	Ref      string             `bson:"_ref,omitempty" json:"_ref,omitempty"`
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Database string             `bson:"_db,omitempty" json:"_db,omitempty"`

	// common model fields
	CreatedAt *time.Time `json:"_created,omitempty" bson:"_created,omitempty"`
	UpdatedAt *time.Time `json:"_updated,omitempty" bson:"_updated,omitempty"`
	Revision  int        `json:"_revision,omitempty" bson:"_revision,omitempty"`

	Type AuthType `json:"type" bson:"type,omitempty"`
	User *User    `json:"user" bson:"user,omitempty"`

	// exists if type is 'password'
	Password string `json:"password" bson:"password,omitempty"`

	// exists if type is 'openid'
	OpenID *OpenIDData `json:"openIdData,omitempty" bson:"openIdData,omitempty"`
}

func (auth *Auth) HashPassword() error {
	if auth.Password == "" {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(auth.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	auth.Password = string(hashedPassword)
	return nil
}

// OpenIDData contains the details returned with ID Token (Open ID Connect).
// See these docs for identity providers.
// Google OpenID Connect
// https://developers.google.com/identity/one-tap/android/idtoken-auth#node.js
// Apple OpenID Connect
// https://developer.apple.com/documentation/sign_in_with_apple/sign_in_with_apple_rest_api/authenticating_users_with_sign_in_with_apple
type OpenIDData struct {
	// Issuer of the ID Token.
	//   https://accounts.google.com or accounts.google.com for Google ID tokens.
	Iss string `json:"iss,omitempty"`

	// Audience of the ID Token. Usually the client ID of the application created in the identity provider.
	Aud string `json:"aud,omitempty"`

	// Subject of the ID Token. It's a unique identifier for the user in the identity provider. Remains the same
	// even if the user has multiple emails or changed their email in the associated identity provider.
	Sub string `json:"sub,omitempty"`

	// User email in the ID Token. Only provided if the proper scope is included in the authentication request.
	Email string `json:"email,omitempty"`
}

func (auth *Auth) AddCreationFields() {
	auth.ID = primitive.NewObjectID()
	now := time.Now()
	auth.CreatedAt = &now
	auth.UpdatedAt = &now
	auth.Revision = 1
}

// PrepareUpdateFields sets the UpdatedAt and deletes the Revision. Giving 0 value to Revision results bson
// ignoring the revision field in $set function. It's incremented by the $inc command
func (auth *Auth) PrepareUpdateFields() {
	auth.Revision = 0
	now := time.Now()
	auth.UpdatedAt = &now
}

func (u *Auth) DBRef(database string) *Auth {
	return &Auth{
		Ref:      CollectionAuths,
		ID:       u.ID,
		Database: database,
	}
}
