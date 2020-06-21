package database_service

import (
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (service DatabaseService) FindAuthOfUser(base string, user *model.User, authType model.AuthType) (*model.Auth, error) {
	result := make([]*model.Auth, 0)
	query := bson.M{"$and": bson.A{
		bson.M{"user.$id": user.ID},
		bson.M{"type": authType},
	}}

	err := service.Database.Query(base, model.CollectionAuths, query, func(cur *mongo.Cursor) error {

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
