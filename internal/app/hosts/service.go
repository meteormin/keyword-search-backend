package hosts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/core/logger"
	"github.com/miniyus/go-fiber/internal/entity"
)

type Service interface {
	All() ([]entity.Host, error)
	GetByUserId(userId uint) ([]HostResponse, error)
	Find(pk uint, userId uint) (*HostResponse, error)
	Create(host *CreateHost) (*HostResponse, error)
	Update(pk uint, userId uint, host *UpdateHost) (*HostResponse, error)
	Patch(pk uint, userId uint, host *PatchHost) (*HostResponse, error)
	Delete(pk uint, userId uint) (bool, error)
	logger.HasLogger
}

type ServiceStruct struct {
	repo Repository
	logger.HasLoggerStruct
}

func NewService(repo Repository) Service {
	return &ServiceStruct{
		repo:            repo,
		HasLoggerStruct: logger.HasLoggerStruct{Logger: repo.GetLogger()},
	}
}

func (s *ServiceStruct) All() ([]entity.Host, error) {
	hosts, err := s.repo.All()
	if err != nil {
		return make([]entity.Host, 0), err
	}

	return hosts, err
}

func (s *ServiceStruct) GetByUserId(userId uint) ([]HostResponse, error) {
	ent, err := s.repo.GetByUserId(userId)

	var dto []HostResponse

	if err != nil {
		return make([]HostResponse, 0), err
	}

	for _, e := range ent {
		dto = append(dto, *ToHostResponse(&e))
	}

	if dto == nil {
		return make([]HostResponse, 0), nil
	}

	return dto, nil
}

func (s *ServiceStruct) Find(pk uint, userId uint) (*HostResponse, error) {
	ent, err := s.repo.Find(pk)
	if err != nil {
		return nil, err
	}

	if ent == nil {
		return nil, fiber.ErrNotFound
	}

	if ent.UserId != userId {
		return nil, fiber.ErrForbidden
	}

	return ToHostResponse(ent), nil
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

	return ToHostResponse(created), err
}

func (s *ServiceStruct) Update(pk uint, userId uint, host *UpdateHost) (*HostResponse, error) {
	exists, err := s.repo.Find(pk)
	if err != nil {
		return nil, err
	}

	if exists == nil {
		return nil, fiber.ErrNotFound
	}

	if exists.UserId != userId {
		return nil, fiber.ErrForbidden
	}

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

	return ToHostResponse(updated), err
}

func (s *ServiceStruct) Patch(pk uint, userId uint, host *PatchHost) (*HostResponse, error) {
	ent, err := s.repo.Find(pk)
	if err != nil {
		return nil, err
	}

	if ent == nil {
		return nil, fiber.ErrNotFound
	}

	if ent.UserId != userId {
		return nil, fiber.ErrForbidden
	}

	if host.Host != nil {
		ent.Host = *host.Host
	}
	if host.Subject != nil {
		ent.Subject = *host.Subject
	}

	if host.Description != nil {
		ent.Description = *host.Description
	}

	if host.Path != nil {
		ent.Path = *host.Path
	}

	if host.Publish != nil {
		ent.Publish = *host.Publish
	}

	updated, err := s.repo.Update(pk, *ent)
	if err != nil {
		return nil, err
	}

	return ToHostResponse(updated), err
}

func (s *ServiceStruct) Delete(pk uint, userId uint) (bool, error) {
	exists, err := s.repo.Find(pk)
	if err != nil {
		return false, err
	}

	if exists == nil {
		return false, fiber.ErrNotFound
	}

	if exists.UserId != userId {
		return false, fiber.ErrForbidden
	}

	return s.repo.Delete(pk)
}
