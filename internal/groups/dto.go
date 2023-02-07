package groups

import (
	"github.com/miniyus/gofiber/permission"
	"github.com/miniyus/gofiber/utils"
	"github.com/miniyus/keyword-search-backend/entity"
)

type CreateAction struct {
	Resource string            `json:"resource" validate:"required"`
	Method   permission.Method `json:"method" validate:"required"`
}

type CreatePermission struct {
	Name    string         `json:"name" validate:"required"`
	Actions []CreateAction `json:"actions" validate:"required"`
}

type CreateGroup struct {
	Name        string             `json:"name" validate:"required"`
	Permissions []CreatePermission `json:"permissions" validate:"required"`
}

type UpdateGroup struct {
	Name        string             `json:"name"`
	Permissions []CreatePermission `json:"permissions"`
}

type ResponseAction struct {
	Resource string            `json:"resource"`
	Method   permission.Method `json:"method"`
}

type ResponseGroup struct {
	Id      uint             `json:"id"`
	Name    string           `json:"name"`
	Actions []ResponseAction `json:"actions"`
}

type ListResponse struct {
	utils.Paginator[ResponseGroup]
	Data []ResponseGroup `json:"data"`
}

func ToResponse(ent *entity.Group) *ResponseGroup {
	if ent == nil {
		return nil
	}

	res := &ResponseGroup{
		Name:    ent.Name,
		Id:      ent.ID,
		Actions: make([]ResponseAction, 0),
	}

	for _, perm := range ent.Permissions {
		for _, action := range perm.Actions {
			res.Actions = append(res.Actions, ResponseAction{
				Method:   permission.Method(action.Method),
				Resource: action.Resource,
			})
		}
	}
	return res
}
