package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"coraza-waf/backend/services"
	"coraza-waf/backend/handlers"
	"coraza-waf/backend/internal/data"
	"coraza-waf/backend/logger"
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

	// 日志导出CSV接口（新增）
	r.GET("/api/logs/export/csv", handlers.HandleLogExportCSV)

	// 日志导出JSON接口（新增）
	r.GET("/api/logs/export/json", handlers.HandleLogExportJSON)

	// 日志聚合统计接口
	r.GET("/api/logs/agg/rule_id", handlers.HandleLogAggByRuleID)
	r.GET("/api/logs/agg/attack_type", handlers.HandleLogAggByAttackType)
	r.GET("/api/logs/agg/src_ip", handlers.HandleLogAggBySourceIP)
	r.GET("/api/logs/agg/dest_ip", handlers.HandleLogAggByDestIP)

	// 规则管理相关依赖注入
	// 获取 MongoDB collection
	client := logger.Client()
	db := client.Database("waf_logs") // 或用配置
	ruleCol := db.Collection("rules")
	ruleRepo := data.NewRuleRepository(ruleCol)

	// 热加载实现，ReloadURL 可通过配置/env/env变量传入
	reloader := &handlers.APIReloader{ReloadURL: "http://127.0.0.1:9090/api/reload"}
	ruleHandler := handlers.NewRuleHandler(ruleRepo, reloader)

	// 规则管理 RESTful API
	r.POST("/api/rules", ruleHandler.CreateRule)
	r.PUT("/api/rules/:id", ruleHandler.UpdateRule)
	r.DELETE("/api/rules/:id", ruleHandler.DeleteRule)
	r.GET("/api/rules/:id", ruleHandler.GetRule)
	r.GET("/api/rules", ruleHandler.ListRules)
	r.PATCH("/api/rules/:id/enable", ruleHandler.EnableRule)

	// 异步导出相关接口
	r.POST("/api/logs/export/async", handlers.HandleLogExportAsync)
	r.GET("/api/logs/export/task/:task_id", handlers.HandleLogExportTaskStatus)
	r.GET("/api/logs/export/download/:task_id", handlers.HandleLogExportTaskDownload)

	// 日志详情接口
	r.GET("/api/logs/:id", handlers.HandleLogDetail)

	// 日志全文检索接口
	r.GET("/api/logs/search", handlers.HandleLogFullTextSearch)
}


