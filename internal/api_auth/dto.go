package api_auth

import (
	"github.com/miniyus/gofiber/utils"
)

type SignUp struct {
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`
	Email           string `json:"email" validate:"required,email"`
}

type SignIn struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenInfo struct {
	Token     string         `json:"token"`
	ExpiresAt utils.JsonTime `json:"expires_at"`
	ExpiresIn int64          `json:"expires_in"`
}

type SignUpResponse struct {
	UserId uint `json:"user_id"`
	TokenInfo
}

type ResetPasswordStruct struct {
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`
}
