package migrations

import (
	"github.com/miniyus/go-fiber/internal/entity"
	"gorm.io/gorm"
	"log"
)

func Migrate(db *gorm.DB) {
	log.Println("Auto Migrate...")
	err := db.AutoMigrate(
		&entity.TestEntity{},
		&entity.User{},
	)

	if err != nil {
		log.Fatalf("Failed Auto Migration")
	}

	log.Println("Success Auto Migrate...")
}
