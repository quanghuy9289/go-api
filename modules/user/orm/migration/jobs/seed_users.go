package jobs

import (
	"api_new/modules/user/orm"
	"api_new/modules/user/orm/models"
	"api_new/utils"
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
	"time"
)

var (
	firstUser *models.User = &models.User{
		Email:        "admin@cookingthebooks.com.au",
		Fullname:     "Administrator",
		Nickname:     "ad",
		Password:     "admin",
		AvatarBase64: "",
		RoleID:       "admin",
	}
)

// SeedUsers inserts the first users
var SeedUsers *gormigrate.Migration = &gormigrate.Migration{
	ID: "SEED_USERS",
	Migrate: func(db *gorm.DB) error {
		return CreateUser(db, firstUser)
	},
	Rollback: func(db *gorm.DB) error {
		return db.Delete(&firstUser).Error
	},
}

func CreateUser(db *gorm.DB, user *models.User) error {
	password, err := utils.HashAndSalt([]byte(user.Password))
	if err == nil {
		orm := &orm.ORM{
			DB: db,
		}
		// Fill in
		contextUser := models.NewContextUser(orm)
		_, errAddUser := contextUser.AddUser(&models.User{
			ID:           utils.GenerateUUID(),
			CreatedOn:    time.Now().Unix(),
			Email:        user.Email,
			Fullname:     user.Fullname,
			Nickname:     user.Nickname,
			IsActive:     true,
			Password:     password,
			AvatarBase64: user.AvatarBase64,
			RoleID:       user.RoleID,
		})

		return errAddUser
	}

	return nil
}
