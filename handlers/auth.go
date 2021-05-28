package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
	"github.com/vahidmostofi/minaria/domain"
)

var ErrCantParseBodyToJson = fmt.Errorf("can't parse body to json")
var ErrUsernamePasswordDontMatch = GenericError{
	Message:        domain.ErrEmailPasswordNotMatch.Error(),
	AdditionalInfo: nil,
	Err:            nil,
	HTTPStatusCode: http.StatusUnauthorized,
}

type Auth struct {
	l       *log.Logger
	usecase domain.UserUsecase
	v       *domain.Validation
}

func (a *Auth) AttachRouter(mr *mux.Router) *mux.Router {
	heathHandler := mr.PathPrefix("/auth").Subrouter()

	heathHandler.HandleFunc("/login", a.Login).Methods(http.MethodPost)
	heathHandler.HandleFunc("/register", a.Register).Methods(http.MethodPost)

	heathHandler.Use(a.postProcessMiddleware)
	return heathHandler
}

// NewAuth returns a new Auth handler
func NewAuth(l *log.Logger, usecase domain.UserUsecase, v *domain.Validation) *Auth {
	return &Auth{l: l, usecase: usecase, v: v}
}

// swagger:route POST /auth/login auth loginUser
// Returns the jwt token for the User if the email or password are correct
// responses:
//	200: jwtDTOResponse
//	400: genericErrorResponse
//  400: validationErrorResponse
//	401: usernamePasswordNotMatchResponse
// 	500: internalErrorResponse

// Login checks the health status
func (a *Auth) Login(rw http.ResponseWriter, r *http.Request) {
	a.l.Debug("Handle login request.")

	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	ld := &domain.LoginDTO{}
	gerr := a.validateDTO(ld, r.Body)

	if gerr != nil {
		rw.WriteHeader(gerr.HTTPStatusCode)
		ToJSON(gerr, rw)
		return
	}

	ds, _ := time.ParseDuration("5s") // TODO
	ctx, cancel = context.WithDeadline(r.Context(), time.Now().Add(ds))
	defer cancel()

	res, err := a.usecase.LoginByEmail(ctx, ld)
	if err == domain.ErrNoUserFound || err == domain.ErrEmailPasswordNotMatch {
		a.l.Info("Username and password don't match.")
		gerr := ErrUsernamePasswordDontMatch
		rw.WriteHeader(gerr.HTTPStatusCode)
		ToJSON(gerr, rw)
		return

	} else if err != nil {
		a.l.Errorf("Error while loging in: %s.", err.Error())
		gerr := GenericError{
			Message:        "internal server error",
			AdditionalInfo: nil,
			Err:            err,
			HTTPStatusCode: http.StatusInternalServerError,
		}
		rw.WriteHeader(gerr.HTTPStatusCode)
		ToJSON(gerr, rw)
		return
	}

	rw.WriteHeader(http.StatusOK)
	ToJSON(res, rw)
}

// swagger:route POST /auth/register auth registerUser
// Stores and registers a new user and then returns
// the jwt token for the newly created user.
// responses:
//	200: jwtDTOResponse
//	400: genericErrorResponse
//  400: validationErrorResponse
// 	500: internalErrorResponse

// Register a new user and return the jwt token
func (a *Auth) Register(rw http.ResponseWriter, r *http.Request) {
	a.l.Debug("Handle register request.")

	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	rd := &domain.RegisterDTO{}
	gerr := a.validateDTO(rd, r.Body)

	if gerr != nil {
		rw.WriteHeader(gerr.HTTPStatusCode)
		ToJSON(gerr, rw)
		return
	}

	ds, _ := time.ParseDuration("5s") // TODO
	ctx, cancel = context.WithDeadline(r.Context(), time.Now().Add(ds))
	defer cancel()

	res, err := a.usecase.Create(ctx, rd)
	if err == domain.ErrPasswordsDoNotMatch || err == domain.ErrUsernameAlreadyTaken || err == domain.ErrEmailAlreadyTaken {

		a.l.Info("Username and password don't match.")
		gerr := GenericError{
			Message:        err.Error(),
			AdditionalInfo: nil,
			Err:            err,
			HTTPStatusCode: http.StatusBadRequest,
		}
		rw.WriteHeader(gerr.HTTPStatusCode)
		ToJSON(gerr, rw)
		return
	} else if err != nil {
		a.l.Errorf("Error while registering in: %s.", err.Error())
		gerr := GenericError{
			Message:        "internal server error",
			AdditionalInfo: nil,
			Err:            err,
			HTTPStatusCode: http.StatusInternalServerError,
		}
		rw.WriteHeader(gerr.HTTPStatusCode)
		ToJSON(gerr, rw)
		return
	}

	rw.WriteHeader(http.StatusOK)
	ToJSON(res, rw)
}

func (a *Auth) validateDTO(in interface{}, r io.Reader) *GenericError {
	err := FromJSON(in, r)

	// check if the DTO can be parsed into json
	if err != nil {
		gerr := &GenericError{
			Message:        ErrCantParseBodyToJson.Error(),
			AdditionalInfo: nil,
			Err:            ErrCantParseBodyToJson,
			HTTPStatusCode: http.StatusBadRequest,
		}

		return gerr
	}

	// validate each field of the DTO
	verrs := a.v.Validate(in)
	if len(verrs) != 0 {
		gerr := &GenericError{
			Message:        "FieldError",
			AdditionalInfo: verrs.FieldsError(),
			HTTPStatusCode: http.StatusBadRequest,
		}
		return gerr
	}

	return nil
}

func (a *Auth) postProcessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
