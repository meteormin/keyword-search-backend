package permission

import (
	"fmt"
	"github.com/miniyus/gofiber/app"
	cLog "github.com/miniyus/gofiber/log"
	"github.com/miniyus/gollection"
	"github.com/miniyus/gorm-extension/gormrepo"
	"github.com/miniyus/keyword-search-backend/entity"
	"github.com/miniyus/keyword-search-backend/repo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	PATCH  Method = "PATCH"
	DELETE Method = "DELETE"
)

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
	gollection.Collection[Permission]
	All() []Permission
	GetByName(name string) (*Permission, error)
	RemoveByName(name string) bool
}

type CollectionStruct struct {
	*gollection.BaseCollection[Permission]
	items []Permission
}

func NewPermissionCollection(perms ...Permission) Collection {
	defaultPerms := make([]Permission, 0)
	if len(perms) == 0 {
		perms = defaultPerms
	}

	base := gollection.NewCollection(perms).(*gollection.BaseCollection[Permission])

	return &CollectionStruct{BaseCollection: base}
}

func (p *CollectionStruct) All() []Permission {
	return p.items
}

func (p *CollectionStruct) RemoveByName(name string) bool {
	filtered := p.Filter(func(v Permission, i int) bool {
		return v.Name == name
	})

	if len(filtered.Items()) == 0 {
		return false
	}

	var rmIndex int
	for i, perm := range p.items {
		if perm.Name == filtered.Items()[0].Name {
			rmIndex = i
		}
	}

	err := p.Remove(rmIndex)
	if err != nil {
		return false
	}

	return true
}

func (p *CollectionStruct) GetByName(name string) (*Permission, error) {
	filtered := p.Filter(func(v Permission, i int) bool {
		return v.Name == name
	})

	if filtered.Count() == 0 {
		return nil, fmt.Errorf("can't found %s Permission", name)
	}

	return &filtered.Items()[0], nil
}

func (perm Permission) ToEntity() entity.Permission {
	var ent entity.Permission
	ent.Permission = perm.Name
	ent.GroupId = perm.GroupId
	for _, action := range perm.Actions {
		for _, method := range action.Methods {
			ent.Actions = append(ent.Actions, entity.Action{
				Resource: action.Resource,
				Method:   string(method),
			})
		}
	}

	return ent
}

func (perm Permission) FromEntity(permission entity.Permission) Permission {
	actions := make([]Action, 0)
	gollection.NewCollection(permission.Actions).For(func(v entity.Action, i int) {
		filtered := gollection.NewCollection(permission.Actions).Filter(func(a entity.Action, j int) bool {
			return a.PermissionId == v.PermissionId && a.Resource == v.Resource
		})

		methods := make([]Method, 0)
		filtered.For(func(f entity.Action, k int) {
			methods = append(methods, Method(f.Method))
		})

		actions = append(actions, Action{
			Resource: v.Resource,
			Methods:  methods,
		})
	})

	return Permission{
		GroupId: permission.GroupId,
		Name:    permission.Permission,
		Actions: actions,
	}
}

func CreateDefaultPermissions(cfg []Config) app.Boot {
	pCfg := cfg
	return func(app app.Application) {
		var db *gorm.DB
		err := app.Resolve(&db)
		if err != nil {
			log.Println(err)
		}

		perms := NewPermissionsFromConfig(pCfg)
		permCollection := NewPermissionCollection(perms...)

		repository := repo.NewPermissionRepository(db)
		var entities []entity.Permission

		permCollection.For(func(perm Permission, i int) {
			entities = append(entities, perm.ToEntity())
		})

		all, err := repository.All()
		if err != nil {
			cLog.GetLogger().Error(err)
			log.Println(err)
		}

		if len(all) != 0 {
			actionRepo := gormrepo.NewGenericRepository(db, entity.Action{})
			for _, ent := range entities {
				canInsert := false
				for _, perm := range all {
					canInsert = perm.Permission == ent.Permission && perm.GroupId == ent.GroupId
					if canInsert {
						for i, action := range ent.Actions {
							action.PermissionId = perm.ID
							ent.Actions[i] = action
						}
						break
					}
				}
				if canInsert {
					actions := ent.Actions
					err = actionRepo.DB().Transaction(func(tx *gorm.DB) error {
						return tx.Clauses(clause.OnConflict{
							Columns: []clause.Column{
								{Name: "permission_id"},
								{Name: "method"},
								{Name: "resource"},
							},
							DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
						}).Create(&actions).Error
					})
					if err != nil {
						cLog.GetLogger().Error(err)
						log.Println(ent.Actions)
						log.Println(err)
					}
				}
			}
			return
		}

		_, err = repository.BatchCreate(entities)
		if err != nil {
			if err != nil {
				cLog.GetLogger().Error(err)
				log.Println(err)
			}
		}
	}
}
