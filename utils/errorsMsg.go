package utils

import (
	"github.com/vektah/gqlparser/gqlerror"
)

var (
	// ErrEmptyAuthHeader can be thrown if authing with a HTTP header, the Auth header needs to be set
	ErrEmptyAuthHeader = gqlerror.Errorf("Auth header is empty")

	// ErrForbidden when HTTP status 403 is given
	ErrForbidden = gqlerror.Errorf("You don't have permission to access this resource")

	// ErrInvalidAuthHeader indicates auth header is invalid, could for example have the wrong Realm name
	ErrInvalidAuthHeader = gqlerror.Errorf("Auth header is invalid")

	// ErrInitMasterORM error when master orm has not set yet
	ErrInitMasterORM = gqlerror.Errorf("Master ORM hasn't initiated yet.")
)
