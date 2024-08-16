package mongods

import (
	"context"

	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (service MongoDataService) AnonymizeUser(ctx context.Context, base string, id primitive.ObjectID) error {
	collection, err := service.Database.ConnectToCollection(base, model.CollectionUsers)
	if err != nil {
		return err
	}

	var result model.User
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{
		"$set": model.User{
			FirstName: core.String("-"),
			LastName:  core.String("-"),
			Email:     core.String("-"),
		},
		"$inc": bson.M{"_revision": 1},
	}).Decode(&result)
	return err
}
