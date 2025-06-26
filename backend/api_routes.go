package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// WAF 规则管理
		api.GET("/rules", GetRules)
		api.POST("/rules", CreateRule)
		api.DELETE("/rules/:id", DeleteRule)
		
		// 日志查询
		api.GET("/logs", GetLogs)
		
		// AI 分析
		api.POST("/analyze", AnalyzeRequest)
	}
}

// 这里添加各个路由的处理函数...

