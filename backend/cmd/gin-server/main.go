package main

import (
	"log"

	"coraza-waf/backend/config"
	"coraza-waf/backend/internal/agent"
	"coraza-waf/backend/pkg/database"
	"github.com/corazawaf/coraza/v3"
)

func main() {
	// 1. 加载配置
	if err := config.Init("config.yaml"); err != nil {
		log.Fatalf("Failed to load config.yaml: %v", err)
	}
	cfg := config.Get()

	// 2. 初始化 Coraza WAF
	waf, err := coraza.NewWAF(coraza.NewWAFConfig())
	if err != nil {
		log.Fatalf("Failed to initialize WAF: %v", err)
	}

	// 3. 初始化 MongoDB 日志服务（地址、库名、集合名可以根据实际修改）
	mongo, err := database.NewMongoService("mongodb://localhost:27017", "wafdb", "waflogs")
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// 4. 创建 Agent Handler
	handler := agent.NewAgent(waf, mongo)

	// 5. 启动 SPOA Server，使用 config.yaml 中 spoa.bind 字段
	if err := agent.StartServer(cfg.SPOA.Bind, handler); err != nil {
		log.Fatalf("Failed to start SPOA server: %v", err)
	}

	log.Println("Coraza SPOE Agent started at", cfg.SPOA.Bind)
}

