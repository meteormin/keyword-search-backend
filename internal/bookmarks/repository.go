package bookmarks

import (
	"github.com/miniyus/keyword-search-backend/entity"
	"gorm.io/gorm"
)

type Repository interface {
	All() ([]entity.BookMark, error)
	Find(pk uint) (*entity.BookMark, error)
	Create(mark entity.BookMark) (*entity.BookMark, error)
	Update(pk uint, mark entity.BookMark) (*entity.BookMark, error)
	Delete(pk uint) (bool, error)
}

type RepositoryStruct struct {
	db *gorm.DB
}

func (repo *RepositoryStruct) All() ([]entity.BookMark, error) {
	marks := make([]entity.BookMark, 0)

	if err := repo.db.Find(&marks).Error; err != nil {
		return nil, err
	}

	return marks, nil
}

func (repo *RepositoryStruct) Find(pk uint) (*entity.BookMark, error) {
	mark := &entity.BookMark{}
	err := repo.db.Find(mark, pk).Error

	if err != nil {
		return nil, err
	}

	return mark, nil
}

func (repo *RepositoryStruct) Create(mark *entity.BookMark) (*entity.BookMark, error) {
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(mark).Error
	})

	if err != nil {
		return nil, err
	}

	return mark, nil
}

func (repo *RepositoryStruct) Update(pk uint, mark *entity.BookMark) (*entity.BookMark, error) {
	exists, err := repo.Find(pk)
	if err != nil {
		return nil, err
	}
	mark.ID = exists.ID

	err = repo.db.Transaction(func(tx *gorm.DB) error {
		return tx.Save(mark).Error
	})

	if err != nil {
		return nil, err
	}

	return mark, nil
}

func (repo *RepositoryStruct) Delete(pk uint) (bool, error) {
	exists, err := repo.Find(pk)
	if err != nil {
		return false, err
	}

	err = repo.db.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(exists).Error
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
