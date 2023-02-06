package log

import (
	"go.uber.org/zap/zapcore"
	"os"
)

type Config struct {
	Name       string
	TimeFormat string
	FilePath   string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	TimeKey    string
	TimeZone   string
	LogLevel   zapcore.Level
}

var defaultConfig = Config{
	Name:       "default",
	TimeFormat: "2006-01-02 15:04:05",
	FilePath:   "",
	Filename:   "",
	MaxSize:    10,
	MaxBackups: 30,
	MaxAge:     30,
	Compress:   true,
	TimeKey:    "timestamp",
	TimeZone:   os.Getenv("TIME_ZONE"),
	LogLevel:   zapcore.DebugLevel,
}

func getDefaultConfig(config ...Config) Config {
	if len(config) < 1 {
		return defaultConfig
	}

	cfg := config[0]

	if cfg.Compress {
		cfg.Compress = defaultConfig.Compress
	}

	if cfg.TimeZone == "" {
		cfg.TimeZone = defaultConfig.TimeZone
	}

	if cfg.LogLevel == 0 {
		cfg.LogLevel = defaultConfig.LogLevel
	}

	if cfg.TimeFormat == "" {
		cfg.TimeFormat = defaultConfig.TimeFormat
	}

	if cfg.TimeKey == "" {
		cfg.TimeFormat = defaultConfig.TimeKey
	}

	if cfg.MaxSize == 0 {
		cfg.MaxSize = defaultConfig.MaxSize
	}

	if cfg.MaxAge == 0 {
		cfg.MaxAge = defaultConfig.MaxAge
	}

	if cfg.MaxBackups == 0 {
		cfg.MaxBackups = defaultConfig.MaxBackups
	}

	if cfg.Filename == "" {
		panic("filename field is required")
	}

	if cfg.FilePath == "" {
		panic("filepath field is required")
	}

	return cfg
}
