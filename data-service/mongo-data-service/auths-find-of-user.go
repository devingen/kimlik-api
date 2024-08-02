package mongods

import (
	"context"

	"github.com/devingen/api-core/database"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (service MongoDataService) FindAuthOfUser(ctx context.Context, base, userId string, authType model.AuthType) (*model.Auth, error) {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Auth, 0)
	query := bson.M{"$and": bson.A{
		bson.M{"user._id": id},
		bson.M{"type": authType},
	}}

	err = service.Database.Find(ctx, base, model.CollectionAuths, query, database.FindOptions{}, func(cur *mongo.Cursor) error {
		var data model.Auth
		err := cur.Decode(&data)
		if err != nil {
			return err
		}
		result = append(result, &data)
		return nil
	})
	if len(result) > 0 {
		return result[0], err
	}
	return nil, err
}

func (service MongoDataService) FindAuthsOfUser(ctx context.Context, base, userId string) ([]*model.Auth, error) {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Auth, 0)
	query := bson.M{"$and": bson.A{
		bson.M{"user._id": id},
	}}

	err = service.Database.Find(ctx, base, model.CollectionAuths, query, database.FindOptions{}, func(cur *mongo.Cursor) error {
		var data model.Auth
		err := cur.Decode(&data)
		if err != nil {
			return err
		}
		result = append(result, &data)
		return nil
	})
	return result, err
}
