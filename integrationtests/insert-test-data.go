package integrationtests

import (
	"github.com/devingen/api-core/database"
	"github.com/devingen/api-core/util"
	"github.com/devingen/kimlik-api/model"
)

func InsertTestData(db *database.Database, databaseName string) {
	InsertAuths(db, databaseName)
	InsertSessions(db, databaseName)
	InsertUsers(db, databaseName)
}

func InsertAuths(db *database.Database, databaseName string) {
	var list []model.Auth
	util.ReadFile("./data/"+model.CollectionAuths+".json", &list)

	for _, item := range list {
		util.InsertData(db, databaseName, model.CollectionAuths, item)
	}
}

func InsertSessions(db *database.Database, databaseName string) {
	var list []model.Session
	util.ReadFile("./data/"+model.CollectionSessions+".json", &list)

	for _, item := range list {
		util.InsertData(db, databaseName, model.CollectionSessions, item)
	}
}

func InsertUsers(db *database.Database, databaseName string) {
	var list []model.User
	util.ReadFile("./data/"+model.CollectionUsers+".json", &list)

	for _, item := range list {
		util.InsertData(db, databaseName, model.CollectionUsers, item)
	}
}
