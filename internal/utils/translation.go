package utils

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type registerTranslationStruct struct {
	tag           string
	trans         ut.Translator
	registerFn    validator.RegisterTranslationsFunc
	translationFn validator.TranslationFunc
}

func translations(translator ut.Translator) []registerTranslationStruct {
	trans := []registerTranslationStruct{
		{
			tag:   "required",
			trans: translator,
			registerFn: func(ut ut.Translator) error {
				return ut.Add("required", "{0} 필드는 필수 입니다.", true)
			},
			translationFn: func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("required", fe.Field())
				return t
			},
		},
		{
			tag:   "email",
			trans: translator,
			registerFn: func(ut ut.Translator) error {
				return ut.Add("email", "{0} 필드는 Email 형식이어야 합니다.", true)
			},
			translationFn: func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("email", fe.Field())
				return t
			},
		},
		{
			tag:   "url",
			trans: translator,
			registerFn: func(ut ut.Translator) error {
				return ut.Add("url", "{0} 필드는 URL 형식이어야 합니다.", true)
			},
			translationFn: func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("url", fe.Field())
				return t
			},
		},
		{
			tag:   "dir",
			trans: translator,
			registerFn: func(ut ut.Translator) error {
				return ut.Add("dir", "{0} 필드는 디렉토리 경로 형식이어야 합니다.", true)
			},
			translationFn: func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("dir", fe.Field())
				return t
			},
		},
		{
			tag:   "boolean",
			trans: translator,
			registerFn: func(ut ut.Translator) error {
				return ut.Add("boolean", "{0} 필드는 boolean 타입이어야 합니다.", true)
			},
			translationFn: func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("boolean", fe.Field())
				return t
			},
		},
	}

	return trans
}
