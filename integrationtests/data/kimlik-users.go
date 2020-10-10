package data

import (
	"github.com/devingen/api-core/util"
	"github.com/devingen/kimlik-api/model"
)

var Users = []interface{}{
	model.User{
		ID:        util.ObjectIdFromHexIgnoreError("507f191e810c19729de860ea"),
		Email:     "user1@devingen.io",
		FirstName: "User",
		LastName:  "One",
	},
}
