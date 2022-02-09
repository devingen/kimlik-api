package mongods

import (
	"context"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (service MongoDataService) FindUserUserWithId(base, id string) (*model.User, error) {
	result := make([]*model.User, 0)

	mId, mErr := primitive.ObjectIDFromHex(id)
	if mErr != nil {
		return nil, mErr
	}
	query := bson.M{"_id": mId}

	err := service.Database.Find(context.TODO(), base, model.CollectionUsers, query, 0, func(cur *mongo.Cursor) error {

		var data model.User
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
