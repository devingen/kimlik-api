package mongods

import (
	"context"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (service MongoDataService) DeleteSAMLConfig(ctx context.Context, base string, id primitive.ObjectID) error {
	collection, err := service.Database.ConnectToCollection(base, model.CollectionSAMLConfigs)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
