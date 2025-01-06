package mongods

import (
	"context"

	"github.com/devingen/api-core/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/devingen/kimlik-api/model"
)

func (service MongoDataService) CreateOAuth2AuthenticationRequest(ctx context.Context, base string, item *model.OAuth2AuthenticationRequest) (*model.OAuth2AuthenticationRequest, error) {
	collection, err := service.Database.ConnectToCollection(base, model.CollectionOAuth2AuthenticationRequests)
	if err != nil {
		return nil, err
	}

	item.AddCreationFields()

	result, err := collection.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}

	item.ID = result.InsertedID.(string)
	return item, nil
}

func (service MongoDataService) FindOAuth2AuthenticationRequests(ctx context.Context, base string, query bson.M) ([]*model.OAuth2AuthenticationRequest, error) {
	result := make([]*model.OAuth2AuthenticationRequest, 0)

	err := service.Database.Find(ctx, base, model.CollectionOAuth2AuthenticationRequests, query, database.FindOptions{}, func(cur *mongo.Cursor) error {
		var data model.OAuth2AuthenticationRequest
		err := cur.Decode(&data)
		if err != nil {
			return err
		}
		result = append(result, &data)
		return nil
	})
	return result, err
}
