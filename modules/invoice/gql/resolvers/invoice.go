package resolvers

import (
	"context"

	generated "api_new/modules/invoice/gql/models"
	// "api_new/logger"

	"api_new/modules/invoice/orm/models"

	"github.com/vektah/gqlparser/gqlerror"
)

type invoiceResolver struct{ *Resolver }

// <QUERY>
func (r *queryResolver) Invoices(ctx context.Context, input generated.QueryInvoice) (*generated.Invoices, error) {
	orm, err := r.GetORMFromContext(ctx)
	if err != nil {
		return nil, gqlerror.Errorf("Error: %v", err)
	}
	defer r.DisposeORM(orm)

	// Fill in
	if (input.Code == nil || len(*input.Code) == 0) && (input.Description == nil || len(*input.Description) == 0) {
		contextInvoice := models.NewContextInvoice(orm)
		allInvoices, count, err := contextInvoice.GetAllInvoices()
		if err == nil {
			return &generated.Invoices{
				Count: count,
				List:  allInvoices,
			}, nil
		}
		return nil, gqlerror.Errorf("Encounter error: %v", err)
	}

	// Search by code
	if input.Code != nil && len(*input.Code) > 0 {
		contextInvoice := models.NewContextInvoice(orm)
		result, err := contextInvoice.GetInvoiceByCode(*input.Code)
		if err == nil {
			return &generated.Invoices{
				Count: 1,
				List: []*models.Invoice{
					&result,
				},
			}, err
		}
		return nil, gqlerror.Errorf("Not found: %v", err)
	}

	// Search by description
	contextInvoice := models.NewContextInvoice(orm)
	result, err := contextInvoice.GetInvoiceByDescription(*input.Description)
	if err == nil {
		return &generated.Invoices{
			Count: 1,
			List: []*models.Invoice{
				&result,
			},
		}, err
	}
	return nil, gqlerror.Errorf("Not found: %v", err)
}

// </QUERY>

// <MUTATION>
// CreateInvoice create a invoice
func (r *mutationResolver) CreateInvoice(ctx context.Context, input generated.InvoiceInput) (*models.Invoice, error) {
	panic("Not implemented")
}

// UpdateInvoice updates a record
func (r *mutationResolver) UpdateInvoice(ctx context.Context, id string, input generated.InvoiceInput) (*models.Invoice, error) {
	// return invoiceCreateUpdate(r, input, true, id)
	panic("Not implemented")
}

// DeleteInvoice deletes a record
func (r *mutationResolver) DeleteInvoice(ctx context.Context, id string) (bool, error) {
	// return invoiceDelete(r, id)
	panic("Not implemented")
}

// </MUTATION>
