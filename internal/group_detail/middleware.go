package group_detail

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/database"
	"github.com/miniyus/gollection"
	"github.com/miniyus/keyword-search-backend/entity"
	"github.com/miniyus/keyword-search-backend/internal/auth"
	"github.com/miniyus/keyword-search-backend/internal/permission"
	"github.com/miniyus/keyword-search-backend/repo"
	"gorm.io/gorm"
)

func HandleCreatedUser(u *entity.User, db *gorm.DB) error {
	repository := repo.NewGroupDetailRepository(db)

	if u.GroupId == nil {
		group := entity.Group{
			Name: fmt.Sprintf("%s_group", u.Username),
		}

		err := db.Transaction(func(tx *gorm.DB) error {
			return tx.Create(&group).Error
		})

		if err != nil {
			return err
		}

		u.GroupId = &group.ID
		err = db.Transaction(func(tx *gorm.DB) error {
			return tx.Save(u).Error
		})

		if err != nil {
			return err
		}
	}

	var findGroup entity.Group
	err := db.Find(&findGroup, *u.GroupId).Error

	if err != nil {
		return err
	}

	groupDetail := entity.GroupDetail{
		GroupId: findGroup.ID,
		UserId:  u.ID,
		Role:    entity.Owner,
	}

	create, err := repository.Create(groupDetail)
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
	repository := repo.NewGroupDetailRepository(database.GetDB())

	if groupId != 0 {
		user, err := auth.GetAuthUser(ctx)

		get, err := repository.GetByUserId(user.Id)
		if err != nil {
			return *user.GroupId == perm.GroupId
		}

		filtered := gollection.NewCollection(get).Filter(func(v entity.GroupDetail, i int) bool {
			return v.GroupId == perm.GroupId
		})

		if filtered.Count() == 0 {
			return false
		}

		return true
	}
	return false
}
