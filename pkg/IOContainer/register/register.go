package register

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	flogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/miniyus/keyword-search-backend/api/gofiber"
	"github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/api_error"
	"github.com/miniyus/keyword-search-backend/internal/permission"
	"github.com/miniyus/keyword-search-backend/internal/resolver"
	"github.com/miniyus/keyword-search-backend/pkg/IOContainer"
	"github.com/miniyus/keyword-search-backend/pkg/jwt"
	"github.com/miniyus/keyword-search-backend/pkg/worker"
	"go.uber.org/zap"
	"log"
)

// boot is High Priority
// container settings
func boot(w IOContainer.Container) {
	app := w.App()
	w.Singleton(app)

	cfg := w.Config()
	w.Singleton(cfg)

	db := w.Database()
	w.Singleton(db)

	redisClient := resolver.MakeRedisClient(cfg.RedisConfig)
	w.Singleton(redisClient)

	var tg jwt.Generator
	jwtGenerator := resolver.MakeJwtGenerator(resolver.JwtGeneratorConfig{
		DataPath: cfg.Path.DataPath,
		Exp:      cfg.Auth.Exp,
	})
	w.Bind(&tg, jwtGenerator)

	var logs *zap.SugaredLogger
	loggerStruct := resolver.MakeLogger(cfg.CustomLogger)
	w.Bind(&logs, loggerStruct)

	var perms permission.Collection
	permissionCollection := resolver.MakePermissionCollection(cfg.Permission)
	w.Bind(&perms, permissionCollection)

	var dispatcher worker.Dispatcher
	jobDispatcherStruct := resolver.MakeJobDispatcher(resolver.JobDispatcherConfig{
		WorkerCfg: cfg.QueueConfig,
		RedisCfg:  cfg.RedisConfig,
	})
	w.Bind(&dispatcher, jobDispatcherStruct)
}

// middlewares register middleware
// fiber app middleware settings
func middlewares(w IOContainer.Container) {
	log.Print(w.Instances())
	w.App().Use(flogger.New(w.Config().Logger))
	w.App().Use(recover.New(recover.Config{
		EnableStackTrace: !w.IsProduction(),
	}))
	w.App().Use(api_error.ErrorHandler)
	w.App().Use(cors.New(w.Config().Cors))

	// Add Context Container
	w.App().Use(config.AddContext("container", w))
	// Add Context Config
	w.App().Use(config.AddContext(config.ConfigsKey, w.Config()))
	// Add Context Logger
	var logger *zap.SugaredLogger
	w.Resolve(&logger)
	w.App().Use(config.AddContext(config.LoggerKey, logger))
	// Add Context Permissions
	var perms permission.Collection
	w.Resolve(&perms)
	w.App().Use(config.AddContext(config.PermissionsKey, perms))
	// Add Context Redis
	w.App().Use(config.AddContext(config.RedisKey, resolver.MakeRedisClient(w.Config().RedisConfig)))
}

// routes register Routes
func routes(w IOContainer.Container) {

}

// Resister
// private 함수들 모아서 순서대로 실행 해주는 public 함수
func Resister(w IOContainer.Container) {
	boot(w)
	middlewares(w)
	routes(w)
}
