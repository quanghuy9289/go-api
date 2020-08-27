package stock

import (
	"fmt"

	"api_new/config"
	"api_new/logger"
	"api_new/modules/stock/handlers"
	"api_new/modules/stock/orm"
	"api_new/modules/stock/orm/migration"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ModuleStock user module
type ModuleStock struct {
}

// RegisterHandlers Register the module's handlers
func (m *ModuleStock) RegisterHandlers(r *gin.Engine, serverCfg *config.ServerConfig, dbCfg *config.DatabaseConfig) {
	db, err := gorm.Open(
		dbCfg.Dialect,
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			dbCfg.Host,
			dbCfg.Port,
			dbCfg.Username,
			dbCfg.Database,
			dbCfg.Password))
	if err != nil {
		logger.Panicf("[ORM] err: %v", err)
	}
	orm := &orm.ORM{
		DB: db,
	}

	// Log every SQL command on dev, @prod: this should be disabled?
	db.LogMode(dbCfg.LogMode)

	// Automigrate tables
	if dbCfg.AutoMigrate {
		err = migration.ServiceAutoMigration(orm.DB)
	}
	logger.Infof("[ORM] Database connection initialized.")

	logger.Infof("######## STOCK MODULE REGISTRATION ########")
	// GraphQL handlers
	// Playground handler
	GQLPlaygroundPath := "/stock/graphql"
	GQLPath := "/stock/query"
	if serverCfg.GQLPlaygroundEnabled {
		r.GET(GQLPlaygroundPath, handlers.PlaygroundHandler(GQLPath))
		logger.Infof("GraphQL Playground @ " + serverCfg.Host + ":" + serverCfg.Port + GQLPlaygroundPath)
	}
	r.POST(GQLPath, handlers.GraphqlHandler(orm))
	logger.Infof("GraphQL @ " + serverCfg.Host + ":" + serverCfg.Port + GQLPath)
}
