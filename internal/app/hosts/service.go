package hosts

import (
	"github.com/miniyus/go-fiber/internal/entity"
)

type Service interface {
	All() ([]*entity.Host, error)
	GetByUserId(userId uint) ([]*HostResponse, error)
	Find(pk uint, userId uint) (*HostResponse, error)
	Create(host *CreateHost) (*HostResponse, error)
	Update(pk uint, userId uint, host *UpdateHost) (*HostResponse, error)
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

func toResponseDto(host *entity.Host) *HostResponse {
	return &HostResponse{
		Id:          host.ID,
		Host:        host.Host,
		Path:        host.Path,
		UserId:      host.UserId,
		Subject:     host.Subject,
		Description: host.Description,
		Publish:     host.Publish,
	}
}

func (s *ServiceStruct) GetByUserId(userId uint) ([]*HostResponse, error) {
	ent, err := s.repo.GetByUserId(userId)
	if err != nil {
		return nil, err
	}

	var dto []*HostResponse
	for _, e := range ent {
		dto = append(dto, toResponseDto(e))
	}

	return dto, nil
}

func (s *ServiceStruct) Find(pk uint, userId uint) (*HostResponse, error) {
	ent, err := s.repo.Find(pk, userId)
	if err != nil {
		return nil, err
	}

	return toResponseDto(ent), nil
}

func (s *ServiceStruct) Create(host *CreateHost) (*HostResponse, error) {
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

	return toResponseDto(created), err
}

func (s *ServiceStruct) Update(pk uint, userId uint, host *UpdateHost) (*HostResponse, error) {
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

	return toResponseDto(updated), err
}

func (s *ServiceStruct) Delete(pk uint, userId uint) (bool, error) {
	return s.repo.Delete(pk, userId)
}
