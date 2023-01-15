package utils

import "github.com/go-playground/validator/v10"

type registerValidationStruct struct {
	tag string
	fn  validator.Func
}

func validations() []registerValidationStruct {
	registerValidations := []registerValidationStruct{
		{
			tag: "custom_test",
			fn: func(fl validator.FieldLevel) bool {
				return fl.Field().String() == "test"
			},
		},
	}

	return registerValidations
}
