package models

import (
	"api_new/modules/user/orm"
	"time"
)

var ()

func init() {

}

// User the user model
type User struct {
	ID                                 string                 `gorm:"PRIMARY_KEY"`       //
	CreatedOn                          int64                  `gorm:""`                  //
	Email                              string                 `gorm:"UNIQUE"`            //
	PhoneNumber                        string                 `gorm:""`                  //
	Password                           string                 `gorm:""`                  //
	Fullname                           string                 `gorm:""`                  //
	Nickname                           string                 `gorm:""`                  //
	AvatarBase64                       string                 `gorm:""`                  //
	RoleID                             string                 `gorm:""`                  // A user has one role
	StartDate                          *int                   `gorm:""`                  //
	MinDailyUnitPointsRequirement      *float64               `gorm:""`                  // Minimum daily unit points requirement for this employee
	StandardNumberOfWorkingDaysPerWeek *float64               `gorm:""`                  // Number of working days per week that is required to be completed by the employee
	Config                             string                 `gorm:""`                  // JSON configuration to be used for user specific configuration
	IsActive                           bool                   `gorm:""`                  //
	AuthenticationTokens               []*AuthenticationToken `gorm:"FOREIGNKEY:UserID"` // A user has many authentication tokens
	InUseDatabase                      string                 `gorm:""`
}

// ContextUser query context
type ContextUser struct {
	ORM *orm.ORM
}

// NewContextUser new context
func NewContextUser(o *orm.ORM) *ContextUser {
	return &ContextUser{
		ORM: o,
	}
}

// GetAllUsers query all users
func (mc *ContextUser) GetAllUsers() (users []*User, count int, err error) {
	err = mc.ORM.GetDB().
		Order("fullname", true).
		Find(&users).Count(&count).Error
	return
}

// AddUser add new user
func (mc *ContextUser) AddUser(r *User) (id string, err error) {
	if err = mc.ORM.GetDB().Create(r).Error; err == nil {
		id = r.ID
	}
	return
}

// DeleteUser delete existing user
func (mc *ContextUser) DeleteUser(r *User) (err error) {
	err = mc.ORM.GetDB().Where("id = ?", r.ID).Delete(r).Error
	return
}

// GetUserByID Get user by email
func (mc *ContextUser) GetUserByID(id string) (user User, err error) {
	err = mc.ORM.GetDB().
		Where("id = ?", id).
		First(&user).Error
	return
}

// GetUserByEmail Get user by email
func (mc *ContextUser) GetUserByEmail(email string) (user User, err error) {
	err = mc.ORM.GetDB().
		Where("email = ?", email).
		First(&user).Error
	return
}

// GetUserAvatar Get user avatar
func (mc *ContextUser) GetUserAvatar(userID string) (user User, err error) {
	err = mc.ORM.GetDB().
		Select("avatar_base64").
		Where("id = ?", userID).
		First(&user).Error
	return
}

// GetUserByAuthenticationToken Get user by authentication
func (mc *ContextUser) GetUserByAuthenticationToken(token string) (user User, err error) {
	var authToken AuthenticationToken
	contextAuthenticationToken := NewContextAuthenticationToken(mc.ORM) // Use same context
	authToken, err = contextAuthenticationToken.GetAuthenticationTokenByToken(token)
	if err == nil {
		if len(authToken.UserID) > 0 {
			err = mc.ORM.GetDB().Where("id = ?", authToken.UserID).First(&user).Error
		}
	}
	return
}

// UpdateUserPassword update user password
func (mc *ContextUser) UpdateUserPassword(id string, password string) (err error) {
	u := User{
		ID: id,
	}
	err = mc.ORM.GetDB().Model(&u).Update("password", password).Error
	return
}

// UpdateUserConfig update user config
func (mc *ContextUser) UpdateUserConfig(id string, config string) (err error) {
	u := User{
		ID: id,
	}
	err = mc.ORM.GetDB().Model(&u).Update("config", config).Error
	return
}

// UpdateUser update user
func (mc *ContextUser) UpdateUser(u *User) (err error) {
	err = mc.ORM.GetDB().Save(&u).Error
	return
}

// GetUserActiveAuthenticationTokens get active authentication tokens of the user
func (mc *ContextUser) GetUserActiveAuthenticationTokens(u *User) (authenticationTokens []*AuthenticationToken, err error) {
	// err = mc.ORM.GetDB().
	// 	Where("shortcode = ?", shortcode).
	// 	First(&p).Error
	err = mc.ORM.GetDB().
		Model(u).
		Where("expired_on > ?", time.Now().Unix()).
		Related(&authenticationTokens, "AuthenticationTokens").
		Error
	return
}

// GetUserByExternalProvider Get user by external authentication - oauth2
func (mc *ContextUser) GetUserByExternalProvider(email string, provider string, externalUserId string) (user User, err error) {
	var authToken AuthenticationToken
	contextAuthenticationToken := NewContextAuthenticationToken(mc.ORM) // Use same context
	authToken, err = contextAuthenticationToken.GetAuthenticationTokenByExternalUser(email, provider, externalUserId)
	if err == nil {
		if len(authToken.ExternalUserID) > 0 {
			err = mc.ORM.GetDB().Where("id = ?", authToken.ExternalUserID).First(&user).Error
		}
	}
	return
}
