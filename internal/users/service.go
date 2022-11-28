package users

type Service interface {
	All() ([]UserResponse, error)
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

	for _, entity := range entities {
		userRes = append(userRes, UserResponse{
			Id:              entity.ID,
			Username:        entity.Username,
			Email:           entity.Email,
			EmailVerifiedAt: entity.EmailVerifiedAt,
			CreatedAt:       entity.CreatedAt,
			UpdatedAt:       entity.UpdatedAt,
		})
	}

	return userRes, nil
}
