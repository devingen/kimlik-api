package data

import (
	coremodel "github.com/devingen/api-core/model"
	"github.com/devingen/api-core/util"
	"github.com/devingen/kimlik-api/model"
)

var Auths = []interface{}{
	model.Auth{
		Password: "$2a$10$UVN567IGdcUPXRIUuII3HeWY6ZnJ/tT7h1DX2VjbdHuna1UDaUQlu",
		Type:     "password",
		User: &coremodel.DBRef{
			Ref:      "kimlik-users",
			ID:       util.ObjectIdFromHexIgnoreError("507f191e810c19729de860ea"),
			Database: "dvn-kimlik-api-integration-test",
		},
	},
}
