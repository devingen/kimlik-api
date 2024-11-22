package mongods

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/api-core/database"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func (service MongoDataService) CreateTenantInfo(ctx context.Context, base string, item *model.TenantInfo) (*model.TenantInfo, error) {
	tenantInfo, err := service.getTenantInfo(ctx, base)
	if err != nil {
		return nil, err
	}

	if tenantInfo != nil {
		return nil, core.NewError(http.StatusConflict, "tenant-info-already-exists")
	}

	collection, err := service.Database.ConnectToCollection(base, model.CollectionTenantInfo)
	if err != nil {
		return nil, err
	}

	item.AddCreationFields()

	result, err := collection.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}

	item.ID = result.InsertedID.(string)
	return item, nil
}

func (service MongoDataService) GetTenantInfo(ctx context.Context, base string) (*model.TenantInfo, error) {
	return service.getTenantInfo(ctx, base)
}

func (service MongoDataService) UpdateTenantInfo(ctx context.Context, base string, item *model.TenantInfo) (*time.Time, int, error) {

	// generate update entry model, ignore the fields that shouldn't be updated
	data := &model.TenantInfo{
		Name:             item.Name,
		LogoURL:          item.LogoURL,
		TermsOfUseURL:    item.TermsOfUseURL,
		PrivacyPolicyURL: item.PrivacyPolicyURL,
		SupportURL:       item.SupportURL,
		SupportEmail:     item.SupportEmail,
	}

	collection, err := service.Database.ConnectToCollection(base, model.CollectionTenantInfo)
	if err != nil {
		return nil, 0, err
	}
	data.PrepareUpdateFields()

	var result model.TenantInfo
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": item.ID}, bson.M{
		"$set": data,
		"$inc": bson.M{"_revision": 1},
	}).Decode(&result)
	if err != nil {
		return nil, 0, err
	}

	return result.UpdatedAt, result.Revision + 1, nil
}

func (service MongoDataService) getTenantInfo(ctx context.Context, base string) (*model.TenantInfo, error) {

	var tenantInfo *model.TenantInfo
	err := service.Database.Find(ctx, base, model.CollectionTenantInfo, bson.M{}, database.FindOptions{Limit: 1}, func(cur *mongo.Cursor) error {
		err := cur.Decode(&tenantInfo)
		if err != nil {
			return err
		}
		return nil
	})

	return tenantInfo, err
}
