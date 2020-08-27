package resolvers

import (
	"context"

	generated "api_new/modules/example/gql/models"
	// "api_new/logger"
	"api_new/modules/example/orm/models"
	"api_new/utils"
	"github.com/vektah/gqlparser/gqlerror"
	"time"
)

type exampleResolver struct{ *Resolver }

// <QUERY>
func (r *queryResolver) Examples(ctx context.Context, input generated.QueryExample) (*generated.Examples, error) {
	// Fill in
	if (input.Code == nil || len(*input.Code) == 0) && (input.Description == nil || len(*input.Description) == 0) {
		contextExample := models.NewContextExample(r.ORM)
		allExamples, count, err := contextExample.GetAllExamples()
		if err == nil {
			return &generated.Examples{
				Count: count,
				List:  allExamples,
			}, nil
		}
		return nil, gqlerror.Errorf("Encounter error: %v", err)
	}

	// Search by code
	if input.Code != nil && len(*input.Code) > 0 {
		contextExample := models.NewContextExample(r.ORM)
		result, err := contextExample.GetExampleByCode(*input.Code)
		if err == nil {
			return &generated.Examples{
				Count: 1,
				List: []*models.Example{
					&result,
				},
			}, err
		}
		return nil, gqlerror.Errorf("Not found: %v", err)
	}

	// Search by description
	contextExample := models.NewContextExample(r.ORM)
	result, err := contextExample.GetExampleByDescription(*input.Description)
	if err == nil {
		return &generated.Examples{
			Count: 1,
			List: []*models.Example{
				&result,
			},
		}, err
	}
	return nil, gqlerror.Errorf("Not found: %v", err)
}

// </QUERY>

// <MUTATION>
// CreateExample create a example
func (r *mutationResolver) CreateExample(ctx context.Context, input generated.ExampleInput) (*models.Example, error) {
	// Fill in
	contextExample := models.NewContextExample(r.ORM)
	newExampleID, errAddExample := contextExample.AddExample(&models.Example{
		ID:          utils.GenerateUUID(),
		CreatedOn:   time.Now().Unix(),
		Code:        input.Code,
		Description: input.Description,
	})
	if errAddExample == nil {
		theExample, errGetExample := contextExample.GetExampleByID(newExampleID)
		if errGetExample == nil {
			return &theExample, nil
		}
		return nil, gqlerror.Errorf("Could not query the example information: %v", errGetExample)
	}
	return nil, gqlerror.Errorf("Could not create the example: %v", errAddExample)
}

// UpdateExample updates a record
func (r *mutationResolver) UpdateExample(ctx context.Context, id string, input generated.ExampleInput) (*models.Example, error) {
	// return exampleCreateUpdate(r, input, true, id)
	panic("Not implemented")
}

// DeleteExample deletes a record
func (r *mutationResolver) DeleteExample(ctx context.Context, id string) (bool, error) {
	// return exampleDelete(r, id)
	panic("Not implemented")
}

// </MUTATION>
