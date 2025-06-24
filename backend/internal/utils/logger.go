package utils

import (
	"go.uber.org/zap"
)

type LogConfig struct {
	Level string
}

func (l LogConfig) NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	switch l.Level {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	return cfg.Build()
}

