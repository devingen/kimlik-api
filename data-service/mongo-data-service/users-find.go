package mongods

import (
	"context"
	"github.com/devingen/api-core/database"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (service MongoDataService) FindUsers(ctx context.Context, base string, query bson.M) ([]*model.User, error) {
	result := make([]*model.User, 0)

	err := service.Database.Find(ctx, base, model.CollectionUsers, query, database.FindOptions{}, func(cur *mongo.Cursor) error {

		var data model.User
		err := cur.Decode(&data)
		if err != nil {
			return err
		}
		result = append(result, &data)
		return nil
	})
	return result, err
}
