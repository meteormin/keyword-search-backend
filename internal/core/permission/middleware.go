package permission

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/internal/core/auth"
	"github.com/miniyus/keyword-search-backend/internal/core/context"
	"github.com/miniyus/keyword-search-backend/internal/entity"
	"github.com/miniyus/keyword-search-backend/internal/utils"
	"strings"
)

// HasPermission
// has permission?
func HasPermission(permissions ...Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pass := false

		currentUser := c.Locals(context.AuthUser).(*auth.User)
		if currentUser.Role == entity.Admin.RoleToString() {
			return c.Next()
		}
		permCollection := c.Locals(context.Permissions).(Collection)

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
	for _, perm := range hasPerm {
		for _, action := range perm.Actions {
			routePath := c.Path()
			if strings.Contains(routePath, action.Resource) {
				method := c.Method()
				if method == "OPTION" {
					method = "GET"
				}

				filtered := utils.Filter(action.Methods, func(v Method, i int) bool {
					return v.ToString() == method
				})

				if len(filtered) != 0 {
					pass = true
				}
			}
		}
	}
	return pass
}
