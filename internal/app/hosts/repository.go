package hosts

import (
	"github.com/miniyus/go-fiber/internal/core/database"
	"github.com/miniyus/go-fiber/internal/core/logger"
	"github.com/miniyus/go-fiber/internal/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository interface {
	All() ([]entity.Host, error)
	GetByUserId(userId uint) ([]entity.Host, error)
	Find(pk uint) (*entity.Host, error)
	Create(host entity.Host) (*entity.Host, error)
	Update(pk uint, host entity.Host) (*entity.Host, error)
	Delete(pk uint) (bool, error)
	logger.HasLogger
}

type RepositoryStruct struct {
	db *gorm.DB
	logger.HasLoggerStruct
}

func NewRepository(db *gorm.DB, log *zap.SugaredLogger) Repository {
	return &RepositoryStruct{
		db:              db,
		HasLoggerStruct: logger.HasLoggerStruct{Logger: log},
	}
}

func (r *RepositoryStruct) All() ([]entity.Host, error) {
	var hosts []entity.Host

	result := r.db.Find(&hosts)
	_, err := database.HandleResult(result)
	if err != nil {
		return nil, err
	}

	return hosts, nil
}

func (r *RepositoryStruct) GetByUserId(userId uint) ([]entity.Host, error) {
	var hosts []entity.Host
	result := r.db.Where(entity.Host{UserId: userId}).Find(&hosts)
	_, err := database.HandleResult(result)
	if err != nil {
		return hosts, err
	}

	return hosts, nil
}

func (r *RepositoryStruct) Find(pk uint) (*entity.Host, error) {
	host := entity.Host{}
	result := r.db.Preload("Search").Find(&host, pk)
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

	if exists == nil {
		return nil, gorm.ErrRecordNotFound
	}

	if host.ID == exists.ID { // patch
		result := r.db.Save(&host)
		_, err = database.HandleResult(result)
		if err != nil {
			return nil, err
		}
	} else { // put
		host.ID = exists.ID
		result := r.db.Save(&host)
		_, err = database.HandleResult(result)
	}

	return &host, nil
}

func (r *RepositoryStruct) Delete(pk uint) (bool, error) {
	exists, err := r.Find(pk)
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
