package register

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	flogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	_ "github.com/miniyus/keyword-search-backend/api/gofiber"
	"github.com/miniyus/keyword-search-backend/internal/core/api_error"
	"github.com/miniyus/keyword-search-backend/internal/core/container"
	"github.com/miniyus/keyword-search-backend/internal/core/context"
	"github.com/miniyus/keyword-search-backend/internal/core/permission"
	"github.com/miniyus/keyword-search-backend/internal/core/register/resolver"
	router "github.com/miniyus/keyword-search-backend/internal/routes"
	"github.com/miniyus/keyword-search-backend/internal/utils"
	"github.com/miniyus/keyword-search-backend/pkg/jwt"
	"go.uber.org/zap"
)

// boot is High Priority
func boot(w container.Container) {
	w.Singleton(context.App, w.App())
	w.Singleton(context.Config, w.Config())
	w.Singleton(context.DB, w.Database())

	var tg jwt.Generator
	jwtGenerator := resolver.MakeJwtGenerator(w)

	w.Bind(&tg, jwtGenerator)
	w.Resolve(&tg)
	w.Singleton(context.JwtGenerator, tg)

	var logs *zap.SugaredLogger
	loggerStruct := resolver.MakeLogger(w)

	w.Bind(&logs, loggerStruct)
	w.Resolve(&logs)
	w.Singleton(context.Logger, logs)

	var permCollect permission.Collection
	permissionCollection := resolver.MakePermissionCollection(w)

	w.Bind(&permCollect, permissionCollection)
	w.Resolve(&permCollect)
	w.Singleton(context.Permissions, permCollect)
}

// middlewares register middleware
func middlewares(w container.Container) {
	w.App().Use(flogger.New(w.Config().Logger))
	w.App().Use(recover.New())

	// Add Context Container
	w.App().Use(resolver.AddContext(context.Container, w))
	// Add Context Config
	w.App().Use(resolver.AddContext(context.Config, w.Config()))
	// Add Context Logger
	w.App().Use(resolver.AddContext(context.Logger, w.Get(context.Logger)))
	// Add Context Permissions
	w.App().Use(resolver.AddContext(context.Permissions, w.Get(context.Permissions)))

	w.App().Use(api_error.ErrorHandler)
	w.App().Use(cors.New(w.Config().Cors))
}

// healthCheck
// @Summary health check your server
// @Description health check your server
// @Success 200 {object} utils.StatusResponse
// @Tags HealthCheck
// @Accept */*
// @Produce json
// @Router /health-check [get]
func healthCheck(ctx *fiber.Ctx) error {

	err := ctx.JSON(utils.StatusResponse{Status: true})
	if err != nil {
		return ctx.JSON(utils.StatusResponse{Status: false})
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

// Resister
// private 함수들 모아서 순서대로 실행 해주는 public 함수
func Resister(w container.Container) {
	boot(w)
	middlewares(w)
	routes(w)
}
