package resolvers

import (
	"context"

	generated "api_new/modules/user/gql/models"
	// "api_new/logger"
	"api_new/modules/user/orm/models"
	"github.com/vektah/gqlparser/gqlerror"
)

// </QUERY>

// <MUTATION>
// CreateUser create a user
func (r *mutationResolver) CreateNewRegistration(ctx context.Context, input generated.RegistrationInput) (*models.Registration, error) {
	// Fill in
	context := models.NewContextRegistration(r.ORM)
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

// </MUTATION>
