package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/miniyus/gofiber/app"
	_ "github.com/miniyus/keyword-search-backend/api"
)

const WebPrefix = "/web"

func Web(router app.Router, app app.Application) {
	router.Route("/", func(router fiber.Router) {
		router.Get("/swagger/*", swagger.HandlerDefault)
	})
}
