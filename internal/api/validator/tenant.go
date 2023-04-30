package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

type UserValidator struct {
	Validate *validator.Validate
}

// Register Custom Validators
func NewUserValidator() UserValidator {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		fieldVal := fl.Field().String()
		if fl.Field().String() != "" {
			if len([]rune(fieldVal)) < 6 {
				return false
			}
			exp1 := regexp.MustCompile(`[0-9]`)
			match1 := exp1.MatchString(fieldVal)
			exp2 := regexp.MustCompile(`[a-z]`)
			match2 := exp2.MatchString(fieldVal)
			exp3 := regexp.MustCompile(`[A-Z]`)
			match3 := exp3.MatchString(fieldVal)
			exp4 := regexp.MustCompile(`[!"#$%&'()*+,-./:;<=>?@[\]^_{|}~]`)
			match4 := exp4.MatchString(fieldVal)
			if !match1 || !match2 || !match3 || !match4 {
				return false
			}
		}
		return true
	})

	return UserValidator{
		Validate: v,
	}
}

func (v UserValidator) validate(s interface{}) error {
	return v.Validate.Struct(s)
}

func (v UserValidator) buildValidationMessage(field string, rule string) (message string) {
	switch rule {
	case "required":
		return fmt.Sprintf("field '%s' is required.", field)
	case "password":
		return fmt.Sprintf("field '%s' should be minimum 6 characters with at least one letter, number and a special character.", field)
	default:
		return fmt.Sprintf("field '%s' is not valid.", field)
	}
}
