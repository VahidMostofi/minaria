package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/strfmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vahidmostofi/minaria/domain"
	"github.com/vahidmostofi/minaria/repositories"
	"github.com/vahidmostofi/minaria/usecase"
)

const desiredContentType = "application/json"

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestLoginSuccessful(t *testing.T) {
	router := getNewRouter()
	testUserDbIdx := 0
	loginDTOSuccessful := &domain.LoginDTO{Email: strfmt.Email(testUserData[testUserDbIdx].Email), Password: "1234567"}
	loginDTOBytesSuccessful, _ := json.Marshal(loginDTOSuccessful)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(loginDTOBytesSuccessful))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()
	jwtDTO := &domain.JWTDTO{}

	if !basicHTTPResponseChecks(t, http.StatusOK, desiredContentType, jwtDTO, resp) {
		return
	}

	assert.NotNil(t, jwtDTO)
	assert.Greater(t, len(jwtDTO.Token), 0)
	claims := jwt.MapClaims{}

	jwt.ParseWithClaims(jwtDTO.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("<YOUR VERIFICATION KEY>"), nil
	})

	assert.Nil(t, claims.Valid())
	assert.True(t, claims.VerifyAudience(testUserData[testUserDbIdx].ID, false))
}

func TestLoginError(t *testing.T) {
	router := getNewRouter()
	testUserDbIdx := 0

	loginDTOInvalidEmail := &domain.LoginDTO{Email: strfmt.Email(testUserData[testUserDbIdx].Username), Password: "1234567"}
	loginDTONoEmail := &domain.LoginDTO{Email: "", Password: "1234567"}
	loginDTOEmailNotFound := &domain.LoginDTO{Email: "vahid@gmail.com", Password: "1234567"}
	loginDTOWrongPassword := &domain.LoginDTO{Email: strfmt.Email(testUserData[testUserDbIdx].Email), Password: "123456"}

	tests := []struct {
		name       string
		dto        *domain.LoginDTO
		errMessage string
		statusCode int
	}{
		{
			name:       "bad request - invalid email",
			dto:        loginDTOInvalidEmail,
			statusCode: http.StatusBadRequest,
			errMessage: "FieldError",
		},
		{
			name:       "bad request - no email",
			dto:        loginDTONoEmail,
			statusCode: http.StatusBadRequest,
			errMessage: "FieldError",
		},
		{
			name:       "unauthorized - email not found",
			dto:        loginDTOEmailNotFound,
			statusCode: http.StatusUnauthorized,
			errMessage: "email and the password don't match",
		},
		{
			name:       "unauthorized - password is wrong",
			dto:        loginDTOWrongPassword,
			statusCode: http.StatusUnauthorized,
			errMessage: "email and the password don't match",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := json.Marshal(tt.dto)
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(b))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			resp := w.Result()
			gerr := &GenericError{}

			if !basicHTTPResponseChecks(t, tt.statusCode, desiredContentType, gerr, resp) {
				return
			}
			assert.Equal(t, tt.errMessage, gerr.Message)
		})
	}

}

func TestRegisterSuccessful(t *testing.T) {
	router := getNewRouter()
	registerDTO := &domain.RegisterDTO{Email: "gholi@gmail.com", Username: "gholi", Password: "1234567", RepeatPassword: "1234567"}
	registerDTOBytes, _ := json.Marshal(registerDTO)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(registerDTOBytes))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()
	jwtDTO := &domain.JWTDTO{}

	if !basicHTTPResponseChecks(t, http.StatusOK, desiredContentType, jwtDTO, resp) {
		return
	}

	assert.NotNil(t, jwtDTO)
	assert.Greater(t, len(jwtDTO.Token), 0)
	claims := jwt.MapClaims{}

	jwt.ParseWithClaims(jwtDTO.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("<YOUR VERIFICATION KEY>"), nil
	})

	assert.Nil(t, claims.Valid())
}

func TestRegisterError(t *testing.T) {
	router := getNewRouter()
	registerDTOInvalidEmail := &domain.RegisterDTO{Email: strfmt.Email("gholigmail.com"), Password: "1234567", RepeatPassword: "1234567"}
	registerDTONoUsernameProvided := &domain.RegisterDTO{Email: "vahid@gmail.com", Password: "1234567", RepeatPassword: "123456"}
	registerDTOPasswordsDontMatch := &domain.RegisterDTO{Email: "vahid@gmail.com", Username: "vahid", Password: "1234567", RepeatPassword: "123456"}
	registerDTOEmailAlreadyTaken := &domain.RegisterDTO{Email: "jack@gmail.com", Username: "vahid", Password: "1234567", RepeatPassword: "1234567"}

	tests := []struct {
		name       string
		dto        *domain.RegisterDTO
		more       interface{}
		errMessage string
		statusCode int
	}{
		{
			name:       "bad request - invalid email",
			dto:        registerDTOInvalidEmail,
			statusCode: http.StatusBadRequest,
			errMessage: "FieldError",
		},
		{
			name:       "bad request - invalid email",
			dto:        registerDTONoUsernameProvided,
			statusCode: http.StatusBadRequest,
			more:       map[string]string{"Username": "Username is a required field"},
			errMessage: "FieldError",
		},
		{
			name:       "bad request - passwords don't match",
			dto:        registerDTOPasswordsDontMatch,
			statusCode: http.StatusBadRequest,
			errMessage: "passwords don't match",
		},
		{
			name:       "bad request - email already taken",
			dto:        registerDTOEmailAlreadyTaken,
			statusCode: http.StatusBadRequest,
			errMessage: "email is already taken",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := json.Marshal(tt.dto)
			req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(b))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			resp := w.Result()
			gerr := &GenericError{}

			if !basicHTTPResponseChecks(t, tt.statusCode, desiredContentType, gerr, resp) {
				return
			}
			assert.Equal(t, tt.errMessage, gerr.Message)
			if tt.more != nil && tt.errMessage == "FieldError" {
				for key, value := range tt.more.(map[string]string) {
					assert.Equal(t, gerr.AdditionalInfo.(map[string]interface{})[key], value)
				}
			}
		})
	}

}

func getNewRouter() *mux.Router {
	l := logrus.New()
	l.SetOutput(ioutil.Discard)
	router := mux.NewRouter()
	ur, _ := repositories.NewUserRepository(
		repositories.InMemoryKind,
		repositories.InMemoryArgs{testUserData},
	)
	uc := usecase.NewUser(l, ur, usecase.UserOptions{})
	ah := NewAuth(l, uc, domain.NewValidation())
	ah.AttachRouter(router)
	return router
}

var testUserData = []*domain.User{{
	ID:       "54215f2a-b752-11eb-8529-0242ac130003",
	Username: "jack",
	Email:    "jack@gmail.com",
	Password: "8bb0cf6eb9b17d0f7d22b456f121257dc1254e1f01665370476383ea776df414",
}, {
	ID:       "5a823a9c-b752-11eb-8529-0242ac130003",
	Username: "john",
	Email:    "john@gmail.com",
	Password: "8bb0cf6eb9b17d0f7d22b456f121257dc1254e1f01665370476383ea776df414",
}, {
	ID:       "601427c2-b752-11eb-8529-0242ac130003",
	Username: "jill",
	Email:    "jill@gmail.com",
	Password: "8bb0cf6eb9b17d0f7d22b456f121257dc1254e1f01665370476383ea776df414",
}}

func basicHTTPResponseChecks(t *testing.T, desiredStatusCode int, desiredContentType string, bodyResult interface{}, resp *http.Response) bool {
	var err error

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		b, _ := ioutil.ReadAll(resp.Body)

		t.Fatal(fmt.Sprintf("desired Content-Type is %s, the server returned %s; %s", desiredContentType, contentType, string(b)))
		return false
	}

	if desiredContentType == "application/json" {
		b, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(b, bodyResult)
		if err != nil {
			t.Log(string(b))
			t.Fatal(fmt.Sprintf("failed to parse json %s: %s", string(b), err))
			return false
		}
	}

	if err != nil {
		t.Fatal(fmt.Sprintf("failed read body content %s", err))
	}

	if desiredStatusCode != resp.StatusCode {
		t.Fatal(fmt.Sprintf("desired status code is %d, the server returned %d", desiredStatusCode, resp.StatusCode))
		return false
	}

	return true
}
