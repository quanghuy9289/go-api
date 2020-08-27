package models

import (
	"api_new/modules/example/orm"
)

var ()

func init() {

}

// Example the example model
type Example struct {
	ID          string `gorm:"PRIMARY_KEY"` //
	CreatedOn   int64  `gorm:""`            //
	Code        string `gorm:"UNIQUE"`      //
	Description string `gorm:"UNIQUE"`      //
	IsActive    bool   `gorm:""`            //
}

// ContextExample query context
type ContextExample struct {
	ORM *orm.ORM
}

// NewContextExample new context
func NewContextExample(o *orm.ORM) *ContextExample {
	return &ContextExample{
		ORM: o,
	}
}

// GetAllExamples query all users
func (mc *ContextExample) GetAllExamples() (users []*Example, count int, err error) {
	err = mc.ORM.GetDB().
		Order("code", true).
		Find(&users).Count(&count).Error
	return
}

// AddExample add new user
func (mc *ContextExample) AddExample(r *Example) (id string, err error) {
	if err = mc.ORM.GetDB().Create(r).Error; err == nil {
		id = r.ID
	}
	return
}

// GetExampleByID Get user by email
func (mc *ContextExample) GetExampleByID(id string) (user Example, err error) {
	err = mc.ORM.GetDB().
		Where("id = ?", id).
		First(&user).Error
	return
}

// GetExampleByCode Get example by code
func (mc *ContextExample) GetExampleByCode(code string) (user Example, err error) {
	err = mc.ORM.GetDB().
		Where("code = ?", code).
		First(&user).Error
	return
}

// GetExampleByDescription Get example by description
func (mc *ContextExample) GetExampleByDescription(description string) (user Example, err error) {
	err = mc.ORM.GetDB().
		Where("description = ?", description).
		First(&user).Error
	return
}

// GetExampleAvatar Get user avatar
func (mc *ContextExample) GetExampleAvatar(userID string) (user Example, err error) {
	err = mc.ORM.GetDB().
		Select("avatar_base64").
		Where("id = ?", userID).
		First(&user).Error
	return
}

// UpdateExampleCode update example code
func (mc *ContextExample) UpdateExampleCode(id string, code string) (err error) {
	u := Example{
		ID: id,
	}
	err = mc.ORM.GetDB().Model(&u).Update("code", code).Error
	return
}

// UpdateExampleDescription update example description
func (mc *ContextExample) UpdateExampleDescription(id string, description string) (err error) {
	u := Example{
		ID: id,
	}
	err = mc.ORM.GetDB().Model(&u).Update("description", description).Error
	return
}

// UpdateExample update example
func (mc *ContextExample) UpdateExample(u *Example) (err error) {
	err = mc.ORM.GetDB().Save(&u).Error
	return
}
