package groups

import (
	"github.com/miniyus/keyword-search-backend/internal/core/database"
	"github.com/miniyus/keyword-search-backend/internal/core/logger"
	"github.com/miniyus/keyword-search-backend/internal/entity"
	"github.com/miniyus/keyword-search-backend/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository interface {
	Count(group entity.Group) (int64, error)
	All(page utils.Page) ([]entity.Group, int64, error)
	Create(group entity.Group) (*entity.Group, error)
	Update(pk uint, group entity.Group) (*entity.Group, error)
	Find(pk uint) (*entity.Group, error)
	FindByName(groupName string) (*entity.Group, error)
	Delete(pk uint) (bool, error)
	logger.HasLogger
}

type RepositoryStruct struct {
	db *gorm.DB
	logger.HasLoggerStruct
}

func NewRepository(db *gorm.DB, log *zap.SugaredLogger) Repository {
	return &RepositoryStruct{db, logger.HasLoggerStruct{Logger: log}}
}

func (r *RepositoryStruct) Count(group entity.Group) (int64, error) {
	var count int64 = 0
	rs := r.db.Model(&group).Count(&count)
	_, err := database.HandleResult(rs)

	return count, err
}

func (r *RepositoryStruct) All(page utils.Page) ([]entity.Group, int64, error) {
	var groups []entity.Group

	count, err := r.Count(entity.Group{})

	if count != 0 {
		result := r.db.Scopes(utils.Paginate(page)).Order("id desc").Find(&groups)
		_, err = database.HandleResult(result)
	}

	return groups, count, err
}

func (r *RepositoryStruct) Create(group entity.Group) (*entity.Group, error) {
	result := r.db.Create(&group)
	_, err := database.HandleResult(result)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (r *RepositoryStruct) Update(pk uint, group entity.Group) (*entity.Group, error) {
	exists, err := r.Find(pk)
	if err != nil {
		return nil, err
	}

	if exists == nil {
		return nil, gorm.ErrRecordNotFound
	}

	if group.ID == exists.ID {
		result := r.db.Save(&group)
		_, err = database.HandleResult(result)
		if err != nil {
			return nil, err
		}
	} else {
		group.ID = exists.ID
		result := r.db.Save(&group)
		_, err = database.HandleResult(result)
	}

	return &group, nil
}

func (r *RepositoryStruct) Find(pk uint) (*entity.Group, error) {
	group := entity.Group{}
	result := r.db.Preload("Permissions.Actions").First(&group, pk)
	_, err := database.HandleResult(result)

	if err != nil {
		return nil, err
	}

	return &group, err
}

func (r *RepositoryStruct) FindByName(groupName string) (*entity.Group, error) {
	group := &entity.Group{}

	result := r.db.Preload("Permissions.Actions").Where(entity.Group{Name: groupName}).First(group)
	_, err := database.HandleResult(result)
	if err != nil {
		return nil, err
	}

	return group, err
}

func (r *RepositoryStruct) Delete(pk uint) (bool, error) {
	exists, err := r.Find(pk)
	if err != nil {
		return false, err
	}

	result := r.db.Delete(exists)
	_, err = database.HandleResult(result)
	if err != nil {
		return false, err
	}

	return true, nil
}
