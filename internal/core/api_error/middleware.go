package api_error

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	configure "github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/core/context"
	"github.com/miniyus/keyword-search-backend/internal/core/register/resolver"
	"go.uber.org/zap"
	"runtime/debug"
)

func OverrideDefaultErrorHandler(ctx *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	var config *configure.Configs
	config, err = resolver.Resolve[*configure.Configs](ctx, config)
	if err != nil {
		errRes := NewErrorResponse(ctx, fiber.StatusInternalServerError, "Can not found context.Config")
		return errRes.Response()
	}

	var logger *zap.SugaredLogger
	logger, err = resolver.Resolve[*zap.SugaredLogger](ctx, logger)

	if err != nil {
		errRes := NewErrorResponse(ctx, fiber.StatusInternalServerError, "Can not found context.Logger")
		return errRes.Response()
	}

	logger.Errorln(err)
	if config.AppEnv != configure.PRD {
		debug.PrintStack()
	}

	errRes := NewFromError(ctx, err)
	return errRes.Response()
}

func ErrorHandler(ctx *fiber.Ctx) error {
	err := ctx.Next()

	if err == nil {
		return nil
	}

	var config *configure.Configs
	config, ok := ctx.Locals(context.Config).(*configure.Configs)

	if !ok {
		errRes := NewErrorResponse(ctx, fiber.StatusInternalServerError, "Can not found context.Config")
		return errRes.Response()
	}

	var logger *zap.SugaredLogger
	logger, ok = ctx.Locals(context.Logger).(*zap.SugaredLogger)

	if !ok {
		errRes := NewErrorResponse(ctx, fiber.StatusInternalServerError, "Can not found context.Logger")
		return errRes.Response()
	}

	logger.Errorln(err)
	if config.AppEnv != configure.PRD {
		debug.PrintStack()
	}

	errRes := NewFromError(ctx, err)

	b, _ := json.Marshal(errRes)

	logger.Errorln(string(b))

	return errRes.Response()
}
