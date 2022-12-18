package search

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/core/database"
	"github.com/miniyus/go-fiber/internal/entity"
	"github.com/miniyus/go-fiber/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	All(page utils.Page) ([]entity.Search, int64, error)
	GetByHostId(hostId uint, page utils.Page) ([]entity.Search, int64, error)
	Find(pk uint) (*entity.Search, error)
	Create(ent entity.Search) (*entity.Search, error)
	BatchCreate(entities []entity.Search) ([]entity.Search, error)
	FindByShortUrl(code string, userId uint) (*entity.Search, error)
	Update(pk uint, ent entity.Search) (*entity.Search, error)
}

type RepositoryStruct struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryStruct{
		db: db,
	}
}

func (r *RepositoryStruct) Count(search *entity.Search) (int64, error) {
	var count int64 = 0

	rs := r.db.Model(search).Count(&count)
	_, err := database.HandleResult(rs)

	return count, err
}

func (r *RepositoryStruct) All(page utils.Page) ([]entity.Search, int64, error) {
	var search []entity.Search
	count, err := r.Count(&entity.Search{})

	if count != 0 {
		rs := r.db.Scopes(utils.Paginate(page)).Find(&search)
		_, err = database.HandleResult(rs)
	}

	return search, count, err
}

func (r *RepositoryStruct) GetByHostId(hostId uint, page utils.Page) ([]entity.Search, int64, error) {
	var search []entity.Search
	count, err := r.Count(&entity.Search{})
	if count != 0 {
		where := entity.Search{HostId: hostId}
		scopes := utils.Paginate(page)

		rs := r.db.Where(where).Scopes(scopes).Find(&search)
		_, err = database.HandleResult(rs)
	}

	return search, count, err
}

func (r *RepositoryStruct) Find(pk uint) (*entity.Search, error) {
	var search *entity.Search
	rs := r.db.Joins("Host", r.db.Where(&entity.Host{Publish: true})).Find(&search, pk)
	_, err := database.HandleResult(rs)
	if search.Host == nil {
		return nil, fiber.ErrForbidden
	}

	return search, err
}

func (r *RepositoryStruct) Create(search entity.Search) (*entity.Search, error) {
	rs := r.db.Create(&search)
	_, err := database.HandleResult(rs)

	return &search, err
}

func (r *RepositoryStruct) BatchCreate(search []entity.Search) ([]entity.Search, error) {
	rs := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"short_url"}),
	}).Create(&search)
	_, err := database.HandleResult(rs)

	return search, err
}

func (r *RepositoryStruct) FindByShortUrl(code string, userId uint) (*entity.Search, error) {
	var search *entity.Search
	rs := r.db.Joins("Host", r.db.Where(&entity.Host{Publish: true})).
		Where(&entity.Search{ShortUrl: &code}).
		Find(&search)

	_, err := database.HandleResult(rs)

	if search.Host.UserId != userId {
		return nil, fiber.ErrForbidden
	}

	return search, err
}

func (r *RepositoryStruct) Update(pk uint, search entity.Search) (*entity.Search, error) {
	exists, err := r.Find(pk)
	if err != nil {
		return nil, err
	}

	search.ID = exists.ID

	result := r.db.Save(&search)

	_, err = database.HandleResult(result)

	if err != nil {
		return nil, err
	}

	return &search, nil
}
