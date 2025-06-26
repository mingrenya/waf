package handlers

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"coraza-waf/backend/logger"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// 公共分页参数解析
func parsePageParams(c *gin.Context) (int, int) {
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
	return page, pageSize
}

// 公共时间区间过滤
func buildTimeFilter(c *gin.Context) bson.M {
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")
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
		return bson.M{"request_time": timeFilter}
	}
	return bson.M{}
}

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
	timeFilter := buildTimeFilter(c)
	for k, v := range timeFilter {
		filter[k] = v
	}
	page, pageSize := parsePageParams(c)
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

// 优化聚合接口，支持时间区间过滤
func aggWithTime(c *gin.Context, match bson.M, group bson.M) {
	timeFilter := buildTimeFilter(c)
	for k, v := range timeFilter {
		match[k] = v
	}
	pipeline := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$sort": bson.M{"count": -1}},
	}
	results, err := logger.AggregateLogs(c.Request.Context(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results})
}

func HandleLogAggByRuleID(c *gin.Context) {
	aggWithTime(c, bson.M{"rule_id": bson.M{"$ne": nil}}, bson.M{"_id": "$rule_id", "count": bson.M{"$sum": 1}})
}

func HandleLogAggByAttackType(c *gin.Context) {
	aggWithTime(c, bson.M{"attack_type": bson.M{"$ne": ""}}, bson.M{"_id": "$attack_type", "count": bson.M{"$sum": 1}})
}

func HandleLogAggBySourceIP(c *gin.Context) {
	aggWithTime(c, bson.M{"src_ip": bson.M{"$ne": ""}}, bson.M{"_id": "$src_ip", "count": bson.M{"$sum": 1}})
}

func HandleLogAggByDestIP(c *gin.Context) {
	aggWithTime(c, bson.M{"request_host": bson.M{"$ne": ""}}, bson.M{"_id": "$request_host", "count": bson.M{"$sum": 1}})
}

// 日志导出接口：按条件导出为CSV
func HandleLogExportCSV(c *gin.Context) {
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
	logs, _, err := logger.QueryLogs(context.Background(), filter, 1, 10000) // 最多导出1万条
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=logs.csv")
	w := csv.NewWriter(c.Writer)
	defer w.Flush()
	// 写表头
	headers := []string{"request_time","src_ip","request_method","request_uri","status_code","request_host","user_agent","referer","rule_id","attack_type"}
	w.Write(headers)
	// 写数据
	for _, log := range logs {
		row := []string{
			toStr(log["request_time"]),
			toStr(log["src_ip"]),
			toStr(log["request_method"]),
			toStr(log["request_uri"]),
			toStr(log["status_code"]),
			toStr(log["request_host"]),
			toStr(log["user_agent"]),
			toStr(log["referer"]),
			toStr(log["rule_id"]),
			toStr(log["attack_type"]),
		}
		w.Write(row)
	}
}

// 日志导出接口：按条件导出为JSON
func HandleLogExportJSON(c *gin.Context) {
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
	logs, _, err := logger.QueryLogs(context.Background(), filter, 1, 10000) // 最多导出1万条
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=logs.json")
	c.JSON(http.StatusOK, logs)
}

// toStr 辅助函数（适配 map[string]interface{}）
func toStr(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case int, int32, int64:
		return fmt.Sprintf("%v", val)
	case float64:
		return fmt.Sprintf("%.0f", val)
	default:
		return strings.Trim(fmt.Sprintf("%v", val), "{}[]")
	}
}

// 导出任务结构体
// 可存MongoDB/文件/内存，这里用本地JSON文件模拟任务表

type ExportTask struct {
	TaskID     string `json:"task_id"`
	Status     string `json:"status"` // pending, running, done, failed
	FilePath   string `json:"file_path"`
	CreatedAt  int64  `json:"created_at"`
	FinishedAt int64  `json:"finished_at"`
	Error      string `json:"error"`
}

// 创建导出任务（异步）
func HandleLogExportAsync(c *gin.Context) {
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
	taskID := genTaskID()
	task := ExportTask{
		TaskID:    taskID,
		Status:    "pending",
		FilePath:  "",
		CreatedAt: time.Now().Unix(),
	}
	taskFile := "/workspaces/waf/backend/export_tasks/" + taskID + ".json"
	// 立即写入任务文件，状态pending
	writeTaskFile(taskFile, &task)
	// 启动异步导出
	go runExportTask(taskFile, filter)
	c.JSON(http.StatusOK, gin.H{"task_id": taskID, "status": "pending"})
}

// 查询导出任务状态
func HandleLogExportTaskStatus(c *gin.Context) {
	taskID := c.Param("task_id")
	taskFile := "/workspaces/waf/backend/export_tasks/" + taskID + ".json"
	task, err := readTaskFile(taskFile)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// 下载导出文件
func HandleLogExportTaskDownload(c *gin.Context) {
	taskID := c.Param("task_id")
	taskFile := "/workspaces/waf/backend/export_tasks/" + taskID + ".json"
	task, err := readTaskFile(taskFile)
	if err != nil || task.Status != "done" {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件未就绪"})
		return
	}
	c.FileAttachment(task.FilePath, "logs_export_"+taskID+".json")
}

// 异步导出任务worker
func runExportTask(taskFile string, filter bson.M) {
	task, _ := readTaskFile(taskFile)
	task.Status = "running"
	writeTaskFile(taskFile, task)
	logs, _, err := logger.QueryLogs(context.Background(), filter, 1, 100000) // 支持10万条
	if err != nil {
		task.Status = "failed"
		task.Error = err.Error()
		task.FinishedAt = time.Now().Unix()
		writeTaskFile(taskFile, task)
		return
	}
	filePath := strings.Replace(taskFile, ".json", ".export.json", 1)
	f, ferr := os.Create(filePath)
	if ferr != nil {
		task.Status = "failed"
		task.Error = ferr.Error()
		task.FinishedAt = time.Now().Unix()
		writeTaskFile(taskFile, task)
		return
	}
	json.NewEncoder(f).Encode(logs)
	f.Close()
	task.Status = "done"
	task.FilePath = filePath
	task.FinishedAt = time.Now().Unix()
	writeTaskFile(taskFile, task)
}

// 工具函数：生成任务ID
func genTaskID() string {
	return fmt.Sprintf("tsk%x", rand.Int63())
}

func writeTaskFile(path string, task *ExportTask) {
	f, _ := os.Create(path)
	defer f.Close()
	json.NewEncoder(f).Encode(task)
}

func readTaskFile(path string) (*ExportTask, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	task := &ExportTask{}
	json.NewDecoder(f).Decode(task)
	return task, nil
}

// 日志详情接口：按ID查询单条日志及关联规则
func HandleLogDetail(c *gin.Context) {
	logID := c.Param("id")
	if logID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少日志ID"})
		return
	}
	// 查询日志
	logDoc, err := logger.FindLogByID(context.Background(), logID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "日志不存在"})
		return
	}
	// 查询关联规则（如有）
	var rule interface{} = nil
	if rid, ok := logDoc["rule_id"].(string); ok && rid != "" {
		rule, _ = logger.FindRuleByID(context.Background(), rid)
	}
	c.JSON(http.StatusOK, gin.H{"log": logDoc, "rule": rule})
}

// 日志全文检索接口：attack_type 精确匹配优先，text 关键字模糊检索
func HandleLogFullTextSearch(c *gin.Context) {
	attackType := c.Query("attack_type")
	text := c.Query("text")
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

	var filter bson.M
	if attackType != "" {
		// attack_type 精确匹配优先
		filter = bson.M{"attack_type": attackType}
	} else if text != "" {
		// text 关键字全文检索，支持多字段
		filter = bson.M{"$text": bson.M{"$search": text}}
	} else {
		filter = bson.M{}
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
