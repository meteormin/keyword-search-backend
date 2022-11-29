package users

import "github.com/miniyus/go-fiber/internal/entity"

type Service interface {
	All() ([]UserResponse, error)
	Get(pk uint) (*UserResponse, error)
	Update(pk uint, user PatchUser) (*UserResponse, error)
	ResetPassword(pk uint, passwordStruct ResetPasswordStruct) (*UserResponse, error)
}

type ServiceStruct struct {
	repo Repository
}

func NewService(repo Repository) *ServiceStruct {
	return &ServiceStruct{repo: repo}
}

func toUserResponseFromEntity(user *entity.User) UserResponse {
	return UserResponse{
		Id:              user.ID,
		Username:        user.Username,
		Email:           user.Email,
		EmailVerifiedAt: user.EmailVerifiedAt,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
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

func (s *ServiceStruct) Update(pk uint, user PatchUser) (*UserResponse, error) {
	rsUser, err := s.repo.Update(pk, entity.User{
		Email: user.Email,
	})

	if err != nil {
		return nil, err
	}

	userRes := toUserResponseFromEntity(rsUser)

	return &userRes, nil
}

func (s *ServiceStruct) ResetPassword(pk uint, passwordStruct ResetPasswordStruct) (*UserResponse, error) {
	rsUser, err := s.repo.Update(pk, entity.User{
		Password: passwordStruct.Password,
	})

	if err != nil {
		return nil, err
	}

	userRes := toUserResponseFromEntity(rsUser)

	return &userRes, nil
}
