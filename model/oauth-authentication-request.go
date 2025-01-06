package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OAuth2AuthenticationRequest struct {
	// DBRef fields
	Ref      string `bson:"_ref,omitempty" json:"_ref,omitempty"`
	ID       string `bson:"_id,omitempty" json:"_id,omitempty"`
	Database string `bson:"_db,omitempty" json:"_db,omitempty"`

	// common model fields
	CreatedAt *time.Time `json:"_created,omitempty" bson:"_created,omitempty"`
	UpdatedAt *time.Time `json:"_updated,omitempty" bson:"_updated,omitempty"`
	Revision  int        `json:"_revision,omitempty" bson:"_revision,omitempty"`

	State        *string `json:"state,omitempty" bson:"state,omitempty"`
	ClientID     *string `json:"clientId,omitempty" bson:"clientId,omitempty"`
	CodeVerifier *string `json:"codeVerifier,omitempty" bson:"codeVerifier,omitempty"`
}

func (oar *OAuth2AuthenticationRequest) AddCreationFields() {
	oar.ID = primitive.NewObjectID().Hex()
	now := time.Now()
	oar.CreatedAt = &now
	oar.UpdatedAt = &now
	oar.Revision = 1
}

// PrepareUpdateFields sets the UpdatedAt and deletes the Revision. Giving 0 value to Revision results bson
// ignoring the revision field in $set function. It's incremented by the $inc command
func (oar *OAuth2AuthenticationRequest) PrepareUpdateFields() {
	oar.Revision = 0
	now := time.Now()
	oar.UpdatedAt = &now
}
