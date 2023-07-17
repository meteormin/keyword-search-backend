package repo

import (
	"github.com/miniyus/gorm-extension/gormrepo"
	"github.com/miniyus/keyword-search-backend/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	gormrepo.GenericRepository[entity.User]
	FindByUsername(username string) (*entity.User, error)
}

type UserRepositoryStruct struct {
	gormrepo.GenericRepository[entity.User]
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryStruct{
		gormrepo.NewGenericRepository[entity.User](db, entity.User{}),
	}
}

func (repo *UserRepositoryStruct) FindByUsername(username string) (*entity.User, error) {
	return repo.FindByEntity(entity.User{Username: username})
}
