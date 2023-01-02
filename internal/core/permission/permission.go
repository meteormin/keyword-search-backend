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

func NewAction(method []Method, resource string) Action {
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
	Add(perm Permission)
	Remove(name string) bool
	Get(name string) (*Permission, error)
	Filter(fn func(p Permission, i int) bool) []Permission
	Map(fn func(p Permission, i int) Permission) []Permission
	Concat(perms []Permission)
}

type CollectionStruct struct {
	permissions []Permission
}

func NewPermissionCollection(perms ...Permission) Collection {
	defaultPerms := make([]Permission, 0)
	if len(perms) == 0 {
		perms = defaultPerms
	}

	return &CollectionStruct{perms}
}

func (p *CollectionStruct) Add(perm Permission) {
	p.permissions = append(p.permissions, perm)
}

func (p *CollectionStruct) Remove(name string) bool {
	filtered := utils.Filter(p.permissions, func(v Permission, i int) bool {
		return v.Name == name
	})

	if len(filtered) == 0 {
		return false
	}

	var rmIndex int
	for i, perm := range p.permissions {
		if perm.Name == filtered[0].Name {
			rmIndex = i
		}
	}

	slice := p.permissions
	p.permissions = append(slice[:rmIndex], slice[rmIndex+1:]...)

	return true
}

func (p *CollectionStruct) Get(name string) (*Permission, error) {
	filtered := utils.Filter(p.permissions, func(v Permission, i int) bool {
		return v.Name == name
	})

	if len(filtered) == 0 {
		return nil, fmt.Errorf("can't found %s Permission", name)

	}

	return &filtered[0], nil
}

func (p *CollectionStruct) Filter(fn func(p Permission, i int) bool) []Permission {
	return utils.Filter(p.permissions, fn)
}

func (p *CollectionStruct) Map(fn func(p Permission, i int) Permission) []Permission {
	return utils.Map(p.permissions, fn)
}

func (p *CollectionStruct) Concat(perms []Permission) {
	if len(perms) == 0 {
		return
	}

	for _, perm := range perms {
		p.Add(perm)
	}

}
