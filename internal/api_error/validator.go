package api_error

import "github.com/go-playground/validator/v10"

func newValidator() *validator.Validate {
	return validator.New()
}

func Validate(err error) map[string]string {
	fields := map[string]string{}

	for _, err := range err.(validator.ValidationErrors) {
		fields[err.Field()] = err.Error()
	}

	return fields
}
