package middleware

import (
	"api_new/modules/user/orm"
	"api_new/modules/user/orm/models"
	"api_new/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func authError(c *gin.Context, err error) {
	e := gin.H{"Error": err.Error()}
	c.AbortWithStatusJSON(http.StatusUnauthorized, e)
}

// AuthMiddleware wraps the request with auth middleware
func AuthMiddleware(orm *orm.ORM) gin.HandlerFunc {
	return func(c *gin.Context) {
		if token, err := JWTFromHeader(c); err == nil {
			contextUser := models.NewContextUser(orm)
			_, err := contextUser.GetUserByAuthenticationToken(token)

			if err != nil {
				authError(c, utils.ErrForbidden)
			} else {
				c.Next()
			}
			// session := sessions.Default(c)
			// tokenSession := session.Get("token")
			// if tokenSession == nil {
			// 	authError(c, utils.ErrForbidden)
			// } else {
			// 	if token != tokenSession {
			// 		authError(c, utils.ErrForbidden)
			// 	}
			// }

			// c.Next()
		} else {
			authError(c, utils.ErrForbidden)
		}
	}
}

// JWTFromHeader get jwt token from header
func JWTFromHeader(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")

	if authHeader == "" {
		return "", utils.ErrEmptyAuthHeader
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", utils.ErrInvalidAuthHeader
	}

	return parts[1], nil
}
