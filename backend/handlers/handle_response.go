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
	}

	// 可根据实际检测逻辑设置拦截信息
	// logData.Blocked = true
	// logData.CorazaRuleID = ptr("942120")
	// logData.CorazaRuleMsg = "SQL Injection Attempt"
	// logData.CorazaAction = "deny"

	if err := logger.InsertLog(logData); err != nil {
		log.Println("写入响应日志失败:", err)
	}
}


