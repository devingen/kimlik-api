package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SessionStatus string

const (
	SessionStatusSuccessful  SessionStatus = "successful"
	SessionStatusFailed                    = "failed"
	SessionStatusInvalidated               = "invalidated"
)

type Session struct {
	// DBRef fields
	Ref      string             `bson:"_ref,omitempty" json:"_ref,omitempty"`
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Database string             `bson:"_db,omitempty" json:"_db,omitempty"`

	// common model fields
	CreatedAt *time.Time `json:"_created,omitempty" bson:"_created,omitempty"`
	UpdatedAt *time.Time `json:"_updated,omitempty" bson:"_updated,omitempty"`
	Revision  int        `json:"_revision,omitempty" bson:"_revision,omitempty"`

	Auth      *Auth         `json:"auth,omitempty" bson:"auth,omitempty"`
	User      *User         `json:"user,omitempty" bson:"user,omitempty"`
	UserAgent string        `json:"userAgent,omitempty" bson:"userAgent,omitempty"`
	Client    string        `json:"client,omitempty" bson:"client,omitempty"`
	Status    SessionStatus `json:"status,omitempty" bson:"status,omitempty"`
	IP        string        `json:"ip,omitempty" bson:"ip,omitempty"`
	Error     string        `json:"error,omitempty" bson:"error,omitempty"`
}

func (s *Session) AddCreationFields() {
	s.ID = primitive.NewObjectID()
	now := time.Now()
	s.CreatedAt = &now
	s.UpdatedAt = &now
	s.Revision = 1
}
