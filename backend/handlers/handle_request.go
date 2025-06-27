package handlers

import (
	"github.com/gin-gonic/gin"
	"coraza-waf/backend/logger"
	"io"
	"bytes"
	"time"
	"log"
)

// HandleRequest 拦截请求并记录日志
func HandleRequest(c *gin.Context) {
	// 保存请求体
	body, _ := c.GetRawData()
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// ====== 规则检测示例（可替换为真实WAF引擎） =====
	// 假设命中规则（实际应调用 wafService 或 coraza 检测）
	var matched bool
	var ruleID, ruleContent, ruleFormat string
	if bytes.Contains(bytes.ToLower(body), []byte("select")) {
		matched = true
		ruleID = "1001"
		ruleContent = "SecRule ARGS \"select\" id:1001,deny,msg:'SQLi'"
		ruleFormat = "modsec"
	}
	// ============================================

	if matched {
		c.Set("matched_rule_id", ruleID)
		c.Set("matched_rule_content", ruleContent)
		c.Set("matched_rule_format", ruleFormat)
		// 联动告警（可扩展为HTTP/钉钉/Prometheus/MCP/AI等）
		// go sendAlarm(ruleID, ruleContent, c.ClientIP())
	}

	logData := logger.WafLog{
		RequestTime:   time.Now().Format(time.RFC3339),
		SourceIP:      c.ClientIP(),
		RequestMethod: c.Request.Method,
		RequestURI:    c.Request.URL.Path,
		UserAgent:     c.Request.UserAgent(),
		RequestHost:   c.Request.Host,
		RequestCookie: c.Request.Header.Get("Cookie"),
		RequestID:     c.GetHeader("X-Request-ID"),
		Blocked:       matched, // 命中规则可直接标记拦截
		RequestBody:   string(body),
		RequestHeaders: c.Request.Header,
		RuleID:        toStrPtr(ruleID),
		RuleContent:   ruleContent,
		RuleFormat:    ruleFormat,
	}

	if err := logger.InsertLog(logData); err != nil {
		log.Printf("写入请求日志失败: %v, data: %+v\n", err, logData)
	}
	c.Next()
}

// sendAlarm 示例（可扩展为实际告警/MCP/AI）
func sendAlarm(ruleID, ruleContent, srcIP string) {
	// log.Printf("[ALARM] 命中规则ID:%s, IP:%s, 规则内容:%s\n", ruleID, srcIP, ruleContent)
	// 可扩展为HTTP/邮件/钉钉/Prometheus/MCP/AI等
}

// toStrPtr、toStr 辅助函数建议移至 utils 包


