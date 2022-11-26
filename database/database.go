package database

import (
	"fmt"
	"github.com/miniyus/go-fiber/config"
	"github.com/miniyus/go-fiber/database/migrations"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func DB(config config.DB) *gorm.DB {
	var sslMode string = "disable"
	if config.SSLMode {
		sslMode = "enable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.Host, config.Username, config.Password, config.Dbname, config.Port, sslMode, config.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed: Connect DB %v", err)
	}

	log.Println("Success: Connect DB")

	if config.AutoMigrate {
		migrations.Migrate(db)
	}

	return db
}
