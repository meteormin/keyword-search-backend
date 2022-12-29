package api_error

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/internal/core/context"
	"go.uber.org/zap"
	"runtime/debug"
)

func OverrideDefaultErrorHandler(ctx *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	logger, ok := ctx.Locals(context.Logger).(*zap.SugaredLogger)
	if !ok {
		debug.PrintStack()
		return err
	}

	logger.Errorln(err)
	debug.PrintStack()

	errRes := NewFromError(ctx, err)
	return errRes.Response()
}

func ErrorHandler(ctx *fiber.Ctx) error {
	err := ctx.Next()

	if err == nil {
		return nil
	}

	logger, ok := ctx.Locals(context.Logger).(*zap.SugaredLogger)
	if !ok {
		debug.PrintStack()
		return err
	}

	logger.Errorln(err)
	debug.PrintStack()

	errRes := NewFromError(ctx, err)

	b, _ := json.Marshal(errRes)

	logger.Errorln(string(b))

	return errRes.Response()
}
