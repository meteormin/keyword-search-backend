package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
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
	"github.com/miniyus/keyword-search-backend/routes"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

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

	a.Register(func(app app.Application) {
		pCfg := &cfg
		app.Bind(&pCfg, func() *config.Configs {
			return pCfg
		})

		var rClient *redis.Client
		rClientMaker := utils.RedisClientMaker(cfg.RedisConfig)

		app.Bind(&rClient, func() *redis.Client {
			return rClientMaker()
		})
	})

	a.Middleware(func(fiberApp *fiber.App, app app.Application) {
		fiberApp.Use(compress.New())
		fiberApp.Use(etag.New())
		fiberApp.Use(requestid.New())
		fiberApp.Use(loginlogs.Middleware(database.GetDB(), fiber.MethodPost, "/api/auth/token"))
	})

	a.Register(func(app app.Application) {
		var db *gorm.DB
		app.Resolve(&db)

		permission.CreateDefaultPermissions(db, cfg.Permission)
	})

	a.Register(entity.RegisterHooks)

	// register routes
	a.Route(routes.ApiPrefix, func(router app.Router, app app.Application) {
		routes.Api(router, app)
	}, "api")

	a.Route(routes.WebPrefix, func(router app.Router, app app.Application) {
		routes.Web(router, app)
	}, "web")

	// print status
	a.Status()
	// run application
	a.Run()
}
