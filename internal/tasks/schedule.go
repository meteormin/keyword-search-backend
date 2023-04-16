package tasks

import (
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/gofiber/log"
	"github.com/miniyus/gofiber/schedule"
	"github.com/miniyus/keyword-search-backend/config"
	"time"
)

func RegisterSchedule(app app.Application) {
	var cfg *config.Configs

	app.Resolve(&cfg)

	log.New(log.Config{
		Name:     "tasks_schedule",
		FilePath: cfg.Path.LogPath,
		Filename: "tasks_schedule.log",
	})

	logger := log.GetLogger("tasks_schedule")

	loc, err := time.LoadLocation(cfg.App.TimeZone)
	if err != nil {
		logger.Error(err)
	}

	scheduleWorker := schedule.NewWorker(schedule.WorkerConfig{
		TagsUnique: true,
		Logger:     logger,
		Location:   loc,
	})

	go scheduleWorker.Run()
}
