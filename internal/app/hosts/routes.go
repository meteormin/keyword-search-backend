package hosts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/core/router"
)

const Prefix = "/hosts"

func Register(handler Handler) router.Register {
	return func(router fiber.Router) {
		router.Post("/", handler.Create).Name("api.hosts.create")
		router.Get("/", handler.All).Name("api.hosts.get")
		router.Get("/:id", handler.Get).Name("api.hosts.id")
		router.Put("/:id", handler.Update).Name("api.hosts.put")
		router.Patch("/:id", handler.Patch).Name("api.hosts.patch")
		router.Delete("/:id", handler.Delete).Name("api.hosts.delete")
	}
}
