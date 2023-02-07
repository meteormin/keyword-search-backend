package group_detail

import (
	"github.com/miniyus/gofiber/database"
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
	var models []entity.GroupDetail
	rs := repo.db.Find(&models)
	rs, err := database.HandleResult(rs)

	if rs.RowsAffected == 0 {
		return make([]entity.GroupDetail, 0), nil
	}

	return models, err
}

func (repo *RepositoryStruct) Find(pk uint) (*entity.GroupDetail, error) {
	var model entity.GroupDetail
	rs := repo.db.First(&model, pk)
	rs, err := database.HandleResult(rs)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (repo *RepositoryStruct) Create(detail entity.GroupDetail) (*entity.GroupDetail, error) {
	rs := repo.db.Create(&detail)
	rs, err := database.HandleResult(rs)
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

	if detail.ID == exists.ID {
		rs := repo.db.Save(&detail)
		rs, err = database.HandleResult(rs)
		if err != nil {
			return nil, err
		}
	} else {
		detail.ID = exists.ID
		rs := repo.db.Save(&detail)
		rs, err = database.HandleResult(rs)
	}

	return &detail, nil
}

func (repo *RepositoryStruct) Delete(pk uint) (bool, error) {
	exists, err := repo.Find(pk)
	if err != nil {
		return false, err
	}

	rs := repo.db.Delete(exists)
	rs, err = database.HandleResult(rs)
	if err != nil {
		return false, err
	}

	return true, err
}

func (repo *RepositoryStruct) GetByGroupId(groupId uint) ([]entity.GroupDetail, error) {
	var groupDetails []entity.GroupDetail
	rs := repo.db.Where(&entity.GroupDetail{GroupId: groupId}).Find(&groupDetails)
	rs, err := database.HandleResult(rs)
	if rs.RowsAffected == 0 {
		return make([]entity.GroupDetail, 0), nil
	}

	return groupDetails, err
}

func (repo *RepositoryStruct) GetByUserId(userId uint) ([]entity.GroupDetail, error) {
	var groupDetails []entity.GroupDetail
	rs := repo.db.Where(&entity.GroupDetail{UserId: userId}).Find(&groupDetails)
	rs, err := database.HandleResult(rs)
	if rs.RowsAffected == 0 {
		return make([]entity.GroupDetail, 0), nil
	}

	return groupDetails, err
}

func (repo *RepositoryStruct) FindByOwner(groupId uint, userId uint) (*entity.GroupDetail, error) {
	var groupDetail entity.GroupDetail
	rs := repo.db.Where(&entity.GroupDetail{
		GroupId: groupId,
		UserId:  userId,
		Role:    entity.Owner,
	}).First(groupDetail)
	rs, err := database.HandleResult(rs)

	if err != nil {
		return nil, err
	}

	return &groupDetail, err
}
