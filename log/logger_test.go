package log_test

import (
	"github.com/miniyus/keyword-search-backend/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"testing"
)

var logger *zap.SugaredLogger

func TestNew(t *testing.T) {
	cfg := log.Config{
		Name:       "default",
		TimeFormat: "2006-01-02 15:04:05",
		FilePath:   "./",
		Filename:   "logger_test.log",
		MaxSize:    10,
		MaxBackups: 30,
		MaxAge:     30,
		Compress:   true,
		TimeKey:    "timestamp",
		TimeZone:   os.Getenv("TIME_ZONE"),
		LogLevel:   zapcore.DebugLevel,
	}
	logger = log.New(cfg)

	logger.Infof("%s", "default")
}

func TestGetLogger(t *testing.T) {
	get := log.GetLogger()
	if get != logger {
		t.Error(get)
	}
}
