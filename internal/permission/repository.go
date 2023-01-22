package permission

import (
	"github.com/miniyus/keyword-search-backend/internal/database"
	"github.com/miniyus/keyword-search-backend/internal/entity"
	"github.com/miniyus/keyword-search-backend/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Save(permission []entity.Permission) ([]entity.Permission, error)
	Get(groupId uint) ([]entity.Permission, error)
}

type RepositoryStruct struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryStruct{
		db: db,
	}
}

func (r *RepositoryStruct) Save(permission []entity.Permission) ([]entity.Permission, error) {
	actions := make([][]entity.Action, 0)

	permission = utils.NewCollection(permission).Map(func(v entity.Permission, i int) entity.Permission {
		if len(v.Actions) != 0 && v.Actions != nil {
			actions = append(actions, v.Actions)
			v.Actions = nil
		} else {
			actions = append(actions, nil)
		}
		return v
	})

	rs := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"permission"}),
	}).Create(&permission)
	_, err := database.HandleResult(rs)

	if err != nil {
		return make([]entity.Permission, 0), err
	}

	var createActions []entity.Action

	utils.NewCollection(permission).For(func(v entity.Permission, i int) {
		if actions[i] != nil {
			for _, action := range actions[i] {
				action.PermissionId = v.ID
				createActions = append(createActions, action)
			}
		}
	})

	rs = r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"method", "resource"}),
	}).Create(&createActions)
	_, err = database.HandleResult(rs)

	if err != nil {
		return permission, err
	}

	created := utils.NewCollection(permission).Map(func(perm entity.Permission, i int) entity.Permission {
		permActions := utils.NewCollection(createActions).Filter(func(action entity.Action, j int) bool {
			return action.PermissionId == perm.ID
		})
		perm.Actions = permActions

		return perm
	})

	return created, nil
}

func (r *RepositoryStruct) Get(groupId uint) ([]entity.Permission, error) {
	permissions := make([]entity.Permission, 0)

	rs := r.db.Preload("Actions").Where(entity.Permission{GroupId: groupId}).Find(&permissions)
	_, err := database.HandleResult(rs)

	return permissions, err
}
