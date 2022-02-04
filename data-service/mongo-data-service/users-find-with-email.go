package mongods

import (
	"context"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (service MongoDataService) FindUserUserWithEmail(base, email string) (*model.User, error) {
	result := make([]*model.User, 0)
	query := bson.M{"email": bson.M{"$regex": `^` + email + `$`, "$options": "i"}}

	err := service.Database.Find(context.TODO(), service.DatabaseName, model.CollectionUsers, query, 0, func(cur *mongo.Cursor) error {

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
