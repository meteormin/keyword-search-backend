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

var gLoggers = make(map[string]*zap.Logger)

func GetLogger(loggerName ...string) *zap.SugaredLogger {
	if gLoggers == nil {
		return New()
	}

	if len(loggerName) == 0 {
		return gLoggers[Default].Named(Default).Sugar()
	}

	return gLoggers[loggerName[0]].Named(loggerName[0]).Sugar()
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
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.LevelKey = "level"
	encoderConfig.StacktraceKey = "stacktrace"
	encoderConfig.CallerKey = "caller"
	encoderConfig.MessageKey = "msg"
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

	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), ws, cfg.LogLevel)
	zapLogger := zap.New(core, zap.AddCaller())
	logger := zapLogger.Named(cfg.Name).Sugar()
	gLoggers[cfg.Name] = zapLogger

	return logger
}
