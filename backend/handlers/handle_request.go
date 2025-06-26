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

	logData := logger.WafLog{
		RequestTime:   time.Now().Format(time.RFC3339),
		SourceIP:      c.ClientIP(),
		RequestMethod: c.Request.Method,
		RequestURI:    c.Request.URL.Path,
		UserAgent:     c.Request.UserAgent(),
		RequestHost:   c.Request.Host,
		RequestCookie: c.Request.Header.Get("Cookie"),
		RequestID:     c.GetHeader("X-Request-ID"),
		Blocked:       false, // 你可根据实际检测逻辑设置
		RequestBody:   string(body),
		RequestHeaders: c.Request.Header,
	}

	// 可根据实际检测逻辑设置拦截信息
	// logData.Blocked = true
	// logData.CorazaRuleID = ptr("942120")
	// logData.CorazaRuleMsg = "SQL Injection Attempt"
	// logData.CorazaAction = "deny"

	if err := logger.InsertLog(logData); err != nil {
		log.Println("写入请求日志失败:", err)
	}

	c.Next()
}

// ptr 帮助函数（用于处理 *string 类型）
func ptr(s string) *string {
	return &s
}


