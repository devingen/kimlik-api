package mongods

import (
	"context"
	"github.com/devingen/kimlik-api/model"
)

func (service MongoDataService) CreateOAuthAccessCode(ctx context.Context, base string, item *model.OAuthAccessCode) (*model.OAuthAccessCode, error) {
	collection, err := service.Database.ConnectToCollection(base, model.CollectionAccessCodes)
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
