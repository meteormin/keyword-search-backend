package login_logs

import (
	"github.com/miniyus/keyword-search-backend/entity"
	"gorm.io/gorm"
)

type Repository interface {
	Create(log entity.LoginLog) (*entity.LoginLog, error)
	GetByUserId(userId uint) ([]entity.LoginLog, error)
}

type RepositoryStruct struct {
	db *gorm.DB
}

func (r RepositoryStruct) Create(log entity.LoginLog) (*entity.LoginLog, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&log).Error
	})

	if err != nil {
		return nil, err
	}

	return &log, nil
}

func (r RepositoryStruct) GetByUserId(userId uint) ([]entity.LoginLog, error) {
	var logs []entity.LoginLog
	if err := r.db.Where(&entity.LoginLog{UserId: userId}).Find(&logs).Error; err != nil {
		return make([]entity.LoginLog, 0), err
	}

	return logs, nil
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryStruct{
		db: db,
	}
}
