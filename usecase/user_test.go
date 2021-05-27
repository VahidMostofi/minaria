package usecase

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/strfmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vahidmostofi/minaria/domain"
	"github.com/vahidmostofi/minaria/repositories"
)

func getUserRepository(t *testing.T) domain.UserRepository {
	ur, err := repositories.NewUserRepository(
		repositories.InMemoryKind,
		repositories.InMemoryArgs{Data: []*domain.User{{Email: "jack@gmail.com"}}},
	)

	if err != nil {
		t.Fatalf("error while creating user repository: %s", err.Error())
	}

	return ur
}

func TestCheckEmailAvailability(t *testing.T) {
	l := logrus.New()
	l.SetOutput(ioutil.Discard)

	ur := getUserRepository(t)

	uc := NewUser(l, ur, UserOptions{})
	err := uc.CheckEmailAvailable(context.TODO(), "jack@gmail.com")
	assert.Equal(t, err, domain.ErrEmailAlreadyTaken)

	err = uc.CheckEmailAvailable(context.TODO(), "gholi@gmail.com")
	assert.Nil(t, err)
}

func TestCheckUsernameAvailability(t *testing.T) {
	l := logrus.New()
	l.SetOutput(ioutil.Discard)

	ur := getUserRepository(t)

	uc := NewUser(l, ur, UserOptions{})
	err := uc.CheckUsernameAvailable(context.TODO(), "jack")
	assert.Equal(t, err, domain.ErrUsernameAlreadyTaken)

	err = uc.CheckUsernameAvailable(context.TODO(), "gholi")
	assert.Nil(t, err)
}

func TestCreateUser(t *testing.T) {
	l := logrus.New()
	l.SetOutput(ioutil.Discard)

	ur := getUserRepository(t)

	uc := NewUser(l, ur, UserOptions{})
	rd := &domain.RegisterDTO{"gholi", strfmt.Email("gholi@gmail.com"), strfmt.Password("password"), strfmt.Password("password")}

	jwtDTO, err := uc.Create(context.TODO(), rd)
	if err != nil {
		t.Fatal(err)
	}

	claims := jwt.MapClaims{}

	jwt.ParseWithClaims(jwtDTO.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("<YOUR VERIFICATION KEY>"), nil
	})

	u, err := ur.GetByEmail(context.TODO(), rd.Email.String())
	if err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, claims.Valid())
	assert.True(t, claims.VerifyAudience(u.ID, false))

	_, err = uc.Create(context.TODO(), rd)
	assert.Equal(t, err, domain.ErrEmailAlreadyTaken)

	rd.Email = "another@gmail.com"
	_, err = uc.Create(context.TODO(), rd)
	assert.Equal(t, err, domain.ErrUsernameAlreadyTaken)

	rd.RepeatPassword = strfmt.Password("another password")
	_, err = uc.Create(context.TODO(), rd)
	assert.Equal(t, err, domain.ErrPasswordsDoNotMatch)
}
