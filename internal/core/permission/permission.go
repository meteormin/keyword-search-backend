package permission

import (
	"fmt"
	"github.com/miniyus/keyword-search-backend/internal/utils"
)

type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	PATCH  Method = "PATCH"
	DELETE Method = "DELETE"
)

func (m Method) ToString() string {
	return string(m)
}

type Action struct {
	Methods  []Method
	Resource string
}

func NewAction(resource string, method []Method) Action {
	return Action{method, resource}
}

type Permission struct {
	GroupId uint
	Name    string
	Actions []Action
}

func NewPermission(groupId uint, name string, actions []Action) Permission {
	return Permission{
		groupId, name, actions,
	}
}

type user struct {
	Id      *uint `json:"id"`
	GroupId *uint `json:"group_id"`
	Role    *uint `json:"role"`
}

type includeUser struct {
	UserId  *uint `json:"user_id"`
	GroupId *uint `json:"group_id"`
	User    *user `json:"user"`
}

type Collection interface {
	utils.Collection[Permission]
	All() []Permission
	RemoveByName(name string) bool
	Get(name string) (*Permission, error)
}

type CollectionStruct struct {
	*utils.BaseCollection[Permission]
	items []Permission
}

func NewPermissionCollection(perms ...Permission) Collection {
	defaultPerms := make([]Permission, 0)
	if len(perms) == 0 {
		perms = defaultPerms
	}

	base := utils.NewCollection(perms).(*utils.BaseCollection[Permission])

	return &CollectionStruct{BaseCollection: base}
}

func (p *CollectionStruct) All() []Permission {
	return p.items
}

func (p *CollectionStruct) RemoveByName(name string) bool {
	filtered := p.Filter(func(v Permission, i int) bool {
		return v.Name == name
	})

	if len(filtered) == 0 {
		return false
	}

	var rmIndex int
	for i, perm := range p.items {
		if perm.Name == filtered[0].Name {
			rmIndex = i
		}
	}

	p.Remove(rmIndex)

	return true
}

func (p *CollectionStruct) Get(name string) (*Permission, error) {
	filtered := p.Filter(func(v Permission, i int) bool {
		return v.Name == name
	})

	if len(filtered) == 0 {
		return nil, fmt.Errorf("can't found %s Permission", name)

	}

	return &filtered[0], nil
}
