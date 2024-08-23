package mongods

import (
	"context"

	"github.com/devingen/api-core/database"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (service MongoDataService) FindSessionWithId(ctx context.Context, base, id string) (*model.Session, error) {

	mId, mErr := primitive.ObjectIDFromHex(id)
	if mErr != nil {
		return nil, mErr
	}
	query := bson.M{"_id": mId}

	var result *model.Session
	err := service.Database.Find(ctx, base, model.CollectionSessions, query, database.FindOptions{}, func(cur *mongo.Cursor) error {

		err := cur.Decode(&result)
		if err != nil {
			return err
		}
		return nil
	})
	return result, err
}

func (service MongoDataService) FindSessions(ctx context.Context, base string, query bson.M) ([]*model.Session, error) {

	result := make([]*model.Session, 0)

	err := service.Database.Find(ctx, base, model.CollectionSessions, query, database.FindOptions{}, func(cur *mongo.Cursor) error {
		var data model.Session
		err := cur.Decode(&data)
		if err != nil {
			return err
		}
		result = append(result, &data)
		return nil
	})
	return result, err
}
