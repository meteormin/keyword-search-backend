package group_detail

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/auth"
	"github.com/miniyus/gofiber/database"
	"github.com/miniyus/gofiber/permission"
	"github.com/miniyus/gofiber/utils"
	"github.com/miniyus/keyword-search-backend/entity"
	"gorm.io/gorm"
)

func HandleCreatedUser(u *entity.User, db *gorm.DB) error {
	repo := NewRepository(db)

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
	repo := NewRepository(database.GetDB())

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
