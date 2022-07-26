package mongods

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func (service MongoDataService) GetAPIKey(ctx context.Context, base, id string) (*model.APIKey, error) {

	result := model.APIKey{}
	err := service.Database.Get(ctx, base, model.CollectionAPIKeys, id, &result)
	if err == mongo.ErrNoDocuments {
		return nil, core.NewError(http.StatusNotFound, "api-key-not-found")
	}
	return &result, err
}
