package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber"
	"github.com/miniyus/gofiber/apierrors"
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/gofiber/database"
	"github.com/miniyus/gofiber/routes"
	"github.com/miniyus/gofiber/utils"
	"github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/entity"
	"github.com/miniyus/keyword-search-backend/internal/loginlogs"
	ksRoutes "github.com/miniyus/keyword-search-backend/routes"
	"github.com/redis/go-redis/v9"
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
	appConfig := cfg.App
	appConfig.FiberConfig.ErrorHandler = apierrors.OverrideDefaultErrorHandler(appConfig.Env)

	a := gofiber.New(cfg)

	a.Register(func(app app.Application) {
		var rClient *redis.Client
		rClientMaker := utils.RedisClientMaker(cfg.RedisConfig)

		app.Bind(&rClient, func() *redis.Client {
			return rClientMaker()
		})
	})

	a.Middleware(func(fiberApp *fiber.App, app app.Application) {
		fiberApp.Use(loginlogs.Middleware(database.GetDB(), fiber.MethodPost, "/api/auth/token"))
	})

	a.Register(func(app app.Application) {
		entity.RegisterHooks(app)
	})

	// register routes
	a.Route(routes.ApiPrefix, func(router app.Router, app app.Application) {
		routes.Api(router, app)
		ksRoutes.Api(router, app)
	}, "api")

	a.Route("/", func(router app.Router, app app.Application) {
		routes.External(router, app)
	}, "external")

	a.Route(ksRoutes.WebPrefix, func(router app.Router, app app.Application) {
		ksRoutes.Web(router, app)
	}, "web")

	// print status
	a.Status()
	// run application
	a.Run()

}
