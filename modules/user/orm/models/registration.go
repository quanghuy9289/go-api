package models

import (
	"api_new/modules/user/orm"
)

var ()

func init() {

}

// Registration the registration model
type Registration struct {
	RegistrationID int64    `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	CompanyName    string   `gorm:""`
	PhoneNo        string   `gorm:""`
	FaxNo          string   `gorm:""`
	Website        string   `gorm:""`
	StreetAddress  string   `gorm:""`
	Database       Database `gorm:"FOREIGNKEY:RegistrationID"`
}

// ContextRegistration query context
type ContextRegistration struct {
	ORM *orm.ORM
}

// NewContextRegistration new context
func NewContextRegistration(o *orm.ORM) *ContextRegistration {
	return &ContextRegistration{
		ORM: o,
	}
}

// CreateRegistration add new registration
func (mc *ContextRegistration) CreateRegistration(r *Registration) (id int64, err error) {
	if err = mc.ORM.GetDB().Create(r).Error; err == nil {
		id = r.RegistrationID
	}
	return
}

// GetRegistrationByID get registration by id
func (mc *ContextRegistration) GetRegistrationByID(id int64) (reg Registration, err error) {
	err = mc.ORM.GetDB().
		Where("registration_id = ?", id).
		First(&reg).Error
	return
}
