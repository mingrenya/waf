package main

import (
	"log"

	"coraza-waf/backend/config"
	"coraza-waf/backend/internal/utils"

	"go.uber.org/zap"
)

const configPath = "config.yaml"

func main() {
	// 读取配置
	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed loading config: %v", err)
	}

	// 创建日志
	logCfg := utils.LogConfig{Level: cfg.LogLevel}
	logger, err := logCfg.NewLogger()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting Coraza-SPOA...")
	logger.Info("Listening address", zap.String("bind", cfg.Bind))

	// 打印每个应用配置
	for _, app := range cfg.Applications {
		logger.Info("Loaded application",
			zap.String("name", app.Name),
			zap.String("address", app.Address),
			zap.Int("workers", app.Workers),
			zap.Bool("response_check", app.ResponseCheck),
			zap.Int("ttl_ms", app.TransactionTTLMS),
		)
	}

	// 热加载配置
	err = config.WatchConfig(configPath, func() {
		newCfg, err := config.ReadConfig(configPath)
		if err != nil {
			logger.Error("Failed to reload config", zap.Error(err))
			return
		}
		logger.Info("Config reloaded successfully", zap.String("default_app", newCfg.DefaultApplication))
	})
	if err != nil {
		logger.Error("Failed to watch config file", zap.Error(err))
	}

	// 保持运行
	select {}
}

