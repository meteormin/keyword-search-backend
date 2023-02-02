package internal

import (
	baseContext "context"
	"encoding/json"
	"github.com/go-redis/redis/v9"
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/api_error"
	jobWorker "github.com/miniyus/keyword-search-backend/pkg/worker"
	"github.com/miniyus/keyword-search-backend/utils"
)

func HandleValidate(c *fiber.Ctx, data interface{}) *api_error.ValidationErrorResponse {
	err := c.BodyParser(data)
	if err != nil {
		errRes := api_error.NewValidationErrorResponse(c, map[string]string{
			"parse_error": err.Error(),
		})
		return errRes
	}

	failed := utils.Validate(data)
	if failed != nil {
		errRes := api_error.NewValidationErrorResponse(c, failed)
		return errRes
	}

	return nil
}

func FindJobFromQueueWorker(jobDispatcher jobWorker.Dispatcher) func(ctx *fiber.Ctx, jobId string, worker ...string) (*jobWorker.Job, error) {
	return func(ctx *fiber.Ctx, jobId string, worker ...string) (*jobWorker.Job, error) {
		workerName := jobWorker.DefaultWorker

		if len(worker) != 0 {
			workerName = worker[0]
		}

		jobDispatcher.SelectWorker(workerName)

		redisClient := jobDispatcher.GetRedis()()

		var convJob *jobWorker.Job
		value, err := redisClient.Get(baseContext.Background(), jobId).Result()
		if err == redis.Nil {
			return nil, nil
		} else if err != nil {
			return nil, err
		} else {

			bytes := []byte(value)
			err = json.Unmarshal(bytes, &convJob)
			if err != nil {
				return nil, err
			}
		}

		return convJob, nil
	}
}
