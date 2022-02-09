package mongods

import (
	"github.com/devingen/api-core/database"
	ds "github.com/devingen/kimlik-api/data-service"
)

// MongoDataService implements IKimlikDataService interface with MongoDB connection
type MongoDataService struct {
	Database *database.Database
}

// New generates new MongoDataService
func New(database *database.Database) ds.IKimlikDataService {
	return MongoDataService{
		Database: database,
	}
}
