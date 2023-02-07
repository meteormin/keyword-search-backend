package config

import (
	loggerMiddleware "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/miniyus/gofiber/log"
	"go.uber.org/zap/zapcore"
	"os"
)

func flogger() loggerMiddleware.Config {
	return loggerMiddleware.Config{
		Format:     "[${time}] ${ip}:${port} | (${pid}) ${status} - ${method} ${path}\n",
		TimeZone:   os.Getenv("TIME_ZONE"),
		TimeFormat: "2006-01-02 15:04:05",
	}
}

func loggerConfig() map[string]log.Config {
	filePath := getPath().LogPath

	return map[string]log.Config{
		"default": {
			Name:       "default",
			TimeFormat: "2006-01-02 15:04:05",
			FilePath:   filePath,
			Filename:   "log.log",
			MaxSize:    10,
			MaxBackups: 30,
			MaxAge:     30,
			Compress:   true,
			TimeKey:    "timestamp",
			TimeZone:   os.Getenv("TIME_ZONE"),
			LogLevel:   zapcore.DebugLevel,
		},
		"default_worker": {
			Name:       "default_worker",
			TimeFormat: "2006-01-02 15:04:05",
			FilePath:   filePath,
			Filename:   "default_worker.log",
			MaxSize:    10,
			MaxBackups: 30,
			MaxAge:     30,
			Compress:   true,
			TimeKey:    "timestamp",
			TimeZone:   os.Getenv("TIME_ZONE"),
			LogLevel:   zapcore.DebugLevel,
		},
	}
}
