package mongods

import (
	"context"
	"net/http"
	"time"

	core "github.com/devingen/api-core"
	"github.com/devingen/api-core/database"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (service MongoDataService) CreateIntegrationSettings(ctx context.Context, base string, item *model.IntegrationSettings) (*model.IntegrationSettings, error) {
	integrationSettings, err := service.getIntegrationSettings(ctx, base)
	if err != nil {
		return nil, err
	}

	if integrationSettings != nil {
		return nil, core.NewError(http.StatusConflict, "integration-settings-already-exists")
	}

	collection, err := service.Database.ConnectToCollection(base, model.CollectionIntegrationSettings)
	if err != nil {
		return nil, err
	}

	item.AddCreationFields()

	result, err := collection.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}

	item.ID = core.String(result.InsertedID.(string))
	return item, nil
}

func (service MongoDataService) GetIntegrationSettings(ctx context.Context, base string) (*model.IntegrationSettings, error) {
	return service.getIntegrationSettings(ctx, base)
}

func (service MongoDataService) UpdateIntegrationSettings(ctx context.Context, base string, item *model.IntegrationSettings) (*time.Time, int, error) {

	// generate update entry model, ignore the fields that shouldn't be updated
	domain := &model.IntegrationSettings{
		Ulak: item.Ulak,
	}

	collection, err := service.Database.ConnectToCollection(base, model.CollectionIntegrationSettings)
	if err != nil {
		return nil, 0, err
	}
	domain.PrepareUpdateFields()

	var result model.IntegrationSettings
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": item.ID}, bson.M{
		"$set": domain,
		"$inc": bson.M{"_revision": 1},
	}).Decode(&result)
	if err != nil {
		return nil, 0, err
	}

	return result.UpdatedAt, result.Revision + 1, nil
}

func (service MongoDataService) getIntegrationSettings(ctx context.Context, base string) (*model.IntegrationSettings, error) {

	var item *model.IntegrationSettings
	err := service.Database.Find(ctx, base, model.CollectionIntegrationSettings, bson.M{}, database.FindOptions{Limit: 1}, func(cur *mongo.Cursor) error {
		err := cur.Decode(&item)
		if err != nil {
			return err
		}
		return nil
	})

	return item, err
}
