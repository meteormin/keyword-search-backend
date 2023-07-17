package repo

import (
	"github.com/miniyus/gorm-extension/gormrepo"
	"github.com/miniyus/keyword-search-backend/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FileRepository struct {
	gormrepo.GenericRepository[entity.File]
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{
		GenericRepository: gormrepo.NewGenericRepository[entity.File](
			db,
			entity.File{},
		),
	}
}

func (r *FileRepository) BatchCreate(files []entity.File) ([]entity.File, error) {
	err := r.DB().Transaction(func(tx *gorm.DB) error {
		return tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "id"},
			},
			DoNothing: true,
		}).Create(&files).Error
	})

	if err != nil {
		return make([]entity.File, 0), err
	}

	return files, err
}
