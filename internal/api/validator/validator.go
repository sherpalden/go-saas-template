package validator

import (
	validator "github.com/go-playground/validator/v10"
	"github.com/sherpalden/go-saas-template/internal/app_error"
)

type customValidator interface {
	buildValidationMessage(string, string) string
	validate(s interface{}) error
}

func GenerateValidation(cv customValidator, s interface{}) (validations []app_error.FieldError, passed bool) {
	errs := cv.validate(s)
	if errs != nil {
		for _, value := range errs.(validator.ValidationErrors) {
			field, rule := value.Field(), value.Tag()
			validation := app_error.FieldError{Field: field, Message: cv.buildValidationMessage(field, rule)}
			validations = append(validations, validation)
		}
	}
	if len(validations) > 0 {
		return validations, false
	} else {
		return validations, true
	}
}
