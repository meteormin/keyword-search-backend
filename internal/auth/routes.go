package auth

import (
	"github.com/gofiber/fiber/v2"
	jwtWare "github.com/gofiber/jwt/v3"
	"github.com/miniyus/gofiber/app"
)

const Prefix = "/auth"

func Register(handler Handler, config jwtWare.Config) app.SubRouter {
	return func(router fiber.Router) {
		router.Post("/register", handler.SignUp).Name("api.auth.register")
		router.Post("/token", handler.SignIn).Name("api.auth.token")

		router.Get("/me", JwtMiddleware(config), Middlewares(), handler.Me).Name("api.auth.me")

		router.Patch("/password", JwtMiddleware(config), Middlewares(), handler.Me).Name("api.auth.password")

		router.Delete("/revoke", JwtMiddleware(config), Middlewares(), handler.Me).Name("api.auth.revoke")
	}
}
