package hosts

import (
	"github.com/miniyus/go-fiber/internal/core/database"
	"github.com/miniyus/go-fiber/internal/entity"
	"gorm.io/gorm"
)

type Repository interface {
	All() ([]*entity.Host, error)
	GetByUserId(userId uint) ([]*entity.Host, error)
	Find(pk uint, userId uint) (*entity.Host, error)
	Create(host entity.Host) (*entity.Host, error)
	Update(pk uint, userId uint, host entity.Host) (*entity.Host, error)
	Delete(pk uint, userId uint) (bool, error)
}

type RepositoryStruct struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *RepositoryStruct {
	return &RepositoryStruct{db: db}
}

func (r *RepositoryStruct) All() ([]*entity.Host, error) {
	var hosts []*entity.Host

	result := r.db.Find(&hosts)
	_, err := database.HandleResult(result)
	if err != nil {
		return nil, err
	}

	return hosts, nil
}

func (r *RepositoryStruct) GetByUserId(userId uint) ([]*entity.Host, error) {
	var hosts []*entity.Host
	result := r.db.Where(entity.Host{UserId: userId}).Find(&hosts)
	_, err := database.HandleResult(result)
	if err != nil {
		return hosts, err
	}

	return hosts, nil
}

func (r *RepositoryStruct) Find(pk uint, userId uint) (*entity.Host, error) {
	host := entity.Host{}
	result := r.db.Preload("Search").Where(entity.Host{UserId: userId}).Find(&host, pk)
	_, err := database.HandleResult(result)

	if err != nil {
		return nil, err
	}

	return &host, nil
}

func (r *RepositoryStruct) Create(host entity.Host) (*entity.Host, error) {
	result := r.db.Create(&host)
	_, err := database.HandleResult(result)

	if err != nil {
		return nil, err
	}

	return &host, nil
}

func (r *RepositoryStruct) Update(pk uint, userId uint, host entity.Host) (*entity.Host, error) {
	exists, err := r.Find(pk, userId)
	if err != nil {
		return nil, err
	}

	host.ID = exists.ID
	result := r.db.Save(&host)
	_, err = database.HandleResult(result)
	if err != nil {
		return nil, err
	}

	return &host, nil
}

func (r *RepositoryStruct) Delete(pk uint, userId uint) (bool, error) {
	exists, err := r.Find(pk, userId)
	if err != nil {
		return false, err
	}

	result := r.db.Delete(exists)
	_, err = database.HandleResult(result)
	if err != nil {
		return false, err
	}

	return true, nil
}
