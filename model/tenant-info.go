package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TenantInfo struct {
	// DBRef fields
	Ref      string `bson:"_ref,omitempty" json:"_ref,omitempty"`
	ID       string `bson:"_id,omitempty" json:"_id,omitempty"`
	Database string `bson:"_db,omitempty" json:"_db,omitempty"`

	// common model fields
	CreatedAt *time.Time `json:"_created,omitempty" bson:"_created,omitempty"`
	UpdatedAt *time.Time `json:"_updated,omitempty" bson:"_updated,omitempty"`
	Revision  int        `json:"_revision,omitempty" bson:"_revision,omitempty"`

	Name                    *string `json:"name,omitempty" bson:"name,omitempty"`
	LogoURL                 *string `json:"logoUrl,omitempty" bson:"logoUrl,omitempty"`
	TermsOfUseURL           *string `json:"termsOfUseUrl,omitempty" bson:"termsOfUseUrl,omitempty"`
	PrivacyPolicyURL        *string `json:"privacyPolicyUrl,omitempty" bson:"privacyPolicyUrl,omitempty"`
	SupportURL              *string `json:"supportUrl,omitempty" bson:"supportUrl,omitempty"`
	SupportEmail            *string `json:"supportEmail,omitempty" bson:"supportEmail,omitempty"`
	OAuth2IssuerIdentifier  *string `json:"oAuth2IssuerIdentifier,omitempty" bson:"oAuth2IssuerIdentifier,omitempty"`
	OAuth2RedirectionURL    *string `json:"oAuth2RedirectionUrl,omitempty" bson:"oAuth2RedirectionUrl,omitempty"`
	OAuth2AuthorizationURL  *string `json:"oAuth2AuthorizationUrl,omitempty" bson:"oAuth2AuthorizationUrl,omitempty"`
	OAuth2SigningPrivateKey *string `json:"oAuth2SigningPrivateKey,omitempty" bson:"oAuth2SigningPrivateKey,omitempty"`
	OAuth2SigningPublicKey  *string `json:"oAuth2SigningPublicKey,omitempty" bson:"oAuth2SigningPublicKey,omitempty"`
}

func (ti *TenantInfo) AddCreationFields() {
	ti.ID = primitive.NewObjectID().Hex()
	now := time.Now()
	ti.CreatedAt = &now
	ti.UpdatedAt = &now
	ti.Revision = 1
}

// PrepareUpdateFields sets the UpdatedAt and deletes the Revision. Giving 0 value to Revision results bson
// ignoring the revision field in $set function. It's incremented by the $inc command
func (ti *TenantInfo) PrepareUpdateFields() {
	ti.Revision = 0
	now := time.Now()
	ti.UpdatedAt = &now
}
