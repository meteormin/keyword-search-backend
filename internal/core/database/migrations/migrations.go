package migrations

import (
	"github.com/miniyus/keyword-search-backend/internal/entity"
	"gorm.io/gorm"
	"log"
)

// Migrate
// db entity 스키마에 맞춰 자동으로 migration
func Migrate(db *gorm.DB) {
	log.Println("Auto Migrate...")
	err := db.AutoMigrate(
		&entity.User{},
		&entity.AccessToken{},
		&entity.Host{},
		&entity.BookMark{},
		&entity.Search{},
		&entity.Tag{},
		&entity.Permission{},
	)

	if err != nil {
		log.Fatalf("Failed Auto Migration")
	}

	log.Println("Success Auto Migrate...")
}
