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
	// 解析查询参数
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")
	srcIP := c.Query("src_ip")
	requestMethod := c.Query("request_method")
	requestURI := c.Query("request_uri")
	statusCodeStr := c.Query("status_code")
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

	filter := bson.M{}
	if srcIP != "" {
		filter["src_ip"] = srcIP
	}
	if requestMethod != "" {
		filter["request_method"] = requestMethod
	}
	if requestURI != "" {
		filter["request_uri"] = requestURI
	}
	if statusCodeStr != "" {
		if code, err := strconv.Atoi(statusCodeStr); err == nil {
			filter["status_code"] = code
		}
	}
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
