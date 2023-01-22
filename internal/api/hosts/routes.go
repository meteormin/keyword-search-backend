package hosts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/internal/app"
)

const Prefix = "/hosts"

func Register(handler Handler) app.SubRouter {
	return func(router fiber.Router) {
		router.Post("/", handler.Create).Name("api.hosts.create")
		router.Get("/", handler.All).Name("api.hosts.all")
		router.Get("/subjects", handler.GetSubjects).Name("api.hosts.subjects")
		router.Get("/:id", handler.Get).Name("api.hosts.get")
		router.Put("/:id", handler.Update).Name("api.hosts.put")
		router.Patch("/:id", handler.Patch).Name("api.hosts.patch")
		router.Delete("/:id", handler.Delete).Name("api.hosts.delete")
	}
}
