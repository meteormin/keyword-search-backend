package group_detail

import (
	"github.com/miniyus/keyword-search-backend/auth"
	"github.com/miniyus/keyword-search-backend/entity"
	"github.com/miniyus/keyword-search-backend/permission"
	"github.com/miniyus/keyword-search-backend/utils"
	"gorm.io/gorm"
)

type FilterParameter struct {
	DB *gorm.DB
}

func FilterFunc(parameter FilterParameter) func(user *auth.User, perm permission.Permission) bool {
	repo := NewRepository(parameter.DB)
	return func(user *auth.User, perm permission.Permission) bool {
		if user.GroupId != nil {
			get, err := repo.GetByUserId(user.Id)
			if err != nil {
				return false
			}

			filtered := utils.NewCollection(get).Filter(func(v entity.GroupDetail, i int) bool {
				return v.GroupId == perm.GroupId
			})

			if len(filtered) == 0 {
				return false
			}

			return true
		}
		return false
	}
}
