package model

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type APIKey struct {
	// DBRef fields
	Ref      string             `bson:"_ref,omitempty" json:"_ref,omitempty"`
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Database string             `bson:"_db,omitempty" json:"_db,omitempty"`

	// common model fields
	CreatedAt *time.Time `json:"_created,omitempty" bson:"_created,omitempty"`
	UpdatedAt *time.Time `json:"_updated,omitempty" bson:"_updated,omitempty"`
	Revision  int        `json:"_revision,omitempty" bson:"_revision,omitempty"`

	CreatedBy *User    `json:"createdBy" bson:"createdBy,omitempty"`
	Hash      *string  `json:"hash,omitempty" bson:"hash,omitempty"`
	Name      *string  `json:"name,omitempty" bson:"name,omitempty"`
	KeyID     *string  `json:"keyId,omitempty" bson:"keyId,omitempty"`
	Scopes    []string `json:"scopes,omitempty" bson:"scopes,omitempty"`
}

func (apiKey *APIKey) AddCreationFields() {
	apiKey.ID = primitive.NewObjectID()
	now := time.Now()
	apiKey.CreatedAt = &now
	apiKey.UpdatedAt = &now
	apiKey.Revision = 1
}

// PrepareUpdateFields sets the UpdatedAt and deletes the Revision. Giving 0 value to Revision results bson
// ignoring the revision field in $set function. It's incremented by the $inc command
func (apiKey *APIKey) PrepareUpdateFields() {
	apiKey.Revision = 0
	now := time.Now()
	apiKey.UpdatedAt = &now
}

func (apiKey *APIKey) ContainsScope(scope string) bool {
	for _, s := range apiKey.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

func (apiKey *APIKey) ValidateScope(scope string) error {
	return apiKey.ValidateScopes([]string{scope})
}

func (apiKey *APIKey) ValidateScopes(scopes []string) error {
	if apiKey.Scopes == nil {
		return errors.New("api-key-has-no-scope")
	}

	allowedScopes := map[string]bool{}
	for _, scope := range apiKey.Scopes {
		allowedScopes[scope] = true
	}

	for _, expectedScope := range scopes {
		if !allowedScopes[expectedScope] {
			return errors.New("api-key-missing-scope:" + expectedScope)
		}
	}
	return nil
}
