package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStatus string

const (
	// UserStatusNotActivated is used when a user is created but s/he never activated the account.
	//   This state can occur if the user is invited by an admin to the system.
	//   The status is updated later to active after s/he activates the account.
	UserStatusNotActivated UserStatus = "not-activated"

	// UserStatusActive is used when user activated the account by agreeing the terms and conditions explicitly
	UserStatusActive UserStatus = "active"
)

type User struct {
	// DBRef fields
	Ref      string             `bson:"_ref,omitempty" json:"_ref,omitempty"`
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Database string             `bson:"_db,omitempty" json:"_db,omitempty"`

	// common model fields
	CreatedAt *time.Time `json:"_created,omitempty" bson:"_created,omitempty"`
	UpdatedAt *time.Time `json:"_updated,omitempty" bson:"_updated,omitempty"`
	Revision  int        `json:"_revision,omitempty" bson:"_revision,omitempty"`

	Status          *UserStatus `json:"status,omitempty" bson:"status,omitempty"`
	IsEmailVerified *bool       `json:"isEmailVerified,omitempty" bson:"isEmailVerified,omitempty"`
	Email           *string     `json:"email,omitempty" bson:"email,omitempty"`
	FirstName       *string     `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName        *string     `json:"lastName,omitempty" bson:"lastName,omitempty"`
}

func (u *User) DBRef(database string) *User {
	return &User{
		Ref:      CollectionUsers,
		ID:       u.ID,
		Database: database,
	}
}

func (u *User) AddCreationFields() {
	u.ID = primitive.NewObjectID()
	now := time.Now()
	u.CreatedAt = &now
	u.UpdatedAt = &now
	u.Revision = 1
}

// PrepareUpdateFields sets the UpdatedAt and deletes the Revision. Giving 0 value to Revision results bson
// ignoring the revision field in $set function. It's incremented by the $inc command
func (u *User) PrepareUpdateFields() {
	u.Revision = 0
	now := time.Now()
	u.UpdatedAt = &now
}
