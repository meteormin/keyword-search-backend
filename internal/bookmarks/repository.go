package bookmarks

import (
	"github.com/miniyus/gofiber/database"
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
	var marks []entity.BookMark

	result := repo.db.Find(&marks)
	_, err := database.HandleResult(result)
	if err != nil {
		return nil, err
	}

	return marks, nil
}

func (repo *RepositoryStruct) Find(pk uint) (*entity.BookMark, error) {
	mark := &entity.BookMark{}
	result := repo.db.Find(mark, pk)
	_, err := database.HandleResult(result)
	if err != nil {
		return nil, err
	}

	return mark, nil
}

func (repo *RepositoryStruct) Create(mark *entity.BookMark) (*entity.BookMark, error) {
	result := repo.db.Create(mark)
	_, err := database.HandleResult(result)
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

	result := repo.db.Save(mark)
	_, err = database.HandleResult(result)

	return mark, err
}

func (repo *RepositoryStruct) Delete(pk uint) (bool, error) {
	exists, err := repo.Find(pk)
	if err != nil {
		return false, err
	}

	result := repo.db.Delete(exists)
	_, err = database.HandleResult(result)
	if err != nil {
		return false, err
	}

	return true, nil
}
