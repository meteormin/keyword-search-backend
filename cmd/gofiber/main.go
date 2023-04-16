package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/miniyus/gofiber"
	"github.com/miniyus/gofiber/apierrors"
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/gofiber/database"
	"github.com/miniyus/gofiber/utils"
	"github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/entity"
	"github.com/miniyus/keyword-search-backend/internal/loginlogs"
	"github.com/miniyus/keyword-search-backend/internal/permission"
	"github.com/miniyus/keyword-search-backend/internal/tasks"
	"github.com/miniyus/keyword-search-backend/routes"
	"github.com/redis/go-redis/v9"
)

func register(cfg *config.Configs) app.Register {
	return func(app app.Application) {
		app.Bind(&cfg, func() *config.Configs {
			return cfg
		})

		var rClient *redis.Client
		rClientMaker := utils.RedisClientMaker(cfg.RedisConfig)

		app.Bind(&rClient, func() *redis.Client {
			return rClientMaker()
		})
	}
}

func middleware(cfg *config.Configs) app.MiddlewareRegister {
	return func(fiberApp *fiber.App, app app.Application) {
		fiberApp.Use(cors.New(cfg.Cors))
		fiberApp.Use(compress.New())
		fiberApp.Use(etag.New())
		fiberApp.Use(requestid.New())
		fiberApp.Use(loginlogs.Middleware(database.GetDB(), fiber.MethodPost, "/api/auth/token"))
	}
}

// @title keyword-search-backend Swagger API Documentation
// @version 1.0.1
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
	appConfig := cfg.App
	appConfig.FiberConfig.ErrorHandler = apierrors.OverrideDefaultErrorHandler(appConfig.Env)
	a := gofiber.New(*cfg.Configs)

	a.Register(register(&cfg))
	a.Register(entity.RegisterHooks)
	a.Register(tasks.RegisterJob)
	a.Register(tasks.RegisterSchedule)
	a.Register(permission.CreateDefaultPermissions(cfg.Permission))

	// register middlewares
	a.Middleware(middleware(&cfg))

	// register routes
	a.Route(routes.ApiPrefix, routes.Api, "api")
	a.Route(routes.WebPrefix, routes.Web, "web")

	// print status
	a.Status()
	// run application
	a.Run()
}
