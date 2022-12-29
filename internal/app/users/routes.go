package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/internal/core/router"
)

const Prefix = "/users"

func Register(handler Handler) router.Register {
	return func(router fiber.Router) {
		router.Get("/", handler.All)
		router.Get("/:id", handler.Get)
		router.Post("/", handler.Create)
		router.Put("/:id", handler.Update)
		router.Patch("/:id", handler.Patch)
		router.Delete("/:id", handler.Delete)
	}
}
