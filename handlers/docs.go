// Package classification Minaria
//
// Documentation for Minaria
//
//	Schemes: http
//	BasePath: /
//	Version: 0.1.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package handlers

import "github.com/vahidmostofi/minaria/domain"

//
// NOTE: Types defined here are purely for documentation purposes
// these types are not used by any of the handers

// No content is returned by this API endpoint
// swagger:response noContentResponse
type noContentResponseWrapper struct {
}

// JWT Data Transfer Object response contains the jwt token string
// swagger:response jwtDTOResponse
type jwtDTOResponseWrapper struct {
	// in: body
	Body domain.JWTDTO
}

// Generic Error respones contains an error object returned
// swagger:response genericErrorResponse
type genericErrorResponseWrapper struct {
	// in: body
	Body GenericError
}

// Validation Error respones contains an error object similar
// to GenericError but the mesage field is "FieldError" and
// the more field contains a map from field to error
// swagger:response validationErrorResponse
type validationErrorResponseWrapper struct {
	// in: body
	Body GenericError
}

// Username Password don't match Error response contains an
// error object returned, the message field is:
// "email and the password don't match".
// swagger:response usernamePasswordNotMatchResponse
type usernamePasswordNotMatchResponseWrapper struct {
	// in: body
	Body GenericError
}

// Internal Server error response contains an error object
// returned, the message field is:
// "internal server error".
// swagger:response internalErrorResponse
type internalErrorResponseWrapper struct {
	// in: body
	Body GenericError
}

//swagger:parameters loginUser
type loginDTOWrapper struct {
	// in: body
	Body domain.LoginDTO
}

//swagger:parameters registerUser
type registerDTOWrapper struct {
	// in: body
	Body domain.RegisterDTO
}
