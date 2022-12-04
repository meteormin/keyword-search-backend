package api_auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/core/auth"
)

func Register(router fiber.Router, handler Handler) {
	authApi := router.Group("/auth")
	authApi.Post("/register", handler.SignUp)
	authApi.Post("/token", handler.SignIn)

	authApi.Get("/me", auth.JwtMiddleware, auth.GetUserFromJWT, auth.CheckExpired, handler.Me)
}
