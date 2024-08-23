package mongods

import (
	"context"
	"time"

	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (service MongoDataService) CreateSession(ctx context.Context, base, client, userAgent, ip, refreshToken, error string, auth *model.Auth, user *model.User) (*model.Session, error) {
	collection, err := service.Database.ConnectToCollection(base, model.CollectionSessions)
	if err != nil {
		return nil, err
	}

	status := model.SessionStatusSuccessful
	if error != "" {
		status = model.SessionStatusFailed
	}
	var authRef *model.Auth
	if auth != nil {
		authRef = auth.DBRef(base)
	}
	item := &model.Session{
		RefreshToken: &refreshToken,
		Auth:         authRef,
		User:         user.DBRef(base),
		UserAgent:    &userAgent,
		Client:       &client,
		Status:       &status,
		IP:           &ip,
		Error:        &error,
	}
	item.AddCreationFields()

	result, err := collection.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}

	item.ID = result.InsertedID.(primitive.ObjectID)
	return item, nil
}

func (service MongoDataService) UpdateSession(ctx context.Context, base string, session *model.Session) (*time.Time, int, error) {
	collection, err := service.Database.ConnectToCollection(base, model.CollectionSessions)
	if err != nil {
		return nil, 0, err
	}
	session.PrepareUpdateFields()

	var result model.Session
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": session.ID}, bson.M{
		"$set": session,
		"$inc": bson.M{"_revision": 1},
	}).Decode(&result)
	if err != nil {
		return nil, 0, err
	}

	return result.UpdatedAt, result.Revision + 1, nil
}
