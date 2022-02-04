package data

import (
	"github.com/devingen/api-core/util"
	"github.com/devingen/kimlik-api/model"
)

var Sessions = []interface{}{
	model.Session{
		Client:       "",
		IP:           "1.2.3.4",
		SessionCount: 0,
		Status:       "successful",
		User: &model.User{
			Ref:      "kimlik-users",
			ID:       util.ObjectIdFromHexIgnoreError("507f191e810c19729de860ea"),
			Database: "dvn-kimlik-api-integration-test",
		},
		UserAgent: "",
	},
}
