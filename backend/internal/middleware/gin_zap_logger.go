package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinZapLogger 返回一个使用 zap 的 Gin 日志中间件
func GinZapLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 请求路径
		path := c.Request.URL.Path
		// 查询参数
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()
		// 延迟
		latency := end.Sub(start)
		// 客户端IP
		clientIP := c.ClientIP()
		// 请求方法
		method := c.Request.Method
		// 状态码
		statusCode := c.Writer.Status()
		// 用户代理
		userAgent := c.Request.UserAgent()

		// 记录日志
		logger.Info("HTTP request",
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", clientIP),
			zap.String("user-agent", userAgent),
			zap.Duration("latency", latency),
		)
	}
}
