package config

import (
	"github.com/miniyus/gofiber/database"
	"github.com/miniyus/keyword-search-backend/entity"
	gormLogger "gorm.io/gorm/logger"
	"os"
	"strconv"
	"time"
)

func databaseConfig() map[string]database.Config {
	autoMigrate, err := strconv.ParseBool(os.Getenv("DB_AUTO_MIGRATE"))

	var autoMigrateEntities []interface{}

	if err != nil {
		autoMigrate = false
	}

	if autoMigrate {
		autoMigrateEntities = []interface{}{
			&entity.AccessToken{},
			&entity.Action{},
			&entity.Permission{},
			&entity.Group{},
			&entity.GroupDetail{},
			&entity.User{},
			&entity.Host{},
			&entity.Search{},
			&entity.BookMark{},
			&entity.Tag{},
			&entity.JobHistory{},
		}
	}

	return map[string]database.Config{
		"default": {
			Name:        "default",
			Driver:      "postgres",
			Host:        os.Getenv("DB_HOST"),
			Dbname:      os.Getenv("DB_DATABASE"),
			Username:    os.Getenv("DB_USERNAME"),
			Password:    os.Getenv("DB_PASSWORD"),
			Port:        os.Getenv("DB_PORT"),
			TimeZone:    os.Getenv("TIME_ZONE"),
			SSLMode:     false,
			AutoMigrate: autoMigrateEntities,
			Logger: gormLogger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  gormLogger.Silent,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
			MaxIdleConn: 10,
			MaxOpenConn: 100,
			MaxLifeTime: time.Hour,
		},
	}
}
