package config

import (
	"github.com/go-redis/redis/v9"
	"github.com/gofiber/fiber/v2"
	fCors "github.com/gofiber/fiber/v2/middleware/cors"
	fCsrf "github.com/gofiber/fiber/v2/middleware/csrf"
	loggerMiddleware "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/miniyus/keyword-search-backend/pkg/worker"
	"log"
	"os"
	"strconv"
)

type Env string

const (
	PRD   Env = "production"
	DEV   Env = "development"
	LOCAL Env = "local"
)

type Configs struct {
	AppEnv       Env
	AppPort      int
	Locale       string
	TimeZone     string
	App          fiber.Config
	Logger       loggerMiddleware.Config
	CustomLogger LoggerConfig
	Database     DB
	Path         Path
	Auth         Auth
	Cors         fCors.Config
	Csrf         fCsrf.Config
	Permission   []PermissionConfig
	CreateAdmin  CreateAdminConfig
	RedisConfig  *redis.Options
	QueueConfig  worker.DispatcherOption
}

func GetConfigs() *Configs {
	port, err := strconv.Atoi(os.Getenv("APP_PORT"))

	if err != nil {
		log.Printf("App Port is not numeric... %v", err)
		port = 8000
	}

	return &Configs{
		AppEnv:       Env(os.Getenv("APP_ENV")),
		AppPort:      port,
		Locale:       os.Getenv("LOCALE"),
		TimeZone:     os.Getenv("TIME_ZONE"),
		App:          app(),
		Logger:       flogger(),
		CustomLogger: logger(),
		Database:     database(),
		Path:         getPath(),
		Auth:         auth(),
		Cors:         cors(),
		Csrf:         csrf(),
		Permission:   getPermissions(),
		CreateAdmin:  createAdminConfig(),
		RedisConfig:  redisConfig(),
		QueueConfig:  QueueConfig(),
	}
}
