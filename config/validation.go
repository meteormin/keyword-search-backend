package config

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	mConfig "github.com/miniyus/gofiber/config"
	"github.com/miniyus/gofiber/pkg/validation"
)

func validationConfig() mConfig.Validation {
	return mConfig.Validation{
		Validations: []validation.Tag{},
		Translations: []validation.TranslationTag{
			{
				Tag: "required",
				RegisterFn: func(ut ut.Translator) error {
					return ut.Add("required", "{0} 필드는 필수 입니다.", true)
				},
				TranslationFn: func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T("required", fe.Field())
					return t
				},
			},
			{
				Tag: "email",
				RegisterFn: func(ut ut.Translator) error {
					return ut.Add("email", "{0} 필드는 Email 형식이어야 합니다.", true)
				},
				TranslationFn: func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T("email", fe.Field())
					return t
				},
			},
			{
				Tag: "url",
				RegisterFn: func(ut ut.Translator) error {
					return ut.Add("url", "{0} 필드는 URL 형식이어야 합니다.", true)
				},
				TranslationFn: func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T("url", fe.Field())
					return t
				},
			},
			{
				Tag: "dir",
				RegisterFn: func(ut ut.Translator) error {
					return ut.Add("dir", "{0} 필드는 디렉토리 경로 형식이어야 합니다.", true)
				},
				TranslationFn: func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T("dir", fe.Field())
					return t
				},
			},
			{
				Tag: "boolean",
				RegisterFn: func(ut ut.Translator) error {
					return ut.Add("boolean", "{0} 필드는 boolean 타입이어야 합니다.", true)
				},
				TranslationFn: func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T("boolean", fe.Field())
					return t
				},
			},
			{
				Tag: "eqfield",
				RegisterFn: func(ut ut.Translator) error {
					return ut.Add("eqfield", "{0} 필드는 {1} 필드와 일치해야 합니다.", true)
				},
				TranslationFn: func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T("eqfield", fe.Field(), fe.Param())
					return t
				},
			},
		},
	}
}
