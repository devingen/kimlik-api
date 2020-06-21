package database_service

import (
	"context"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (service DatabaseService) CreateSession(base, client, userAgent, ip string, user *model.User) (*model.Session, error) {
	collection, err := service.Database.ConnectToCollection(base, model.CollectionSessions)
	if err != nil {
		return nil, err
	}

	item := &model.Session{
		User:         user.DBRef(base),
		UserAgent:    userAgent,
		Client:       client,
		SessionCount: 1,
		Status:       model.SessionStatusSuccessful,
		IP:           ip,
	}

	result, err := collection.InsertOne(context.TODO(), item)
	if err != nil {
		return nil, err
	}

	item.ID = result.InsertedID.(primitive.ObjectID)
	return item, nil
}
