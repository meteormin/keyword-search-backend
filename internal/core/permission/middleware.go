package permission

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/internal/core/auth"
	"github.com/miniyus/keyword-search-backend/internal/core/container"
	"github.com/miniyus/keyword-search-backend/internal/core/context"
	"github.com/miniyus/keyword-search-backend/internal/entity"
	"github.com/miniyus/keyword-search-backend/internal/utils"
	"gorm.io/gorm"
	"strings"
)

// HasPermission check has permissions middleware
func HasPermission(permissions ...Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pass := false

		currentUser := c.Locals(context.AuthUser).(*auth.User)
		if currentUser.Role == string(entity.Admin) {
			return c.Next()
		}

		var permCollection Collection

		db, ok := c.Locals(context.DB).(*gorm.DB)
		if !ok {
			return fiber.NewError(fiber.StatusInternalServerError, "can't find context.DB")
		}

		repo := NewRepository(db)

		get, err := repo.Get(*currentUser.GroupId)
		if err == nil {
			permCollection = NewPermissionCollection()
			utils.NewCollection(get).For(func(v entity.Permission, i int) {
				permCollection.Add(EntityToPermission(v))
			})
		}

		if permCollection == nil {
			permCollection, ok = c.Locals(context.Permissions).(Collection)
			if !ok {
				permCollection = nil
				containerContext := c.Locals(context.Container).(container.Container)
				permCollection, ok = containerContext.Resolve(permCollection).(Collection)
				if !ok {
					return fiber.NewError(fiber.StatusInternalServerError, "can not found context permissions")
				}
			}

			SavePermission(permCollection.Items()...)
		}

		if len(permissions) != 0 {
			permCollection.Concat(permissions)
		}

		userHasPerm := permCollection.Filter(func(p Permission, i int) bool {
			if currentUser.GroupId != nil {
				return currentUser.GroupId == &p.GroupId
			}

			return false
		})

		pass = checkPermissionFromCtx(userHasPerm, c)

		if pass {
			return c.Next()
		}

		return fiber.ErrForbidden
	}
}

func checkPermissionFromCtx(hasPerm []Permission, c *fiber.Ctx) bool {
	pass := false
	utils.NewCollection(hasPerm).For(func(perm Permission, i int) {
		utils.NewCollection(perm.Actions).For(func(action Action, j int) {
			routePath := c.Path()
			if strings.Contains(routePath, action.Resource) {
				method := c.Method()
				if method == "OPTION" {
					method = "GET"
				}

				filtered := utils.NewCollection(action.Methods).Filter(func(v Method, i int) bool {
					return string(v) == method
				})

				if len(filtered) != 0 {
					pass = true
				}
			}
		})
	})

	return pass
}

func SavePermission(permissions ...Permission) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authUser, ok := ctx.Locals(context.AuthUser).(*auth.User)
		if !ok {
			return ctx.Next()
		}

		if authUser == nil {
			return ctx.Next()
		}

		db, ok := ctx.Locals(context.DB).(*gorm.DB)
		if !ok {
			return fiber.NewError(fiber.StatusInternalServerError, "can't find context.DB")
		}

		repo := NewRepository(db)

		entities := make([]entity.Permission, 0)
		for _, perm := range permissions {
			entities = append(entities, ToPermissionEntity(perm))
		}

		_, err := repo.Save(entities)
		if err != nil {
			return err
		}

		return ctx.Next()
	}
}
