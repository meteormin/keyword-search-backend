package migrations

import (
	"github.com/miniyus/keyword-search-backend/entity"
	"gorm.io/gorm"
	"log"
)

// Migrate
// db entity 스키마에 맞춰 자동으로 migration
func Migrate(db *gorm.DB) {
	log.Println("Auto Migrate...")
	err := db.AutoMigrate(
		&entity.User{},
		&entity.Group{},
		&entity.Permission{},
		&entity.Action{},
		&entity.AccessToken{},
		&entity.Host{},
		&entity.BookMark{},
		&entity.Search{},
		&entity.Tag{},
		&entity.GroupDetail{},
	)

	if err != nil {
		log.Fatalf("Failed Auto Migration")
	}

	log.Println("Success Auto Migrate...")
}
