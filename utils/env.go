package utils

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	// LoadEnvironmentFile("config/.env")
	// Do nothing for now
}

// LoadEnvironmentFile read your env file(s) and load them into ENV for this process.
func LoadEnvironmentFile(envfile string) error {
	// load .env file
	err := godotenv.Overload(envfile)
	return err

	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %s", envfile)
	// } else {
	// 	log.Printf("Loaded .env file: %s\n", envfile)
	// }
}

// MustGet will return the env or return error if it is not present
func MustGet(k string) (string, error) {
	v := os.Getenv(k)
	if v == "" {
		return v, fmt.Errorf("ENV missing, key: %s", k)
	}
	return v, nil
}

// MustGetBool will return the env as boolean or return error if it is not present
func MustGetBool(k string, defaultValue bool) (bool, error) {
	v := os.Getenv(k)
	if v == "" {
		return defaultValue, fmt.Errorf("ENV missing, key: %s", k)
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return defaultValue, err
	}
	return b, nil
}
