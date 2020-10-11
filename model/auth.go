package model

import (
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
	// DBRef fields
	Ref      string             `bson:"_ref,omitempty" json:"_ref,omitempty"`
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Database string             `bson:"_db,omitempty" json:"_db,omitempty"`

	// common model fields
	CreatedAt *time.Time `json:"_created,omitempty" bson:"_created,omitempty"`
	UpdatedAt *time.Time `json:"_updated,omitempty" bson:"_updated,omitempty"`
	Revision  int        `json:"_revision,omitempty" bson:"_revision,omitempty"`

	Password string   `json:"password" bson:"password,omitempty"`
	Type     AuthType `json:"type" bson:"type,omitempty"`
	User     *User    `json:"user" bson:"user,omitempty"`
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
