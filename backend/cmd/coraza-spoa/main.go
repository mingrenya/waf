package main

import (
	"log"

	"github.com/corazawaf/coraza/v3"
	"coraza-waf/backend/internal/spoa"
	"coraza-waf/backend/pkg/database"
)

func main() {
	// 初始化 WAF
	waf, err := coraza.NewWAF(coraza.NewWAFConfig())
	if err != nil {
		log.Fatalf("Failed to create WAF: %v", err)
	}

	// 初始化 MongoDB（你可能用的是 localhost:27017）
	mongo, err := database.NewMongoService("mongodb://localhost:27017", "wafdb", "waflogs")
	if err != nil {
		log.Fatalf("Failed to connect MongoDB: %v", err)
	}

	// 启动 SPOE Server（注意传入 addr）
	server := spoa.NewServer("127.0.0.1:8080", waf, mongo)

	// 正确调用 Run()
	if err := server.Run(); err != nil {
		log.Fatalf("SPOE Server error: %v", err)
	}
}

