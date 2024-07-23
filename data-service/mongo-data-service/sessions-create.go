package mongods

import (
	"context"

	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (service MongoDataService) CreateSession(ctx context.Context, base, client, userAgent, ip, error string, auth *model.Auth, user *model.User) (*model.Session, error) {
	collection, err := service.Database.ConnectToCollection(base, model.CollectionSessions)
	if err != nil {
		return nil, err
	}

	status := model.SessionStatusSuccessful
	if error != "" {
		status = model.SessionStatusFailed
	}
	item := &model.Session{
		Auth:      auth.DBRef(base),
		User:      user.DBRef(base),
		UserAgent: userAgent,
		Client:    client,
		Status:    status,
		IP:        ip,
		Error:     error,
	}
	item.AddCreationFields()

	result, err := collection.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}

	item.ID = result.InsertedID.(primitive.ObjectID)
	return item, nil
}
