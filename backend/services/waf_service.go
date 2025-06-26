package services

import (
	"fmt"

	"github.com/corazawaf/coraza/v3"
	"github.com/corazawaf/coraza/v3/debuglog"
	"github.com/corazawaf/coraza/v3/types"
	"go.uber.org/zap"
)

type WAFService struct {
	waf    coraza.WAF
	logger *zap.Logger
}

// NewWAFService 创建新的WAF服务实例
func NewWAFService(logger *zap.Logger, directives string) (*WAFService, error) {
	// 创建WAF配置
	config := coraza.NewWAFConfig().
		WithErrorCallback(logError(logger)).
		WithDebugLogger(debuglog.Default())
	
	// 初始化WAF实例
	waf, err := coraza.NewWAF(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create WAF: %w", err)
	}

	// 直接使用WithDirectives加载规则
	waf, err = coraza.NewWAF(
		config.WithDirectives(`
			SecRuleEngine On
			SecRequestBodyAccess On
			SecDebugLog /var/log/coraza_debug.log
			SecRule ARGS "(union\s+select|sleep\(|\bselect\b.*\bfrom\b)" \
				"id:1001,phase:2,deny,status:403,msg:'SQLi detected'"
			SecRule ARGS|REQUEST_HEADERS "<script>|javascript:" \
				"id:1002,phase:2,deny,msg:'XSS attack detected'"
			SecRule REQUEST_HEADERS:User-Agent "(nikto|w3af|sqlmap)" \
				"id:1003,deny,msg:'Scanner detected'"
			` + directives),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load WAF rules: %w", err)
	}

	return &WAFService{
		waf:    waf,
		logger: logger,
	}, nil
}

func logError(logger *zap.Logger) func(types.MatchedRule) {
	return func(mr types.MatchedRule) {
		logger.Error("WAF blocked request",
			zap.String("message", mr.Message()),
			zap.Int("rule_id", mr.Rule().ID()), // 正确使用 zap.Int
		)
	}
}

// ProcessRequest 处理HTTP请求 (保持不变)
func (s *WAFService) ProcessRequest(clientIP, host, method, path, proto string, headers map[string][]string) (bool, *types.Interruption) {
	tx := s.waf.NewTransaction()
	defer tx.ProcessLogging()
	
	tx.ProcessConnection(clientIP, 0, host, 0)
	tx.ProcessURI(path, method, proto)
	
	for k, values := range headers {
		for _, v := range values {
			tx.AddRequestHeader(k, v)
		}
	}
	
	if it := tx.ProcessRequestHeaders(); it != nil {
		return false, it
	}
	
	return true, nil
}

