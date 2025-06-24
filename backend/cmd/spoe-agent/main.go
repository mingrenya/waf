package main

import (
	"coraza-waf/backend/internal/agent"
	"coraza-waf/backend/pkg/waf"
	"log"
)

func main() {
	wafInstance, err := waf.NewWAF([]string{"./rules/*.conf"}, "")
	if err != nil {
		log.Fatalf("failed to initialize WAF: %v", err)
	}

	handler := agent.NewHandler(wafInstance)
	if err := agent.StartServer(12345, handler); err != nil {
		log.Fatalf("failed to start SPOE agent: %v", err)
	}
}
