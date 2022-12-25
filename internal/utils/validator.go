package utils

import (
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

func newValidator() *validator.Validate {
	return validator.New()
}

func newTranslator(locale locales.Translator) (trans ut.Translator) {
	uni := ut.New(locale, locale)
	trans, found := uni.GetTranslator("en")
	if found {
		return trans
	}

	return nil
}

func Validate(data interface{}) map[string]string {
	validate := newValidator()

	trans := newTranslator(en.New())
	var transErr error
	if trans != nil {
		transErr = enTranslations.RegisterDefaultTranslations(validate, trans)
		if transErr == nil {
			registerTranslation(validate, trans)
		}
	}

	fields := map[string]string{}
	errs := validate.Struct(data)

	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			if err != nil {
				if transErr == nil && trans != nil {
					fields[err.Field()] = err.Translate(trans)
				} else {
					fields[err.Field()] = err.Error()
				}
			}
		}

		return fields
	}

	return nil
}

func registerTranslation(validate *validator.Validate, trans ut.Translator) {
	for _, t := range translations(trans) {
		_ = validate.RegisterTranslation(
			t.tag,
			t.trans,
			t.registerFn,
			t.translationFn,
		)
	}
}
