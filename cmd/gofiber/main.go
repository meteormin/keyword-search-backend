package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	flogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/api_error"
	"github.com/miniyus/keyword-search-backend/internal/app"
	"github.com/miniyus/keyword-search-backend/internal/resolver"
	"github.com/miniyus/keyword-search-backend/internal/routes"
)

// @title keyword-search-backend Swagger API Documentation
// @version 1.0.0
// @description keyword-search-backend API
// @contact.name miniyus
// @contact.url https://miniyus.github.io
// @contact.email miniyu97@gmail.com
// @license.name MIT
// @host localhost:9090
// @BasePath /
// @schemes http
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				   Bearer token type
func main() {
	a := app.New()

	a.Middleware(func(fiberApp *fiber.App, application app.Application) {
		configure := application.Config()

		fiberApp.Use(flogger.New(configure.Logger))
		fiberApp.Use(recover.New(recover.Config{
			EnableStackTrace: !application.IsProduction(),
		}))
		fiberApp.Use(api_error.ErrorHandler)
		fiberApp.Use(cors.New(configure.Cors))

		// Add Context Config
		fiberApp.Use(config.AddContext(config.ConfigsKey, configure))
		// Add Context Logger
		fiberApp.Use(config.AddContext(config.LoggerKey, resolver.MakeLogger(configure.CustomLogger)))
	})

	a.Route(routes.ApiPrefix, routes.Api, "api")
	a.Route("/", routes.External, "external")

	a.Stats()
	a.Run()
}
