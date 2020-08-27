package auth

import (
	"fmt"

	"api_new/config"
	"api_new/logger"
	"api_new/modules/auth/handlers"
	"api_new/modules/user/orm"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

// ModuleUser user module
type ModuleAuth struct {
}

// RegisterHandlers Register the module's handlers
func (m *ModuleAuth) RegisterHandlers(r *gin.Engine, serverCfg *config.ServerConfig, dbCfg *config.DatabaseConfig) {
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

	logger.Infof("[ORM] Database connection initialized.")

	logger.Infof("######## AUTH MODULE REGISTRATION ########")

	// init auth providers
	InitalizeAuthProviders(serverCfg)

	// register callback handler for auth provider
	r.GET("/auth/provider/:provider", handlers.Begin())
	r.GET("/auth/provider/:provider/callback", handlers.Callback(serverCfg, orm))

	GQLPlaygroundPath := "/auth/graphql"
	GQLPath := "/auth/query"
	// Playground handler
	if serverCfg.GQLPlaygroundEnabled {
		r.GET(GQLPlaygroundPath, handlers.PlaygroundHandler(GQLPath))
		logger.Infof("GraphQL Playground @ " + serverCfg.Host + ":" + serverCfg.Port + GQLPlaygroundPath)
	}
	// GraphQL handlers
	r.POST(GQLPath, handlers.GraphqlHandler(orm, dbCfg))
	logger.Infof("GraphQL @ " + serverCfg.Host + ":" + serverCfg.Port + GQLPath)
}

// InitalizeAuthProviders does just that, with Goth providers
func InitalizeAuthProviders(cfg *config.ServerConfig) error {
	providers := []goth.Provider{}
	// Initialize Goth providers
	for _, p := range cfg.AuthProviders {
		switch p.Provider {
		case "google":
			providers = append(
				providers,
				google.New(p.ClientKey, p.Secret, "http://"+cfg.Host+":"+cfg.Port+"/auth/provider/"+p.Provider+"/callback", p.Scopes...),
			)
		}
	}

	goth.UseProviders(providers...)
	return nil
}
