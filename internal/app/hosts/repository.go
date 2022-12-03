package hosts

import (
	"github.com/miniyus/go-fiber/database"
	"github.com/miniyus/go-fiber/internal/entity"
	"gorm.io/gorm"
)

type Repository interface {
	All() ([]*entity.Host, error)
	Find(pk uint) (*entity.Host, error)
	Create(host entity.Host) (*entity.Host, error)
	Update(pk uint, host entity.Host) (*entity.Host, error)
	Delete(pk uint) (bool, error)
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

func (r *RepositoryStruct) Find(pk uint) (*entity.Host, error) {
	host := entity.Host{}
	result := r.db.Find(&host, pk)
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

func (r *RepositoryStruct) Update(pk uint, host entity.Host) (*entity.Host, error) {
	exists, err := r.Find(pk)
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

func (r *RepositoryStruct) Delete(pk uint) (bool, error) {
	exists, err := r.Find(pk)
	if err != nil {
		return false, nil
	}

	result := r.db.Delete(exists)
	_, err = database.HandleResult(result)
	if err != nil {
		return false, nil
	}

	return true, nil
}
