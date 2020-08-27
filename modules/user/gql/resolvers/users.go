package resolvers

import (
	"context"

	generated "api_new/modules/user/gql/models"
	// "api_new/logger"
	"api_new/middleware"
	"api_new/modules/user/orm/models"
	"api_new/utils"
	"fmt"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/vektah/gqlparser/gqlerror"
)

type userResolver struct{ *Resolver }

// <QUERY>
func (r *userResolver) AuthenticationTokens(ctx context.Context, obj *models.User) ([]*models.AuthenticationToken, error) {
	// Fill in
	contextUser := models.NewContextUser(r.ORM)
	return contextUser.GetUserActiveAuthenticationTokens(obj)
}

func (r *queryResolver) Users(ctx context.Context, input generated.QueryUser) (*generated.Users, error) {
	// Fill in
	if len(input.Email) == 0 {
		contextUser := models.NewContextUser(r.ORM)
		allUsers, count, err := contextUser.GetAllUsers()
		if err == nil {
			return &generated.Users{
				Count: count,
				List:  allUsers,
			}, nil
		}
		return nil, gqlerror.Errorf("Encounter error: %v", err)
	}
	contextUser := models.NewContextUser(r.ORM)
	userByEmail, err := contextUser.GetUserByEmail(input.Email)
	if err == nil {
		return &generated.Users{
			Count: 1,
			List: []*models.User{
				&userByEmail,
			},
		}, err
	}
	return nil, gqlerror.Errorf("Not found: %v", err)
}

// </QUERY>

// <MUTATION>
// CreateUser create a user
func (r *mutationResolver) CreateUser(ctx context.Context, input generated.UserInput) (*models.User, error) {
	password, errHashPassword := utils.HashAndSalt([]byte(input.Password))
	if errHashPassword == nil {
		ctxDB := models.NewContextDatabase(r.ORM)
		_, err := ctxDB.GetDatabaseByName(input.InUseDatabase)
		if err == nil {
			contextUser := models.NewContextUser(r.ORM)
			newUserID, errAddUser := contextUser.AddUser(&models.User{
				ID:            utils.GenerateUUID(),
				CreatedOn:     time.Now().Unix(),
				Email:         input.Email,
				Fullname:      input.Fullname,
				Nickname:      input.Nickname,
				IsActive:      true,
				Password:      password,
				AvatarBase64:  input.AvatarBase64,
				RoleID:        input.RoleID,
				InUseDatabase: input.InUseDatabase,
			})
			if errAddUser == nil {
				theUser, errGetUser := contextUser.GetUserByID(newUserID)
				if errGetUser == nil {
					return &theUser, nil
				}
				return nil, gqlerror.Errorf("Could not query the user information: %v", errGetUser)
			}
			return nil, gqlerror.Errorf("Could not create the user: %v", errAddUser)
		}
		return nil, gqlerror.Errorf("Could not create new user: In use database not exist")
	}
	return nil, gqlerror.Errorf("Could not generate hashed password: %v", errHashPassword)
}

// UpdateUser updates a record
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input generated.UserInfo) (*models.User, error) {
	// return userCreateUpdate(r, input, true, id)
	contextUser := models.NewContextUser(r.ORM)
	theUser, errGetUser := contextUser.GetUserByID(id)
	if errGetUser == nil {
		updatedUser := &models.User{
			ID:           id,
			CreatedOn:    theUser.CreatedOn,
			Email:        input.Email,
			Fullname:     input.Fullname,
			Nickname:     input.Nickname,
			IsActive:     theUser.IsActive,
			Password:     theUser.Password,
			AvatarBase64: input.AvatarBase64,
			RoleID:       theUser.RoleID,
		}
		errUpdate := contextUser.UpdateUser(updatedUser)

		if errUpdate == nil {
			return updatedUser, nil
		}
		return nil, gqlerror.Errorf("Could not update the user information: %v", errUpdate)

	}
	return nil, gqlerror.Errorf("Could not update the user information: %v", errGetUser)
}

// DeleteUser deletes a record
func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*models.User, error) {
	// return userDelete(r, id)
	// panic("Not implemented")
	contextUser := models.NewContextUser(r.ORM)
	theUser, errGetUser := contextUser.GetUserByID(id)
	if errGetUser == nil {
		deleteUser := &models.User{
			ID: id,
		}
		errDelete := contextUser.DeleteUser(deleteUser)

		if errDelete == nil {
			return &theUser, nil
		}
		return nil, gqlerror.Errorf("Could not delete the user information: %v", errDelete)

	}
	return nil, gqlerror.Errorf("Could not delete the user: %v", errGetUser)
}

// Logout user
func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	ginContext, err := middleware.GinContextFromContext(ctx)

	if err != nil {
		return false, gqlerror.Errorf("Wrong context")
	}

	session := sessions.Default(ginContext)
	tokenSession := session.Get("token")

	// expired token in db
	contextAuth := models.NewContextAuthenticationToken(r.ORM)
	contextAuth.SetExpiredAuthenticationToken(fmt.Sprintf("%v", tokenSession))

	// clear cache
	session.Clear()
	session.Save()

	return true, nil
}

// Logout all services for user
func (r *mutationResolver) LogoutAll(ctx context.Context) (bool, error) {
	ginContext, err := middleware.GinContextFromContext(ctx)

	if err != nil {
		return false, gqlerror.Errorf("Wrong context")
	}

	session := sessions.Default(ginContext)
	email := session.Get("email")

	// expired token in db
	contextAuth := models.NewContextAuthenticationToken(r.ORM)
	// contextAuth.SetExpiredAuthenticationToken(fmt.Sprintf("%v", tokenSession))
	contextAuth.SetExpiredAllToken(fmt.Sprintf("%v", email))

	// clear cache
	session.Clear()
	session.Save()

	return true, nil
}

// </MUTATION>
