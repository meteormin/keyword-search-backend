package search

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/pagination"
	"github.com/miniyus/gorm-extension/gormrepo"
	"github.com/miniyus/keyword-search-backend/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Filter struct {
	pagination.Page
	Query    *string
	QueryKey *string
	Publish  *bool
	SortBy   *string
	OrderBy  *string
}

func (f Filter) FillEntity(ent *entity.Search) {
	if f.QueryKey != nil {
		ent.QueryKey = *f.QueryKey
	}

	if f.Query != nil {
		ent.Query = *f.Query
	}

	if f.Publish != nil {
		ent.Publish = *f.Publish
	}
}

type Repository interface {
	gormrepo.GenericRepository[entity.Search]
	AllWithPage(page pagination.Page) ([]entity.Search, int64, error)
	GetByHostId(hostId uint, filter Filter) ([]entity.Search, int64, error)
	GetDescriptionsByHostId(hostId uint, page pagination.Page) ([]entity.Search, int64, error)
	BatchCreate(entities []entity.Search) ([]entity.Search, error)
	FindByShortUrl(code string, userId uint) (*entity.Search, error)
	HasHost(hostId uint, userId uint) bool
}

type RepositoryStruct struct {
	gormrepo.GenericRepository[entity.Search]
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryStruct{
		gormrepo.NewGenericRepository(db, entity.Search{}),
	}
}

func (r *RepositoryStruct) HasHost(hostId uint, userId uint) bool {
	var host entity.Host
	err := r.DB().Where(&entity.Host{UserId: userId}).First(&host, hostId).Error
	if err != nil {
		return false
	}

	return true
}

func (r *RepositoryStruct) Count(search entity.Search) (int64, error) {
	var count int64 = 0

	err := r.DB().Model(&search).Count(&count).Error

	return count, err
}

func (r *RepositoryStruct) AllWithPage(page pagination.Page) ([]entity.Search, int64, error) {
	var search []entity.Search
	count, err := r.Count(entity.Search{})

	if count != 0 {
		err = r.DB().Scopes(pagination.Paginate(page)).Find(&search).Error
	}

	if err != nil || count == 0 {
		return make([]entity.Search, 0), 0, err
	}

	return search, count, err
}

func (r *RepositoryStruct) GetByHostId(hostId uint, filter Filter) ([]entity.Search, int64, error) {
	var search []entity.Search
	var count int64

	where := entity.Search{HostId: hostId}
	filter.FillEntity(&where)
	err := r.DB().Model(&entity.Search{}).Where(where).Count(&count).Error

	if count != 0 {
		scopes := pagination.Paginate(filter.Page)
		var order string
		if filter.SortBy != nil && filter.OrderBy != nil {
			order = fmt.Sprintf("%s %s, id desc", *filter.SortBy, *filter.OrderBy)
		} else {
			order = "id desc"
		}
		err = r.DB().Where(where).Scopes(scopes).Order(order).Find(&search).Error
	}

	if err != nil || count == 0 {
		return make([]entity.Search, 0), 0, err
	}

	return search, count, err
}

func (r *RepositoryStruct) GetDescriptionsByHostId(hostId uint, page pagination.Page) ([]entity.Search, int64, error) {
	var search []entity.Search
	var count int64

	where := &entity.Search{HostId: hostId}
	err := r.DB().Model(&entity.Search{}).Where(where).Count(&count).Error

	if err != nil {
		return make([]entity.Search, 0), 0, err
	}

	if count != 0 {
		scopes := pagination.Paginate(page)

		err = r.DB().Select("id", "description", "short_url").
			Where(where).
			Scopes(scopes).
			Order("id desc").
			Find(&search).Error
	}

	if err != nil || count == 0 {
		return make([]entity.Search, 0), 0, err
	}

	return search, count, err
}

func (r *RepositoryStruct) Find(pk uint) (*entity.Search, error) {
	var search entity.Search
	err := r.DB().Joins("Host", r.DB().Where(&entity.Host{Publish: true})).First(&search, pk).Error

	if err != nil {
		return nil, err
	}

	if search.Host == nil {
		return nil, fiber.ErrForbidden
	}

	return &search, nil
}

func (r *RepositoryStruct) BatchCreate(search []entity.Search) ([]entity.Search, error) {
	err := r.DB().Transaction(func(tx *gorm.DB) error {
		return tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{"short_url"}),
		}).Create(&search).Error
	})

	if err != nil {
		return make([]entity.Search, 0), err
	}

	return search, err
}

func (r *RepositoryStruct) FindByShortUrl(code string, userId uint) (*entity.Search, error) {
	var search entity.Search
	err := r.DB().Joins("Host", r.DB().Where(&entity.Host{Publish: true})).
		Where(&entity.Search{ShortUrl: &code}).
		First(&search).Error

	if err != nil {
		return nil, err
	}

	if search.Host.UserId != userId {
		return nil, fiber.ErrForbidden
	}

	return &search, err
}
