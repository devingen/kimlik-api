package model

import (
	coremodel "github.com/devingen/api-core/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthType string

const (
	AuthTypePassword AuthType = "password"
	AuthTypeOAuth2            = "oauth2"
)

type Auth struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `json:"_created" bson:"_created"`
	UpdatedAt time.Time          `json:"_updated" bson:"_updated"`
	Revision  int                `json:"_revision" bson:"_revision,omitempty"`

	Password string           `json:"password" bson:"password,omitempty"`
	Type     AuthType         `json:"type" bson:"type,omitempty"`
	User     *coremodel.DBRef `json:"user" bson:"user,omitempty"`
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

func (auth *Auth) AddCreationFields() {
	auth.ID = primitive.NewObjectID()
	auth.CreatedAt = time.Now()
	auth.UpdatedAt = time.Now()
	auth.Revision = 1
}

// PrepareUpdateFields sets the UpdatedAt and deletes the Revision. Giving 0 value to Revision results bson
// ignoring the revision field in $set function. It's incremented by the $inc command
func (auth *Auth) PrepareUpdateFields() {
	auth.Revision = 0
	auth.UpdatedAt = time.Now()
}
