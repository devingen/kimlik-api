package mongods

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (service MongoDataService) CreateAPIKey(ctx context.Context, base, name string, scopes []string, keyID, hash string) (*model.APIKey, error) {
	collection, err := service.Database.ConnectToCollection(base, model.CollectionAPIKeys)
	if err != nil {
		return nil, err
	}

	item := &model.APIKey{
		Hash:   core.String(hash),
		Name:   core.String(name),
		KeyID:  core.String(keyID),
		Scopes: scopes,
	}
	item.AddCreationFields()

	result, err := collection.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}

	item.ID = result.InsertedID.(primitive.ObjectID)
	return item, nil
}
