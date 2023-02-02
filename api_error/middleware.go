package api_error

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	configure "github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/utils"
	"go.uber.org/zap"
	"runtime/debug"
)

func OverrideDefaultErrorHandler(cfg *configure.Configs) fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {

		if err == nil {
			return nil
		}

		var config *configure.Configs
		var insideErr error
		if cfg == nil {
			config, insideErr = utils.GetContext[*configure.Configs](ctx, utils.ConfigsKey)
			if err != nil {
				errRes := NewFromError(ctx, insideErr)
				return errRes.Response()
			}

		} else {
			config = cfg
		}

		var logger *zap.SugaredLogger
		logger, insideErr = utils.GetContext[*zap.SugaredLogger](ctx, utils.LoggerKey)

		if insideErr != nil {
			errRes := NewFromError(ctx, insideErr)
			return errRes.Response()
		}

		logger.Errorln(err)
		if config.AppEnv != configure.PRD {
			debug.PrintStack()
		}

		errRes := NewFromError(ctx, err)
		return errRes.Response()
	}
}

func ErrorHandler(cfg *configure.Configs) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := ctx.Next()

		if err == nil {
			return nil
		}

		var config *configure.Configs
		var insideErr error
		if cfg == nil {
			config, insideErr = utils.GetContext[*configure.Configs](ctx, utils.ConfigsKey)
			if err != nil {
				errRes := NewFromError(ctx, insideErr)
				return errRes.Response()
			}
		} else {
			config = cfg
		}

		var logger *zap.SugaredLogger
		logger, insideErr = utils.GetContext[*zap.SugaredLogger](ctx, utils.LoggerKey)

		if insideErr != nil {
			errRes := NewFromError(ctx, insideErr)
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
}
