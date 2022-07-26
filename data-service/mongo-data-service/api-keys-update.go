package mongods

import (
	"context"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (service MongoDataService) UpdateAPIKey(ctx context.Context, base string, item *model.APIKey) (*time.Time, int, error) {

	// generate update entry model, ignore the fields that shouldn't be updated
	data := &model.APIKey{
		Name:   item.Name,
		Scopes: item.Scopes,
	}

	collection, err := service.Database.ConnectToCollection(base, model.CollectionAPIKeys)
	if err != nil {
		return nil, 0, err
	}
	data.PrepareUpdateFields()

	var result model.APIKey
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": item.ID}, bson.M{
		"$set": data,
		"$inc": bson.M{"_revision": 1},
	}).Decode(&result)
	if err != nil {
		return nil, 0, err
	}

	return result.UpdatedAt, result.Revision + 1, nil
}
