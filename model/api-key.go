package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ApiKey struct {
	// DBRef fields
	Ref      string             `bson:"_ref,omitempty" json:"_ref,omitempty"`
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Database string             `bson:"_db,omitempty" json:"_db,omitempty"`

	// common model fields
	CreatedAt *time.Time `json:"_created,omitempty" bson:"_created,omitempty"`
	UpdatedAt *time.Time `json:"_updated,omitempty" bson:"_updated,omitempty"`
	Revision  int        `json:"_revision,omitempty" bson:"_revision,omitempty"`

	CreatedBy *User    `json:"user" bson:"user,omitempty"`
	Hash      string   `json:"hash,omitempty" bson:"hash,omitempty"`
	Name      string   `json:"name,omitempty" bson:"name,omitempty"`
	KeyPrefix string   `json:"keyPrefix,omitempty" bson:"keyPrefix,omitempty"`
	ProductId string   `json:"productId,omitempty" bson:"productId,omitempty"`
	Scopes    []string `json:"scopes,omitempty" bson:"scopes,omitempty"`
}

func (apiKey *ApiKey) AddCreationFields() {
	apiKey.ID = primitive.NewObjectID()
	now := time.Now()
	apiKey.CreatedAt = &now
	apiKey.UpdatedAt = &now
	apiKey.Revision = 1
}

// PrepareUpdateFields sets the UpdatedAt and deletes the Revision. Giving 0 value to Revision results bson
// ignoring the revision field in $set function. It's incremented by the $inc command
func (apiKey *ApiKey) PrepareUpdateFields() {
	apiKey.Revision = 0
	now := time.Now()
	apiKey.UpdatedAt = &now
}
