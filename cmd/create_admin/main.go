package main

import (
	"github.com/miniyus/gofiber/database"
	"github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/admin"
)

func main() {
	cfg := config.GetConfigs()
	db := database.New(cfg.Database["default"])
	admin.CreateAdmin(db, &cfg)
}
