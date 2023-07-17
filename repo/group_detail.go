package repo

import (
	"github.com/miniyus/gorm-extension/gormrepo"
	"github.com/miniyus/keyword-search-backend/entity"
	"gorm.io/gorm"
)

type GroupDetailRepository interface {
	gormrepo.GenericRepository[entity.GroupDetail]
	GetByGroupId(groupId uint) ([]entity.GroupDetail, error)
	GetByUserId(userId uint) ([]entity.GroupDetail, error)
	FindByOwner(groupId uint, userId uint) (*entity.GroupDetail, error)
}

type GroupDetailRepositoryStruct struct {
	gormrepo.GenericRepository[entity.GroupDetail]
}

func NewGroupDetailRepository(db *gorm.DB) GroupDetailRepository {
	return &GroupDetailRepositoryStruct{
		gormrepo.NewGenericRepository(db, entity.GroupDetail{}),
	}
}

func (repo *GroupDetailRepositoryStruct) GetByGroupId(groupId uint) ([]entity.GroupDetail, error) {
	return repo.GetByEntity(entity.GroupDetail{GroupId: groupId})
}

func (repo *GroupDetailRepositoryStruct) GetByUserId(userId uint) ([]entity.GroupDetail, error) {
	return repo.GetByEntity(entity.GroupDetail{UserId: userId})
}

func (repo *GroupDetailRepositoryStruct) FindByOwner(groupId uint, userId uint) (*entity.GroupDetail, error) {
	return repo.FindByEntity(entity.GroupDetail{
		GroupId: groupId,
		UserId:  userId,
		Role:    entity.Owner,
	})
}
