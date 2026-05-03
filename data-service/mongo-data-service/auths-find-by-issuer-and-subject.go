package mongods

import (
	"context"

	"github.com/devingen/api-core/database"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (service MongoDataService) FindOIDCAuthByIssuerAndSubject(ctx context.Context, base, issuer, subject string) (*model.Auth, error) {
	result := make([]*model.Auth, 0)
	query := bson.M{"$and": bson.A{
		bson.M{"openIdData.iss": issuer},
		bson.M{"openIdData.sub": subject},
		bson.M{"type": model.AuthTypeOpenID},
	}}

	err := service.Database.Find(ctx, base, model.CollectionAuths, query, database.FindOptions{}, func(cur *mongo.Cursor) error {
		var data model.Auth
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
