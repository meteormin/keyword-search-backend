package config

import (
	mConfig "github.com/miniyus/gofiber/config"
	"github.com/miniyus/keyword-search-backend/internal/permission"
)

var cfg *Configs

type Configs struct {
	*mConfig.Configs
	Permission  []permission.Config
	CreateAdmin CreateAdminConfig
}

func init() {
	cfg = &Configs{
		Configs: &mConfig.Configs{
			App:            appConfig(),
			Logger:         flogger(),
			CustomLogger:   loggerConfig(),
			Database:       databaseConfig(),
			Path:           getPath(),
			Auth:           auth(),
			Cors:           cors(),
			Csrf:           csrf(),
			RedisConfig:    redisConfig(),
			JobQueueConfig: jobQueueConfig(),
			Validation:     validationConfig(),
		},
		Permission:  permissionConfig(),
		CreateAdmin: createAdminConfig(),
	}
}

func GetConfigs() Configs {
	return *cfg
}
