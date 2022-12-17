package migrations

import (
	"github.com/miniyus/go-fiber/internal/entity"
	"gorm.io/gorm"
	"log"
)

func Migrate(db *gorm.DB) {
	log.Println("Auto Migrate...")
	err := db.AutoMigrate(
		&entity.User{},
		&entity.AccessToken{},
		&entity.Host{},
		&entity.BookMark{},
		&entity.Search{},
		&entity.Tag{},
	)

	if err != nil {
		log.Fatalf("Failed Auto Migration")
	}

	log.Println("Success Auto Migrate...")
}
