package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Email     string             `json:"email"`
}

func (u *User) DBRef(database string) *DBRef {
	return &DBRef{
		Ref:      CollectionUsers,
		ID:       u.ID,
		Database: database,
	}
}
