package usecase

import (
	"context"
	"crypto"
	"crypto/subtle"
	"encoding/hex"

	"github.com/go-openapi/strfmt"

	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vahidmostofi/minaria/common"
	"github.com/vahidmostofi/minaria/domain"
	"github.com/vahidmostofi/minaria/repositories"
)

type UserOptions struct {
	// default is 31 days
	jwtExpiresAfter *time.Duration

	// default is SHA256
	hashMethod *crypto.Hash
}

type User struct {
	l               *log.Logger
	r               domain.UserRepository
	jwtExpiresAfter time.Duration
	hashMethod      crypto.Hash
}

func NewUser(l *log.Logger, r domain.UserRepository, opts UserOptions) domain.UserUsecase {
	u := &User{}
	u.l = l
	u.r = r

	if opts.jwtExpiresAfter != nil {
		u.jwtExpiresAfter = *opts.jwtExpiresAfter
	} else {
		d, _ := time.ParseDuration("744h")
		u.jwtExpiresAfter = time.Duration(d)
	}
	if opts.hashMethod != nil {
		u.hashMethod = *opts.hashMethod
	} else {
		u.hashMethod = crypto.SHA256
	}

	return u
}

func (uc *User) LoginByEmail(ctx context.Context, ld *domain.LoginDTO) (*domain.JWTDTO, error) {
	user, err := uc.r.GetByEmail(ctx, ld.Email.String())

	if err != nil {
		if err == repositories.ErrNoUserFound {
			return nil, domain.ErrNoUserFound
		}
	}

	hashedPassword := uc.hash([]byte(ld.Password))

	currentedHashedPassword, err := hex.DecodeString(user.Password)
	if err != nil {
		return nil, fmt.Errorf("error while decoding hex string: %w", err)
	}

	if subtle.ConstantTimeCompare(hashedPassword, currentedHashedPassword) == 1 {
		token, err := uc.generateJWT(user.ID, user.Username)
		if err != nil {
			return nil, fmt.Errorf("error while generating jwt token: %w", err)
		}
		return &domain.JWTDTO{Token: token}, nil
	}

	return nil, domain.ErrEmailPasswordNotMatch
}

func (uc *User) CheckEmailAvailable(ctx context.Context, email string) error {
	_, err := uc.r.GetByEmail(ctx, email)

	if err != nil {
		if err == repositories.ErrNoUserFound {
			return nil
		} else {
			return err
		}
	}

	return domain.ErrEmailAlreadyTaken
}

func (uc *User) CheckUsernameAvailable(ctx context.Context, username string) error {
	_, err := uc.r.GetByUsername(ctx, username)

	if err != nil {
		if err == repositories.ErrNoUserFound {
			return nil
		} else {
			return err
		}
	}

	return domain.ErrUsernameAlreadyTaken
}

func (uc *User) Create(ctx context.Context, r *domain.RegisterDTO) (*domain.JWTDTO, error) {
	rawPassword := r.Password

	if r.Password.String() != r.RepeatPassword.String() {
		return nil, domain.ErrPasswordsDoNotMatch
	}

	if uc.CheckEmailAvailable(ctx, r.Email.String()) != nil {
		return nil, domain.ErrEmailAlreadyTaken
	}

	if uc.CheckUsernameAvailable(ctx, r.Username) != nil {
		return nil, domain.ErrUsernameAlreadyTaken
	}

	u := domain.User{Username: r.Username, Password: hex.EncodeToString(uc.hash([]byte(r.Password))), Email: r.Email.String()}

	usr, err := uc.r.Store(ctx, &u)
	if err != nil {
		return nil, err
	}

	return uc.LoginByEmail(ctx, &domain.LoginDTO{strfmt.Email(usr.Email), strfmt.Password(rawPassword)})
}

func (uc *User) hash(in []byte) []byte {
	h := uc.hashMethod.New()
	h.Write(in)
	return h.Sum(nil)
}

func (uc *User) generateJWT(ID, Username string) (string, error) {
	signKey := []byte(viper.GetString(common.JWT_SIGN_KEY))
	claims := &jwt.StandardClaims{Id: ID, ExpiresAt: time.Now().Add(uc.jwtExpiresAfter).Unix()}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signKey)
	if err != nil {
		return "", fmt.Errorf("error while signing the token: %w", err)
	}

	return ss, nil
}
