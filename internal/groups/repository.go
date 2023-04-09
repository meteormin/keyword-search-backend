package groups

import (
	"github.com/miniyus/gofiber/pagination"
	"github.com/miniyus/gorm-extension/gormrepo"
	"github.com/miniyus/keyword-search-backend/entity"
	"gorm.io/gorm"
)

type Repository interface {
	gormrepo.GenericRepository[entity.Group]
	Count(group entity.Group) (int64, error)
	AllWithPage(page pagination.Page) ([]entity.Group, int64, error)
	FindByName(groupName string) (*entity.Group, error)
}

type RepositoryStruct struct {
	gormrepo.GenericRepository[entity.Group]
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryStruct{gormrepo.NewGenericRepository(db, entity.Group{})}
}

func (r *RepositoryStruct) Count(group entity.Group) (int64, error) {
	var count int64 = 0
	err := r.DB().Transaction(func(tx *gorm.DB) error {
		return tx.Model(&group).Count(&count).Error
	})

	if err != nil {
		return 0, err
	}

	return count, err
}

func (r *RepositoryStruct) AllWithPage(page pagination.Page) ([]entity.Group, int64, error) {
	var groups []entity.Group

	count, err := r.Count(entity.Group{})

	if count != 0 {
		err = r.DB().Model(&entity.Group{}).
			Preload("Permissions.Actions").
			Scopes(pagination.Paginate(page)).
			Order("id desc").
			Find(&groups).Error

		if err != nil {
			return make([]entity.Group, 0), 0, err
		}
	} else {
		return make([]entity.Group, 0), 0, nil
	}

	return groups, count, err
}

func (r *RepositoryStruct) FindByName(groupName string) (*entity.Group, error) {
	group := &entity.Group{}

	if err := r.DB().Preload("Permissions.Actions").Where(entity.Group{Name: groupName}).First(group).Error; err != nil {
		return nil, err
	}

	return group, nil
}
