package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var l *zap.Logger

func InitLogger(level, format string) error {
	var config zap.Config

	if format == "json" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	logger, err := config.Build()
	if err != nil {
		return err
	}
	l = logger
	zap.ReplaceGlobals(logger)
	return nil
}

func Sync() {
	_ = l.Sync()
}
