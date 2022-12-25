package host_search

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/core/auth"
	"github.com/miniyus/go-fiber/internal/core/router"
)

const Prefix = "/hosts/:id/search"

func Register(handler Handler) router.Register {
	return func(router fiber.Router) {
		router.Get("/", handler.GetByHostId)
		router.Post("/", auth.HasPerm(), handler.BatchCreate)
	}
}
