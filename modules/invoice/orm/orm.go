// Package orm provides `GORM` helpers for the creation, migration and access
// on the project's database
package orm

import (
	//Imports the database dialect of choice
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
)

// ORM struct to holds the gorm pointer to db
type ORM struct {
	DB *gorm.DB
}

// GetDB Clone new database connection without search conditions
func (o *ORM) GetDB() *gorm.DB {
	// return o.DB.New()
	// CONTINUE: Implement connection pools
	return o.DB
}
