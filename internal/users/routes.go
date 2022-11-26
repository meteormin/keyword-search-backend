package users

import "github.com/gofiber/fiber/v2"

func SetRoutes(router fiber.Router, handler Handler) {
	usersApi := router.Group("/users")

	usersApi.Get("/", handler.All)
	usersApi.Get("/:id", handler.Get)
	usersApi.Post("/", handler.Create)
	usersApi.Put("/:id", handler.Update)
	usersApi.Patch("/:id", handler.Patch)
	usersApi.Delete("/:id", handler.Delete)
}
