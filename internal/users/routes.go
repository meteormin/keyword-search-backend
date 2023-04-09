package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/app"
)

const Prefix = "/users"

func Register(handler Handler) app.SubRouter {
	return func(router fiber.Router) {
		router.Get("/", handler.All).Name("api.users.all")
		router.Get("/:id", handler.Get).Name("api.users.get")
		router.Post("/", handler.Create).Name("api.users.create")
		router.Patch("/", handler.PatchForMe).Name("api.users.update")
		router.Patch("/:id", handler.PatchForAdm).Name("api.users.patch")
		router.Delete("/:id", handler.Delete).Name("api.users.delete")
	}
}
