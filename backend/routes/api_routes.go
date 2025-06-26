package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"coraza-waf/backend/services"
	"coraza-waf/backend/handlers"
)

func RegisterAPIRoutes(r *gin.Engine, wafService *services.WAFService) {
	// 健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// WAF 检测路由
	r.POST("/inspect", func(c *gin.Context) {
		allowed, _ := wafService.ProcessRequest(
			c.ClientIP(),
			c.Request.Host,
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.Proto,
			c.Request.Header,
		)
		
		if !allowed {
			c.AbortWithStatusJSON(403, gin.H{"error": "request blocked"})
			return
		}
		
		c.JSON(200, gin.H{"status": "allowed"})
	})

	// Prometheus Metrics 端点（新增）
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// 日志查询接口
	r.GET("/api/logs", handlers.HandleLogQuery)
}

