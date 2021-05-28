package domain

import (
	"log"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// ValidationError wraps the validators FieldError
type ValidationError struct {
	validator.FieldError
	translator ut.Translator
}

func (v ValidationError) Error() string {
	return v.Translate(v.translator)
}

// ValidationErrors is a collection of ValidationError
type ValidationErrors []ValidationError

func (v ValidationErrors) FieldsError() map[string]string {
	mapErr := make(map[string]string)
	for _, err := range v {
		mapErr[err.Field()] = err.Error()
	}
	return mapErr
}

func (v ValidationErrors) Errors() []string {
	errs := []string{}
	for _, err := range v {
		errs = append(errs, err.Error())
	}
	return errs
}

type Validation struct {
	validate   *validator.Validate
	translator ut.Translator
}

func NewValidation() *Validation {
	translator, _ := ut.New(en.New(), en.New()).GetTranslator("en")

	validate := validator.New()

	if err := en_translations.RegisterDefaultTranslations(validate, translator); err != nil {
		log.Fatalf("failed to register a default translation to validator: %s", err)
	}

	return &Validation{validate, translator}
}

// Validate an object
func (v *Validation) Validate(i interface{}) ValidationErrors {
	errs, ok := v.validate.Struct(i).(validator.ValidationErrors)

	if !ok || len(errs) == 0 {
		return nil
	}

	var returnErrs []ValidationError
	for _, err := range errs {
		ve := ValidationError{err.(validator.FieldError), v.translator}
		returnErrs = append(returnErrs, ve)
	}
	return returnErrs
}
