package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"
	"coraza-waf/backend/logger"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// 日志查询接口
func HandleLogQuery(c *gin.Context) {
	params := map[string]string{
		"src_ip":         c.Query("src_ip"),
		"request_method": c.Query("request_method"),
		"request_uri":    c.Query("request_uri"),
		"status_code":    c.Query("status_code"),
		"request_host":   c.Query("request_host"),
		"user_agent":     c.Query("user_agent"),
		"referer":        c.Query("referer"),
		"http_version":   c.Query("http_version"),
	}
	filter := BuildLogFilter(params)

	// 时间区间单独处理
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")
	if startTimeStr != "" || endTimeStr != "" {
		timeFilter := bson.M{}
		if startTimeStr != "" {
			if t, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
				timeFilter["$gte"] = t
			}
		}
		if endTimeStr != "" {
			if t, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
				timeFilter["$lte"] = t
			}
		}
		if len(timeFilter) > 0 {
			filter["request_time"] = timeFilter
		}
	}

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	logs, total, err := logger.QueryLogs(context.Background(), filter, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"page": page,
		"page_size": pageSize,
		"logs": logs,
	})
}

// 动态参数映射辅助函数
func BuildLogFilter(params map[string]string) bson.M {
	filter := bson.M{}
	for k, v := range params {
		if v == "" {
			continue
		}
		switch k {
		case "status_code":
			if code, err := strconv.Atoi(v); err == nil {
				filter[k] = code
			}
		case "start_time", "end_time":
			// 跳过，时间区间单独处理
		default:
			filter[k] = v
		}
	}
	return filter
}
