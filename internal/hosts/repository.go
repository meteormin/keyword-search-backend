package hosts

import (
	"github.com/miniyus/gofiber/pkg/gormrepo"
	"github.com/miniyus/gofiber/utils"
	"github.com/miniyus/keyword-search-backend/entity"
	"gorm.io/gorm"
)

type Repository interface {
	gormrepo.GenericRepository[entity.Host]
	AllWithPage(page utils.Page) ([]entity.Host, int64, error)
	GetByUserId(userId uint, page utils.Page) ([]entity.Host, int64, error)
	GetByGroupId(groupId uint, page utils.Page) ([]entity.Host, int64, error)
	GetSubjectsByUserId(userId uint, page utils.Page) ([]entity.Host, int64, error)
	GetSubjectsByGroupId(groupId uint, page utils.Page) ([]entity.Host, int64, error)
	Update(pk uint, ent entity.Host) (*entity.Host, error)
}

type RepositoryStruct struct {
	gormrepo.GenericRepository[entity.Host]
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryStruct{
		gormrepo.NewGenericRepository(db, entity.Host{}),
	}
}

func (r *RepositoryStruct) AllWithPage(page utils.Page) ([]entity.Host, int64, error) {
	var hosts []entity.Host
	count, err := r.Count(entity.Host{})

	if count != 0 {
		err = r.DB().Scopes(utils.Paginate(page)).Find(&hosts).Error
	}

	if err != nil || count == 0 {
		return make([]entity.Host, 0), 0, err
	}

	return hosts, count, err
}

func (r *RepositoryStruct) Count(host entity.Host) (int64, error) {
	var count int64 = 0

	err := r.DB().Model(&host).Count(&count).Error

	return count, err
}

func (r *RepositoryStruct) GetByUserId(userId uint, page utils.Page) (host []entity.Host, count int64, e error) {
	var hosts []entity.Host
	var cnt int64 = 0

	err := r.DB().Model(&entity.Host{}).Where(&entity.Host{UserId: userId}).Count(&cnt).Error

	if cnt != 0 {
		err = r.DB().Scopes(utils.Paginate(page)).
			Where(&entity.Host{UserId: userId}).
			Order("id desc").
			Find(&hosts).Error
	}

	if err != nil || cnt == 0 {
		return make([]entity.Host, 0), cnt, err
	}

	return hosts, cnt, err
}

func (r *RepositoryStruct) GetByGroupId(groupId uint, page utils.Page) ([]entity.Host, int64, error) {
	var group entity.Group
	var count int64
	hosts := make([]entity.Host, 0)

	if err := r.DB().Preload("Users").Find(&group, groupId).Error; err != nil {
		return hosts, 0, err
	}

	userIds := make([]int, 0)
	for _, user := range group.Users {
		userIds = append(userIds, int(user.ID))
	}

	if err := r.DB().Model(&entity.Host{}).Where("user_id IN ?", userIds).Count(&count).Error; err != nil {
		return hosts, 0, err
	}

	if err := r.DB().Scopes(utils.Paginate(page)).Where("user_id IN ?", userIds).Find(&hosts).Error; err != nil {
		return hosts, 0, err
	}

	return hosts, count, nil
}

func (r *RepositoryStruct) GetSubjectsByUserId(userId uint, page utils.Page) ([]entity.Host, int64, error) {
	var hosts []entity.Host
	var cnt int64 = 0

	err := r.DB().Model(&entity.Host{}).Where(&entity.Host{UserId: userId}).Count(&cnt).Error
	if cnt != 0 {
		err = r.DB().Select("id", "subject").Scopes(utils.Paginate(page)).
			Where(&entity.Host{UserId: userId}).
			Order("id desc").
			Find(&hosts).Error
	}

	if err != nil || cnt == 0 {
		return make([]entity.Host, 0), cnt, err
	}

	return hosts, cnt, err
}

func (r *RepositoryStruct) GetSubjectsByGroupId(groupId uint, page utils.Page) ([]entity.Host, int64, error) {
	var group entity.Group
	hosts := make([]entity.Host, 0)

	if err := r.DB().Preload("Users").Find(&group, groupId).Error; err != nil {
		return hosts, 0, err
	}

	userIds := make([]int, 0)
	for _, user := range group.Users {
		userIds = append(userIds, int(user.ID))
	}

	err := r.DB().Select("id", "subject").
		Scopes(utils.Paginate(page)).
		Where("user_id IN ?", userIds).
		Find(&hosts).Error

	if err != nil {
		return hosts, 0, err
	}

	return hosts, int64(len(hosts)), err
}

func (r *RepositoryStruct) Update(pk uint, ent entity.Host) (*entity.Host, error) {
	find, err := r.Find(pk)
	if err != nil {
		return nil, err
	}

	ent.ID = find.ID

	return r.Save(ent)
}
