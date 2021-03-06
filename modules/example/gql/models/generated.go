// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gql

import (
	"api_new/modules/example/orm/models"
)

type ExampleInput struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type Examples struct {
	Count int               `json:"count"`
	List  []*models.Example `json:"list"`
}

type QueryExample struct {
	Code        *string `json:"code"`
	Description *string `json:"description"`
}
