package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/core/auth"
)

func Register(router fiber.Router, handler Handler) {

	usersApi := router.Group("/users", auth.JwtMiddleware)

	usersApi.Get("/", handler.All)
	usersApi.Get("/:id", handler.Get)
	usersApi.Post("/", handler.Create)
	usersApi.Put("/:id", handler.Update)
	usersApi.Patch("/:id", handler.Patch)
	usersApi.Delete("/:id", handler.Delete)
}
