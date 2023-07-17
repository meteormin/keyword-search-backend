package repo

import (
	"github.com/miniyus/gorm-extension/gormrepo"
	"github.com/miniyus/keyword-search-backend/entity"
	"gorm.io/gorm"
)

type LoginLogRepository interface {
	gormrepo.GenericRepository[entity.LoginLog]
	Create(log entity.LoginLog) (*entity.LoginLog, error)
	GetByUserId(userId uint) ([]entity.LoginLog, error)
}

type LoginLogRepositoryStruct struct {
	gormrepo.GenericRepository[entity.LoginLog]
}

func (r LoginLogRepositoryStruct) GetByUserId(userId uint) ([]entity.LoginLog, error) {
	return r.GetByEntity(entity.LoginLog{UserId: userId})
}

func NewLoginLogRepository(db *gorm.DB) LoginLogRepository {
	return &LoginLogRepositoryStruct{
		gormrepo.NewGenericRepository(db, entity.LoginLog{}),
	}
}
