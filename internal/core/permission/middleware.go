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

		currentUser, err := context.GetContext[*auth.User](c, context.AuthUser)

		if currentUser.Role == string(entity.Admin) {
			return c.Next()
		}

		var permCollection Collection

		db, err := context.GetContext[*gorm.DB](c, context.DB)
		if err != nil {
			return err
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
			permCollection, err = context.GetContext[Collection](c, context.Permissions)

			if err != nil {
				permCollection = nil
				containerContext, err := context.GetContext[container.Container](c, context.Container)
				if err != nil {
					return err
				}

				containerContext.Resolve(&permCollection)
			}

			entities := make([]entity.Permission, 0)
			permCollection.For(func(perm Permission, i int) {
				entities = append(entities, ToPermissionEntity(perm))
			})

			_, err = repo.Save(entities)
			if err != nil {
				return err
			}
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
