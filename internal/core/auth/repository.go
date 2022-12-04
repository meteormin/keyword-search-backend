package auth

import (
	"github.com/miniyus/go-fiber/database"
	"github.com/miniyus/go-fiber/internal/entity"
	"gorm.io/gorm"
)

type Repository interface {
	All() ([]*entity.AccessToken, error)
	Create(token entity.AccessToken) (*entity.AccessToken, error)
	Find(pk uint) (*entity.AccessToken, error)
	FindByToken(token string) (*entity.AccessToken, error)
	FindByUserId(userId uint) (*entity.AccessToken, error)
	Update(token entity.AccessToken) (*entity.AccessToken, error)
	Delete(pk uint) (bool, error)
}

type RepositoryStruct struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *RepositoryStruct {
	return &RepositoryStruct{db}
}

func (repo *RepositoryStruct) All() ([]*entity.AccessToken, error) {
	var tokens []*entity.AccessToken
	result := repo.db.Find(&tokens)
	_, err := database.HandleResult(result)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (repo *RepositoryStruct) Create(token entity.AccessToken) (*entity.AccessToken, error) {
	result := repo.db.Create(&token)

	_, err := database.HandleResult(result)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (repo *RepositoryStruct) Find(pk uint) (*entity.AccessToken, error) {
	token := entity.AccessToken{}
	result := repo.db.Find(&token, pk)

	_, err := database.HandleResult(result)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (repo *RepositoryStruct) Update(token entity.AccessToken) (*entity.AccessToken, error) {
	result := repo.db.Save(&token)

	_, err := database.HandleResult(result)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (repo *RepositoryStruct) Delete(pk uint) (bool, error) {
	token, err := repo.Find(pk)

	if token != nil && err == nil {
		result := repo.db.Delete(&token)
		_, err := database.HandleResult(result)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, err
}

func (repo *RepositoryStruct) FindByToken(token string) (*entity.AccessToken, error) {
	ent := entity.AccessToken{}

	result := repo.db.Where(&entity.AccessToken{
		Token: token,
	}).First(&ent)

	_, err := database.HandleResult(result)
	if err != nil {
		return nil, err
	}

	return &ent, nil
}

func (repo *RepositoryStruct) FindByUserId(userId uint) (*entity.AccessToken, error) {
	ent := entity.AccessToken{}

	result := repo.db.Where(&entity.AccessToken{
		UserId: userId,
	}).Last(&ent)

	_, err := database.HandleResult(result)
	if err != nil {
		return nil, err
	}

	return &ent, nil
}
