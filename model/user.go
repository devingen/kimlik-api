package model

import (
	coremodel "github.com/devingen/api-core/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `json:"_created" bson:"_created"`
	UpdatedAt time.Time          `json:"_updated" bson:"_updated"`
	Revision  int                `json:"_revision" bson:"_revision,omitempty"`

	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func (u *User) DBRef(database string) *coremodel.DBRef {
	return &coremodel.DBRef{
		Ref:      CollectionUsers,
		ID:       u.ID,
		Database: database,
	}
}

func (u *User) AddCreationFields() *User {
	u.ID = primitive.NewObjectID()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.Revision = 1
	return u
}
