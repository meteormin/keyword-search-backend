package api_auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/core/auth"
)

func Register(router fiber.Router, handler Handler) {
	authApi := router.Group("/auth")
	authApi.Post("/register", handler.SignUp)
	authApi.Post("/token", handler.SignIn)

	authMiddlewares := auth.Middlewares()

	meHandlers := append(authMiddlewares, handler.Me)
	authApi.Get("/me", meHandlers...)

	resetPassHandlers := append(authMiddlewares, handler.ResetPassword)
	authApi.Patch("/password", resetPassHandlers...)

	revokeHandlers := append(authMiddlewares, handler.RevokeToken)
	authApi.Delete("/revoke", revokeHandlers...)
}
