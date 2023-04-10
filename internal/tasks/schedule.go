package tasks

import (
	"github.com/go-co-op/gocron"
	"github.com/miniyus/gofiber/app"
	"log"
	"time"
)

func RegisterSchedule() app.Register {
	return func(app app.Application) {
		cfg := app.Config()
		location, err := time.LoadLocation(cfg.TimeZone)
		if err != nil {
			panic(err)
		}

		scheduler := gocron.NewScheduler(location)
		scheduler.TagsUnique()
		_, err = scheduler.Tag("health-check").Cron("* 0 * * *").Do(func() {
			log.Println("scheduler health check")
		})
		if err != nil {
			log.Println("scheduler:", err)
		}
	}
}
