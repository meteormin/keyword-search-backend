package users

import (
	"github.com/miniyus/gofiber/database"
	"github.com/miniyus/gofiber/entity"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user entity.User) (*entity.User, error)
	Find(pk uint) (*entity.User, error)
	All() ([]entity.User, error)
	Update(pk uint, user entity.User) (*entity.User, error)
	Delete(pk uint) (bool, error)
	FindByUsername(username string) (*entity.User, error)
	FindByEntity(user entity.User) (*entity.User, error)
}

type RepositoryStruct struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryStruct{db}
}

func (repo *RepositoryStruct) Create(user entity.User) (*entity.User, error) {
	result := repo.db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &user, nil
}

func (repo *RepositoryStruct) Find(pk uint) (*entity.User, error) {
	user := entity.User{}

	result := repo.db.First(&user, pk)
	_, err := database.HandleResult(result)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *RepositoryStruct) All() ([]entity.User, error) {
	var users []entity.User
	result := repo.db.Find(&users)
	_, err := database.HandleResult(result)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *RepositoryStruct) Update(pk uint, user entity.User) (*entity.User, error) {
	exists, err := repo.Find(pk)
	if err != nil {
		return nil, err
	}

	user.ID = exists.ID

	result := repo.db.Save(&user)

	_, err = database.HandleResult(result)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *RepositoryStruct) Delete(pk uint) (bool, error) {

	user, err := repo.Find(pk)
	if err != nil {
		return false, err
	}

	result := repo.db.Delete(user)

	_, err = database.HandleResult(result)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (repo *RepositoryStruct) FindByUsername(username string) (*entity.User, error) {
	var user entity.User

	result := repo.db.Where(&entity.User{Username: username}).First(&user)
	_, err := database.HandleResult(result)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *RepositoryStruct) FindByEntity(user entity.User) (*entity.User, error) {
	var rsUser entity.User
	result := repo.db.Where(&user).First(&rsUser)

	_, err := database.HandleResult(result)
	if err != nil {
		return nil, err
	}

	return &rsUser, nil
}
