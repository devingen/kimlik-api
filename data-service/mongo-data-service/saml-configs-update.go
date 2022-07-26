package mongods

import (
	"context"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (service MongoDataService) UpdateSAMLConfig(ctx context.Context, base string, item *model.SAMLConfig) (*time.Time, int, error) {

	// generate update entry model, ignore the fields that shouldn't be updated
	data := &model.SAMLConfig{
		Name:                        item.Name,
		MetadataURL:                 item.MetadataURL,
		MetadataContent:             item.MetadataContent,
		AssertionConsumerServiceURL: item.AssertionConsumerServiceURL,
		AudienceURI:                 item.AudienceURI,
		ServiceProviderIssuer:       item.ServiceProviderIssuer,
		AttributeKeyMappingTemplate: item.AttributeKeyMappingTemplate,
		MetaAttributeKeyMapping:     item.MetaAttributeKeyMapping,
	}

	collection, err := service.Database.ConnectToCollection(base, model.CollectionSAMLConfigs)
	if err != nil {
		return nil, 0, err
	}
	data.PrepareUpdateFields()

	var result model.SAMLConfig
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": item.ID}, bson.M{
		"$set": data,
		"$inc": bson.M{"_revision": 1},
	}).Decode(&result)
	if err != nil {
		return nil, 0, err
	}

	return result.UpdatedAt, result.Revision + 1, nil
}
