package users

import (
	"github.com/miniyus/go-fiber/internal/entity"
	"github.com/miniyus/go-fiber/internal/utils"
)

type Service interface {
	All() ([]UserResponse, error)
	Get(pk uint) (*UserResponse, error)
	Update(pk uint, user *PatchUser) (*UserResponse, error)
}

type ServiceStruct struct {
	repo Repository
}

func NewService(repo Repository) *ServiceStruct {
	return &ServiceStruct{repo: repo}
}

func toUserResponseFromEntity(user *entity.User) UserResponse {
	createdAt := utils.TimeIn(user.CreatedAt, "Asia/Seoul")
	updatedAt := utils.TimeIn(user.UpdatedAt, "Asia/Seoul")

	var emailVerifiedAt *string
	if user.EmailVerifiedAt == nil {
		emailVerifiedAt = nil
	} else {
		formatString := user.EmailVerifiedAt.Format("2006-01-02 15:04:05")
		emailVerifiedAt = &formatString
	}

	return UserResponse{
		Id:              user.ID,
		Username:        user.Username,
		Email:           user.Email,
		EmailVerifiedAt: emailVerifiedAt,
		CreatedAt:       createdAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       updatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (s *ServiceStruct) All() ([]UserResponse, error) {
	var userRes []UserResponse

	entities, err := s.repo.All()

	if err != nil {
		return userRes, err
	}

	for _, ent := range entities {
		userRes = append(userRes, toUserResponseFromEntity(ent))
	}

	return userRes, nil
}

func (s *ServiceStruct) Get(pk uint) (*UserResponse, error) {
	user, err := s.repo.Find(pk)
	if err != nil {
		return nil, err
	}

	userRes := toUserResponseFromEntity(user)

	return &userRes, nil

}

func (s *ServiceStruct) Update(pk uint, user *PatchUser) (*UserResponse, error) {
	rsUser, err := s.repo.Update(pk, entity.User{
		Email: user.Email,
	})

	if err != nil {
		return nil, err
	}

	userRes := toUserResponseFromEntity(rsUser)

	return &userRes, nil
}
