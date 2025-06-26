package waf

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/corazawaf/coraza/v3"
	"github.com/corazawaf/coraza/v3/debuglog"
	"github.com/corazawaf/coraza/v3/experimental/plugins"
	"github.com/corazawaf/coraza/v3/types"
	"github.com/fsnotify/fsnotify"
)

const (
	defaultRulesDir = "./rules"       // 默认规则目录
	reloadCooldown  = 10 * time.Second // 重载冷却时间
)

// WAF管理器结构体
type WAFManager struct {
	waf     coraza.WAF            // Coraza核心实例
	mu      sync.RWMutex          // 读写锁
	logger  *slog.Logger          // 结构化日志
	watcher *fsnotify.Watcher     // 文件监听器
}

// 初始化选项
type Option func(*config)

type config struct {
	rulesDir      string                   // 规则目录路径
	debugLog      bool                     // 是否启用调试日志
	errorCallback func(types.MatchedRule)  // 自定义错误回调
}

// 初始化方法
func NewWAFManager(opts ...Option) (*WAFManager, error) {
	// 1. 初始化配置
	cfg := &config{
		rulesDir: defaultRulesDir,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	// 2. 创建WAF配置
	wafConfig := coraza.NewWAFConfig().
		WithErrorCallback(cfg.errorCallback)
	
	if cfg.debugLog {
		wafConfig = wafConfig.WithDebugLogger(&logAdapter{
			logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		})
	}

	// 3. 创建WAF实例
	instance, err := coraza.NewWAF(wafConfig)
	if err != nil {
		return nil, fmt.Errorf("WAF初始化失败: %w", err)
	}

	// 4. 创建文件监听器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("文件监听器创建失败: %w", err)
	}

	wm := &WAFManager{
		waf:     instance,
		logger:  slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		watcher: watcher,
	}

	// 5. 加载规则
	if err := wm.loadRules(cfg.rulesDir); err != nil {
		return nil, err
	}

	// 6. 启动规则监听协程
	go wm.watchRuleFiles(cfg.rulesDir)

	return wm, nil
}

// 核心方法实现
func (wm *WAFManager) loadRules(dir string) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	parser, err := crz_rules.NewParser(wm.waf)
	if err != nil {
		return fmt.Errorf("规则解析器创建失败: %w", err)
	}

	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".conf" {
			return nil
		}

		if err := parser.FromFile(path); err != nil {
			return fmt.Errorf("规则加载失败[%s]: %w", path, err)
		}

		wm.logger.Info("规则文件已加载", "path", path)
		return nil
	})
}

// 文件监听方法
func (wm *WAFManager) watchRuleFiles(dir string) {
	var lastReload time.Time
	for {
		select {
		case event, ok := <-wm.watcher.Events:
			if !ok || event.Op != fsnotify.Write {
				continue
			}
			
			if time.Since(lastReload) > reloadCooldown {
				wm.logger.Info("检测到规则变更", "file", event.Name)
				if err := wm.Reload(); err != nil {
					wm.logger.Error("规则热重载失败", "error", err)
				}
				lastReload = time.Now()
			}

		case err, ok := <-wm.watcher.Errors:
			if !ok {
				return
			}
			wm.logger.Error("文件监听错误", "error", err)
		}
	}
}

// API方法
func (wm *WAFManager) EvaluateRequest(ctx context.Context, req *Request) (bool, error) {
	tx := wm.waf.NewTransaction()
	defer tx.Close()

	if err := tx.ProcessRequest(req); err != nil {
		return false, fmt.Errorf("请求处理失败: %w", err)
	}

	return tx.Interrupted(), nil
}

func (wm *WAFManager) AddRule(rule string) error {
	// ...解析并添加单条规则...
}

func (wm *WAFManager) Close() error {
	return wm.watcher.Close()
}

