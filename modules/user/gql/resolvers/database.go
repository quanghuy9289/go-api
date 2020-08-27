package resolvers

import (
	"api_new/logger"
	"api_new/services"
	"context"
	"fmt"
	"github.com/jinzhu/gorm"

	generated "api_new/modules/user/gql/models"
	// "api_new/logger"
	"api_new/modules/user/orm/models"
	"github.com/vektah/gqlparser/gqlerror"
	"strconv"
)

// </QUERY>

// <MUTATION>
// CreateNewDatabase create a database for client
func (r *mutationResolver) CreateNewDatabase(ctx context.Context, input generated.DatabaseInput) (*models.Database, error) {
	ctxRegistration := models.NewContextRegistration(r.ORM)
	_, err := ctxRegistration.GetRegistrationByID(int64(input.RegistrationID))
	if err == nil {
		context := models.NewContextDatabase(r.ORM)
		db, errGetDB := context.GetDatabaseByRegistrationID(int64(input.RegistrationID))
		if gorm.IsRecordNotFoundError(errGetDB) { // db not exist, so create new database
			logger.Info(db.DatabaseName)
			newDatabaseID, err := context.CreateDatabase(&models.Database{
				ExpiryDate:     int64(input.ExpiryDate),
				DatabaseName:   "CTB_User_" + strconv.Itoa(input.RegistrationID) + "_" + input.DatabaseName,
				RegistrationID: int64(input.RegistrationID),
				IsActive:       true,
			})
			if err == nil {
				database, errGetDatabase := context.GetDatabaseByID(newDatabaseID)
				if errGetDatabase == nil {
					// persist new database
					sqlQuery := fmt.Sprintf("CREATE DATABASE %s", database.DatabaseName)
					newDB := r.ORM.GetDB().Exec(sqlQuery)
					if newDB.Error != nil {
						return nil, gqlerror.Errorf("Unable to create DB: %v", newDB.Error)
					}

					// run migration for new database
					err = services.MigrationClientDabatase(database.DatabaseName, r.DBConfig)
					if err != nil {
						return nil, gqlerror.Errorf("Migration database error: %v", err)
					}

					return &database, nil
				}
				return nil, gqlerror.Errorf("Could not query the Database information: %v", errGetDatabase)
			}
			return nil, gqlerror.Errorf("Could not create the Database: %v", err)
		}
		return nil, gqlerror.Errorf("Database for registration %v already exist", input.RegistrationID)
	}
	return nil, gqlerror.Errorf("RegistrationID %v is not exist", input.RegistrationID)

}

// </MUTATION>
