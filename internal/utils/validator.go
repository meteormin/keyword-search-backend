package utils

import (
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

var validate *validator.Validate
var transErr error
var trans ut.Translator

func init() {
	validate = newValidator()
}

func newValidator() *validator.Validate {
	v := validator.New()
	registerValidation(v)
	trans = newTranslator(en.New())
	if trans != nil {
		transErr = enTranslations.RegisterDefaultTranslations(v, trans)
		if transErr == nil {
			registerTranslation(v, trans)
		}
	}

	return v
}

func newTranslator(locale locales.Translator) (t ut.Translator) {
	uni := ut.New(locale, locale)
	t, found := uni.GetTranslator("en")
	if found {
		return t
	}

	return nil
}

func registerValidation(validate *validator.Validate) {
	for _, v := range validations() {
		_ = validate.RegisterValidation(
			v.tag,
			v.fn,
		)
	}
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

func Validate(data interface{}) map[string]string {
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
