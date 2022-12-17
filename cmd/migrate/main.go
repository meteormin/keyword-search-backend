package main

import (
	"github.com/miniyus/go-fiber/config"
	"github.com/miniyus/go-fiber/internal/core/database"
	"github.com/miniyus/go-fiber/internal/core/database/migrations"
)

func main() {
	configure := config.GetConfigs()

	db := database.DB(configure.Database)

	migrations.Migrate(db)
}
