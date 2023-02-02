package config

import (
	"fmt"
	loggerMiddleware "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/miniyus/keyword-search-backend/logger"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

func flogger() loggerMiddleware.Config {

	return loggerMiddleware.Config{
		Format:     "[${time}] ${ip}:${port} | (${pid}) ${status} - ${method} ${path}\n",
		TimeZone:   os.Getenv("TIME_ZONE"),
		TimeFormat: "2006-01-02 15:04:05",
	}
}

func loggerConfig() logger.Config {
	loc, err := time.LoadLocation(os.Getenv("TIME_ZONE"))
	var today string
	if err != nil {
		today = time.Now().Format("2006-01-02")
	} else {
		today = time.Now().In(loc).Format("2006-01-02")
	}

	filePath := getPath().LogPath
	filename := fmt.Sprintf("log-%s.log", today)
	return logger.Config{
		TimeFormat: "2006-01-02 15:04:05",
		FilePath:   filePath,
		Filename:   filename,
		MaxSize:    10,
		MaxBackups: 30,
		MaxAge:     30,
		Compress:   true,
		TimeKey:    "timestamp",
		TimeZone:   os.Getenv("TIME_ZONE"),
		LogLevel:   zapcore.DebugLevel,
	}
}
