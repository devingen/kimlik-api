package mongods

import (
	"context"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (service MongoDataService) CreateAPIKey(base, name, productId string, scopes []string, keyPrefix, hash string, user *model.User) (*model.ApiKey, error) {
	collection, err := service.Database.ConnectToCollection(base, model.CollectionApiKeys)
	if err != nil {
		return nil, err
	}

	item := &model.ApiKey{
		CreatedBy: user.DBRef(base),
		Hash:      hash,
		Name:      name,
		KeyPrefix: keyPrefix,
		ProductId: productId,
		Scopes:    scopes,
	}
	item.AddCreationFields()

	result, err := collection.InsertOne(context.Background(), item)
	if err != nil {
		return nil, err
	}

	item.ID = result.InsertedID.(primitive.ObjectID)
	return item, nil
}
