package test_api

import "gorm.io/gorm"

type Repository interface {
	create()
}

type RepositoryStruct struct {
	model *gorm.DB
}

func NewRepository(db *gorm.DB) RepositoryStruct {
	return RepositoryStruct{db}
}

func (repo RepositoryStruct) create() {

}
