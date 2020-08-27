package models

import (
	"api_new/modules/user/orm"
	"time"
)

var ()

func init() {

}

// AuthenticationToken the authenticationtoken model
type AuthenticationToken struct {
	Token          string `gorm:"PRIMARY_KEY"`
	CreatedOn      int64  `gorm:""`
	ExpiredOn      int64  `gorm:""`
	UserID         string `gorm:""` // OneToMany relationship
	DeviceID       string `gorm:""`
	Provider       string `gorm:""`
	ExternalUserID string `gorm:""` // UserId from provider
	Email          string `gorm:""`
	RefreshToken   string `gorm:""`
}

// ContextAuthenticationToken query context
type ContextAuthenticationToken struct {
	ORM *orm.ORM
}

// NewContextAuthenticationToken new context
func NewContextAuthenticationToken(o *orm.ORM) *ContextAuthenticationToken {
	return &ContextAuthenticationToken{
		ORM: o,
	}
}

// AddAuthenticationToken add new authenticationtoken
func (mc *ContextAuthenticationToken) AddAuthenticationToken(r *AuthenticationToken) (token string, err error) {
	if err = mc.ORM.GetDB().Create(r).Error; err == nil {
		token = r.Token
	}
	return
}

// GetAuthenticationTokenByToken Get authenticationtoken by token
func (mc *ContextAuthenticationToken) GetAuthenticationTokenByToken(token string) (authenticationtoken AuthenticationToken, err error) {
	err = mc.ORM.GetDB().
		Where("token = ? AND expired_on > ?", token, time.Now().Unix()).
		First(&authenticationtoken).Error
	return
}

// GetAuthenticationTokenByDeviceId Get authenticationtoken by device id
func (mc *ContextAuthenticationToken) GetAuthenticationTokenByDeviceId(deviceId string) (authenticationtoken AuthenticationToken, err error) {
	err = mc.ORM.GetDB().
		Where("device_id = ? AND expired_on > ?", deviceId, time.Now().Unix()).
		First(&authenticationtoken).Error
	return
}

// SetExpiredAuthenticationToken Set authenticationtoken expired
func (mc *ContextAuthenticationToken) SetExpiredAuthenticationToken(token string) (err error) {
	auth := AuthenticationToken{
		Token: token,
	}
	err = mc.ORM.GetDB().Model(&auth).Update("expired_on", time.Now().Unix()).Error
	return
}

// SetExpiredAllToken Set expired all authentication tokens for current user
func (mc *ContextAuthenticationToken) SetExpiredAllToken(email string) (err error) {

	var authTokens []AuthenticationToken
	err = mc.ORM.GetDB().
		Where("email = ? AND expired_on > ?", email, time.Now().Unix()).Find(&authTokens).Error

	if err == nil {
		for _, item := range authTokens {
			// set expired all token for this user
			mc.SetExpiredAuthenticationToken(item.Token)
		}
	}
	return
}

// GetAuthenticationTokenByExternalUser Get authenticationtoken by external user
func (mc *ContextAuthenticationToken) GetAuthenticationTokenByExternalUser(email string, provider string, externalUserId string) (authenticationtoken AuthenticationToken, err error) {
	err = mc.ORM.GetDB().
		Where("email = ? AND provider = ? AND external_user_id = ? AND expired_on > ?", email, provider, externalUserId, time.Now().Unix()).
		First(&authenticationtoken).Error
	return
}
