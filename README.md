# WAF 后端服务

## 目录结构说明

```
backend/
├── main.go                # 程序入口
├── routes/                # 路由注册
├── handlers/              # 业务 handler（如日志查询）
├── logger/                # 日志与日志相关数据库操作
├── services/              # 业务服务层（如 WAF 逻辑）
├── middlewares/           # Gin 中间件
├── pkg/                   # 通用包（如数据库、工具等）
├── config/                # 配置相关
├── internal/              # 内部实现细节
└── ...                    # 其他目录
```

## 主要接口文档

### 1. 日志查询
- **接口**：`GET /api/logs`
- **参数（Query）**：
  - `src_ip`：源 IP（可选）
  - `request_method`：请求方法（可选）
  - `request_uri`：请求 URI（可选）
  - `status_code`：状态码（可选）
  - `request_host`：请求 Host（可选）
  - `user_agent`：User-Agent（可选）
  - `referer`：Referer（可选）
  - `http_version`：HTTP 版本（可选）
  - `start_time`：起始时间（RFC3339，可选）
  - `end_time`：结束时间（RFC3339，可选）
  - `page`：页码（默认1）
  - `page_size`：每页数量（默认20，最大100）
- **返回**：
```json
{
  "total": 100,
  "page": 1,
  "page_size": 20,
  "logs": [
    {
      "_id": "...",
      "request_time": "2025-06-26T12:54:24+08:00",
      "src_ip": "127.0.0.1",
      "request_method": "GET",
      "request_uri": "/inspect",
      "http_version": "HTTP/1.1",
      "request_host": "127.0.0.1:8080",
      "user_agent": "...",
      "status_code": 200
    }
  ]
}
```

### 2. 健康检查
- **接口**：`GET /health`
- **返回**：`{"status": "ok"}`

### 3. WAF 检测
- **接口**：`POST /inspect`
- **参数**：请求体内容由业务自定义
- **返回**：
  - 允许：`{"status": "allowed"}`
  - 拦截：`{"error": "request blocked"}`

### 4. Prometheus 监控
- **接口**：`GET /metrics`
- **返回**：Prometheus 指标数据

## 索引建议

建议在 MongoDB 日志集合建立如下索引提升查询性能：
```js
db.logs.createIndex({request_time: -1})
db.logs.createIndex({src_ip: 1})
db.logs.createIndex({request_method: 1})
db.logs.createIndex({status_code: 1})
```

## 扩展建议

- 如需用户认证、权限管理，可新建 `auth/` 或 `users/` 目录。
- 其他业务功能可按模块独立分层，保持结构清晰。
