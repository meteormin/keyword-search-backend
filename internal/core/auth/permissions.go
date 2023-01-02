package auth

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/utils"
	"strings"
)

type Action struct {
	Method   config.Methods
	Resource string
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

type PermissionCollection interface {
	Add(perm Permission)
	Remove(name string) bool
	Get(name string) (*Permission, error)
	Filter(fn func(p Permission, i int) bool) []Permission
	Map(fn func(p Permission, i int) Permission) []Permission
}

type PermissionCollectionStruct struct {
	permissions []Permission
}

func NewPermissionCollection(perms ...Permission) PermissionCollection {
	defaultPerms := make([]Permission, 0)
	if len(perms) == 0 {
		perms = defaultPerms
	}

	return &PermissionCollectionStruct{perms}
}

func (p *PermissionCollectionStruct) Add(perm Permission) {
	p.permissions = append(p.permissions, perm)
}

func (p *PermissionCollectionStruct) Remove(name string) bool {
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

func (p *PermissionCollectionStruct) Get(name string) (*Permission, error) {
	filtered := utils.Filter(p.permissions, func(v Permission, i int) bool {
		return v.Name == name
	})

	if len(filtered) == 0 {
		return nil, fmt.Errorf("can't found %s Permission", name)

	}

	return &filtered[0], nil
}

func (p *PermissionCollectionStruct) Filter(fn func(p Permission, i int) bool) []Permission {
	return utils.Filter(p.permissions, fn)
}

func (p *PermissionCollectionStruct) Map(fn func(p Permission, i int) Permission) []Permission {
	return utils.Map(p.permissions, fn)
}

func CheckPermissionFromCtx(hasPerm []Permission, c *fiber.Ctx) bool {
	pass := false
	for _, perm := range hasPerm {
		for _, action := range perm.Actions {
			routePath := c.Path()
			if strings.Contains(routePath, action.Resource) {
				if action.Method.ToString() == c.Method() {
					pass = true
				}
			}
		}
	}
	return pass
}

func NewPermissionsFromConfig(permission config.Permission) {
	for _, p := range permission {
		permission.Resources
	}
}
