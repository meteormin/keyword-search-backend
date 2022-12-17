package search

import (
	"github.com/miniyus/go-fiber/internal/core/database"
	"github.com/miniyus/go-fiber/internal/entity"
	"github.com/miniyus/go-fiber/internal/utils"
	"gorm.io/gorm"
)

type Repository interface {
	All(page utils.Page) ([]*entity.Search, int64, error)
	GetByHostId(hostId uint, page utils.Page) ([]*entity.Search, int64, error)
	Find(pk uint) (*entity.Search, error)
	Create(ent *entity.Search) (*entity.Search, error)
}

type RepositoryStruct struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *RepositoryStruct {
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

func (r *RepositoryStruct) All(page utils.Page) ([]*entity.Search, int64, error) {
	var search []*entity.Search
	count, err := r.Count(&entity.Search{})

	if count != 0 {
		rs := r.db.Scopes(utils.Paginate(page)).Find(&search)
		_, err = database.HandleResult(rs)
	}

	return search, count, err
}

func (r *RepositoryStruct) GetByHostId(hostId uint, page utils.Page) ([]*entity.Search, int64, error) {
	var search []*entity.Search
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
	rs := r.db.Find(&search, pk)
	_, err := database.HandleResult(rs)

	return search, err
}

func (r *RepositoryStruct) Create(search entity.Search) (*entity.Search, error) {
	rs := r.db.Create(&search)
	_, err := database.HandleResult(rs)

	return &search, err
}
