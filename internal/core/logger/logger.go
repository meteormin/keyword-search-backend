package logger

import (
	"fmt"
	"github.com/miniyus/go-fiber/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"path"
	"time"
)

var log *zap.SugaredLogger

func init() {
	today := time.Now().Format("2006-01-02")
	logFilename := path.Join(config.GetPath().LogPath, fmt.Sprintf("log-%s.log", today))

	ll := &lumberjack.Logger{
		Filename:   logFilename,
		MaxSize:    10,
		MaxBackups: 30,
		MaxAge:     30,
		Compress:   true,
	}

	ws := zapcore.AddSync(ll)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"

	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		type appendTimeEncoder interface {
			AppendTimeLayout(time.Time, string)
		}

		if enc, ok := enc.(appendTimeEncoder); ok {
			enc.AppendTimeLayout(t, "2006-01-02 15:04:05")
			return
		}

		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.StacktraceKey = ""

	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), ws, zap.InfoLevel)
	logger := zap.New(core, zap.AddCaller())
	log = logger.Sugar()
}

func GetLogger() *zap.SugaredLogger {
	return log
}
