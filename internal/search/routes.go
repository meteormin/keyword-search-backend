package search

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/app"
)

const Prefix = "/search"

func Register(handler Handler) app.SubRouter {
	return func(router fiber.Router) {
		router.Get("/all", handler.All).Name("api.search.all")
		router.Get("/:id", handler.Get).Name("api.search.get")
		router.Post("/", handler.Create).Name("api.search.create")
	}

}
