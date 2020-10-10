package database_service

import (
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (service DatabaseService) FindAuthOfUser(base, userId string, authType model.AuthType) (*model.Auth, error) {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Auth, 0)
	query := bson.M{"$and": bson.A{
		bson.M{"user.$id": id},
		bson.M{"type": authType},
	}}

	err = service.Database.Query(base, model.CollectionAuths, query, func(cur *mongo.Cursor) error {
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
