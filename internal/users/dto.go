package users

import (
	"github.com/miniyus/gofiber/utils"
	"github.com/miniyus/keyword-search-backend/entity"
)

type CreateUser struct {
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`
	Email           string `json:"email" validate:"required,email"`
}

type PatchUser struct {
	Email *string `json:"email" validate:"email"`
	Role  *string `json:"role" validate:"string"`
}

type UserResponse struct {
	Id              uint    `json:"id"`
	Username        string  `json:"username"`
	Email           string  `json:"email"`
	EmailVerifiedAt *string `json:"email_verified_at"`
	Role            string  `json:"role"`
	GroupId         *uint   `json:"group_id"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

func (cr CreateUser) ToEntity() entity.User {
	res := entity.User{
		Username: cr.Username,
		Password: cr.Password,
		Email:    cr.Email,
	}

	return res
}

func (ur UserResponse) FromEntity(user entity.User) UserResponse {
	createdAt := utils.TimeIn(user.CreatedAt, "Asia/Seoul")
	updatedAt := utils.TimeIn(user.UpdatedAt, "Asia/Seoul")

	var emailVerifiedAt *string
	if user.EmailVerifiedAt == nil {
		emailVerifiedAt = nil
	} else {
		formatString := user.EmailVerifiedAt.Format(utils.DefaultDateLayout)
		emailVerifiedAt = &formatString
	}

	return UserResponse{
		Id:              user.ID,
		Username:        user.Username,
		Role:            string(user.Role),
		Email:           user.Email,
		EmailVerifiedAt: emailVerifiedAt,
		CreatedAt:       createdAt.Format(utils.DefaultDateLayout),
		UpdatedAt:       updatedAt.Format(utils.DefaultDateLayout),
	}
}
