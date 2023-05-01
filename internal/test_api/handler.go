package test_api

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/gofiber/log"
	"github.com/miniyus/gofiber/utils"
	jobWorker "github.com/miniyus/goworker"
	"github.com/redis/go-redis/v9"
	"time"
)

const Prefix = "/test"

// Register
// @Summary test api
// @Description test api
// @Tags Test
// @Router /api/test [post]
// @Success 200 {object} utils.StatusResponse
func Register(dispatcher jobWorker.Dispatcher, client *redis.Client) app.SubRouter {
	logger := log.GetLogger()
	testContext := context.Background()

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
				client.Set(testContext, "TEST.API", time.Now(), time.Minute)
				return nil
			})

			logger.Infof("run dispatcher")

			if err != nil {
				return err
			}
			client.Get(testContext, "TEST.API")
			return ctx.JSON(utils.StatusResponse{Status: true})
		})
	}
}
