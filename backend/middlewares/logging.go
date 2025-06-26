package middlewares

import (
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "time"
)

func LoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        logger.Info("HTTP Request",
            zap.String("method", c.Request.Method),
            zap.String("path", c.Request.URL.Path),
            zap.Int("status", c.Writer.Status()),
            zap.String("ip", c.ClientIP()),
            zap.Duration("latency", time.Since(start)),
            zap.String("user_agent", c.Request.UserAgent()),
        )
    }
}

