package register

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	flogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/miniyus/go-fiber/internal/context"
	"github.com/miniyus/go-fiber/internal/core/api_error"
	"github.com/miniyus/go-fiber/internal/core/container"
	"github.com/miniyus/go-fiber/internal/core/logger"
	router "github.com/miniyus/go-fiber/internal/routes"
	"github.com/miniyus/go-fiber/pkg/jwt"
	rsGen "github.com/miniyus/go-fiber/pkg/rs256"
	"go.uber.org/zap"
	"path"
)

// Boot is High Priority
func boot(w container.Container) {
	w.Inject("app", w.App())
	w.Inject("config", w.Config())
	w.Inject("db", w.Database())

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
	w.Inject(context.JwtGenerator, tg)

	var logs *zap.SugaredLogger
	loggerConfig := w.Config().CustomLogger
	w.Bind(&logs, logger.NewLogger(logger.Config{
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
	w.Inject(context.Logger, logs)
}

// Middlewares register middleware
func middlewares(w container.Container) {
	w.App().Use(flogger.New(w.Config().Logger))
	w.App().Use(recover.New())
	w.App().Use(func(ctx *fiber.Ctx) error {
		ctx.Locals(context.Container, w)
		return ctx.Next()
	})

	w.App().Use(func(ctx *fiber.Ctx) error {
		ctx.Locals(context.Config, w.Config())
		return ctx.Next()
	})

	w.App().Use(func(ctx *fiber.Ctx) error {
		zLogger := w.Get(context.Logger).(*zap.SugaredLogger)

		ctx.Locals(context.Logger, zLogger)
		return ctx.Next()
	})

	w.App().Use(api_error.ErrorHandler)
	w.App().Use(cors.New(w.Config().Cors))
}

// Routes register Routes
func routes(w container.Container) {
	router.SetRoutes(w)
	w.App().Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))
}

func Resister(w container.Container) {
	boot(w)
	middlewares(w)
	routes(w)
}
