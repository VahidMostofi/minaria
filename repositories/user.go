package repositories

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/vahidmostofi/minaria/domain"
)

// ErrUnknownRepository ...
var ErrUnknownRepository = fmt.Errorf("no user repository found with provided kind")

// ErrNoUserFound ...
var ErrNoUserFound = fmt.Errorf("no user found")

// ErrUsernameNotUnique ...
var ErrUsernameNotUnique = fmt.Errorf("username is not unique, it already exists")

// ErrEmailNotUnique ...
var ErrEmailNotUnique = fmt.Errorf("email is not unique, it already exists")

const InMemoryKind string = "InMemory"

type InMemoryArgs struct {
	Data []*domain.User
}

func NewUserRepository(kind string, args interface{}) (domain.UserRepository, error) {

	switch kind {
	case InMemoryKind:
		if ima, ok := args.(*InMemoryArgs); !ok {
			return newInMemoryUserRepository(ima), nil
		} else {
			return newInMemoryUserRepository(nil), nil
		}

	}

	return nil, errors.Wrap(ErrUnknownRepository, fmt.Sprintf("kind: %s", kind))
}
