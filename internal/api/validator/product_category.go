package validator

import (
	"fmt"
	"reflect"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

type ProductCategoryValidator struct {
	Validate *validator.Validate
}

// Register Custom Validators
func NewProductCategoryValidator() ProductCategoryValidator {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return ProductCategoryValidator{
		Validate: v,
	}
}

func (v ProductCategoryValidator) validate(s interface{}) error {
	return v.Validate.Struct(s)
}

func (v ProductCategoryValidator) buildValidationMessage(field string, rule string) (message string) {
	switch rule {
	default:
		return fmt.Sprintf("field '%s' is not valid.", field)
	case "required":
		return fmt.Sprintf("field '%s' is required.", field)
	}
}
