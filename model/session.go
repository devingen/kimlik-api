package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SessionStatus string

const (
	SessionStatusSuccessful  SessionStatus = "successful"
	SessionStatusFailed                    = "failed"
	SessionStatusInvalidated               = "invalidated"
)

type Session struct {
	// DBRef fields
	Ref      string             `bson:"_ref,omitempty" json:"ref,omitempty"`
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Database string             `bson:"_db,omitempty" json:"db,omitempty"`

	// common model fields
	CreatedAt *time.Time `json:"_created,omitempty" bson:"_created,omitempty"`
	UpdatedAt *time.Time `json:"_updated,omitempty" bson:"_updated,omitempty"`
	Revision  int        `json:"_revision,omitempty" bson:"_revision,omitempty"`

	User         *User         `json:"user,omitempty" bson:"user,omitempty"`
	UserAgent    string        `json:"userAgent,omitempty" bson:"userAgent,omitempty"`
	Client       string        `json:"client,omitempty" bson:"client,omitempty"`
	SessionCount float64       `json:"sessionCount,omitempty" bson:"sessionCount,omitempty"`
	Status       SessionStatus `json:"status,omitempty" bson:"status,omitempty"`
	IP           string        `json:"ip,omitempty" bson:"ip,omitempty"`
}

func (s *Session) AddCreationFields() {
	s.ID = primitive.NewObjectID()
	now := time.Now()
	s.CreatedAt = &now
	s.UpdatedAt = &now
	s.Revision = 1
}
