package repo

import (
	"github.com/miniyus/gorm-extension/gormrepo"
	"github.com/miniyus/keyword-search-backend/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PermissionRepository interface {
	gormrepo.GenericRepository[entity.Permission]
	BatchCreate(permission []entity.Permission) ([]entity.Permission, error)
	GetByGroupId(groupId uint) ([]entity.Permission, error)
}

type PermissionRepositoryStruct struct {
	gormrepo.GenericRepository[entity.Permission]
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &PermissionRepositoryStruct{
		gormrepo.NewGenericRepository(db, entity.Permission{}),
	}
}

func (r *PermissionRepositoryStruct) BatchCreate(permission []entity.Permission) ([]entity.Permission, error) {
	err := r.DB().Transaction(func(tx *gorm.DB) error {
		return tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "permission"},
				{Name: "group_id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
		}).Create(&permission).Error
	})

	if err != nil {
		return make([]entity.Permission, 0), err
	}

	return permission, nil
}

func (r *PermissionRepositoryStruct) GetByGroupId(groupId uint) ([]entity.Permission, error) {
	permissions := make([]entity.Permission, 0)
	err := r.DB().Preload("Actions").Where(entity.Permission{GroupId: groupId}).Find(&permissions).Error

	return permissions, err
}
