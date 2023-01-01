package hosts

import (
	"github.com/miniyus/keyword-search-backend/internal/core/database"
	"github.com/miniyus/keyword-search-backend/internal/core/logger"
	"github.com/miniyus/keyword-search-backend/internal/entity"
	"github.com/miniyus/keyword-search-backend/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository interface {
	All(page utils.Page) ([]entity.Host, int64, error)
	GetByUserId(userId uint, page utils.Page) ([]entity.Host, int64, error)
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

func (r *RepositoryStruct) Count(host *entity.Host) (int64, error) {
	var count int64 = 0

	rs := r.db.Model(host).Count(&count)
	_, err := database.HandleResult(rs)

	return count, err
}

func (r *RepositoryStruct) All(page utils.Page) ([]entity.Host, int64, error) {
	var hosts []entity.Host
	count, err := r.Count(&entity.Host{})

	if count != 0 {
		result := r.db.Scopes(utils.Paginate(page)).Find(&hosts)
		_, err = database.HandleResult(result)
	}

	return hosts, count, err
}

func (r *RepositoryStruct) GetByUserId(userId uint, page utils.Page) (host []entity.Host, count int64, e error) {
	var hosts []entity.Host
	var cnt int64 = 0

	result := r.db.Where(entity.Host{UserId: userId}).Find(&hosts).Count(&cnt)
	_, err := database.HandleResult(result)

	if cnt == 0 {
		return make([]entity.Host, 0), cnt, err
	}

	result = r.db.Scopes(utils.Paginate(page)).Where(entity.Host{UserId: userId}).Find(&hosts)
	_, err = database.HandleResult(result)

	return hosts, cnt, err
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
