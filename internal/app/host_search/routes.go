package host_search

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/internal/core/permission"
	"github.com/miniyus/keyword-search-backend/internal/core/router"
)

const Prefix = "/hosts/:id/search"

func Register(handler Handler) router.Register {
	return func(router fiber.Router) {
		router.Get("/", handler.GetByHostId).Name("api.hosts.id.search")
		router.Get("/descriptions", handler.GetDescriptionsByHostId).Name("api.hosts.id.search.descriptions")
		router.Post("/", permission.HasPermission(), handler.BatchCreate).Name("api.hosts.batch-create")
	}
}
