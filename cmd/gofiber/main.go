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
	"github.com/miniyus/keyword-search-backend/create_admin"
	"github.com/miniyus/keyword-search-backend/database"
	"github.com/miniyus/keyword-search-backend/logger"
	"github.com/miniyus/keyword-search-backend/pkg/IOContainer"
	"github.com/miniyus/keyword-search-backend/routes"
	"github.com/miniyus/keyword-search-backend/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
	cfg := config.GetConfigs()
	a := app.New(cfg)
	create_admin.CreateAdmin(a)

	a.RegisterContainer(bindings(cfg))
	a.Middleware(middlewares)
	a.Route(routes.ApiPrefix, routes.Api, "api")
	a.Route("/", routes.External, "external")
	a.Status()
	a.Run()
}

func bindings(configs *config.Configs) app.RegisterContainer {
	return func(c IOContainer.Container) {
		var cfg *config.Configs
		c.Bind(&cfg, func() *config.Configs {
			return configs
		})

		var db *gorm.DB
		c.Bind(&db, func() *gorm.DB {
			return database.DB(configs.Database)
		})

		var zLogger *zap.SugaredLogger
		c.Bind(&zLogger, func() *zap.SugaredLogger {
			return logger.New(configs.CustomLogger)
		})
	}
}

func middlewares(fiberApp *fiber.App, application app.Application) {
	var cfg *config.Configs

	application.Resolve(&cfg)

	fiberApp.Use(flogger.New(cfg.Logger))
	fiberApp.Use(recover.New(recover.Config{
		EnableStackTrace: !application.IsProduction(),
	}))

	fiberApp.Use(api_error.ErrorHandler(cfg))
	fiberApp.Use(cors.New(cfg.Cors))

	// Add Context Config
	fiberApp.Use(utils.AddContext(utils.ConfigsKey, cfg))
	// Add Context Logger
	var zLogger *zap.SugaredLogger
	application.Resolve(&zLogger)

	fiberApp.Use(utils.AddContext(utils.LoggerKey, zLogger))
}
