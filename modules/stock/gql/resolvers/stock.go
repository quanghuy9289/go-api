package resolvers

import (
	"context"

	generated "api_new/modules/stock/gql/models"
	// "api_new/logger"
	"api_new/modules/stock/orm/models"
	"api_new/utils"
	"github.com/vektah/gqlparser/gqlerror"
	"time"
)

type stockResolver struct{ *Resolver }

// <QUERY>
func (r *queryResolver) Stocks(ctx context.Context, input generated.QueryStock) (*generated.Stocks, error) {
	// Fill in
	if (input.Code == nil || len(*input.Code) == 0) && (input.Description == nil || len(*input.Description) == 0) {
		contextStock := models.NewContextStock(r.ORM)
		allStocks, count, err := contextStock.GetAllStocks()
		if err == nil {
			return &generated.Stocks{
				Count: count,
				List:  allStocks,
			}, nil
		}
		return nil, gqlerror.Errorf("Encounter error: %v", err)
	}

	// Search by code
	if input.Code != nil && len(*input.Code) > 0 {
		contextStock := models.NewContextStock(r.ORM)
		result, err := contextStock.GetStockByCode(*input.Code)
		if err == nil {
			return &generated.Stocks{
				Count: 1,
				List: []*models.Stock{
					&result,
				},
			}, err
		}
		return nil, gqlerror.Errorf("Not found: %v", err)
	}

	// Search by description
	contextStock := models.NewContextStock(r.ORM)
	result, err := contextStock.GetStockByDescription(*input.Description)
	if err == nil {
		return &generated.Stocks{
			Count: 1,
			List: []*models.Stock{
				&result,
			},
		}, err
	}
	return nil, gqlerror.Errorf("Not found: %v", err)
}

// </QUERY>

// <MUTATION>
// CreateStock create a stock
func (r *mutationResolver) CreateStock(ctx context.Context, input generated.StockInput) (*models.Stock, error) {
	// Fill in
	contextStock := models.NewContextStock(r.ORM)
	newStockID, errAddStock := contextStock.AddStock(&models.Stock{
		ID:          utils.GenerateUUID(),
		CreatedOn:   time.Now().Unix(),
		Code:        input.Code,
		Description: input.Description,
	})
	if errAddStock == nil {
		theStock, errGetStock := contextStock.GetStockByID(newStockID)
		if errGetStock == nil {
			return &theStock, nil
		}
		return nil, gqlerror.Errorf("Could not query the stock information: %v", errGetStock)
	}
	return nil, gqlerror.Errorf("Could not create the stock: %v", errAddStock)
}

// UpdateStock updates a record
func (r *mutationResolver) UpdateStock(ctx context.Context, id string, input generated.StockInput) (*models.Stock, error) {
	// return stockCreateUpdate(r, input, true, id)
	panic("Not implemented")
}

// DeleteStock deletes a record
func (r *mutationResolver) DeleteStock(ctx context.Context, id string) (bool, error) {
	// return stockDelete(r, id)
	panic("Not implemented")
}

// </MUTATION>
