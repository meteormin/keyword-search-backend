package searchs

import "gorm.io/gorm"

type Repository interface {
}

type RepositoryStruct struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *RepositoryStruct {
	return &RepositoryStruct{
		db: db,
	}
}
