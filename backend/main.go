package main

import (
    "os"
    "time"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "coraza-waf/backend/logger"
    "context"
    "go.mongodb.org/mongo-driver/v2/mongo/readpref"
    "bytes"
    "github.com/prometheus/client_golang/prometheus"
    "coraza-waf/backend/services"
    "coraza-waf/backend/routes"
)

func main() {
    // 1. 检查 MongoDB 服务是否运行（例如通过 Docker 或本地 mongod）

    // 2. 从环境变量读取 MongoDB 配置（可选）
    mongoUri := os.Getenv("MONGO_URI")
    mongoDbName := os.Getenv("MONGO_DB_NAME")

    // 3. 设置默认值
    if mongoUri == "" {
        mongoUri = "mongodb://localhost:27017"
    }
    if mongoDbName == "" {
        mongoDbName = "waf_logs"
    }

    // 4. 初始化 Zap 日志
    zapLogger, err := zap.NewProduction()
    if err != nil {
        panic(err)
    }
    defer zapLogger.Sync()

    // 5. 初始化 MongoDB
    zapLogger.Info("开始初始化 MongoDB...", zap.String("uri", mongoUri), zap.String("db", mongoDbName))
    if err := logger.InitializeMongoDB(mongoUri, mongoDbName); err != nil {
        zapLogger.Fatal("MongoDB 初始化失败", zap.Error(err))
    }

    // 6. 验证 MongoDB 连接
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := logger.Client().Ping(ctx, readpref.Primary()); err != nil {
        zapLogger.Fatal("MongoDB 连接失败", zap.Error(err))
    }

    // 7. 初始化 Prometheus 指标（可选）
    requestCounter := prometheus.NewCounter(prometheus.CounterOpts{
        Name: "http_requests_total",
        Help: "Total number of HTTP requests made",
    })
    prometheus.MustRegister(requestCounter)

    // 8. 初始化 WAF 服务（确保传入 zapLogger）
    wafService, err := services.NewWAFService(zapLogger, `SecRuleEngine On`)
    if err != nil {
        zapLogger.Fatal("WAF 初始化失败", zap.Error(err))
    }

    // 9. 配置 Gin 引擎
    r := gin.Default()
    r.SetTrustedProxies([]string{"127.0.0.1", "192.168.1.0/24", "10.0.0.0/8"})

    // 10. 注册 API 路由
    routes.RegisterAPIRoutes(r, wafService)

    // 11. 定义日志中间件（处理响应体）
    r.Use(func(c *gin.Context) {
        var buf bytes.Buffer
        w := &logger.ResponseBodyWriter{
            ResponseWriter: c.Writer,
            Body:           &buf,
        }
        c.Writer = w

        // 记录请求日志
        logger.HandleRequest(c)

        // 记录开始时间
        start := time.Now()

        // 继续处理请求
        c.Next()

        // 记录响应耗时（毫秒）
        duration := time.Since(start).Milliseconds()

        // 记录响应日志
        logger.HandleResponse(c, buf.String(), duration)
    })

    // 12. 启动服务
    if err := r.Run("0.0.0.0:8080"); err != nil {
        zapLogger.Fatal("启动服务器失败", zap.Error(err))
    }
}
