package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OAuth2Config contains the configuration for the identity provider for OAuth2 integration.
// See references:
//
//	https://developer.okta.com/docs/guides/add-an-external-idp/openidconnect/main/
type OAuth2Config struct {
	// DBRef fields
	Ref      string             `bson:"_ref,omitempty" json:"_ref,omitempty"`
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Database string             `bson:"_db,omitempty" json:"_db,omitempty"`

	// common model fields
	CreatedAt *time.Time `json:"_created,omitempty" bson:"_created,omitempty"`
	UpdatedAt *time.Time `json:"_updated,omitempty" bson:"_updated,omitempty"`
	Revision  int        `json:"_revision,omitempty" bson:"_revision,omitempty"`

	Name   *string  `json:"name" bson:"name,omitempty" validate:"required"`
	Scopes []string `json:"scopes" bson:"scopes,omitempty" validate:"required"`

	// ClientID of the Service Provider that matches the records in Identity Provider
	ClientID     *string `json:"clientId" bson:"clientId,omitempty" validate:"required"`
	ClientSecret *string `json:"clientSecret" bson:"clientSecret,omitempty"`
	Issuer       *string `json:"issuer" bson:"issuer,omitempty"`

	// AuthorizationEndpoint is the UI in Identity Provider where user consents views the authorization page
	AuthorizationEndpoint *string `json:"authorizationEndpoint" bson:"authorizationEndpoint,omitempty"`

	// TokenEndpoint is the endpoint of the Identity Provider api to exchange the code with the token
	TokenEndpoint    *string `json:"tokenEndpoint" bson:"tokenEndpoint,omitempty"`
	JWKSEndpoint     *string `json:"jwksEndpoint" bson:"jwksEndpoint,omitempty"`
	UserinfoEndpoint *string `json:"userinfoEndpoint" bson:"userinfoEndpoint,omitempty"`
}

func (sc *OAuth2Config) AddCreationFields() {
	sc.ID = primitive.NewObjectID()
	now := time.Now()
	sc.CreatedAt = &now
	sc.UpdatedAt = &now
	sc.Revision = 1
}

// PrepareUpdateFields sets the UpdatedAt and deletes the Revision. Giving 0 value to Revision results bson
// ignoring the revision field in $set function. It's incremented by the $inc command
func (sc *OAuth2Config) PrepareUpdateFields() {
	sc.Revision = 0
	now := time.Now()
	sc.UpdatedAt = &now
}
