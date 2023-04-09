package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/swagger"
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/gofiber/utils"
	_ "github.com/miniyus/keyword-search-backend/api"
)

const WebPrefix = "/web"

func Web(router app.Router, app app.Application) {
	router.Route("/", func(router fiber.Router) {
		router.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))
		router.Get("/swagger/*", swagger.HandlerDefault)
		router.Get("/health-check", healthCheck)
	})
}

// healthCheck
// @Summary health check your server
// @Description health check your server
// @Success 200 {object} utils.StatusResponse
// @Tags HealthCheck
// @Accept */*
// @Produce json
// @Router /health-check [get]
func healthCheck(ctx *fiber.Ctx) error {

	err := ctx.JSON(utils.StatusResponse{Status: true})
	if err != nil {
		return ctx.JSON(utils.StatusResponse{Status: false})
	}

	return err
}
