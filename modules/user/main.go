package user

import (
	"fmt"

	"api_new/config"
	"api_new/logger"
	"api_new/modules/user/handlers"
	"api_new/modules/user/orm"
	"api_new/modules/user/orm/migration"

	"api_new/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ModuleUser user module
type ModuleUser struct {
}

// RegisterHandlers Register the module's handlers
func (m *ModuleUser) RegisterHandlers(r *gin.Engine, serverCfg *config.ServerConfig, dbCfg *config.DatabaseConfig) {
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
	orm.SharedORM = &orm.ORM{
		DB: db,
	}

	// TODO: save cache orm

	// Log every SQL command on dev, @prod: this should be disabled?
	db.LogMode(dbCfg.LogMode)

	// Automigrate tables
	if dbCfg.AutoMigrate {
		err = migration.ServiceAutoMigration(orm.SharedORM.DB)
	}
	logger.Infof("[ORM] Database connection initialized.")

	logger.Infof("######## USER MODULE REGISTRATION ########")

	// GraphQL handlers
	// Playground handler
	GQLPlaygroundPath := "/user/graphql"
	GQLPath := "/user/query"
	if serverCfg.GQLPlaygroundEnabled {
		r.GET(GQLPlaygroundPath, handlers.PlaygroundHandler(GQLPath))
		logger.Infof("GraphQL Playground @ " + serverCfg.Host + ":" + serverCfg.Port + GQLPlaygroundPath)
	}

	// add authentication middleware to graphql query
	r.POST(GQLPath, middleware.AuthMiddleware(orm.SharedORM), handlers.GraphqlHandler(orm.SharedORM, dbCfg))
	logger.Infof("GraphQL @ " + serverCfg.Host + ":" + serverCfg.Port + GQLPath)
}
