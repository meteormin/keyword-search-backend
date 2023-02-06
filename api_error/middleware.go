package api_error

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/app"
	"github.com/miniyus/keyword-search-backend/log"
	"runtime/debug"
)

func OverrideDefaultErrorHandler(appEnv app.Env) fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {

		if err == nil {
			return nil
		}

		logger := log.GetLogger()

		logger.Errorln(err)
		if appEnv != app.PRD {
			debug.PrintStack()
		}

		errRes := NewFromError(ctx, err)
		return errRes.Response()
	}
}

func ErrorHandler(appEnv app.Env) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := ctx.Next()

		if err == nil {
			return nil
		}

		logger := log.GetLogger()

		logger.Errorln(err)
		if appEnv != app.PRD {
			debug.PrintStack()
		}

		errRes := NewFromError(ctx, err)

		b, _ := json.Marshal(errRes)

		logger.Errorln(string(b))

		return errRes.Response()
	}
}
