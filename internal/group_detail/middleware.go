package group_detail

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/auth"
	"github.com/miniyus/keyword-search-backend/database"
	"github.com/miniyus/keyword-search-backend/entity"
	"github.com/miniyus/keyword-search-backend/permission"
	"github.com/miniyus/keyword-search-backend/utils"
	"gorm.io/gorm"
)

type CreateParameter struct {
	DB *gorm.DB
}

func HandleCreatedUser(parameter CreateParameter) func(u *entity.User, tx *gorm.DB) error {
	repo := NewRepository(parameter.DB)

	return func(u *entity.User, tx *gorm.DB) error {
		if u.GroupId == nil {
			group := entity.Group{
				Name: fmt.Sprintf("%s_group", u.Username),
			}

			rs := tx.Create(&group)
			rs, err := database.HandleResult(rs)
			if err != nil {
				return err
			}

			u.GroupId = &group.ID
			rs = tx.Save(u)
			rs, err = database.HandleResult(rs)
			if err != nil {
				return err
			}
		}

		var findGroup entity.Group
		rs := tx.Find(&findGroup, *u.GroupId)
		rs, err := database.HandleResult(rs)
		if err != nil {
			return err
		}

		groupDetail := entity.GroupDetail{
			GroupId: findGroup.ID,
			UserId:  u.ID,
			Role:    entity.Owner,
		}

		create, err := repo.Create(groupDetail)
		if err != nil {
			return err
		}

		if create == nil {
			return fiber.NewError(
				fiber.StatusInternalServerError,
				fmt.Sprintf("Can not Create group_detail(user_id:%d, group_id:%d)", u.ID, findGroup.ID),
			)
		}

		return nil
	}
}

type FilterParameter struct {
	DB *gorm.DB
}

func FilterFunc(parameter FilterParameter) func(user *auth.User, perm permission.Permission) bool {
	repo := NewRepository(parameter.DB)
	return func(user *auth.User, perm permission.Permission) bool {
		if user.GroupId != nil {
			get, err := repo.GetByUserId(user.Id)
			if err != nil {
				return *user.GroupId == perm.GroupId
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
