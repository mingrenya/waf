package handlers

import (
	"github.com/gin-gonic/gin"
	"coraza-waf/backend/logger"
	"io"
	"bytes"
)

// HandleRequest 拦截请求并记录日志（使用结构体）
func HandleRequest(c *gin.Context) {
	// 保存请求体
	body, _ := c.GetRawData()
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// 构造 RequestLog 实例（新增：字段映射）
	var logData logger.RequestLog
	logData.RequestTime = time.Now().Format(time.RFC3339)
	logData.SourceIP = c.ClientIP()
	logData.RequestMethod = c.Request.Method
	logData.RequestURI = c.Request.URL.Path
	logData.UserAgent = c.Request.UserAgent()
	logData.RequestHeaders = c.Request.Header
	logData.RequestBody = string(body)

	// **新增：模拟 Coraza 的规则和拦截信息**
	if blocked {
		logData.Blocked = true
		logData.CorazaRuleID = ptr("942120")         // 示例 SQL 注入规则 ID
		logData.CorazaRuleMsg = "SQL Injection Attempt"
		logData.CorazaAction = "deny"
	}

	// 写入 MongoDB（使用结构体）
	if err := logger.LogRequest(&logData); err != nil {
		log.Println("写入请求日志失败:", err)
	}

	c.Next()
}

// ptr 帮助函数（用于处理 *string 类型）
func ptr(s string) *string {
	return &s
}

