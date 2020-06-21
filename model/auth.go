package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthType string

const (
	AuthTypePassword AuthType = "password"
	AuthTypeOAuth2            = "oauth2"
)

type Auth struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Password string             `json:"password"`
	Type     AuthType           `json:"type"`
	User     *DBRef             `json:"user"`
}
