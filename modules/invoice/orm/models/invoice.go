package models

import (
	"api_new/modules/invoice/orm"
)

var ()

func init() {

}

// Invoice the invoice model
type Invoice struct {
	ID          string `gorm:"PRIMARY_KEY"` //
	CreatedOn   int64  `gorm:""`            //
	Code        string `gorm:"UNIQUE"`      //
	Description string `gorm:"UNIQUE"`      //
	IsActive    bool   `gorm:""`            //
}

// ContextInvoice query context
type ContextInvoice struct {
	ORM *orm.ORM
}

// NewContextInvoice new context
func NewContextInvoice(o *orm.ORM) *ContextInvoice {
	return &ContextInvoice{
		ORM: o,
	}
}

// GetAllInvoices query all users
func (mc *ContextInvoice) GetAllInvoices() (users []*Invoice, count int, err error) {
	err = mc.ORM.GetDB().
		Order("code", true).
		Find(&users).Count(&count).Error
	return
}

// AddInvoice add new user
func (mc *ContextInvoice) AddInvoice(r *Invoice) (id string, err error) {
	if err = mc.ORM.GetDB().Create(r).Error; err == nil {
		id = r.ID
	}
	return
}

// GetInvoiceByID Get user by email
func (mc *ContextInvoice) GetInvoiceByID(id string) (user Invoice, err error) {
	err = mc.ORM.GetDB().
		Where("id = ?", id).
		First(&user).Error
	return
}

// GetInvoiceByCode Get invoice by code
func (mc *ContextInvoice) GetInvoiceByCode(code string) (user Invoice, err error) {
	err = mc.ORM.GetDB().
		Where("code = ?", code).
		First(&user).Error
	return
}

// GetInvoiceByDescription Get invoice by description
func (mc *ContextInvoice) GetInvoiceByDescription(description string) (user Invoice, err error) {
	err = mc.ORM.GetDB().
		Where("description = ?", description).
		First(&user).Error
	return
}

// GetInvoiceAvatar Get user avatar
func (mc *ContextInvoice) GetInvoiceAvatar(userID string) (user Invoice, err error) {
	err = mc.ORM.GetDB().
		Select("avatar_base64").
		Where("id = ?", userID).
		First(&user).Error
	return
}

// UpdateInvoiceCode update invoice code
func (mc *ContextInvoice) UpdateInvoiceCode(id string, code string) (err error) {
	u := Invoice{
		ID: id,
	}
	err = mc.ORM.GetDB().Model(&u).Update("code", code).Error
	return
}

// UpdateInvoiceDescription update invoice description
func (mc *ContextInvoice) UpdateInvoiceDescription(id string, description string) (err error) {
	u := Invoice{
		ID: id,
	}
	err = mc.ORM.GetDB().Model(&u).Update("description", description).Error
	return
}

// UpdateInvoice update invoice
func (mc *ContextInvoice) UpdateInvoice(u *Invoice) (err error) {
	err = mc.ORM.GetDB().Save(&u).Error
	return
}
