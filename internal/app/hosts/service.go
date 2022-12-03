package hosts

import "github.com/miniyus/go-fiber/internal/entity"

type Service interface {
	All() ([]*entity.Host, error)
	AllFromUser(userId uint) ([]*entity.Host, error)
	Find(pk uint) (*entity.Host, error)
	Create(host *CreateHost) (*entity.Host, error)
	Update(pk uint, host *UpdateHost) (*entity.Host, error)
	Delete(pk uint) (bool, error)
}

type ServiceStruct struct {
	repo Repository
}

func NewService(repo Repository) *ServiceStruct {
	return &ServiceStruct{repo: repo}
}

func (s *ServiceStruct) All() ([]*entity.Host, error) {
	hosts, err := s.repo.All()
	if err != nil {
		return nil, err
	}

	return hosts, nil
}

func (s *ServiceStruct) AllFromUser(userId uint) ([]*entity.Host, error) {
	return s.repo.AllFromUser(userId)
}

func (s *ServiceStruct) Find(pk uint) (*entity.Host, error) {
	return s.repo.Find(pk)
}

func (s *ServiceStruct) Create(host *CreateHost) (*entity.Host, error) {
	ent := entity.Host{
		Host:        host.Host,
		Subject:     host.Subject,
		Description: host.Description,
		Path:        host.Path,
		Publish:     host.Publish,
	}

	created, err := s.repo.Create(ent)
	if err != nil {
		return nil, err
	}

	return created, err
}

func (s *ServiceStruct) Update(pk uint, host *UpdateHost) (*entity.Host, error) {
	ent := entity.Host{
		Subject:     host.Subject,
		Description: host.Description,
		Path:        host.Path,
		Publish:     host.Publish,
	}

	updated, err := s.repo.Update(pk, ent)
	if err != nil {
		return nil, err
	}

	return updated, err
}

func (s *ServiceStruct) Delete(pk uint) (bool, error) {
	return s.repo.Delete(pk)
}
