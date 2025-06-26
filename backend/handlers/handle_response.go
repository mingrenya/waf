package handlers

import (
	"github.com/gin-gonic/gin"
	"coraza-waf/backend/logger"
	"time"
)

// HandleResponse 记录响应日志（使用结构体）
func HandleResponse(c *gin.Context) {
	// **新增：计算请求耗时**
	responseTime := time.Since(start).Seconds() // 需要在请求中定义 start 变量

	// 构造 ResponseLog 实例（新增：字段映射）
	var logData logger.ResponseLog
	logData.StatusCode = c.Writer.Status()
	logData.ResponseTime = int64(responseTime * 1000) // 转换为毫秒
	logData.ResponseHeaders = c.Writer.Header()
	logData.ResponseBody = string(responseBody) // 假设你已经捕获了响应内容

	// **新增：同步 Coraza 的拦截状态**
	if blocked {
		logData.Blocked = true
		logData.CorazaRuleID = ptr("942120")
		logData.CorazaRuleMsg = "SQL Injection Attempt"
		logData.CorazaAction = "deny"
	}

	// 写入 MongoDB
	if err := logger.LogResponse(&logData); err != nil {
		log.Println("写入响应日志失败:", err)
	}

	c.Next()
}

