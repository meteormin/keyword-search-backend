package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	flogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/miniyus/keyword-search-backend/api/gofiber"
	"github.com/miniyus/keyword-search-backend/api_error"
	"github.com/miniyus/keyword-search-backend/app"
	"github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/resolver"
	"github.com/miniyus/keyword-search-backend/routes"
)

// @title keyword-search-backend Swagger API Documentation
// @version 1.0.0
// @description keyword-search-backend API
// @contact.name miniyus
// @contact.url https://miniyus.github.io
// @contact.email miniyu97@gmail.com
// @license.name MIT
// @host localhost:8080
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
		fiberApp.Use(api_error.ErrorHandler(configure))
		fiberApp.Use(cors.New(configure.Cors))

		// Add Context Config
		fiberApp.Use(config.AddContext(config.ConfigsKey, configure))
		// Add Context Logger
		logger := resolver.MakeLogger(configure.CustomLogger)
		fiberApp.Use(config.AddContext(config.LoggerKey, logger()))
	})

	a.Route(routes.ApiPrefix, routes.Api, "api")
	a.Route("/", routes.External, "external")

	a.Stats()
	a.Run()
}
