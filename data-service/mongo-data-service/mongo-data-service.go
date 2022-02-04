package mongods

import (
	"github.com/devingen/api-core/database"
	ds "github.com/devingen/kimlik-api/data-service"
)

// MongoDataService implements IKimlikDataService interface with MongoDB connection
type MongoDataService struct {
	DatabaseName string
	Database     *database.Database
}

// New generates new MongoDataService
func New(databaseName string, database *database.Database) ds.IKimlikDataService {
	return MongoDataService{
		DatabaseName: databaseName,
		Database:     database,
	}
}
