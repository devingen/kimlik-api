package mongods

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func (service MongoDataService) CreateSAMLConfig(ctx context.Context, base string, item *model.SAMLConfig) (*model.SAMLConfig, error) {
	collection, err := service.Database.ConnectToCollection(base, model.CollectionSAMLConfigs)
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

func (service MongoDataService) GetSAMLConfig(ctx context.Context, base, id string) (*model.SAMLConfig, error) {

	result := model.SAMLConfig{}
	err := service.Database.Get(ctx, base, model.CollectionSAMLConfigs, id, &result)
	if err == mongo.ErrNoDocuments {
		return nil, core.NewError(http.StatusNotFound, "saml-config-not-found")
	}
	return &result, err
}
