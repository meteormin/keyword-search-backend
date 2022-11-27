package api_error

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/config"
	"go.uber.org/zap"
)

func OverrideDefaultErrorHandler(ctx *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}
	errRes := NewFromError(err)
	return errRes.Response(ctx)
}

func ErrorHandler(ctx *fiber.Ctx) error {
	err := ctx.Next()

	if err == nil {
		return nil
	}

	logger, ok := ctx.Locals(config.Logger).(*zap.SugaredLogger)
	if !ok {
		return OverrideDefaultErrorHandler(ctx, err)
	}

	logger.Errorln(err)

	errRes := NewFromError(err)

	b, _ := json.Marshal(errRes)

	logger.Errorln(string(b))

	return errRes.Response(ctx)
}
