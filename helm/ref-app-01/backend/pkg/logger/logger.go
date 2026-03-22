package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New builds a Zap logger from level string ("debug", "info", "warn", "error").
func New(level string) (*zap.Logger, error) {
	var cfg zap.Config
	if level == "debug" {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
	}
	atomic := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if parsed, err := zapcore.ParseLevel(level); err == nil {
		atomic = zap.NewAtomicLevelAt(parsed)
	}
	cfg.Level = atomic
	return cfg.Build()
}
