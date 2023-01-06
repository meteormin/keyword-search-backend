package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/internal/core/api_error"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	DefaultDateLayout = "2006-01-02 15:04:05"
)

func HandleValidate(c *fiber.Ctx, data interface{}) *api_error.ValidationErrorResponse {
	failed := Validate(data)
	if failed != nil {
		errRes := api_error.NewValidationErrorResponse(c, failed)
		return errRes
	}

	return nil
}

func TimeIn(t time.Time, tz string) time.Time {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		panic(err)
	}

	return t.In(loc)
}

func HashPassword(password string) (string, error) {
	fromPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(fromPassword), err
}

func HashCheck(hashPass string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password))
	return err == nil
}
