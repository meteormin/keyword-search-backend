package short_url

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/core/auth"
)

func Register(router fiber.Router, handler Handler) {
	shortUrlApi := router.Group("redirect/", auth.Middlewares()...)
	shortUrlApi.Get("/:code", handler.Redirect)
}
