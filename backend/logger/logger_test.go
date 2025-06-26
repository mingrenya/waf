package logger

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
    "go.mongodb.org/mongo-driver/v2/mongo/readpref"
    "context"
    "bson"
)

func TestInsertLog(t *testing.T) {
    // 连接本地 MongoDB
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
    assert.NoError(t, err)
    defer client.Disconnect(context.TODO())

    db := client.Database("waf_logs")
    collection := db.Collection("logs")

    // 清空集合（可选）
    _, err = collection.DeleteMany(context.TODO(), bson.M{})
    assert.NoError(t, err)

    // 构造测试日志对象
    log := WafLog{
        RequestTime:       "2025-03-03T00:00:00Z",
        SourceIP:          "192.168.1.1",
        RequestMethod:     "GET",
        RequestURI:         "/ping",
        HTTPVersion:        "HTTP/1.1",
        RequestHost:        "localhost",
        UserAgent:          "Mozilla/5.0",
        Referer:            "http://localhost",
        RequestCookie:      "session=abc123",
        RequestID:          "req-1234",
        Blocked:            true,
        CorazaRuleID:       &"rule123",
        CorazaRuleMsg:      "SQL Injection",
        CorazaAction:       "deny",
        AttackType:         "SQL Injection",
        BusinessCategory:   "User Management",
        Scene:              "login",
        ClientCountry:      "CN",
        ClientRegion:       "Beijing",
        OriginIPAddress:    "192.168.1.1",
        ResponseTime:      12345,
        ResponseLength:    123,
        ResponseBody:      "pong",
        StatusCode:         200,
    }

    // 插入日志
    if err := insertLog(log); err != nil {
        t.Fatalf("日志插入失败: %v", err)
    }

    // 从数据库中读取日志
    var result WafLog
    err = collection.FindOne(context.TODO(), bson.M{"request_id": log.RequestID}).Decode(&result)
    assert.NoError(t, err)
    assert.Equal(t, "pong", result.ResponseBody, "响应体不匹配")
    assert.Equal(t, true, result.Blocked, "Blocked 标志不匹配")
    assert.Equal(t, "rule123", *result.CorazaRuleID, "Rules ID 不匹配")
}

