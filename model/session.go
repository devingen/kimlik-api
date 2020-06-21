package model

import (
	coremodel "github.com/devingen/api-core/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SessionStatus string

const (
	SessionStatusSuccessful  SessionStatus = "successful"
	SessionStatusFailed                    = "failed"
	SessionStatusInvalidated               = "invalidated"
)

type Session struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	User         *coremodel.DBRef   `json:"user"`
	UserAgent    string             `json:"userAgent"`
	Client       string             `json:"client"`
	SessionCount float64            `json:"sessionCount"`
	Status       SessionStatus      `json:"status"`
	IP           string             `json:"ip"`
}
