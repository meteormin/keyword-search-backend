package logger

import (
	"github.com/miniyus/go-fiber/internal/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"path"
	"time"
)

type Config struct {
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

func NewLogger(config Config) *zap.SugaredLogger {
	logFilename := path.Join(config.FilePath, config.Filename)

	ll := &lumberjack.Logger{
		Filename:   logFilename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	ws := zapcore.AddSync(ll)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = config.TimeKey

	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		t = utils.TimeIn(t, config.TimeZone)
		type appendTimeEncoder interface {
			AppendTimeLayout(time.Time, string)
		}

		if enc, ok := enc.(appendTimeEncoder); ok {
			enc.AppendTimeLayout(t, config.TimeFormat)
			return
		}

		enc.AppendString(t.Format(config.TimeFormat))
	}

	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.StacktraceKey = ""

	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), ws, config.LogLevel)
	logger := zap.New(core, zap.AddCaller())
	return logger.Sugar()
}
