package handlers

import (
	"api_new/middleware"
	"time"

	// "github.com/dgrijalva/jwt-go"
	"api_new/config"
	"api_new/modules/user/orm"
	"api_new/modules/user/orm/models"

	"github.com/gin-contrib/sessions"

	"api_new/logger"

	"github.com/gin-gonic/gin"

	"api_new/utils"

	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/markbates/goth/gothic"
	// "time"
)

// Begin login with the auth provider
func Begin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// You have to add value context with provider name to get provider name in GetProviderName method
		c.Request = middleware.AddProviderToContext(c, c.Param("provider"))
		// try to get the user without re-authenticating
		if gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err != nil {
			gothic.BeginAuthHandler(c.Writer, c.Request)
		} else {
			logger.Debugf("user: %#v", gothUser)
		}
	}
}

// Callback callback to complete auth provider flow
func Callback(cfg *config.ServerConfig, orm *orm.ORM) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add value context with provider name to get provider name in GetProviderName method
		c.Request = middleware.AddProviderToContext(c, c.Param("provider"))
		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// find user in db
		contextUser := models.NewContextUser(orm)
		u, err := contextUser.GetUserByExternalProvider(user.Email, user.Provider, user.UserID)
		// logger.Debugf("gothUser: %#v", user)
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				// add external user
				contextUser.AddUser(&models.User{
					ID:        user.UserID,
					CreatedOn: time.Now().Unix(),
					Email:     user.Email,
					Fullname:  user.Name,
					Nickname:  user.NickName,
					IsActive:  true,
				})
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}

		logger.Debug("[Auth.CallBack.UserLoggedIn]: ", u.ID)

		// generate token
		token, err := utils.GenerateJwtToken(user.Email)
		if err == nil {
			contextAuth := models.NewContextAuthenticationToken(orm)

			// get device id from request
			var deviceID = c.Request.RemoteAddr
			currentAuthToken, err := contextAuth.GetAuthenticationTokenByDeviceId(deviceID)

			if err == nil {
				// set expired for current token
				contextAuth.SetExpiredAuthenticationToken(currentAuthToken.Token)
			}

			authToken := &models.AuthenticationToken{
				Token:          token,
				UserID:         user.UserID,
				CreatedOn:      time.Now().UTC().Unix(),
				ExpiredOn:      time.Now().Add(time.Hour * time.Duration(24)).Unix(), // using expired setting by our system instead of 3rd party
				ExternalUserID: user.UserID,
				DeviceID:       deviceID,
				Email:          user.Email,
				Provider:       user.Provider,
			}
			contextAuth.AddAuthenticationToken(authToken)

			// add to session
			session := sessions.Default(c)
			session.Set("userId", user.UserID)
			session.Set("email", user.Email)
			session.Set("token", token)
			session.Save()

			json := gin.H{
				"type":          "Bearer",
				"token":         token,
				"refresh_token": user.RefreshToken,
			}
			c.JSON(http.StatusOK, json)
		}
	}
}
