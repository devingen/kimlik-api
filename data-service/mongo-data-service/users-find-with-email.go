package mongods

import (
	"context"
	"regexp"

	"github.com/devingen/api-core/database"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (service MongoDataService) FindUserWithEmail(ctx context.Context, base, email string) (*model.User, error) {
	result := make([]*model.User, 0)
	query := bson.M{"email": bson.M{"$regex": `^` + regexp.QuoteMeta(email) + `$`, "$options": "i"}}

	err := service.Database.Find(ctx, base, model.CollectionUsers, query, database.FindOptions{}, func(cur *mongo.Cursor) error {

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
