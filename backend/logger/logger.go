package logger

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"io"
	"log"
	"time"
)

const (
	mongoQueryTimeout   = 10 * time.Second
	mongoAggTimeout     = 15 * time.Second
	mongoInsertTimeout  = 5 * time.Second
)

// WafLog 定义
// 结构体字段分组注释，便于协作
// 基本请求信息
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
	RequestBody      string              `bson:"request_body,omitempty"`
	RequestHeaders   map[string][]string `bson:"request_headers,omitempty"`
	ResponseTime     int64               `bson:"response_time,omitempty"`
	ResponseLength   int64               `bson:"response_length,omitempty"`
	ResponseHeaders  map[string][]string `bson:"response_headers,omitempty"`
	ResponseBody     string              `bson:"response_body,omitempty"`
	StatusCode       int                 `bson:"status_code,omitempty"`
	RuleID           *string             `bson:"rule_id,omitempty" json:"rule_id,omitempty"`
	RuleContent      string              `bson:"rule_content,omitempty" json:"rule_content,omitempty"`
	RuleFormat       string              `bson:"rule_format,omitempty" json:"rule_format,omitempty"`
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
	ctx, cancel := context.WithTimeout(context.Background(), mongoInsertTimeout)
	defer cancel()
	_, err := collection.InsertOne(ctx, log)
	if err != nil {
		log.Printf("插入日志失败: %v\n", err)
		return err
	}
	return nil
}

// 请求日志处理
func HandleRequest(ctx *gin.Context) {
	headers := make(map[string][]string)
	for k, v := range ctx.Request.Header {
		headers[k] = v
	}
	body := ""
	if ctx.Request.Body != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(ctx.Request.Body)
		body = buf.String()
		ctx.Request.Body = io.NopCloser(bytes.NewBufferString(body))
	}
	log := WafLog{
		RequestTime:    time.Now().Format(time.RFC3339),
		SourceIP:       ctx.ClientIP(),
		RequestMethod:  ctx.Request.Method,
		RequestURI:     ctx.Request.RequestURI,
		HTTPVersion:    ctx.Request.Proto,
		RequestHost:    ctx.Request.Host,
		UserAgent:      ctx.Request.Header.Get("User-Agent"),
		Referer:        ctx.Request.Header.Get("Referer"),
		RequestCookie:  ctx.Request.Header.Get("Cookie"),
		RequestID:      ctx.GetHeader("X-Request-ID"),
		Blocked:        false,
		RequestBody:    body,
		RequestHeaders: headers,
	}
	if err := InsertLog(log); err != nil {
		log.Printf("请求日志写入失败: %v\n", err)
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
		log.Printf("响应日志写入失败: %v\n", err)
	}
}

// QueryLogs 支持全文检索和 attack_type 精确匹配，分页返回
func QueryLogs(ctx context.Context, filter bson.M, page, pageSize int) ([]bson.M, int64, error) {
	if collection == nil {
		return nil, 0, fmt.Errorf("MongoDB collection 未初始化")
	}
	ctx, cancel := context.WithTimeout(ctx, mongoQueryTimeout)
	defer cancel()
	findOpts := options.Find().SetSort(bson.M{"request_time": -1}).SetSkip(int64((page-1)*pageSize)).SetLimit(int64(pageSize))
	cursor, err := collection.Find(ctx, filter, findOpts)
	if err != nil {
		log.Printf("[QueryLogs] 查询失败: %v\n", err)
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	var results []bson.M
	for cursor.Next(ctx) {
		var m bson.M
		if err := cursor.Decode(&m); err == nil {
			results = append(results, m)
		}
	}
	total, _ := collection.CountDocuments(ctx, filter)
	return results, total, nil
}

// AggregateLogs 增加 context 超时和错误日志
func AggregateLogs(ctx context.Context, pipeline interface{}) ([]map[string]interface{}, error) {
	if collection == nil {
		return nil, fmt.Errorf("MongoDB collection 未初始化")
	}
	ctx, cancel := context.WithTimeout(ctx, mongoAggTimeout)
	defer cancel()
	cur, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("[AggregateLogs] 聚合失败: %v\n", err)
		return nil, err
	}
	defer cur.Close(ctx)
	var results []map[string]interface{}
	for cur.Next(ctx) {
		var m map[string]interface{}
		if err := cur.Decode(&m); err == nil {
			results = append(results, m)
		}
	}
	return results, nil
}

// 按ID查找单条日志
func FindLogByID(ctx context.Context, id string) (map[string]interface{}, error) {
	if collection == nil {
		return nil, fmt.Errorf("MongoDB collection 未初始化")
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("[FindLogByID] ID格式错误: %v\n", err)
		return nil, err
	}
	var result map[string]interface{}
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		log.Printf("[FindLogByID] 查询失败: %v\n", err)
	}
	return result, err
}

// 按ID查找规则（假设规则collection为rules）
func FindRuleByID(ctx context.Context, id string) (map[string]interface{}, error) {
	if db == nil {
		return nil, fmt.Errorf("MongoDB未初始化")
	}
	col := db.Collection("rules")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("[FindRuleByID] ID格式错误: %v\n", err)
		return nil, err
	}
	var result map[string]interface{}
	err = col.FindOne(ctx, bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		log.Printf("[FindRuleByID] 查询失败: %v\n", err)
	}
	return result, err
}

// ResponseBodyWriter 用于中间件捕获响应体，便于日志记录
// Write/WriteString 增加空指针保护
type ResponseBodyWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w *ResponseBodyWriter) Write(b []byte) (int, error) {
	if w.Body == nil {
		w.Body = &bytes.Buffer{}
	}
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *ResponseBodyWriter) WriteString(s string) (int, error) {
	if w.Body == nil {
		w.Body = &bytes.Buffer{}
	}
	if ws, ok := w.ResponseWriter.(interface{ WriteString(string) (int, error) }); ok {
		w.Body.WriteString(s)
		return ws.WriteString(s)
	}
	// fallback
	return w.Write([]byte(s))
}

