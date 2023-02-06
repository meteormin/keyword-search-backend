package log

import (
	"github.com/miniyus/keyword-search-backend/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"path"
	"time"
)

const Default string = "default"

var loggers map[string]*zap.SugaredLogger

func GetLogger(loggerName ...string) *zap.SugaredLogger {
	if loggers == nil {
		return New()
	}

	if len(loggerName) == 0 {
		return loggers[Default]
	}

	return loggers[loggerName[0]]
}

func New(config ...Config) *zap.SugaredLogger {
	cfg := getDefaultConfig(config...)

	logFilename := path.Join(cfg.FilePath, cfg.Filename)

	ll := &lumberjack.Logger{
		Filename:   logFilename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	ws := zapcore.AddSync(ll)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = cfg.TimeKey

	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		t = utils.TimeIn(t, cfg.TimeZone)
		type appendTimeEncoder interface {
			AppendTimeLayout(time.Time, string)
		}

		if enc, ok := enc.(appendTimeEncoder); ok {
			enc.AppendTimeLayout(t, cfg.TimeFormat)
			return
		}

		enc.AppendString(t.Format(cfg.TimeFormat))
	}

	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.StacktraceKey = ""

	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), ws, cfg.LogLevel)
	zapLogger := zap.New(core, zap.AddCaller())
	logger := zapLogger.Sugar()
	loggers[cfg.Name] = logger

	return logger
}
