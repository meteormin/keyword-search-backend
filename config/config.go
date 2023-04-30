package config

import (
	fCors "github.com/gofiber/fiber/v2/middleware/cors"
	mConfig "github.com/miniyus/gofiber/config"
	"github.com/miniyus/keyword-search-backend/internal/permission"
)

var cfg *Configs

type Configs struct {
	*mConfig.Configs
	Auth        Auth
	Cors        fCors.Config
	Permission  []permission.Config
	CreateAdmin CreateAdminConfig
}

func init() {
	cfg = &Configs{
		Auth: auth(),
		Configs: &mConfig.Configs{
			App:            appConfig(),
			Logger:         flogger(),
			CustomLogger:   loggerConfig(),
			Database:       databaseConfig(),
			Path:           getPath(),
			RedisConfig:    redisConfig(),
			JobQueueConfig: jobQueueConfig(),
			Validation:     validationConfig(),
		},
		Cors:        cors(),
		Permission:  permissionConfig(),
		CreateAdmin: createAdminConfig(),
	}
}

func GetConfigs() Configs {
	return *cfg
}
