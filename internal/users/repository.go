package users

import "gorm.io/gorm"

type Repository interface {
	create()
}

type RepositoryImpl struct {
	model *gorm.DB
}

func NewRepository(db *gorm.DB) *RepositoryImpl {
	return &RepositoryImpl{db}
}

func (repo *RepositoryImpl) create() {

}
