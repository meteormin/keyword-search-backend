package permission

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/auth"
	"github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/entity"
	"github.com/miniyus/keyword-search-backend/utils"
	"gorm.io/gorm"
	"strings"
)

type HasPermissionParameter struct {
	DB           *gorm.DB
	DefaultPerms Collection
	FilterFunc   func(currentUser *auth.User, p Permission) bool
}

// HasPermission check has permissions middleware
func HasPermission(parameter HasPermissionParameter, permissions ...Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pass := false

		currentUser, err := config.GetContext[*auth.User](c, config.AuthUserKey)

		if currentUser.Role == string(entity.Admin) {
			return c.Next()
		}

		var permCollection Collection
		var db *gorm.DB

		if parameter.DB == nil {
			db, err = config.GetContext[*gorm.DB](c, config.DBKey)
			if err != nil {
				return err
			}
		} else {
			db = parameter.DB
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
			if parameter.DefaultPerms == nil {
				permCollection, err = config.GetContext[Collection](c, config.PermissionsKey)

				if err != nil {
					return err
				}
			} else {

				permCollection = parameter.DefaultPerms
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
			if parameter.FilterFunc != nil {
				return parameter.FilterFunc(currentUser, p)
			}

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
	if len(hasPerm) == 0 {
		return false
	}

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
