package api_auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/core/auth"
)

const Prefix = "/auth"

func Register(handler Handler) func(router fiber.Router) {
	return func(router fiber.Router) {
		router.Post("/register", handler.SignUp)
		router.Post("/token", handler.SignIn)

		authMiddlewares := auth.Middlewares()

		meHandlers := append(authMiddlewares, handler.Me)
		router.Get("/me", meHandlers...)

		resetPassHandlers := append(authMiddlewares, handler.ResetPassword)
		router.Patch("/password", resetPassHandlers...)

		revokeHandlers := append(authMiddlewares, handler.RevokeToken)
		router.Delete("/revoke", revokeHandlers...)
	}
}
