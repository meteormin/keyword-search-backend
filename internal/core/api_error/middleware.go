package api_error

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/internal/core/context"
	"go.uber.org/zap"
	"os"
	"runtime/debug"
)

var appEnv = os.Getenv("APP_ENV")

func OverrideDefaultErrorHandler(ctx *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	var logger *zap.SugaredLogger
	logger, ok := ctx.Locals(context.Logger).(*zap.SugaredLogger)

	if !ok {
		errRes := NewErrorResponse(ctx, fiber.StatusInternalServerError, "Can not found context.Logger")
		return errRes.Response()
	}

	logger.Errorln(err)
	if appEnv != "production" {
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

	var logger *zap.SugaredLogger
	logger, ok := ctx.Locals(context.Logger).(*zap.SugaredLogger)

	if !ok {
		errRes := NewErrorResponse(ctx, fiber.StatusInternalServerError, "Can not found context.Logger")
		return errRes.Response()
	}

	logger.Errorln(err)
	if appEnv != "production" {
		debug.PrintStack()
	}

	errRes := NewFromError(ctx, err)

	b, _ := json.Marshal(errRes)

	logger.Errorln(string(b))

	return errRes.Response()
}
