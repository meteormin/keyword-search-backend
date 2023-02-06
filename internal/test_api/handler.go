package test_api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/app"
	"github.com/miniyus/keyword-search-backend/log"
	jobWorker "github.com/miniyus/keyword-search-backend/pkg/worker"
	"github.com/miniyus/keyword-search-backend/utils"
	"time"
)

const Prefix = "/test"

// Register
// @Summary test api
// @Description test api
// @Tags Test
// @Router /api/test [post]
// @Success 200 {object} utils.StatusResponse
func Register(dispatcher jobWorker.Dispatcher) app.SubRouter {
	logger := log.GetLogger()
	return func(router fiber.Router) {
		router.Post("/", func(ctx *fiber.Ctx) error {
			logger.Infof(ctx.Path())
			err := dispatcher.SelectWorker(jobWorker.DefaultWorker).Dispatch("test", func(j *jobWorker.Job) error {
				jStr, jErr := j.Marshal()
				if jErr != nil {
					return jErr
				}

				logger.Infof("job: %s", jStr)
				time.Sleep(time.Second * 3)
				logger.Infof("job: %s", jStr)
				return nil
			})

			logger.Infof("run dispatcher")

			if err != nil {
				return err
			}

			return ctx.JSON(utils.StatusResponse{Status: true})
		})
	}
}
