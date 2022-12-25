package register

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	flogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	_ "github.com/miniyus/go-fiber/api/gofiber"
	"github.com/miniyus/go-fiber/internal/core/api_error"
	"github.com/miniyus/go-fiber/internal/core/container"
	"github.com/miniyus/go-fiber/internal/core/context"
	"github.com/miniyus/go-fiber/internal/core/logger"
	router "github.com/miniyus/go-fiber/internal/routes"
	"github.com/miniyus/go-fiber/pkg/jwt"
	rsGen "github.com/miniyus/go-fiber/pkg/rs256"
	"go.uber.org/zap"
	"path"
)

// Boot is High Priority
func boot(w container.Container) {
	w.Singleton(context.App, w.App())
	w.Singleton(context.Config, w.Config())
	w.Singleton(context.Db, w.Database())

	jwtGenerator := func() *jwt.GeneratorStruct {
		dataPath := w.Config().Path.DataPath
		privateKey := rsGen.PrivatePemDecode(path.Join(dataPath, "secret/private.pem"))

		return &jwt.GeneratorStruct{
			PrivateKey: privateKey,
			PublicKey:  privateKey.Public(),
			Exp:        w.Config().Auth.Exp,
		}
	}

	var tg jwt.Generator
	w.Bind(&tg, jwtGenerator)
	w.Resolve(&tg)
	w.Singleton(context.JwtGenerator, tg)

	var logs *zap.SugaredLogger
	loggerConfig := w.Config().CustomLogger
	w.Bind(&logs, logger.New(logger.Config{
		TimeFormat: loggerConfig.TimeFormat,
		FilePath:   loggerConfig.FilePath,
		Filename:   loggerConfig.Filename,
		MaxAge:     loggerConfig.MaxAge,
		MaxBackups: loggerConfig.MaxBackups,
		MaxSize:    loggerConfig.MaxSize,
		Compress:   loggerConfig.Compress,
		TimeKey:    loggerConfig.TimeKey,
		TimeZone:   loggerConfig.TimeZone,
		LogLevel:   loggerConfig.LogLevel,
	}))
	w.Resolve(&logs)
	w.Singleton(context.Logger, logs)
}

// Middlewares register middleware
func middlewares(w container.Container) {
	w.App().Use(flogger.New(w.Config().Logger))
	w.App().Use(recover.New())

	// Add Context Container
	w.App().Use(func(ctx *fiber.Ctx) error {
		ctx.Locals(context.Container, w)
		return ctx.Next()
	})

	// Add Context Config
	w.App().Use(func(ctx *fiber.Ctx) error {
		ctx.Locals(context.Config, w.Config())
		return ctx.Next()
	})

	// Add Context Logger
	w.App().Use(func(ctx *fiber.Ctx) error {
		zLogger := w.Get(context.Logger).(*zap.SugaredLogger)

		ctx.Locals(context.Logger, zLogger)
		return ctx.Next()
	})

	w.App().Use(api_error.ErrorHandler)
	w.App().Use(cors.New(w.Config().Cors))
}

// healthCHeckRes health check response
type healthCheckRes struct {
	Status bool `json:"status"`
}

// healthCheck
// @Summary health check your server
// @Description health check your server
// @Success 200 {object} healthCheckRes
// @Tags healthCheck
// @Accept */*
// @Produce json
// @Router /health-check [get]
func healthCheck(ctx *fiber.Ctx) error {

	err := ctx.JSON(healthCheckRes{Status: true})
	if err != nil {
		return ctx.JSON(healthCheckRes{Status: false})
	}

	return err
}

// routes register Routes
func routes(w container.Container) {
	router.Api(w)
	w.App().Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))
	w.App().Get("/health-check", healthCheck)
	w.App().Get("/swagger/*", swagger.HandlerDefault)

}

func Resister(w container.Container) {
	boot(w)
	middlewares(w)
	routes(w)
}
