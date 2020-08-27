package example

import (
	"fmt"

	"api_new/config"
	"api_new/logger"
	"api_new/modules/example/handlers"
	"api_new/modules/example/orm"
	"api_new/modules/example/orm/migration"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ModuleExample user module
type ModuleExample struct {
}

// RegisterHandlers Register the module's handlers
func (m *ModuleExample) RegisterHandlers(r *gin.Engine, serverCfg *config.ServerConfig, dbCfg *config.DatabaseConfig) {
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

	logger.Infof("######## EXAMPLE MODULE REGISTRATION ########")
	// GraphQL handlers
	// Playground handler
	GQLPlaygroundPath := "/example/graphql"
	GQLPath := "/example/query"
	if serverCfg.GQLPlaygroundEnabled {
		r.GET(GQLPlaygroundPath, handlers.PlaygroundHandler(GQLPath))
		logger.Infof("GraphQL Playground @ " + serverCfg.Host + ":" + serverCfg.Port + GQLPlaygroundPath)
	}
	r.POST(GQLPath, handlers.GraphqlHandler(orm))
	logger.Infof("GraphQL @ " + serverCfg.Host + ":" + serverCfg.Port + GQLPath)
}
