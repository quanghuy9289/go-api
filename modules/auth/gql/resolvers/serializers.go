package resolvers

import (
	generated "api_new/modules/auth/gql/models"
	"api_new/modules/user/orm/models"
)

type UserSerializer struct {
}

func (self *UserSerializer) Response(user models.User) *generated.UserOutput {
	userModel := &generated.UserOutput{
		ID:                                 user.ID,
		Email:                              user.Email,
		PhoneNumber:                        user.PhoneNumber,
		Fullname:                           user.Fullname,
		Nickname:                           user.Nickname,
		AvatarBase64:                       user.AvatarBase64,
		RoleID:                             user.RoleID,
		StartDate:                          user.StartDate,
		MinDailyUnitPointsRequirement:      user.MinDailyUnitPointsRequirement,
		StandardNumberOfWorkingDaysPerWeek: user.StandardNumberOfWorkingDaysPerWeek,
		Config:                             user.Config,
	}
	return userModel
}
