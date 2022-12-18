package search

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/core/auth"
)

func Register(router fiber.Router, handler Handler) {
	searchFromHostApi := router.Group("/hosts/:id/search", auth.Middlewares()...)
	searchFromHostApi.Get("/", handler.GetByHostId)
	searchFromHostApi.Post("/", auth.HasPerm(), handler.BatchCreate)

	searchApi := router.Group("/search", auth.Middlewares()...)
	searchApi.Get("/all", handler.All)
	searchApi.Get("/:id", handler.Get)
	searchApi.Post("/", handler.Create)

}
