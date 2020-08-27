package models

import (
	"api_new/modules/user/orm"
)

var ()

func init() {

}

// Database the database model
type Database struct {
	DatabaseID     int64  `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	ExpiryDate     int64  `gorm:""`
	DatabaseName   string `gorm:"UNIQUE"`
	RegistrationID int64  `gorm:""`
	IsActive       bool   `gorm:""`
}

// ContextDatabase query context
type ContextDatabase struct {
	ORM *orm.ORM
}

// NewContextDatabase new context
func NewContextDatabase(o *orm.ORM) *ContextDatabase {
	return &ContextDatabase{
		ORM: o,
	}
}

// CreateDatabase add new Database
func (mc *ContextDatabase) CreateDatabase(r *Database) (id int64, err error) {
	if err = mc.ORM.GetDB().Create(r).Error; err == nil {
		id = r.DatabaseID
	}
	return
}

// GetDatabaseByID get Database by id
func (mc *ContextDatabase) GetDatabaseByID(id int64) (res Database, err error) {
	err = mc.ORM.GetDB().
		Where("database_id = ?", id).
		First(&res).Error
	return
}

// GetDatabaseByName get Database by name
func (mc *ContextDatabase) GetDatabaseByName(name string) (res Database, err error) {
	err = mc.ORM.GetDB().
		Where("database_name = ?", name).
		First(&res).Error
	return
}

// GetDatabaseByRegistrationID get Database by registration id
func (mc *ContextDatabase) GetDatabaseByRegistrationID(id int64) (res Database, err error) {
	err = mc.ORM.GetDB().
		Where("registration_id = ? AND is_active = ?", id, true).
		First(&res).Error
	return
}
