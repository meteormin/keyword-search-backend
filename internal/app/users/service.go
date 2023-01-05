package users

import (
	"github.com/miniyus/keyword-search-backend/internal/entity"
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

func (s *ServiceStruct) All() ([]UserResponse, error) {
	var userRes []UserResponse

	entities, err := s.repo.All()

	if err != nil {
		return userRes, err
	}

	for _, ent := range entities {
		userRes = append(userRes, ToUserResponse(&ent))
	}

	return userRes, nil
}

func (s *ServiceStruct) Get(pk uint) (*UserResponse, error) {
	user, err := s.repo.Find(pk)
	if err != nil {
		return nil, err
	}

	userRes := ToUserResponse(user)

	return &userRes, nil
}

func (s *ServiceStruct) Update(pk uint, user *PatchUser) (*UserResponse, error) {
	rsUser, err := s.repo.Update(pk, entity.User{
		Email: user.Email,
	})

	if err != nil {
		return nil, err
	}

	userRes := ToUserResponse(rsUser)

	return &userRes, nil
}
