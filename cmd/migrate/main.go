package main

import (
	_ "github.com/miniyus/gofiber/app"
	"github.com/miniyus/gofiber/database"
	"github.com/miniyus/gofiber/database/migrations"
	"github.com/miniyus/keyword-search-backend/config"
)

func main() {
	configure := config.GetConfigs()

	db := database.New(configure.Database["default"])

	migrations.Migrate(db)
}
