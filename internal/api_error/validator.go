package api_error

import (
	"github.com/go-playground/validator/v10"
)

func newValidator() *validator.Validate {
	return validator.New()
}

func Validate(data interface{}) map[string]string {
	validate := newValidator()

	fields := map[string]string{}
	errs := validate.Struct(data)

	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			if err != nil {
				fields[err.Field()] = err.Error()
			}
		}

		return fields
	}

	return nil
}
