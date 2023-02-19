package hosts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/utils"
	"github.com/miniyus/keyword-search-backend/entity"
)

type Service interface {
	All(page utils.Page) (utils.Paginator[entity.Host], error)
	GetByUserId(userId uint, page utils.Page) (utils.Paginator[HostResponse], error)
	GetSubjectsByGroupId(groupId uint, page utils.Page) (utils.Paginator[Subjects], error)
	Find(pk uint, userId uint) (*HostResponse, error)
	Create(host *CreateHost) (*HostResponse, error)
	Update(pk uint, userId uint, host *UpdateHost) (*HostResponse, error)
	Patch(pk uint, userId uint, host *PatchHost) (*HostResponse, error)
	Delete(pk uint, userId uint) (bool, error)
}

type ServiceStruct struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceStruct{
		repo: repo,
	}
}

func (s *ServiceStruct) All(page utils.Page) (utils.Paginator[entity.Host], error) {
	hosts, count, err := s.repo.AllWithPage(page)
	if err != nil {
		hosts = make([]entity.Host, 0)
		return utils.Paginator[entity.Host]{
			Page:       page,
			TotalCount: 0,
			Data:       hosts,
		}, err
	}

	return utils.Paginator[entity.Host]{
		Page:       page,
		TotalCount: count,
		Data:       hosts,
	}, err
}

func (s *ServiceStruct) GetByUserId(userId uint, page utils.Page) (utils.Paginator[HostResponse], error) {
	ent, count, err := s.repo.GetByUserId(userId, page)

	var dto []HostResponse

	if err != nil {
		dto = make([]HostResponse, 0)
		return utils.Paginator[HostResponse]{
			Page:       page,
			TotalCount: 0,
			Data:       dto,
		}, err
	}

	var hr HostResponse
	for _, e := range ent {
		dto = append(dto, hr.FromEntity(e))
	}

	if dto == nil {
		dto = make([]HostResponse, 0)
		return utils.Paginator[HostResponse]{
			Page:       page,
			TotalCount: 0,
			Data:       dto,
		}, err
	}

	return utils.Paginator[HostResponse]{
		Page:       page,
		TotalCount: count,
		Data:       dto,
	}, err
}

func (s *ServiceStruct) GetSubjectsByGroupId(groupId uint, page utils.Page) (utils.Paginator[Subjects], error) {
	ent, count, err := s.repo.GetSubjectsByGroupId(groupId, page)

	var dto []Subjects
	if err != nil {
		dto = make([]Subjects, 0)
		return utils.Paginator[Subjects]{
			Page:       page,
			TotalCount: 0,
			Data:       dto,
		}, err
	}

	for _, e := range ent {
		dto = append(dto, Subjects{
			Id:      e.ID,
			Subject: e.Subject,
		})
	}

	if dto == nil {
		dto = make([]Subjects, 0)
		return utils.Paginator[Subjects]{
			Page:       page,
			TotalCount: 0,
			Data:       dto,
		}, err
	}

	return utils.Paginator[Subjects]{
		Page:       page,
		TotalCount: count,
		Data:       dto,
	}, err
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

	var hr HostResponse
	res := hr.FromEntity(*ent)

	return &res, nil
}

func (s *ServiceStruct) Create(host *CreateHost) (*HostResponse, error) {
	ent := host.ToEntity()

	created, err := s.repo.Create(ent)
	if err != nil {
		return nil, err
	}

	var hr HostResponse
	res := hr.FromEntity(*created)
	return &res, nil
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

	ent := host.ToEntity()
	ent.UserId = userId

	updated, err := s.repo.Update(pk, ent)
	if err != nil {
		return nil, err
	}

	var hr HostResponse
	res := hr.FromEntity(*updated)

	return &res, nil
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

	toEntity := host.ToEntity()
	toEntity.UserId = userId

	updated, err := s.repo.Update(pk, toEntity)
	if err != nil {
		return nil, err
	}

	var hr HostResponse
	res := hr.FromEntity(*updated)

	return &res, nil
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
