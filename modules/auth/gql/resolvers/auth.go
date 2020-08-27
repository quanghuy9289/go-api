package resolvers

import (
	"api_new/config"
	"api_new/middleware"
	generated "api_new/modules/auth/gql/models"
	"api_new/modules/user/orm"
	"api_new/modules/user/orm/models"
	"api_new/services"
	"api_new/utils"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/vektah/gqlparser/gqlerror"
)

type authResolver struct{ *Resolver }

// <QUERY>

func (r *queryResolver) SampleAuth(ctx context.Context, input string) (bool, error) {
	return true, nil
}

// </QUERY>

// <MUTATION>

func (r *mutationResolver) Register(ctx context.Context, input generated.RegisterInput) (*generated.RegisterResult, error) {
	password, errHashPassword := utils.HashAndSalt([]byte(input.Password))
	if errHashPassword == nil {
		// 1. Create a new registration
		registration, err := CreateRegistration(input)

		// 2. Create database for this new registration in step 1
		if err == nil {
			database, err := CreateDatabase(registration.RegistrationID, input.Email, r.DBConfig)

			if err == nil {
				// 3. Create user that own the created database in step 2
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
					RoleID:        "user",
					InUseDatabase: database.DatabaseName,
				})
				if errAddUser == nil {
					theUser, errGetUser := contextUser.GetUserByID(newUserID)
					if errGetUser == nil {
						res := &generated.RegisterResult{
							ID:            theUser.ID,
							Email:         theUser.Email,
							Fullname:      theUser.Fullname,
							Nickname:      theUser.Nickname,
							AvatarBase64:  theUser.AvatarBase64,
							CompanyName:   registration.CompanyName,
							PhoneNo:       registration.PhoneNo,
							FaxNo:         registration.FaxNo,
							Website:       registration.Website,
							StreetAddress: registration.StreetAddress,
							InUseDatabase: database.DatabaseName,
						}
						return res, nil
					}
					return nil, gqlerror.Errorf("Could not query the user information: %v", errGetUser)
				}
				return nil, gqlerror.Errorf("Could not create the user: %v", errAddUser)
			}
			return nil, err
		}
		return nil, err
	}
	return nil, gqlerror.Errorf("Could not generate hashed password: %v", errHashPassword)
}

func (r *mutationResolver) Login(ctx context.Context, input generated.LoginInput) (*generated.LoginResult, error) {
	if len(input.Email) == 0 || len(input.Password) == 0 {
		return nil, gqlerror.Errorf("Authentication failed")
	}

	ginContext, err := middleware.GinContextFromContext(ctx)

	contextUser := models.NewContextUser(r.ORM)
	user, err := contextUser.GetUserByEmail(input.Email)
	if err != nil {
		return nil, gqlerror.Errorf("Invalid email or password")
	}

	match := utils.ComparePasswords(user.Password, []byte(input.Password))
	if match {
		// generate token
		token, err := utils.GenerateJwtToken(input.Email)
		if err == nil {
			contextAuth := models.NewContextAuthenticationToken(r.ORM)

			// get device id from request
			var deviceID = ginContext.Request.RemoteAddr
			currentAuthToken, err := contextAuth.GetAuthenticationTokenByDeviceId(deviceID)

			if err == nil {
				// set expired for current token
				contextAuth.SetExpiredAuthenticationToken(currentAuthToken.Token)
			}

			authToken := &models.AuthenticationToken{
				Token:     token,
				CreatedOn: time.Now().UTC().Unix(),
				ExpiredOn: time.Now().Add(time.Hour * time.Duration(24)).Unix(),
				UserID:    user.ID,
				DeviceID:  deviceID,
				Email:     user.Email,
			}
			contextAuth.AddAuthenticationToken(authToken)

			// add to session
			session := sessions.Default(ginContext)
			session.Set("userId", user.ID)
			session.Set("email", user.Email)
			session.Set("token", token)
			session.Save()

			serializer := UserSerializer{}

			return &generated.LoginResult{
				Authtoken: token,
				User:      serializer.Response(user),
			}, nil
		}
	}

	return nil, gqlerror.Errorf("Invalid email or password")
}

// </MUTATION>

// CreateRegistration create a new registration for free user
func CreateRegistration(input generated.RegisterInput) (*models.Registration, error) {
	context := models.NewContextRegistration(orm.SharedORM)
	newRegistrationID, err := context.CreateRegistration(&models.Registration{
		CompanyName:   input.CompanyName,
		PhoneNo:       input.PhoneNo,
		FaxNo:         input.FaxNo,
		Website:       input.Website,
		StreetAddress: input.StreetAddress,
	})
	if err == nil {
		registration, errGetRegistration := context.GetRegistrationByID(newRegistrationID)
		if errGetRegistration == nil {
			return &registration, nil
		}
		return nil, gqlerror.Errorf("Could not query the registration information: %v", errGetRegistration)
	}
	return nil, gqlerror.Errorf("Could not create the registration: %v", err)
}

// CreateDatabase for free user
func CreateDatabase(registrationID int64, email string, dbConfig *config.DatabaseConfig) (*models.Database, error) {
	context := models.NewContextDatabase(orm.SharedORM)
	db, _ := context.GetDatabaseByRegistrationID(int64(registrationID))
	if db.DatabaseID == 0 { // db not exist, so create new database
		newDatabaseID, err := context.CreateDatabase(&models.Database{
			ExpiryDate:     time.Now().Add(time.Hour * time.Duration(720)).Unix(), // expired time will be configured later
			DatabaseName:   "CTB_User_" + strconv.Itoa(int(registrationID)) + "_" + utils.FreeUser,
			RegistrationID: int64(registrationID),
			IsActive:       true,
		})
		if err == nil {
			database, errGetDatabase := context.GetDatabaseByID(newDatabaseID)
			if errGetDatabase == nil {
				// persist new database
				sqlQuery := fmt.Sprintf("CREATE DATABASE %s", database.DatabaseName)
				newDB := orm.SharedORM.GetDB().Exec(sqlQuery)
				if newDB.Error != nil {
					return nil, gqlerror.Errorf("Unable to create DB: %v", newDB.Error)
				}

				// run migration for new database
				err = services.MigrationClientDabatase(database.DatabaseName, dbConfig)
				if err != nil {
					return nil, gqlerror.Errorf("Migration database error: %v", err)
				}

				return &database, nil
			}
			return nil, gqlerror.Errorf("Could not query the Database information: %v", errGetDatabase)
		}
		return nil, gqlerror.Errorf("Could not create the Database: %v", err)
	}
	return nil, gqlerror.Errorf("Database for registration %v already exist", registrationID)
}
