package search

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/pkg/gormrepo"
	"github.com/miniyus/gofiber/utils"
	"github.com/miniyus/keyword-search-backend/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	gormrepo.GenericRepository[entity.Search]
	AllWithPage(page utils.Page) ([]entity.Search, int64, error)
	GetByHostId(hostId uint, page utils.Page) ([]entity.Search, int64, error)
	GetDescriptionsByHostId(hostId uint, page utils.Page) ([]entity.Search, int64, error)
	BatchCreate(entities []entity.Search) ([]entity.Search, error)
	FindByShortUrl(code string, userId uint) (*entity.Search, error)
	Update(pk uint, ent entity.Search) (*entity.Search, error)
}

type RepositoryStruct struct {
	gormrepo.GenericRepository[entity.Search]
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryStruct{
		gormrepo.NewGenericRepository(db, entity.Search{}),
	}
}

func (r *RepositoryStruct) Count(search entity.Search) (int64, error) {
	var count int64 = 0

	err := r.DB().Model(&search).Count(&count).Error

	return count, err
}

func (r *RepositoryStruct) AllWithPage(page utils.Page) ([]entity.Search, int64, error) {
	var search []entity.Search
	count, err := r.Count(entity.Search{})

	if count != 0 {
		err = r.DB().Scopes(utils.Paginate(page)).Find(&search).Error
	}

	if err != nil || count == 0 {
		return make([]entity.Search, 0), 0, err
	}

	return search, count, err
}

func (r *RepositoryStruct) GetByHostId(hostId uint, page utils.Page) ([]entity.Search, int64, error) {
	var search []entity.Search
	var count int64

	where := &entity.Search{HostId: hostId}
	err := r.DB().Model(&entity.Search{}).Where(where).Count(&count).Error

	if count != 0 {
		scopes := utils.Paginate(page)

		err = r.DB().Where(where).Scopes(scopes).Order("id desc").Find(&search).Error
	}

	if err != nil || count == 0 {
		return make([]entity.Search, 0), 0, err
	}

	return search, count, err
}

func (r *RepositoryStruct) GetDescriptionsByHostId(hostId uint, page utils.Page) ([]entity.Search, int64, error) {
	var search []entity.Search
	var count int64

	where := &entity.Search{HostId: hostId}
	err := r.DB().Model(&entity.Search{}).Where(where).Count(&count).Error

	if err != nil {
		return make([]entity.Search, 0), 0, err
	}

	if count != 0 {
		scopes := utils.Paginate(page)

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

func (r *RepositoryStruct) Update(pk uint, ent entity.Search) (*entity.Search, error) {
	find, err := r.Find(pk)
	if err != nil {
		return nil, err
	}

	ent.ID = find.ID

	return r.Save(ent)
}
