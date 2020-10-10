package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	// DBRef fields
	Ref      string             `bson:"_ref,omitempty" json:"ref,omitempty"`
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Database string             `bson:"_db,omitempty" json:"db,omitempty"`

	// common model fields
	CreatedAt *time.Time `json:"_created,omitempty" bson:"_created,omitempty"`
	UpdatedAt *time.Time `json:"_updated,omitempty" bson:"_updated,omitempty"`
	Revision  int        `json:"_revision,omitempty" bson:"_revision,omitempty"`

	FirstName string `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty" bson:"lastName,omitempty"`
	Email     string `json:"email,omitempty" bson:"email,omitempty"`
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
