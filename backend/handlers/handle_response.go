package handlers

import (
	"github.com/gin-gonic/gin"
	"coraza-waf/backend/logger"
	"log"
)

// HandleResponse 记录响应日志
func HandleResponse(c *gin.Context, responseBody string, responseTime int64) {
	headers := make(map[string][]string)
	for k, v := range c.Writer.Header() {
		headers[k] = v
	}
	// 从 context 获取规则命中信息（如有）
	ruleID, _ := c.Get("matched_rule_id")
	ruleContent, _ := c.Get("matched_rule_content")
	ruleFormat, _ := c.Get("matched_rule_format")

	logData := logger.WafLog{
		RequestTime:     "", // 可选，或传入 time.Now().Format(time.RFC3339)
		SourceIP:        c.ClientIP(),
		RequestMethod:   c.Request.Method,
		RequestURI:      c.Request.URL.Path,
		HTTPVersion:     c.Request.Proto,
		RequestHost:     c.Request.Host,
		UserAgent:       c.Request.UserAgent(),
		Referer:         c.Request.Header.Get("Referer"),
		RequestCookie:   c.Request.Header.Get("Cookie"),
		RequestID:       c.GetHeader("X-Request-ID"),
		ResponseBody:    responseBody,
		ResponseLength:  int64(len(responseBody)),
		ResponseHeaders: headers,
		StatusCode:      c.Writer.Status(),
		ResponseTime:    responseTime,
		Blocked:         false, // 可根据实际检测逻辑设置
		RuleID:          toStrPtr(ruleID),
		RuleContent:     toStr(ruleContent),
		RuleFormat:      toStr(ruleFormat),
	}

	if err := logger.InsertLog(logData); err != nil {
		log.Printf("写入响应日志失败: %v, data: %+v\n", err, logData)
	}
}


