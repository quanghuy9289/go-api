package resolvers

import (
	"api_new/config"
	"api_new/logger"
	"api_new/middleware"
	"api_new/modules/invoice/gql"
	"api_new/modules/invoice/orm"
	"net/http"
	"strings"

	userORM "api_new/modules/user/orm"
	userModels "api_new/modules/user/orm/models"
	"context"
	"fmt"

	// "api_new/modules/invoice/orm/migration"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/vektah/gqlparser/gqlerror"
)

var isInit = false

// Resolver is a modifable struct that can be used to pass on properties used
// in the resolvers, such as DB access
type Resolver struct {
	MasterDBConfig *config.DatabaseConfig
}

// Mutation exposes mutation methods
func (r *Resolver) Mutation() gql.MutationResolver {
	return &mutationResolver{r}
}

// Query exposes query methods
func (r *Resolver) Query() gql.QueryResolver {
	return &queryResolver{r}
}

// GetORMFromContext get orm from request context
func (r *Resolver) GetORMFromContext(ctx context.Context) (o *orm.ORM, err error) {
	ginContext, err := middleware.GinContextFromContext(ctx)

	if err != nil {
		return nil, gqlerror.Errorf("Wrong context")
	}
	authToken, err := middleware.JWTFromHeader(ginContext)
	if err != nil {
		ginContext.AbortWithError(http.StatusUnauthorized, gqlerror.Errorf("Unauthorized"))
		return nil, gqlerror.Errorf("Unauthorized")
	}

	o, err = r.GetORMFromToken(ginContext, authToken)
	return
}

// GetORMFromToken get orm from authentication token
func (r *Resolver) GetORMFromToken(ctx *gin.Context, token string) (o *orm.ORM, err error) {
	// get database infor from token
	context := userModels.NewContextUser(userORM.SharedORM)
	user, err := context.GetUserByAuthenticationToken(token)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, gqlerror.Errorf("Unauthorized"))
		return nil, gqlerror.Errorf("Unauthorized")
	}

	databaseName := user.InUseDatabase
	// host := "127.0.0.1"
	// port := "5432"

	db, err := gorm.Open(
		r.MasterDBConfig.Dialect,
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			r.MasterDBConfig.Host,
			r.MasterDBConfig.Port,
			r.MasterDBConfig.Username,
			strings.ToLower(databaseName),
			r.MasterDBConfig.Password))
	if err != nil {
		logger.Errorf("[ORM] err: %v", err)
		return nil, gqlerror.Errorf("Database " + databaseName + " does not exist")
	}

	o = &orm.ORM{
		DB: db,
	}
	return
}

// DisposeORM dispose orm database connection
func (r *Resolver) DisposeORM(o *orm.ORM) {
	o.GetDB().Close()
}

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }
