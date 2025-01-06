package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AppIntegration contains information about 3rd party app that uses SSO authentication.
type AppIntegration struct {
	// DBRef fields
	Ref      string             `bson:"_ref,omitempty" json:"_ref,omitempty"`
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Database string             `bson:"_db,omitempty" json:"_db,omitempty"`

	// common model fields
	CreatedAt *time.Time `json:"_created,omitempty" bson:"_created,omitempty"`
	UpdatedAt *time.Time `json:"_updated,omitempty" bson:"_updated,omitempty"`
	Revision  int        `json:"_revision,omitempty" bson:"_revision,omitempty"`

	ClientID         *string                     `json:"clientId" bson:"clientId,omitempty" validate:"required"`
	Name             *string                     `json:"name" bson:"name,omitempty" validate:"required"`
	LogoURL          *string                     `json:"logoUrl,omitempty" bson:"logoUrl,omitempty"`
	TermsOfUseURL    *string                     `json:"termsOfUseUrl,omitempty" bson:"termsOfUseUrl,omitempty"`
	PrivacyPolicyURL *string                     `json:"privacyPolicyUrl,omitempty" bson:"privacyPolicyUrl,omitempty"`
	SupportURL       *string                     `json:"supportUrl,omitempty" bson:"supportUrl,omitempty"`
	SupportEmail     *string                     `json:"supportEmail,omitempty" bson:"supportEmail,omitempty"`
	OAuth2Config     *AppIntegrationOAuth2Config `json:"oAuth2Config" bson:"oAuth2Config,omitempty"`
}

type AppIntegrationOAuth2Config struct {
	RedirectURLs []string `json:"redirectUrls,omitempty" bson:"redirectUrls,omitempty"`
	Scopes       []string `json:"scopes" bson:"scopes,omitempty" validate:"required"`
}

func (ai *AppIntegration) AddCreationFields() {
	ai.ID = primitive.NewObjectID()
	now := time.Now()
	ai.CreatedAt = &now
	ai.UpdatedAt = &now
	ai.Revision = 1
}

// PrepareUpdateFields sets the UpdatedAt and deletes the Revision. Giving 0 value to Revision results bson
// ignoring the revision field in $set function. It's incremented by the $inc command
func (ai *AppIntegration) PrepareUpdateFields() {
	ai.Revision = 0
	now := time.Now()
	ai.UpdatedAt = &now
}
