package data

import (
	"github.com/devingen/api-core/util"
	"github.com/devingen/kimlik-api/model"
)

var Auths = []interface{}{
	model.Auth{
		// encrypted version of '123456'
		Password: "$2a$10$2HzFa8e0kYLm20RrGTRg.uADleRhs393FdTugRZW0c/8cFQsc022W",
		Type:     "password",
		User: &model.User{
			Ref:      "kimlik-users",
			ID:       util.ObjectIdFromHexIgnoreError("507f191e810c19729de860ea"),
			Database: "dvn-kimlik-api-integration-test",
		},
	},
}
