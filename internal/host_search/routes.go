package host_search

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/gofiber/permission"
)

const Prefix = "/hosts/:id/search"

func Register(handler Handler, hasPerm permission.HasPermissionParameter) app.SubRouter {
	return func(router fiber.Router) {
		router.Get("/", handler.GetByHostId).Name("api.hosts.id.search")
		router.Get("/descriptions", handler.GetDescriptionsByHostId).Name("api.hosts.id.search.descriptions")
		router.Post("/",
			permission.HasPermission(hasPerm),
			handler.BatchCreate,
		).Name("api.hosts.batch-create")
	}
}
