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
		router.Get("/me", meHandlers...).Name("api.auth.me")

		resetPassHandlers := append(authMiddlewares, handler.ResetPassword)
		router.Patch("/password", resetPassHandlers...).Name("api.auth.password")

		revokeHandlers := append(authMiddlewares, handler.RevokeToken)
		router.Delete("/revoke", revokeHandlers...).Name("api.auth.revoke")
	}
}
