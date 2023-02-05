package main

import (
	"fmt"
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
	"github.com/miniyus/keyword-search-backend/job_queue"
	"github.com/miniyus/keyword-search-backend/logger"
	"github.com/miniyus/keyword-search-backend/pkg/IOContainer"
	"github.com/miniyus/keyword-search-backend/pkg/worker"
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
	// bindings in Container
	a.RegisterContainer(bind(cfg))
	// register middlewares
	a.Middleware(middleware)
	// register boot
	a.Register(boot)

	// register routes
	a.Route(routes.ApiPrefix, routes.Api, "api")
	a.Route("/", routes.External, "external")

	// print status
	a.Status()
	// run application
	a.Run()

}

// bind
// container에 객체 주입
func bind(configs *config.Configs) app.RegisterContainer {
	return func(c IOContainer.Container) {
		cfg := configs
		c.Singleton(func() *config.Configs {
			return cfg
		})

		db := database.DB(configs.Database)
		c.Singleton(func() *gorm.DB {
			return db
		})

		opts := configs.JobQueueConfig
		opts.Redis = utils.RedisClientMaker(configs.RedisConfig)

		opts.WorkerOptions = utils.NewCollection(opts.WorkerOptions).Map(func(v worker.Option, i int) worker.Option {
			wLoggerCfg := configs.CustomLogger
			wLoggerCfg.Filename = fmt.Sprintf("worker-%s.log", v.Name)
			v.Logger = logger.New(wLoggerCfg)

			return v
		}).Items()

		dispatcher := worker.NewDispatcher(opts)

		var jDispatcher worker.Dispatcher
		// bind를 이용했지만 singleton처럼 사용
		c.Bind(&jDispatcher, func() worker.Dispatcher {
			return dispatcher
		})

		var zLogger *zap.SugaredLogger
		c.Bind(&zLogger, func() *zap.SugaredLogger {
			return logger.New(configs.CustomLogger)
		})
	}
}

// middleware
// 미들웨어 등록
func middleware(fiberApp *fiber.App, application app.Application) {
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

// boot
// 등록 과정이 끝난 후 실행되는 로직이나 사전 작업
func boot(a app.Application) {

	var dispatcher worker.Dispatcher
	a.Resolve(&dispatcher)

	var db *gorm.DB
	a.Resolve(&db)

	var configs *config.Configs
	a.Resolve(&configs)

	create_admin.CreateAdmin(db, configs)
	job_queue.RecordHistory(dispatcher, db)
	dispatcher.Run()
}
