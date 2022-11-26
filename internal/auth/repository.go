package auth

import (
	"github.com/miniyus/go-fiber/internal/entity"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func (repo *Repository) Create(token entity.AccessToken) (*entity.AccessToken, error) {
	result := repo.db.Create(&token)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected != 0 {
		return &token, nil
	}

	return &token, nil
}

func (repo *Repository) Find(pk uint) (*entity.AccessToken, error) {
	token := entity.AccessToken{}
	result := repo.db.Find(&token, pk)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected != 0 {
		return &token, nil
	}

	return nil, nil
}

func (repo *Repository) Update(token entity.AccessToken) (*entity.AccessToken, error) {
	result := repo.db.Save(&token)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected != 0 {
		return &token, nil
	}

	return nil, nil
}

func (repo *Repository) Delete(pk uint) (bool, error) {
	token, err := repo.Find(pk)

	if token != nil && err == nil {
		result := repo.db.Delete(&token)
		if result.Error != nil {
			return false, err
		}

		if result.RowsAffected != 0 {
			return true, err
		}

	}

	return false, err
}
