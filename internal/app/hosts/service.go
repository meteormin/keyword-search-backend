package hosts

import (
	"github.com/miniyus/go-fiber/internal/entity"
)

type Service interface {
	All() ([]*entity.Host, error)
	GetByUserId(userId uint) ([]*entity.Host, error)
	Find(pk uint, userId uint) (*entity.Host, error)
	Create(host *CreateHost) (*entity.Host, error)
	Update(pk uint, userId uint, host *UpdateHost) (*entity.Host, error)
	Delete(pk uint, userId uint) (bool, error)
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

func (s *ServiceStruct) GetByUserId(userId uint) ([]*entity.Host, error) {
	return s.repo.GetByUserId(userId)
}

func (s *ServiceStruct) Find(pk uint, userId uint) (*entity.Host, error) {
	return s.repo.Find(pk, userId)
}

func (s *ServiceStruct) Create(host *CreateHost) (*entity.Host, error) {
	ent := entity.Host{
		UserId:      host.UserId,
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

func (s *ServiceStruct) Update(pk uint, userId uint, host *UpdateHost) (*entity.Host, error) {
	ent := entity.Host{
		Subject:     host.Subject,
		Description: host.Description,
		Path:        host.Path,
		Publish:     host.Publish,
	}

	updated, err := s.repo.Update(pk, userId, ent)
	if err != nil {
		return nil, err
	}

	return updated, err
}

func (s *ServiceStruct) Delete(pk uint, userId uint) (bool, error) {
	return s.repo.Delete(pk, userId)
}
