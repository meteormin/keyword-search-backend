package tests

import (
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/gofiber/database"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	if app.App() == nil {
		app.New(
			app.Config{
				Env: app.TEST,
			},
		)
	}

	db = database.New(
		database.Config{
			Name:   "sqlite",
			Driver: "sqlite",
			Dbname: "test",
		},
	).Debug()
}

func Migration(entities ...interface{}) {
	err := db.AutoMigrate(entities...)
	if err != nil {
		panic(err)
	}
}
