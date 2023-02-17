package group_detail

import (
	"github.com/miniyus/keyword-search-backend/entity"
	"gorm.io/gorm"
)

type Repository interface {
	gormrepo.GenericRepository[entity.GroupDetail]
	GetByGroupId(groupId uint) ([]entity.GroupDetail, error)
	GetByUserId(userId uint) ([]entity.GroupDetail, error)
	FindByOwner(groupId uint, userId uint) (*entity.GroupDetail, error)
}

type RepositoryStruct struct {
	gormrepo.GenericRepository[entity.GroupDetail]
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryStruct{
		gormrepo.NewGenericRepository(db, entity.GroupDetail{}),
	}
}

func (repo *RepositoryStruct) GetByGroupId(groupId uint) ([]entity.GroupDetail, error) {
	return repo.GetByEntity(entity.GroupDetail{GroupId: groupId})
}

func (repo *RepositoryStruct) GetByUserId(userId uint) ([]entity.GroupDetail, error) {
	return repo.GetByEntity(entity.GroupDetail{UserId: userId})
}

func (repo *RepositoryStruct) FindByOwner(groupId uint, userId uint) (*entity.GroupDetail, error) {
	return repo.FindByEntity(entity.GroupDetail{
		GroupId: groupId,
		UserId:  userId,
		Role:    entity.Owner,
	})
}
