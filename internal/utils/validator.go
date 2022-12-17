package utils

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

func ValidationTranslate(validatedData map[string]string, translate map[string]string) map[string]string {
	for key, _ := range validatedData {
		validatedData[key] = translate[key]
	}

	return validatedData
}
