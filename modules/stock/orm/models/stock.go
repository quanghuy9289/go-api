package models

import (
	"api_new/modules/stock/orm"
)

var ()

func init() {

}

// Stock the stock model
type Stock struct {
	ID          string `gorm:"PRIMARY_KEY"` //
	CreatedOn   int64  `gorm:""`            //
	Code        string `gorm:"UNIQUE"`      //
	Description string `gorm:"UNIQUE"`      //
	IsActive    bool   `gorm:""`            //
}

// ContextStock query context
type ContextStock struct {
	ORM *orm.ORM
}

// NewContextStock new context
func NewContextStock(o *orm.ORM) *ContextStock {
	return &ContextStock{
		ORM: o,
	}
}

// GetAllStocks query all users
func (mc *ContextStock) GetAllStocks() (users []*Stock, count int, err error) {
	err = mc.ORM.GetDB().
		Order("code", true).
		Find(&users).Count(&count).Error
	return
}

// AddStock add new user
func (mc *ContextStock) AddStock(r *Stock) (id string, err error) {
	if err = mc.ORM.GetDB().Create(r).Error; err == nil {
		id = r.ID
	}
	return
}

// GetStockByID Get user by email
func (mc *ContextStock) GetStockByID(id string) (user Stock, err error) {
	err = mc.ORM.GetDB().
		Where("id = ?", id).
		First(&user).Error
	return
}

// GetStockByCode Get stock by code
func (mc *ContextStock) GetStockByCode(code string) (user Stock, err error) {
	err = mc.ORM.GetDB().
		Where("code = ?", code).
		First(&user).Error
	return
}

// GetStockByDescription Get stock by description
func (mc *ContextStock) GetStockByDescription(description string) (user Stock, err error) {
	err = mc.ORM.GetDB().
		Where("description = ?", description).
		First(&user).Error
	return
}

// GetStockAvatar Get user avatar
func (mc *ContextStock) GetStockAvatar(userID string) (user Stock, err error) {
	err = mc.ORM.GetDB().
		Select("avatar_base64").
		Where("id = ?", userID).
		First(&user).Error
	return
}

// UpdateStockCode update stock code
func (mc *ContextStock) UpdateStockCode(id string, code string) (err error) {
	u := Stock{
		ID: id,
	}
	err = mc.ORM.GetDB().Model(&u).Update("code", code).Error
	return
}

// UpdateStockDescription update stock description
func (mc *ContextStock) UpdateStockDescription(id string, description string) (err error) {
	u := Stock{
		ID: id,
	}
	err = mc.ORM.GetDB().Model(&u).Update("description", description).Error
	return
}

// UpdateStock update stock
func (mc *ContextStock) UpdateStock(u *Stock) (err error) {
	err = mc.ORM.GetDB().Save(&u).Error
	return
}
