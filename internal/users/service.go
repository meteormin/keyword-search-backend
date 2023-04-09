package users

import (
	"github.com/miniyus/keyword-search-backend/entity"
)

type Service interface {
	Create(user CreateUser) (*UserResponse, error)
	All() ([]UserResponse, error)
	Get(pk uint) (*UserResponse, error)
	Update(pk uint, user PatchUser) (*UserResponse, error)
	Delete(pk uint) (bool, error)
}

type ServiceStruct struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceStruct{repo: repo}
}

func (s *ServiceStruct) Create(user CreateUser) (*UserResponse, error) {
	ent := user.ToEntity()
	create, err := s.repo.Create(ent)
	if err != nil {
		return nil, err
	}

	var ur UserResponse
	userRes := ur.FromEntity(*create)

	return &userRes, nil
}

func (s *ServiceStruct) All() ([]UserResponse, error) {
	var userRes []UserResponse

	entities, err := s.repo.All()

	if err != nil {
		return userRes, err
	}

	var ur UserResponse
	for _, ent := range entities {
		userRes = append(userRes, ur.FromEntity(ent))
	}

	return userRes, nil
}

func (s *ServiceStruct) Get(pk uint) (*UserResponse, error) {
	user, err := s.repo.Find(pk)
	if err != nil {
		return nil, err
	}

	var ur UserResponse
	userRes := ur.FromEntity(*user)

	return &userRes, nil
}

func (s *ServiceStruct) Update(pk uint, user PatchUser) (*UserResponse, error) {
	ent := entity.User{}
	if user.Email != nil {
		ent.Email = *user.Email
	}

	if user.Role != nil {
		ent.Role = entity.UserRole(*user.Role)
	}

	rsUser, err := s.repo.Update(pk, ent)

	if err != nil {
		return nil, err
	}

	var ur UserResponse
	userRes := ur.FromEntity(*rsUser)

	return &userRes, nil
}

func (s *ServiceStruct) Delete(pk uint) (bool, error) {
	return s.repo.Delete(pk)
}
