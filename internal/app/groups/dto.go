package groups

import (
	"github.com/miniyus/keyword-search-backend/internal/core/permission"
	"github.com/miniyus/keyword-search-backend/internal/entity"
	"github.com/miniyus/keyword-search-backend/internal/utils"
)

type CreateAction struct {
	Resource string            `json:"resource"`
	Method   permission.Method `json:"method"`
}

type CreatePermission struct {
	Name    string         `json:"name"`
	Actions []CreateAction `json:"actions"`
}

type CreateGroup struct {
	Name        string `json:"name"`
	Permissions []CreatePermission
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
	Name    string           `json:"name"`
	Actions []ResponseAction `json:"actions"`
}

type ListResponse struct {
	utils.Paginator
	Data []ResponseGroup `json:"data"`
}

func ToResponse(ent *entity.Group) *ResponseGroup {
	var res *ResponseGroup
	if ent == nil {
		return res
	}

	res.Name = ent.Name
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
