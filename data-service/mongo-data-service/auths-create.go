package mongods

import (
	"context"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (service MongoDataService) CreateAuthWithPassword(base, password string, user *model.User) (*model.Auth, error) {
	collection, err := service.Database.ConnectToCollection(base, model.CollectionAuths)
	if err != nil {
		return nil, err
	}

	item := &model.Auth{
		Password: password,
		Type:     model.AuthTypePassword,
		User:     user.DBRef(base),
	}
	item.AddCreationFields()

	err = item.HashPassword()
	if err != nil {
		return nil, err
	}

	result, err := collection.InsertOne(context.Background(), item)
	if err != nil {
		return nil, err
	}

	item.ID = result.InsertedID.(primitive.ObjectID)
	return item, nil
}
