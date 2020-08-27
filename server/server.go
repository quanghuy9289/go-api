package server

import (
	"fmt"
	"time"

	"api_new/config"
	"api_new/logger"
	"api_new/modules"
	"api_new/modules/auth"
	"api_new/modules/example"
	"api_new/modules/invoice"
	"api_new/modules/stock"
	"api_new/modules/user"

	"github.com/gin-contrib/sessions"
	// "github.com/gin-contrib/sessions/cookie"
	"api_new/middleware"
	// "github.com/gin-contrib/sessions/redis"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

var serverCfg *config.ServerConfig
var dbCfg *config.DatabaseConfig
var sharedServer *Server

func init() {
	// Do nothing for now
}

// Server Server struct
type Server struct {
	ServerConfig   *config.ServerConfig
	DatabaseConfig *config.DatabaseConfig
	Modules        chan modules.Module
}

// NewServer Create new server instance
func NewServer() *Server {
	var err error
	serverCfg, err = config.LoadServerConfigFromEnvironment()
	if err != nil {
		logger.Panicf("Server config error: %v", err)
		panic(fmt.Sprintf("Server config error: %v", err))
	}

	dbCfg, err = config.LoadDatabaseConfigFromEnvironment()
	if err != nil {
		logger.Panicf("Database config error: %v", err)
		panic(fmt.Sprintf("Database config error: %v", err))
	}

	return &Server{
		ServerConfig:   serverCfg,
		DatabaseConfig: dbCfg,
		Modules:        make(chan modules.Module),
	}
}

// Run Run the web server
func (s *Server) Run() {
	logger.Infof("Server is running...")
	r := gin.Default()

	// register global middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://ctb.com", "http://api.ctb.com"},
		AllowMethods:     []string{"PUT", "POST"},
		AllowHeaders:     []string{"Origin", "authorization", "content-type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "http://localhost:3000"
		// },
		MaxAge: 12 * time.Hour,
	}))

	r.Use(middleware.GinContextToContextMiddleware())

	//session handler
	// store := cookie.NewStore([]byte("secret"))
	store := memstore.NewStore([]byte("secret"))
	// store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))

	r.Use(sessions.Sessions("CTBSession", store))

	// Register modules (demo registering module on the fly)
	go func(s *Server) {
		time.Sleep(2 * time.Second)
		s.Modules <- &example.ModuleExample{}
	}(s)

	go func(s *Server) {
		s.Modules <- &user.ModuleUser{}
		s.Modules <- &stock.ModuleStock{}
		s.Modules <- &auth.ModuleAuth{}

		// Client
		s.Modules <- &invoice.ModuleInvoice{}
	}(s)

	// Register the modules on the fly
	go func(s *Server) {
		for {
			select {
			case module := <-s.Modules:
				module.RegisterHandlers(r, s.ServerConfig, s.DatabaseConfig)
			case <-time.After(3 * time.Second):
				{
					logger.Infof("NO MORE MODULES")
					return
				}
			}
		}
	}(s)
	// user.Register(r, s.ServerConfig, s.DatabaseConfig)
	// stock.Register(r, s.ServerConfig, s.DatabaseConfig)

	// logger.Fatalf("Server.Run: %v", r.Run(serverCfg.Host+":"+serverCfg.Port))
	logger.Fatalf("Server.Run: %v", r.Run(":"+serverCfg.Port))
}

// Run Run the web server with specified port (optional) and environment file (optional)
func Run(port *int, envfile *string) {
	// Load new env file if specified
	loadedEnvironmentFile := false
	if envfile != nil {
		err := config.LoadEnvironmentFile(*envfile)
		if err != nil {
			logger.Errorf("Error loading .env file: %v", err)
		} else {
			loadedEnvironmentFile = true
			logger.Infof("Loaded .env file: %s", *envfile)
		}
	}

	// Load default environment file if can't load the specified environment file
	if !loadedEnvironmentFile {
		defaultEnvFile := ".env"
		logger.Infof("Will try loading default .env file: %s ...", defaultEnvFile)
		err := config.LoadEnvironmentFile(defaultEnvFile)
		if err != nil {
			logger.Fatalf("Error loading default .env file: %s, err = %v", defaultEnvFile, err)
		} else {
			logger.Infof("Loaded default .env file: %s\n", defaultEnvFile)
		}
	}

	// Create server instance
	sharedServer = NewServer()

	// Overwrite the port if parameter is provided
	if port != nil {
		sharedServer.ServerConfig.Port = fmt.Sprintf("%v", *port)
	}

	sharedServer.Run()
}

// GetSharedServer Get shared server
func GetSharedServer() *Server {
	return sharedServer
}
