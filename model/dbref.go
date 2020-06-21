package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type DBRef struct {
	Ref      string             `bson:"$ref" json:"ref"`
	ID       primitive.ObjectID `bson:"$id" json:"id"`
	Database string             `bson:"$db" json:"db"`
}
