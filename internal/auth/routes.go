package auth

import (
	"github.com/gofiber/fiber/v2"
)

func Register(router fiber.Router, handler Handler) {
	authApi := router.Group("/auth")
	authApi.Post("/register", handler.SignUp)
	authApi.Post("/token", handler.SignIn)
}
