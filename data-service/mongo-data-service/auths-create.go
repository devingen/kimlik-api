package mongods

import (
	"context"
	"errors"

	"github.com/devingen/kimlik-api/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (service MongoDataService) CreateAuthWithPassword(ctx context.Context, base, password string, user *model.User) (*model.Auth, error) {
	collection, err := service.Database.ConnectToCollection(base, model.CollectionAuths)
	if err != nil {
		return nil, err
	}

	item := &model.Auth{
		Type:     model.AuthTypePassword,
		User:     user.DBRef(base),
		Password: password,
	}
	item.AddCreationFields()

	err = item.HashPassword()
	if err != nil {
		return nil, err
	}

	result, err := collection.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}

	item.ID = result.InsertedID.(primitive.ObjectID)
	return item, nil
}

func (service MongoDataService) CreateAuthWithIDToken(ctx context.Context, base string, claims map[string]interface{}, user *model.User) (*model.Auth, error) {

	issuer, ok := claims["iss"].(string)
	if !ok {
		return nil, errors.New("issuer-missing-in-token-claims")
	}

	audience, ok := claims["aud"].(string)
	if !ok {
		return nil, errors.New("audience-missing-in-token-claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("subject-missing-in-token-claims")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("email-missing-in-token-claims")
	}

	collection, err := service.Database.ConnectToCollection(base, model.CollectionAuths)
	if err != nil {
		return nil, err
	}

	item := &model.Auth{
		Type: model.AuthTypeOpenID,
		User: user.DBRef(base),
		OpenID: &model.OpenIDData{
			Iss:   issuer,
			Aud:   audience,
			Sub:   subject,
			Email: email,
		},
	}
	item.AddCreationFields()

	result, err := collection.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}

	item.ID = result.InsertedID.(primitive.ObjectID)
	return item, nil
}
