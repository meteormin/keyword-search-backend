package logger

import (
	"github.com/miniyus/go-fiber/internal/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"path"
	"time"
)

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
	logger := zap.New(core, zap.AddCaller())
	return logger.Sugar()
}
