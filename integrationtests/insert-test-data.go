package integrationtests

import (
	"github.com/devingen/api-core/database"
	"github.com/devingen/api-core/util"
	"github.com/devingen/kimlik-api/integrationtests/data"
	"github.com/devingen/kimlik-api/model"
)

func InsertTestData(db *database.Database, databaseName string) {
	util.InsertGoData(db, databaseName, model.CollectionAuths, data.Auths)
	util.InsertGoData(db, databaseName, model.CollectionSessions, data.Sessions)
	util.InsertGoData(db, databaseName, model.CollectionUsers, data.Users)
}
