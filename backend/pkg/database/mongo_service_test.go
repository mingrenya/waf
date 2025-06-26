package database

import (
    "testing"
)

func TestMongoService(t *testing.T) {
    // 替换为你的 MongoDB 地址、数据库名、集合名
    ms, err := NewMongoService("mongodb://localhost:27017", "wafdb", "waflogs")
    if err != nil {
        t.Fatalf("连接 MongoDB 失败: %v", err)
    }
    defer ms.Close()

    // 测试插入
    doc := map[string]interface{}{
        "test_field": "test_value",
    }
    if err := ms.InsertLog(doc); err != nil {
        t.Fatalf("插入文档失败: %v", err)
    }
}
