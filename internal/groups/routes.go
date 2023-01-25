package groups

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/app"
)

const Prefix = "/groups"

func Register(handler Handler) app.SubRouter {
	return func(router fiber.Router) {
		router.Get("/", handler.All).Name("api.groups.all")
		router.Get("/:id", handler.Find).Name("api.groups.find")
		router.Get("/name/:name", handler.FindByName).Name("api.groups.find-by-name")
		router.Put("/:id", handler.Update).Name("api.groups.update")
		router.Patch("/:id", handler.Patch).Name("api.groups.patch")
		router.Delete("/:id", handler.Delete).Name("api.groups.delete")
	}
}
