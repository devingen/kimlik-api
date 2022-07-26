package mongods

import (
	"context"
	"github.com/devingen/api-core/database"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (service MongoDataService) FindSAMLConfigs(ctx context.Context, base string, query bson.M) ([]*model.SAMLConfig, error) {
	result := make([]*model.SAMLConfig, 0)

	err := service.Database.Find(ctx, base, model.CollectionSAMLConfigs, query, database.FindOptions{}, func(cur *mongo.Cursor) error {
		var data model.SAMLConfig
		err := cur.Decode(&data)
		if err != nil {
			return err
		}
		result = append(result, &data)
		return nil
	})
	return result, err
}
