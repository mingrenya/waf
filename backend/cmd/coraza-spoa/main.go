package main

import (
	"log"

	"coraza-waf/backend/internal/agent"
	"coraza-waf/backend/pkg/database"

	"github.com/corazawaf/coraza/v3"
)

func main() {
	waf, err := coraza.NewWAF(coraza.NewWAFConfig())
	if err != nil {
		log.Fatalf("Failed to create WAF: %v", err)
	}

	mongo, err := database.NewMongoService("mongodb://localhost:27017", "wafdb", "waflogs")
	if err != nil {
		log.Fatalf("Failed to connect MongoDB: %v", err)
	}

	handler := agent.NewAgent(waf, mongo)

	//err = agent.StartServer("127.0.0.1:12345", handler)
	//if err != nil {
	//	log.Fatalf("Failed to start SPOE server: %v", err)
	//}
	//log.Println("WAF Agent started (SPOE Server start code not included)")
}

