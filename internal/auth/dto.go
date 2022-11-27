package auth

import "time"

type User struct {
	Id        uint
	GroupId   uint
	Username  string
	Email     string
	CreatedAt time.Time
}

type SignUp struct {
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
}

type SignIn struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
