package permission

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/database"
	"github.com/miniyus/gollection"
	"github.com/miniyus/keyword-search-backend/entity"
	"github.com/miniyus/keyword-search-backend/internal/auth"
	"github.com/miniyus/keyword-search-backend/repo"
	"gorm.io/gorm"
	"strings"
)

type HasPermissionParameter struct {
	DB           *gorm.DB
	DefaultPerms []Config
	FilterFunc   func(ctx *fiber.Ctx, groupId uint, p Permission) bool
}

// HasPermission
// check has permissions middleware
func HasPermission(parameter HasPermissionParameter) func(permissions ...Permission) fiber.Handler {
	db := parameter.DB
	defaultPerms := NewPermissionCollection(
		NewPermissionsFromConfig(parameter.DefaultPerms)...,
	)
	filter := parameter.FilterFunc
	return func(permissions ...Permission) fiber.Handler {
		return func(c *fiber.Ctx) error {
			pass := false

			var permCollection Collection

			if db == nil {
				db = database.GetDB()
			}

			authUser, err := auth.GetAuthUser(c)

			repository := repo.NewPermissionRepository(db)

			get, err := repository.GetByGroupId(*authUser.GroupId)
			if err == nil {
				permCollection = NewPermissionCollection()
				var perm Permission
				gollection.NewCollection(get).For(func(v entity.Permission, i int) {
					permCollection.Add(perm.FromEntity(v))
				})
			}

			if permCollection == nil {
				permCollection = defaultPerms
			}

			if len(permissions) != 0 {
				permCollection.Concat(permissions)
			}

			userHasPerm := permCollection.Filter(func(p Permission, i int) bool {
				if filter != nil {
					return filter(c, *authUser.GroupId, p)
				}

				if authUser.GroupId != nil {
					return *authUser.GroupId == p.GroupId
				}

				return false
			})

			pass = checkPermissionFromCtx(userHasPerm.Items(), c)

			if pass {
				return c.Next()
			}

			return fiber.ErrForbidden
		}
	}
}

func checkPermissionFromCtx(hasPerm []Permission, c *fiber.Ctx) bool {
	if len(hasPerm) == 0 {
		return false
	}

	pass := false
	gollection.NewCollection(hasPerm).For(func(perm Permission, i int) {
		gollection.NewCollection(perm.Actions).For(func(action Action, j int) {
			routePath := c.Path()
			if strings.Contains(routePath, action.Resource) {
				method := c.Method()
				if strings.ToUpper(method) == "OPTION" {
					method = "GET"
				}

				filtered := gollection.NewCollection(action.Methods).Filter(func(v Method, i int) bool {
					return strings.ToUpper(string(v)) == strings.ToUpper(method)
				})

				if filtered.Count() != 0 {
					pass = true
				}
			}
		})
	})

	return pass
}
