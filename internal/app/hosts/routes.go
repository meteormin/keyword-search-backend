package hosts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/core/router"
)

const Prefix = "/hosts"

func Register(handler Handler) router.Register {
	return func(router fiber.Router) {
		router.Post("/", handler.Create)
		router.Get("/", handler.All)
		router.Get("/:id", handler.Get)
		router.Put("/:id", handler.Update)
		router.Patch("/:id", handler.Patch)
		router.Delete("/:id", handler.Delete)
	}
}
