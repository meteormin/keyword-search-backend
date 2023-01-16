package router_test

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/internal/core/router"
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	r := router.New(fiber.New(), "/", "test")

	r.Route("/", func(router fiber.Router) {
		router.Get("/", func(ctx *fiber.Ctx) error {
			return ctx.JSON("hi")
		}).Name("tty")
	})

	for _, route := range r.GetRoutes() {
		log.Print(route)
	}
}
