package logger

import (
    "context"
    "fmt"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
    "go.mongodb.org/mongo-driver/v2/mongo/readpref"
    "time"
    "bytes"
)

// WafLog 定义
type WafLog struct {
    RequestTime      string              `bson:"request_time,omitempty"`
    SourceIP         string              `bson:"src_ip,omitempty"`
    RequestMethod    string              `bson:"request_method,omitempty"`
    RequestURI       string              `bson:"request_uri,omitempty"`
    HTTPVersion      string              `bson:"http_version,omitempty"`
    RequestHost      string              `bson:"request_host,omitempty"`
    UserAgent        string              `bson:"user_agent,omitempty"`
    Referer          string              `bson:"referer,omitempty"`
    RequestCookie    string              `bson:"request_cookie,omitempty"`
    RequestID        string              `bson:"request_id,omitempty"`
    Blocked          bool                `bson:"is_blocked,omitempty"`
    CorazaRuleID     *string             `bson:"coraza_rule_id,omitempty"`
    CorazaRuleMsg    string              `bson:"coraza_rule_msg,omitempty"`
    CorazaAction     string              `bson:"coraza_action,omitempty"`
    AttackType       string              `bson:"attack_type,omitempty"`
    BusinessCategory string              `bson:"business_category,omitempty"`
    Scene            string              `bson:"scene,omitempty"`
    ClientCountry    string              `bson:"client_country,omitempty"`
    ClientRegion     string              `bson:"client_region,omitempty"`
    OriginIPAddress  string              `bson:"origin_ip,omitempty"`
    ResponseTime     int64               `bson:"response_time,omitempty"`
    ResponseLength   int64               `bson:"response_length,omitempty"`
    ResponseHeaders  map[string][]string `bson:"response_headers,omitempty"`
    ResponseBody     string              `bson:"response_body,omitempty"`
    StatusCode       int                 `bson:"status_code,omitempty"`
}

var client *mongo.Client
var db *mongo.Database
var collection *mongo.Collection

// MongoDB 初始化
func InitializeMongoDB(uri, dbName string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    opts := options.Client().ApplyURI(uri)
    c, err := mongo.Connect(opts) // v2 版本 Connect 只接收 opts
    if err != nil {
        return err
    }
    if err := c.Ping(ctx, readpref.Primary()); err != nil {
        return fmt.Errorf("MongoDB 连接失败: %w", err)
    }
    client = c
    db = client.Database(dbName)
    collection = db.Collection("logs")
    return nil
}

// 导出 client 以便 main.go 使用
func Client() *mongo.Client {
    return client
}

// 关闭 MongoDB 连接
func CloseMongoDB() {
    if client != nil {
        if err := client.Disconnect(context.TODO()); err != nil {
            fmt.Printf("关闭 MongoDB 连接失败: %v\n", err)
        }
    }
}

// 插入日志（导出，便于外部调用）
func InsertLog(log WafLog) error {
    if collection == nil {
        return fmt.Errorf("MongoDB collection 未初始化")
    }
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    _, err := collection.InsertOne(ctx, log)
    if err != nil {
        fmt.Printf("插入日志失败: %v\n", err)
        return err
    }
    return nil
}

// 请求日志处理
func HandleRequest(ctx *gin.Context) {
    log := WafLog{
        RequestTime:   time.Now().Format(time.RFC3339),
        SourceIP:      ctx.ClientIP(),
        RequestMethod: ctx.Request.Method,
        RequestURI:    ctx.Request.RequestURI,
        HTTPVersion:   ctx.Request.Proto,
        RequestHost:   ctx.Request.Host,
        UserAgent:     ctx.Request.Header.Get("User-Agent"),
        Referer:       ctx.Request.Header.Get("Referer"),
        RequestCookie: ctx.Request.Header.Get("Cookie"),
        RequestID:     ctx.GetHeader("X-Request-ID"),
        Blocked:       false,
    }
    if err := InsertLog(log); err != nil {
        fmt.Printf("请求日志写入失败: %v\n", err)
    }
}

// 响应日志处理
func HandleResponse(ctx *gin.Context, responseBody string, responseTime int64) {
    headers := make(map[string][]string)
    for k, v := range ctx.Writer.Header() {
        headers[k] = v
    }
    log := WafLog{
        RequestTime:     time.Now().Format(time.RFC3339),
        SourceIP:        ctx.ClientIP(),
        RequestMethod:   ctx.Request.Method,
        RequestURI:      ctx.Request.RequestURI,
        HTTPVersion:     ctx.Request.Proto,
        RequestHost:     ctx.Request.Host,
        UserAgent:       ctx.Request.Header.Get("User-Agent"),
        Referer:         ctx.Request.Header.Get("Referer"),
        RequestCookie:   ctx.Request.Header.Get("Cookie"),
        RequestID:       ctx.GetHeader("X-Request-ID"),
        ResponseBody:    responseBody,
        ResponseLength:  int64(len(responseBody)),
        ResponseHeaders: headers,
        StatusCode:      ctx.Writer.Status(),
        ResponseTime:    responseTime,
    }
    if err := InsertLog(log); err != nil {
        fmt.Printf("响应日志写入失败: %v\n", err)
    }
}

// 响应捕获器（导出类型名）
type ResponseBodyWriter struct {
    gin.ResponseWriter
    Body *bytes.Buffer
}

func (w *ResponseBodyWriter) Write(b []byte) (int, error) {
    w.Body.Write(b)
    return w.ResponseWriter.Write(b)
}

func (w *ResponseBodyWriter) WriteString(s string) (int, error) {
    if ws, ok := w.ResponseWriter.(interface{ WriteString(string) (int, error) }); ok {
        w.Body.WriteString(s)
        return ws.WriteString(s)
    }
    // fallback
    return w.Write([]byte(s))
}
