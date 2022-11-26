package test_api

import (
	"github.com/gofiber/fiber/v2"
)

func SetRoutes(router fiber.Router, handler Handler) {
	testApi := router.Group("/test-api")
	testApi.Get("/", handler.GetTest)
}
