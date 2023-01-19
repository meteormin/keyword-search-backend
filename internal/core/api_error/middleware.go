package api_error

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	configure "github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/core/context"
	"go.uber.org/zap"
	"runtime/debug"
)

func OverrideDefaultErrorHandler(ctx *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	var config *configure.Configs
	config, err = context.GetContext[*configure.Configs](ctx, context.Config)

	if err != nil {
		errRes := NewFromError(ctx, err)
		return errRes.Response()
	}

	var logger *zap.SugaredLogger
	logger, err = context.GetContext[*zap.SugaredLogger](ctx, context.Logger)

	if err != nil {
		errRes := NewFromError(ctx, err)
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
	config, err = context.GetContext[*configure.Configs](ctx, context.Config)

	if err != nil {
		errRes := NewFromError(ctx, err)
		return errRes.Response()
	}

	var logger *zap.SugaredLogger
	logger, err = context.GetContext[*zap.SugaredLogger](ctx, context.Logger)

	if err != nil {
		errRes := NewFromError(ctx, err)
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
