package loginlogs

import (
	"github.com/miniyus/gofiber/pkg/gormrepo"
	"github.com/miniyus/keyword-search-backend/entity"
	"gorm.io/gorm"
)

type Repository interface {
	gormrepo.GenericRepository[entity.LoginLog]
	Create(log entity.LoginLog) (*entity.LoginLog, error)
	GetByUserId(userId uint) ([]entity.LoginLog, error)
}

type RepositoryStruct struct {
	gormrepo.GenericRepository[entity.LoginLog]
}

func (r RepositoryStruct) GetByUserId(userId uint) ([]entity.LoginLog, error) {
	return r.GetByEntity(entity.LoginLog{UserId: userId})
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryStruct{
		gormrepo.NewGenericRepository(db, entity.LoginLog{}),
	}
}
