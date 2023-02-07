package main

import (
	"github.com/miniyus/gofiber/create_admin"
	"github.com/miniyus/gofiber/database"
	"github.com/miniyus/keyword-search-backend/config"
)

func main() {
	cfg := config.GetConfigs()
	db := database.New(cfg.Database["default"])

	create_admin.CreateAdmin(db, &cfg)
}
