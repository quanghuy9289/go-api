package modules

import (
	"api_new/config"
	"github.com/gin-gonic/gin"
)

// Module module interface
type Module interface {
	RegisterHandlers(r *gin.Engine, serverCfg *config.ServerConfig, dbCfg *config.DatabaseConfig)
}

