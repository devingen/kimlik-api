package mongods

import (
	"context"

	core "github.com/devingen/api-core"
	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (service MongoDataService) CreateUser(ctx context.Context, base, firstName, lastName, email string, status model.UserStatus, isEmailVerified bool) (*model.User, error) {
	collection, err := service.Database.ConnectToCollection(base, model.CollectionUsers)
	if err != nil {
		return nil, err
	}

	item := &model.User{
		FirstName:       core.String(firstName),
		LastName:        core.String(lastName),
		Name:            core.String(firstName + " " + lastName),
		Email:           core.String(email),
		Status:          &status,
		IsEmailVerified: &isEmailVerified,
	}
	item.AddCreationFields()

	result, err := collection.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}

	item.ID = result.InsertedID.(primitive.ObjectID)
	return item, nil
}
