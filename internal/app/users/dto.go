package users

import (
	"github.com/miniyus/go-fiber/internal/entity"
	"github.com/miniyus/go-fiber/internal/utils"
)

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

type PatchUser struct {
	Email string `json:"email" validate:"email"`
	Role  string `json:"role" validate:"string"`
}

func ToUserResponse(user *entity.User) UserResponse {
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
		Role:            user.Role,
		Email:           user.Email,
		EmailVerifiedAt: emailVerifiedAt,
		CreatedAt:       createdAt.Format(utils.DefaultDateLayout),
		UpdatedAt:       updatedAt.Format(utils.DefaultDateLayout),
	}
}
