package group_detail

import (
	"github.com/miniyus/keyword-search-backend/entity"
	"gorm.io/gorm"
)

type Repository interface {
	All() ([]entity.GroupDetail, error)
	Find(pk uint) (*entity.GroupDetail, error)
	Create(detail entity.GroupDetail) (*entity.GroupDetail, error)
	Update(pk uint, detail entity.GroupDetail) (*entity.GroupDetail, error)
	Delete(pk uint) (bool, error)
	GetByGroupId(groupId uint) ([]entity.GroupDetail, error)
	GetByUserId(userId uint) ([]entity.GroupDetail, error)
	FindByOwner(groupId uint, userId uint) (*entity.GroupDetail, error)
}

type RepositoryStruct struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryStruct{
		db,
	}
}

func (repo *RepositoryStruct) All() ([]entity.GroupDetail, error) {
	models := make([]entity.GroupDetail, 0)

	if err := repo.db.Find(&models).Error; err != nil {
		return make([]entity.GroupDetail, 0), err
	}

	return models, nil
}

func (repo *RepositoryStruct) Find(pk uint) (*entity.GroupDetail, error) {
	var model entity.GroupDetail
	if err := repo.db.First(&model, pk).Error; err != nil {
		return nil, err
	}

	return &model, nil
}

func (repo *RepositoryStruct) Create(detail entity.GroupDetail) (*entity.GroupDetail, error) {
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&detail).Error
	})

	if err != nil {
		return nil, err
	}

	return &detail, err
}

func (repo *RepositoryStruct) Update(pk uint, detail entity.GroupDetail) (*entity.GroupDetail, error) {
	exists, err := repo.Find(pk)
	if err != nil {
		return nil, err
	}

	err = repo.db.Transaction(func(tx *gorm.DB) error {
		if detail.ID == exists.ID {
			return tx.Save(&detail).Error
		} else {
			detail.ID = exists.ID
			return tx.Save(&detail).Error
		}

	})

	if err != nil {
		return nil, err
	}

	return &detail, nil
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

	return true, err
}

func (repo *RepositoryStruct) GetByGroupId(groupId uint) ([]entity.GroupDetail, error) {
	groupDetails := make([]entity.GroupDetail, 0)
	if err := repo.db.Where(&entity.GroupDetail{GroupId: groupId}).Find(&groupDetails).Error; err != nil {
		return make([]entity.GroupDetail, 0), nil
	}

	return groupDetails, nil
}

func (repo *RepositoryStruct) GetByUserId(userId uint) ([]entity.GroupDetail, error) {
	groupDetails := make([]entity.GroupDetail, 0)
	if err := repo.db.Where(&entity.GroupDetail{UserId: userId}).Find(&groupDetails).Error; err != nil {
		return make([]entity.GroupDetail, 0), nil
	}

	return groupDetails, nil
}

func (repo *RepositoryStruct) FindByOwner(groupId uint, userId uint) (*entity.GroupDetail, error) {
	var groupDetail entity.GroupDetail
	err := repo.db.Where(&entity.GroupDetail{
		GroupId: groupId,
		UserId:  userId,
		Role:    entity.Owner,
	}).First(groupDetail).Error

	if err != nil {
		return nil, err
	}

	return &groupDetail, err
}
