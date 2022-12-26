package short_url

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/core/router"
)

const Prefix = "/redirect"

func Register(handler Handler) router.Register {

	return func(router fiber.Router) {
		router.Get("/:code", handler.Redirect).Name("api.redirect.code")
	}

}
