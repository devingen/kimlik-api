package mongods

import (
	"context"
	"time"

	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (service MongoDataService) UpdateUser(ctx context.Context, base string, user *model.User) (*time.Time, int, error) {
	collection, err := service.Database.ConnectToCollection(base, model.CollectionUsers)
	if err != nil {
		return nil, 0, err
	}
	user.PrepareUpdateFields()
	user.Name = core.String(*user.FirstName + " " + *user.LastName)

	var result model.User
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": user.ID}, bson.M{
		"$set": user,
		"$inc": bson.M{"_revision": 1},
	}).Decode(&result)
	if err != nil {
		return nil, 0, err
	}

	return result.UpdatedAt, result.Revision + 1, nil
}
