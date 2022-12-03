package hosts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/core/auth"
)

func Register(router fiber.Router, handler Handler) {
	hostsApi := router.Group("/hosts", auth.Middlewares()...)

	hostsApi.Post("/", handler.Create)
	hostsApi.Get("/", handler.All)
	hostsApi.Get("/:id", handler.Get)
	hostsApi.Put("/:id", handler.Update)
	hostsApi.Delete("/:id", handler.Delete)
}
