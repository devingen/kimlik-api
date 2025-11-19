package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IntegrationSettings struct {
	// DBRef fields
	Ref      string `bson:"_ref,omitempty" json:"_ref,omitempty"`
	ID       string `bson:"_id,omitempty" json:"_id,omitempty"`
	Database string `bson:"_db,omitempty" json:"_db,omitempty"`

	// common model fields
	CreatedAt *time.Time `json:"_created,omitempty" bson:"_created,omitempty"`
	UpdatedAt *time.Time `json:"_updated,omitempty" bson:"_updated,omitempty"`
	Revision  int        `json:"_revision,omitempty" bson:"_revision,omitempty"`

	Ulak *Ulak `json:"ulak,omitempty" bson:"ulak,omitempty"`
}

type Ulak struct {
	ProductID             *string `json:"productId,omitempty" bson:"productId,omitempty"`
	APIKey                *string `json:"apiKey,omitempty" bson:"apiKey,omitempty"`
	SenderConfigurationID *string `json:"senderConfigurationId,omitempty" bson:"senderConfigurationId,omitempty"`
	SenderEmail           *string `json:"senderEmail,omitempty" bson:"senderEmail,omitempty"`
	SenderName            *string `json:"senderName,omitempty" bson:"senderName,omitempty"`
}

func (p *IntegrationSettings) AddCreationFields() {
	p.ID = primitive.NewObjectID().Hex()
	now := time.Now()
	p.CreatedAt = &now
	p.UpdatedAt = &now
	p.Revision = 1
}

// PrepareUpdateFields sets the UpdatedAt and deletes the Revision. Giving 0 value to Revision results bson
// ignoring the revision field in $set function. It's incremented by the $inc command
func (p *IntegrationSettings) PrepareUpdateFields() {
	p.Revision = 0
	now := time.Now()
	p.UpdatedAt = &now
}
