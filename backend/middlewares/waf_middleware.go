package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func WAFMiddleware(wafService *services.WAFService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		allowed, interruption := wafService.ProcessRequest(c)
		if !allowed {
			logger.Warn("Request blocked by WAF",
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", interruption.Status),
				zap.String("action", interruption.Action),
				zap.String("data", interruption.Data),
			)
			
			c.AbortWithStatusJSON(interruption.Status, gin.H{
				"message": "Request blocked by WAF",
				"action":  interruption.Action,
				"data":    interruption.Data,
			})
			return
		}
		
		c.Next()
	}
}

