package search

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/core/router"
)

const Prefix = "/search"

func Register(handler Handler) router.Register {
	return func(router fiber.Router) {
		router.Get("/all", handler.All)
		router.Get("/:id", handler.Get)
		router.Post("/", handler.Create)
	}

}
