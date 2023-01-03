package groups

import (
	"github.com/miniyus/keyword-search-backend/internal/core/logger"
	"github.com/miniyus/keyword-search-backend/internal/entity"
	"github.com/miniyus/keyword-search-backend/internal/utils"
)

type Service interface {
	Create(group *CreateGroup) (*ResponseGroup, error)
	Update(groupId uint, group *UpdateGroup) (*ResponseGroup, error)
	Delete(pk uint) (bool, error)
	All(page utils.Page) (utils.Paginator, error)
	Find(pk uint) (*ResponseGroup, error)
	FindByName(groupName string) (*ResponseGroup, error)
}

type ServiceStruct struct {
	repo Repository
	logger.HasLoggerStruct
}

func NewService(repo Repository) Service {
	return &ServiceStruct{
		repo: repo,
		HasLoggerStruct: logger.HasLoggerStruct{
			Logger: repo.GetLogger(),
		},
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

	return ToResponse(create), err
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

	return ToResponse(update), err
}

func (s *ServiceStruct) Delete(pk uint) (bool, error) {
	return s.repo.Delete(pk)
}

func (s *ServiceStruct) All(page utils.Page) (utils.Paginator, error) {
	res := make([]ResponseGroup, 0)

	entities, count, err := s.repo.All(page)
	paginator := utils.Paginator{
		TotalCount: count,
		Page:       page,
		Data:       res,
	}

	if err != nil {
		return paginator, err
	}

	for _, ent := range entities {
		resGroup := ToResponse(&ent)
		if resGroup != nil {
			res = append(res, *resGroup)
		}
	}

	paginator.Data = res

	return paginator, nil
}

func (s *ServiceStruct) Find(pk uint) (*ResponseGroup, error) {
	ent, err := s.repo.Find(pk)
	if err != nil {
		return nil, err
	}

	return ToResponse(ent), nil
}

func (s *ServiceStruct) FindByName(groupName string) (*ResponseGroup, error) {
	ent, err := s.repo.FindByName(groupName)
	if err != nil {
		return nil, err
	}

	return ToResponse(ent), nil
}
