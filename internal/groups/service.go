package groups

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/pagination"
	"github.com/miniyus/keyword-search-backend/entity"
	"github.com/miniyus/keyword-search-backend/repo"
)

type Service interface {
	Create(group *CreateGroup) (*ResponseGroup, error)
	Update(groupId uint, group *UpdateGroup) (*ResponseGroup, error)
	Delete(pk uint) (bool, error)
	All(page pagination.Page) (pagination.Paginator[ResponseGroup], error)
	Find(pk uint) (*ResponseGroup, error)
	FindByName(groupName string) (*ResponseGroup, error)
}

type ServiceStruct struct {
	repo repo.GroupRepository
}

func NewService(repo repo.GroupRepository) Service {
	return &ServiceStruct{
		repo: repo,
	}
}

func (s *ServiceStruct) Create(group *CreateGroup) (*ResponseGroup, error) {
	var ent entity.Group
	ent.Name = group.Name
	for _, perm := range group.Permissions {
		entPerm := entity.Permission{
			Permission: perm.Name,
		}
		for _, action := range perm.Actions {
			entPerm.Actions = append(entPerm.Actions, entity.Action{
				Resource: action.Resource,
				Method:   string(action.Method),
			})
		}

		ent.Permissions = append(ent.Permissions, entPerm)
	}

	create, err := s.repo.Create(ent)
	if err != nil {
		return nil, err
	}

	var resGroup ResponseGroup
	res := resGroup.FromEntity(*create)
	return &res, err
}

func (s *ServiceStruct) Update(groupId uint, group *UpdateGroup) (*ResponseGroup, error) {
	var ent entity.Group
	ent.Name = group.Name
	for _, perm := range group.Permissions {
		entPerm := entity.Permission{
			Permission: perm.Name,
			GroupId:    groupId,
		}
		for _, action := range perm.Actions {
			entPerm.Actions = append(entPerm.Actions, entity.Action{
				Resource: action.Resource,
				Method:   string(action.Method),
			})
		}

		ent.Permissions = append(ent.Permissions, entPerm)
	}

	update, err := s.repo.Update(groupId, ent)
	if err != nil {
		return nil, err
	}

	var resGroup ResponseGroup

	res := resGroup.FromEntity(*update)

	return &res, err
}

func (s *ServiceStruct) Delete(pk uint) (bool, error) {
	return s.repo.Delete(pk)
}

func (s *ServiceStruct) All(page pagination.Page) (pagination.Paginator[ResponseGroup], error) {
	res := make([]ResponseGroup, 0)

	entities, count, err := s.repo.AllWithPage(page)

	paginator := pagination.Paginator[ResponseGroup]{
		TotalCount: count,
		Page:       page,
		Data:       res,
	}

	if err != nil {
		return paginator, err
	}

	var resG ResponseGroup
	for _, ent := range entities {
		resGroup := resG.FromEntity(ent)
		res = append(res, resGroup)
	}

	paginator.Data = res

	return paginator, nil
}

func (s *ServiceStruct) Find(pk uint) (*ResponseGroup, error) {
	ent, err := s.repo.Find(pk)
	if err != nil {
		return nil, err
	}

	if ent == nil {
		return nil, fiber.ErrNotFound
	}

	var resGroup ResponseGroup
	res := resGroup.FromEntity(*ent)
	return &res, nil
}

func (s *ServiceStruct) FindByName(groupName string) (*ResponseGroup, error) {
	ent, err := s.repo.FindByName(groupName)
	if err != nil {
		return nil, err
	}

	if ent == nil {
		return nil, fiber.ErrNotFound
	}

	var resGroup ResponseGroup
	res := resGroup.FromEntity(*ent)
	return &res, nil
}
