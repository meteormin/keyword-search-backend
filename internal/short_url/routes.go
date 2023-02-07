package short_url

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/app"
)

const Prefix = "/short-url"

func Register(handler Handler) app.SubRouter {

	return func(router fiber.Router) {
		router.Get("/:code/redirect", handler.Redirect).Name("api.shor-url.code.redirect")
		router.Get("/:code", handler.FindUrlByCode).Name("api.short-url.find-url-by-code")
	}

}
