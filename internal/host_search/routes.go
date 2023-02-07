package host_search

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/app"
)

const Prefix = "/hosts/:id/search"

func Register(handler Handler) app.SubRouter {
	return func(router fiber.Router) {
		router.Get("/", handler.GetByHostId).Name("api.hosts.id.search")
		router.Get("/descriptions", handler.GetDescriptionsByHostId).Name("api.hosts.id.search.descriptions")
		router.Post("/",
			handler.BatchCreate,
		).Name("api.hosts.batch-create")
	}
}
