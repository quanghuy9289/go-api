package config

import (
	"api_new/utils"
)

// AuthProvider defines the configuration for the Goth config
type AuthProvider struct {
	Provider  string
	ClientKey string
	Secret    string
	Domain    string // If needed, like with auth0
	Scopes    []string
}

// ServerConfig Server configuration
type ServerConfig struct {
	Host                 string
	Port                 string
	GQLPlaygroundEnabled bool
	AuthProviders        []AuthProvider
	JWTSecret            string
	JWTAlgorithm         string
	// GQLPath              string
	// GQLPlaygroundPath    string
}

// LoadEnvironmentFile load environment file to overwrite current ENV
func LoadEnvironmentFile(envfile string) error {
	return utils.LoadEnvironmentFile(envfile)
}

// LoadServerConfigFromEnvironment load server configuration from environment
func LoadServerConfigFromEnvironment() (c *ServerConfig, err error) {
	c = &ServerConfig{}

	c.Host, err = utils.MustGet("SERVER_HOST")
	c.Port, err = utils.MustGet("SERVER_PORT")
	// c.GQLPath, err = utils.MustGet("GQL_SERVER_GRAPHQL_PATH")
	// c.GQLPlaygroundPath, err = utils.MustGet("GQL_SERVER_GRAPHQL_PLAYGROUND_PATH")
	c.GQLPlaygroundEnabled, err = utils.MustGetBool("GQL_SERVER_GRAPHQL_PLAYGROUND_ENABLED", false)

	c.Host, err = utils.MustGet("SERVER_HOST")
	c.Port, err = utils.MustGet("SERVER_PORT")

	c.JWTSecret, err = utils.MustGet("JWT_SECRET")
	c.JWTAlgorithm, err = utils.MustGet("JWT_SIGNING_ALGORITHM")

	googleClientKey, err := utils.MustGet("PROVIDER_GOOGLE_KEY")
	googleClientSecret, err := utils.MustGet("PROVIDER_GOOGLE_SECRET")

	googleProvider := AuthProvider{
		Provider:  "google",
		ClientKey: googleClientKey,
		Secret:    googleClientSecret,
	}

	c.AuthProviders = []AuthProvider{googleProvider}

	if err != nil {
		return nil, err
	}

	return c, nil
}

// DatabaseConfig database configuration
type DatabaseConfig struct {
	Dialect     string
	Host        string
	Port        string
	Username    string
	Password    string
	Database    string
	SeedDB      bool // ?
	LogMode     bool // ?
	AutoMigrate bool // ?
}

// LoadDatabaseConfigFromEnvironment load database configuration from environment
func LoadDatabaseConfigFromEnvironment() (c *DatabaseConfig, err error) {
	c = &DatabaseConfig{}

	c.Dialect, err = utils.MustGet("GORM_DIALECT")
	c.Host, err = utils.MustGet("DB_SERVER_HOST")
	c.Port, err = utils.MustGet("DB_SERVER_PORT")
	c.Username, err = utils.MustGet("DB_SERVER_USER")
	c.Password, err = utils.MustGet("DB_SERVER_PASS")
	c.Database, err = utils.MustGet("DB_NAME")
	c.SeedDB, err = utils.MustGetBool("GORM_SEED_DB", false)
	c.LogMode, err = utils.MustGetBool("GORM_LOGMODE", false)
	c.AutoMigrate, err = utils.MustGetBool("GORM_AUTOMIGRATE", false)

	if err != nil {
		return nil, err
	}

	return c, nil
}
