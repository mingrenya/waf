package routes

import (
	"github.com/gin-gonic/gin"
	"coraza-waf/backend/internal/gin/handlers"
	"coraza-waf/backend/internal/middleware"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/", handlers.IndexHandler)
	r.GET("/api/test", handlers.TestHandler)
	r.POST("/login", handlers.LoginHandler)
}

func Setup(router *gin.Engine) {
	// 健康检查
	router.GET("/health", handlers.HealthCheck)
	
	// API路由
	api := router.Group("/api")
	{
		api.POST("/analyze", handlers.AnalyzeHandler)
		api.GET("/logs", handlers.LogsHandler)
	}
	
	// WAF管理
	admin := router.Group("/admin", middleware.AuthMiddleware())
	{
		admin.GET("/rules", handlers.ListRulesHandler)
		admin.POST("/rules", handlers.AddRuleHandler)
		admin.DELETE("/rules/:id", handlers.DeleteRuleHandler)
	}
}
