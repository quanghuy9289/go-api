package resolvers

import (
	"api_new/config"
	"api_new/modules/user/gql"
	"api_new/modules/user/orm"
)

// Resolver is a modifable struct that can be used to pass on properties used
// in the resolvers, such as DB access
type Resolver struct {
	ORM      *orm.ORM
	DBConfig *config.DatabaseConfig
}

// Mutation exposes mutation methods
func (r *Resolver) Mutation() gql.MutationResolver {
	return &mutationResolver{r}
}

// Query exposes query methods
func (r *Resolver) Query() gql.QueryResolver {
	return &queryResolver{r}
}

// User resolve for User
func (r *Resolver) User() gql.UserResolver {
	return &userResolver{r}
}

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }
