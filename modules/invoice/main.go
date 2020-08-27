package invoice

import (
	"api_new/config"
	"api_new/logger"
	"api_new/modules/invoice/handlers"

	"github.com/gin-gonic/gin"
)

// ModuleInvoice user module
type ModuleInvoice struct {
}

// RegisterHandlers Register the module's handlers
func (m *ModuleInvoice) RegisterHandlers(r *gin.Engine, serverCfg *config.ServerConfig, dbCfg *config.DatabaseConfig) {
	// db, err := gorm.Open(
	// 	dbCfg.Dialect,
	// 	fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
	// 		dbCfg.Host,
	// 		dbCfg.Port,
	// 		dbCfg.Username,
	// 		dbCfg.Database,
	// 		dbCfg.Password))
	// if err != nil {
	// 	logger.Panicf("[ORM] err: %v", err)
	// }
	// orm := &orm.ORM{
	// 	DB: db,
	// }

	// // Log every SQL command on dev, @prod: this should be disabled?
	// db.LogMode(dbCfg.LogMode)

	// // Automigrate tables
	// if dbCfg.AutoMigrate {
	// 	err = migration.ServiceAutoMigration(orm.DB)
	// }
	// logger.Infof("[ORM] Database connection initialized.")

	logger.Infof("######## INVOICE MODULE REGISTRATION ########")
	// GraphQL handlers
	// Playground handler
	GQLPlaygroundPath := "/invoice/graphql"
	GQLPath := "/invoice/query"
	if serverCfg.GQLPlaygroundEnabled {
		r.GET(GQLPlaygroundPath, handlers.PlaygroundHandler(GQLPath))
		logger.Infof("GraphQL Playground @ " + serverCfg.Host + ":" + serverCfg.Port + GQLPlaygroundPath)
	}
	r.POST(GQLPath, handlers.GraphqlHandler(dbCfg))
	logger.Infof("GraphQL @ " + serverCfg.Host + ":" + serverCfg.Port + GQLPath)
}
