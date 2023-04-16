package tasks

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/gofiber/schedule"
)

const Prefix = "tasks"

func Register() app.SubRouter {
	return func(router fiber.Router) {
		router.Get("scheduler", getScheduler)
		router.Get("jobs", getJobs)
		router.Post("jobs/:jobId", handleJob)
	}
}

func getScheduler(ctx *fiber.Ctx) error {
	worker := schedule.GetWorker()
	stats := worker.Stats()
	return ctx.JSON(stats)
}

func getJobs(ctx *fiber.Ctx) error {
	js := GetJobs()
	stats := js.container.Status()

	var jobIds []string
	for jobId := range js.container.Jobs() {
		jobIds = append(jobIds, jobId)
	}

	return ctx.JSON(map[string]interface{}{
		"jobs":   jobIds,
		"status": stats,
	})
}

func handleJob(ctx *fiber.Ctx) error {
	params := ctx.AllParams()
	jobId := params["jobId"]

	isFind := false
	for jId := range GetJobs().container.Jobs() {
		if jobId == jId {
			jobId = jId
			isFind = true
		}
	}

	if !isFind {
		return fiber.ErrNotFound
	}

	dispatch, err := GetJobs().SyncDispatch(jobId)
	if err != nil {
		return err
	}

	return ctx.JSON(dispatch)
}
