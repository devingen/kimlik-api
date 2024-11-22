package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type OAuthAccessCode struct {
	// DBRef fields
	Ref      string `bson:"_ref,omitempty" json:"_ref,omitempty"`
	ID       string `bson:"_id,omitempty" json:"_id,omitempty"`
	Database string `bson:"_db,omitempty" json:"_db,omitempty"`

	// common model fields
	CreatedAt *time.Time `json:"_created,omitempty" bson:"_created,omitempty"`
	UpdatedAt *time.Time `json:"_updated,omitempty" bson:"_updated,omitempty"`
	Revision  int        `json:"_revision,omitempty" bson:"_revision,omitempty"`

	CreatedBy   *User   `json:"createdBy" bson:"createdBy,omitempty"`
	Code        *string `json:"code,omitempty" bson:"code,omitempty"`
	ClientID    *string `json:"clientId,omitempty" bson:"clientId,omitempty"`
	RedirectURI *string `json:"redirectUri,omitempty" bson:"redirectUri,omitempty"`
	Scope       *string `json:"scope,omitempty" bson:"scope,omitempty"`
}

func (oac *OAuthAccessCode) AddCreationFields() {
	oac.ID = primitive.NewObjectID().Hex()
	now := time.Now()
	oac.CreatedAt = &now
	oac.UpdatedAt = &now
	oac.Revision = 1
}

// PrepareUpdateFields sets the UpdatedAt and deletes the Revision. Giving 0 value to Revision results bson
// ignoring the revision field in $set function. It's incremented by the $inc command
func (oac *OAuthAccessCode) PrepareUpdateFields() {
	oac.Revision = 0
	now := time.Now()
	oac.UpdatedAt = &now
}
