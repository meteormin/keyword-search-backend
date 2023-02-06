package group_detail

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/auth"
	"github.com/miniyus/keyword-search-backend/database"
	"github.com/miniyus/keyword-search-backend/entity"
	"github.com/miniyus/keyword-search-backend/internal"
	"github.com/miniyus/keyword-search-backend/permission"
	"github.com/miniyus/keyword-search-backend/utils"
	"gorm.io/gorm"
)

func HandleCreatedUser(u *entity.User, tx *gorm.DB) error {
	repo := NewRepository(internal.DB())

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

func FilterFunc(ctx *fiber.Ctx, groupId uint, perm permission.Permission) bool {
	repo := NewRepository(internal.DB())

	if groupId != 0 {
		user, err := auth.GetAuthUser(ctx)

		get, err := repo.GetByUserId(user.Id)
		if err != nil {
			return *user.GroupId == perm.GroupId
		}

		filtered := utils.NewCollection(get).Filter(func(v entity.GroupDetail, i int) bool {
			return v.GroupId == perm.GroupId
		})

		if filtered.Count() == 0 {
			return false
		}

		return true
	}
	return false
}
