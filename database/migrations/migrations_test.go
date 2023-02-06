package migrations_test

import (
	"github.com/miniyus/keyword-search-backend/database"
	"github.com/miniyus/keyword-search-backend/database/migrations"
	gormLogger "gorm.io/gorm/logger"
	"testing"
	"time"
)

func TestMigrate(t *testing.T) {
	cfg := database.Config{
		Name:        "default",
		Driver:      "postgres",
		Host:        "localhost",
		Dbname:      "go_fiber",
		Username:    "smyoo",
		Password:    "smyoo",
		Port:        "5432",
		TimeZone:    "Asia/Seoul",
		SSLMode:     false,
		AutoMigrate: false,
		Logger: gormLogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormLogger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
		MaxIdleConn: 10,
		MaxOpenConn: 10,
		MaxLifeTime: time.Hour,
	}

	db := database.New(cfg)

	migrations.Migrate(db)
}
