package mongods

import (
	"context"
	"time"

	"github.com/devingen/api-core/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/devingen/kimlik-api/model"
)

func (service MongoDataService) CreateOAuth2Config(ctx context.Context, base string, item *model.OAuth2Config) (*model.OAuth2Config, error) {
	collection, err := service.Database.ConnectToCollection(base, model.CollectionOAuth2Configs)
	if err != nil {
		return nil, err
	}

	item.AddCreationFields()

	result, err := collection.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}

	item.ID = result.InsertedID.(primitive.ObjectID)
	return item, nil
}

func (service MongoDataService) FindOAuth2Configs(ctx context.Context, base string, query bson.M) ([]*model.OAuth2Config, error) {
	result := make([]*model.OAuth2Config, 0)

	err := service.Database.Find(ctx, base, model.CollectionOAuth2Configs, query, database.FindOptions{}, func(cur *mongo.Cursor) error {
		var data model.OAuth2Config
		err := cur.Decode(&data)
		if err != nil {
			return err
		}
		result = append(result, &data)
		return nil
	})
	return result, err
}

func (service MongoDataService) UpdateOAuth2Config(ctx context.Context, base string, item *model.OAuth2Config) (*time.Time, int, error) {

	// generate update entry model, ignore the fields that shouldn't be updated
	data := &model.OAuth2Config{
		Name:                  item.Name,
		Scopes:                item.Scopes,
		ClientID:              item.ClientID,
		ClientSecret:          item.ClientSecret,
		Issuer:                item.Issuer,
		AuthorizationEndpoint: item.AuthorizationEndpoint,
		TokenEndpoint:         item.TokenEndpoint,
		JWKSEndpoint:          item.JWKSEndpoint,
		UserinfoEndpoint:      item.UserinfoEndpoint,
	}

	collection, err := service.Database.ConnectToCollection(base, model.CollectionOAuth2Configs)
	if err != nil {
		return nil, 0, err
	}
	data.PrepareUpdateFields()

	var result model.OAuth2Config
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": item.ID}, bson.M{
		"$set": data,
		"$inc": bson.M{"_revision": 1},
	}).Decode(&result)
	if err != nil {
		return nil, 0, err
	}

	return result.UpdatedAt, result.Revision + 1, nil
}

func (service MongoDataService) DeleteOAuth2Config(ctx context.Context, base string, id primitive.ObjectID) error {
	collection, err := service.Database.ConnectToCollection(base, model.CollectionOAuth2Configs)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
