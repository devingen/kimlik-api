package database_service

import (
	"context"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func (service DatabaseService) CreateAuth(base, password string, user *model.User) (*model.Auth, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	collection, err := service.Database.ConnectToCollection(base, model.CollectionAuths)
	if err != nil {
		return nil, err
	}

	item := &model.Auth{
		Password: string(hashedPassword),
		Type:     model.AuthTypePassword,
		User:     user.DBRef(base),
	}

	result, err := collection.InsertOne(context.TODO(), item)
	if err != nil {
		return nil, err
	}

	item.ID = result.InsertedID.(primitive.ObjectID)
	return item, nil
}
