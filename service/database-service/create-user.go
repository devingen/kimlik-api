package database_service

import (
	"context"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (service DatabaseService) CreateUser(base, firstName, lastName, email string) (*model.User, error) {
	collection, err := service.Database.ConnectToCollection(base, model.CollectionUsers)
	if err != nil {
		return nil, err
	}

	item := &model.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	result, err := collection.InsertOne(context.TODO(), item)
	if err != nil {
		return nil, err
	}

	item.ID = result.InsertedID.(primitive.ObjectID)
	return item, nil
}
