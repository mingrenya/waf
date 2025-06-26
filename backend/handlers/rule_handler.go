package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"coraza-waf/backend/internal/data"
	"coraza-waf/backend/models"
)

type Reloader interface {
	ReloadRules() error
}

type APIReloader struct {
	ReloadURL string
}

func (r *APIReloader) ReloadRules() error {
	resp, err := http.Post(r.ReloadURL, "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("reload failed: %s", resp.Status)
	}
	return nil
}

type RuleHandler struct {
	Repo     data.RuleRepository
	Reloader Reloader
}

func NewRuleHandler(repo data.RuleRepository, reloader Reloader) *RuleHandler {
	return &RuleHandler{Repo: repo, Reloader: reloader}
}

// 校验规则内容格式
func validateRuleFormat(format, content string) error {
	switch strings.ToLower(format) {
	case "modsec":
		if !strings.HasPrefix(content, "SecRule") {
			return fmt.Errorf("ModSecurity 规则必须以 SecRule 开头")
		}
	case "json":
		var js map[string]interface{}
		if err := json.Unmarshal([]byte(content), &js); err != nil {
			return fmt.Errorf("JSON 格式非法: %v", err)
		}
	case "expr":
		if len(content) == 0 {
			return fmt.Errorf("表达式内容不能为空")
		}
	default:
		return fmt.Errorf("不支持的规则格式: %s", format)
	}
	return nil
}

// 创建规则
func (h *RuleHandler) CreateRule(c *gin.Context) {
	var rule models.Rule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}
	if err := validateRuleFormat(rule.Format, rule.Content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "规则内容格式校验失败: " + err.Error()})
		return
	}
	if err := h.Repo.CreateRule(c.Request.Context(), &rule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 规则变更后热加载
	if h.Reloader != nil {
		h.Reloader.ReloadRules()
	}
	c.JSON(http.StatusOK, gin.H{"message": "规则创建成功", "id": rule.ID})
}

// 更新规则
func (h *RuleHandler) UpdateRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效ID"})
		return
	}
	var update bson.M
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}
	// 校验格式
	if format, ok := update["format"].(string); ok {
		content, _ := update["content"].(string)
		if err := validateRuleFormat(format, content); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "规则内容格式校验失败: " + err.Error()})
			return
		}
	}
	if err := h.Repo.UpdateRule(c.Request.Context(), id, update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if h.Reloader != nil {
		h.Reloader.ReloadRules()
	}
	c.JSON(http.StatusOK, gin.H{"message": "规则更新成功"})
}

// 删除规则
func (h *RuleHandler) DeleteRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效ID"})
		return
	}
	if err := h.Repo.DeleteRule(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if h.Reloader != nil {
		h.Reloader.ReloadRules()
	}
	c.JSON(http.StatusOK, gin.H{"message": "规则删除成功"})
}

// 查询单条规则
func (h *RuleHandler) GetRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效ID"})
		return
	}
	rule, err := h.Repo.GetRule(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rule)
}

// 分页查询规则
func (h *RuleHandler) ListRules(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	filter := bson.M{}
	if name := c.Query("name"); name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if ruleType := c.Query("rule_type"); ruleType != "" {
		filter["rule_type"] = ruleType
	}
	if enabled := c.Query("enabled"); enabled != "" {
		filter["enabled"] = enabled == "true"
	}
	if tag := c.Query("tag"); tag != "" {
		filter["tags"] = tag
	}
	rules, total, err := h.Repo.ListRules(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": total, "rules": rules})
}

// 启用/禁用规则
func (h *RuleHandler) EnableRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效ID"})
		return
	}
	enabled := c.DefaultQuery("enabled", "true") == "true"
	if err := h.Repo.EnableRule(c.Request.Context(), id, enabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if h.Reloader != nil {
		h.Reloader.ReloadRules()
	}
	c.JSON(http.StatusOK, gin.H{"message": "规则状态已更新"})
}
