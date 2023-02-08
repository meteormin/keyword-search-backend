package config

import (
	mConfig "github.com/miniyus/gofiber/config"
)

var cfg *mConfig.Configs

func init() {
	cfg = &mConfig.Configs{
		App:            appConfig(),
		Logger:         flogger(),
		CustomLogger:   loggerConfig(),
		Database:       databaseConfig(),
		Path:           getPath(),
		Auth:           auth(),
		Cors:           cors(),
		Csrf:           csrf(),
		Permission:     permissionConfig(),
		CreateAdmin:    createAdminConfig(),
		RedisConfig:    redisConfig(),
		JobQueueConfig: jobQueueConfig(),
		Validation:     validationConfig(),
	}
}

func GetConfigs() mConfig.Configs {
	return *cfg
}
