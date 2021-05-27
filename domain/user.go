package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"
)

var ErrNoUserFound = fmt.Errorf("no user found")
var ErrEmailPasswordNotMatch = fmt.Errorf("email and the password don't match")
var ErrEmailAlreadyTaken = fmt.Errorf("email is already taken")
var ErrUsernameAlreadyTaken = fmt.Errorf("username is already taken")
var ErrPasswordsDoNotMatch = fmt.Errorf("passwords don't match")

// User ...
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginDTO struct {
	// the email address for this user
	//
	// required: true
	// example: john@provider.net
	Email strfmt.Email `json:"email" validate:"required,email"`

	// the password for this user
	//
	// required: true
	Password strfmt.Password `json:"password" validate:"required"`
}

type RegisterDTO struct {
	// the username fo the new user
	//
	// required: true
	// example: john
	Username string `json:"username" validate:"required,min=5"`

	// the email fo the new user
	//
	// required: true
	// example: john@provider.net
	Email strfmt.Email `json:"email" validate:"email,required"`

	// the password for the new user
	//
	// required: true
	// example: $tR0n@p@$SW0rD
	Password strfmt.Password `json:"password" validate:"required,min=5"`

	// the repeat of the password field
	//
	// required: true
	// example: $tR0n@p@$SW0rD
	RepeatPassword strfmt.Password `json:"repeatPassword" validate:"required,min=5"`
}

type JWTDTO struct {
	// the jwt token for the logged in user
	Token string `json:"token"`
}

// UserUsecase interface represents the user's usecases
type UserUsecase interface {

	// LoginByEmail logs a user in with the email and password and returns a valid jwt for the newly registered user
	LoginByEmail(ctx context.Context, ld *LoginDTO) (*JWTDTO, error)

	// Create Registers the user and return a valid jwt for the newly registered user
	Create(ctx context.Context, u *RegisterDTO) (*JWTDTO, error)

	// CheckEmailAvailable returns EmailAlreadyTakenErr error if the email is not available
	CheckEmailAvailable(ctx context.Context, email string) error

	// CheckEmailAvailable returns EmailAlreadyTakenErr error if the email is not available
	CheckUsernameAvailable(ctx context.Context, username string) error
}

// UserRepository represents the user's repository contract
type UserRepository interface {
	// GetByUsername ...
	GetByUsername(ctx context.Context, username string) (*User, error)

	// GetByEmail ...
	GetByEmail(ctx context.Context, email string) (*User, error)

	// Store ...
	Store(ctx context.Context, u *User) (*User, error)
}
