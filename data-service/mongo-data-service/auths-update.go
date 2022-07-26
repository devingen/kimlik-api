package mongods

import (
	"context"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (service MongoDataService) UpdateAuth(ctx context.Context, base string, auth *model.Auth) (*time.Time, int, error) {
	collection, err := service.Database.ConnectToCollection(base, model.CollectionAuths)
	if err != nil {
		return nil, 0, err
	}
	auth.PrepareUpdateFields()

	err = auth.HashPassword()
	if err != nil {
		return nil, 0, err
	}

	var result model.Auth
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": auth.ID}, bson.M{
		"$set": auth,
		"$inc": bson.M{"_revision": 1},
	}).Decode(&result)
	if err != nil {
		return nil, 0, err
	}

	return result.UpdatedAt, result.Revision + 1, nil
}
