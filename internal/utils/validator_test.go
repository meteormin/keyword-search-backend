package utils_test

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/miniyus/keyword-search-backend/internal/utils"
	"testing"
)

type TestData struct {
	Id   string `validate:"required"`
	Pass string `validate:"required"`
}

func TestValidate(t *testing.T) {
	data := TestData{
		Id: "test",
	}

	validated := utils.Validate(data)

	testValidated := map[string]string{
		"Pass": "Key: 'TestData.Pass' Error:Field validation for 'Pass' failed on the 'required' tag",
	}

	t.Log(validated["Pass"])
	assert.Equal(t, validated, testValidated)
}

func TestValidationTranslate(t *testing.T) {
	data := TestData{
		Id: "test",
	}
	locale := en.New()
	uni := ut.New(locale, locale)
	trans, _ := uni.GetTranslator("en")
	val := validator.New()
	err := enTranslations.RegisterDefaultTranslations(val, trans)
	if err != nil {
		t.Error(err)
	}

	err = val.Struct(data)
	if err != nil {
		errs := err.(validator.ValidationErrors)

		fmt.Println(errs.Translate(trans))
	}

}
