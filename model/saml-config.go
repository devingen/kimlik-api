package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SAMLConfig struct {
	// DBRef fields
	Ref      string             `bson:"_ref,omitempty" json:"_ref,omitempty"`
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Database string             `bson:"_db,omitempty" json:"_db,omitempty"`

	// common model fields
	CreatedAt *time.Time `json:"_created,omitempty" bson:"_created,omitempty"`
	UpdatedAt *time.Time `json:"_updated,omitempty" bson:"_updated,omitempty"`
	Revision  int        `json:"_revision,omitempty" bson:"_revision,omitempty"`

	MetadataURL                 string   `json:"metadataURL" bson:"metadataURL,omitempty" validate:"required"`
	AssertionConsumerServiceURL string   `json:"assertionConsumerServiceURL" bson:"assertionConsumerServiceURL,omitempty" validate:"required"`
	AudienceURI                 string   `json:"audienceURI" bson:"audienceURI,omitempty" validate:"required"`
	ServiceProviderIssuer       string   `json:"serviceProviderIssuer" bson:"serviceProviderIssuer,omitempty" validate:"required"`
	SAMLResponseValues          []string `json:"samlResponseValues" bson:"samlResponseValues,omitempty"`
}

func (sc *SAMLConfig) AddCreationFields() {
	sc.ID = primitive.NewObjectID()
	now := time.Now()
	sc.CreatedAt = &now
	sc.UpdatedAt = &now
	sc.Revision = 1
}

// PrepareUpdateFields sets the UpdatedAt and deletes the Revision. Giving 0 value to Revision results bson
// ignoring the revision field in $set function. It's incremented by the $inc command
func (sc *SAMLConfig) PrepareUpdateFields() {
	sc.Revision = 0
	now := time.Now()
	sc.UpdatedAt = &now
}
