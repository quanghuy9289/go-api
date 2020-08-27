package services

import (
	"api_new/config"
	"api_new/logger"
	"fmt"
	"strings"

	invoiceMigration "api_new/modules/invoice/orm/migration"
	"github.com/jinzhu/gorm"
)

// MigrationClientDabatase run migration for client db
func MigrationClientDabatase(databaseName string, dbConfig *config.DatabaseConfig) error {
	db, err := gorm.Open(
		dbConfig.Dialect,
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.Username,
			strings.ToLower(databaseName),
			dbConfig.Password))
	if err != nil {
		logger.Panicf("[ORM] err: %v", err)
		return err
	}
	defer db.Close()

	// run migration
	err = invoiceMigration.ServiceAutoMigration(db)

	return err
}
