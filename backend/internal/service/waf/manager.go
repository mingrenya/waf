package waf

import "github.com/corazawaf/coraza/v3"

// NewWAF initializes Coraza with inline SecLang rules
func NewWAF(directives string) (coraza.WAF, error) {
	cfg := coraza.NewWAFConfig().WithDirectives(directives)
	return coraza.NewWAF(cfg)
}

