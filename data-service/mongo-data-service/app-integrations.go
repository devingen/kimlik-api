package mongods

import (
	"context"
	"time"

	"github.com/devingen/api-core/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/devingen/kimlik-api/model"
)

func (service MongoDataService) CreateAppIntegration(ctx context.Context, base string, item *model.AppIntegration) (*model.AppIntegration, error) {
	collection, err := service.Database.ConnectToCollection(base, model.CollectionAppIntegrations)
	if err != nil {
		return nil, err
	}

	item.AddCreationFields()

	result, err := collection.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}

	item.ID = result.InsertedID.(primitive.ObjectID)
	return item, nil
}

func (service MongoDataService) FindAppIntegrations(ctx context.Context, base string, query bson.M) ([]*model.AppIntegration, error) {
	result := make([]*model.AppIntegration, 0)

	err := service.Database.Find(ctx, base, model.CollectionAppIntegrations, query, database.FindOptions{}, func(cur *mongo.Cursor) error {
		var data model.AppIntegration
		err := cur.Decode(&data)
		if err != nil {
			return err
		}
		result = append(result, &data)
		return nil
	})
	return result, err
}

func (service MongoDataService) UpdateAppIntegration(ctx context.Context, base string, item *model.AppIntegration) (*time.Time, int, error) {

	// generate update entry model, ignore the fields that shouldn't be updated
	data := &model.AppIntegration{
		//ClientID:         item.ClientID,
		Name:             item.Name,
		LogoURL:          item.LogoURL,
		TermsOfUseURL:    item.TermsOfUseURL,
		PrivacyPolicyURL: item.PrivacyPolicyURL,
		SupportURL:       item.SupportURL,
		SupportEmail:     item.SupportEmail,
		OAuth2Config:     item.OAuth2Config,
	}

	collection, err := service.Database.ConnectToCollection(base, model.CollectionAppIntegrations)
	if err != nil {
		return nil, 0, err
	}
	data.PrepareUpdateFields()

	var result model.AppIntegration
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": item.ID}, bson.M{
		"$set": data,
		"$inc": bson.M{"_revision": 1},
	}).Decode(&result)
	if err != nil {
		return nil, 0, err
	}

	return result.UpdatedAt, result.Revision + 1, nil
}

func (service MongoDataService) DeleteAppIntegration(ctx context.Context, base string, id primitive.ObjectID) error {
	collection, err := service.Database.ConnectToCollection(base, model.CollectionAppIntegrations)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
