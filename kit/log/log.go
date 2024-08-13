package log

import (
	"Backend/kit/enum"
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"sync"
)

type ctxKey struct{}

var (
	once        sync.Once
	Logger      *zap.Logger
	logsFolder  = enum.LogFolder
	logFileName = enum.LogFileName
)

// Get initializes a zap.Logger instance if it has not been initialized
// already and returns the same instance for subsequent calls.
func Get(path string) *zap.Logger {
	once.Do(func() {
		stdout := zapcore.AddSync(os.Stdout)

		// Get the absolute path to the logs directory

		logsDir := filepath.Join(path, logsFolder)
		logFilePath := filepath.Join(logsDir, logFileName)
		logFile := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logFilePath,
			MaxSize:    5, // megabytes
			MaxBackups: 3,
			MaxAge:     14, // days
			Compress:   true,
		})

		logLevel := zap.NewAtomicLevelAt(zap.InfoLevel)

		productionCfg := zap.NewProductionEncoderConfig()
		productionCfg.TimeKey = enum.LogTimeKey
		productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

		developmentCfg := zap.NewDevelopmentEncoderConfig()
		developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

		consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
		fileEncoder := zapcore.NewJSONEncoder(productionCfg)

		// Log to multiple destinations (console and file)
		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, stdout, logLevel),
			zapcore.NewCore(fileEncoder, logFile, logLevel),
		)

		Logger = zap.New(core)
	})

	return Logger
}

// FromCtx returns the Logger associated with the ctx. If no logger
// is associated, the default logger is returned, unless it is nil
// in which case a disabled logger is returned.
func FromCtx(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return l
	} else if l := Logger; l != nil {
		return l
	}

	return zap.NewNop()
}

func StartUpError(l *zap.Logger, message string) {
	l.Fatal(message)
}
