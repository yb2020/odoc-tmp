package plugins

import (
	"context"
	"fmt"

	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/ratelimit"
)

func init() {
	// 注册白名单插件工厂
	ratelimit.RegisterPluginFactory("whitelist", NewWhitelistPlugin)
}

// WhitelistPlugin 白名单插件
type WhitelistPlugin struct {
	whitelist map[string]bool
	logger    logging.Logger
}

// NewWhitelistPlugin 创建白名单插件
func NewWhitelistPlugin(params map[string]interface{}, logger logging.Logger) (ratelimit.RateLimiterPlugin, error) {
	// 从参数中获取白名单
	whitelistParam, ok := params["whitelist"]
	if !ok {
		return nil, fmt.Errorf("缺少必要参数: whitelist")
	}

	// 转换参数类型
	whitelistSlice, ok := whitelistParam.([]interface{})
	if !ok {
		return nil, fmt.Errorf("whitelist 参数类型错误，应为字符串数组")
	}

	// 构建白名单映射
	whitelist := make(map[string]bool)
	for _, item := range whitelistSlice {
		if key, ok := item.(string); ok {
			whitelist[key] = true
		}
	}

	return &WhitelistPlugin{
		whitelist: whitelist,
		logger:    logger,
	}, nil
}

// Name 插件名称
func (p *WhitelistPlugin) Name() string {
	return "whitelist_plugin"
}

// DoFilter 过滤方法
func (p *WhitelistPlugin) DoFilter(ctx context.Context, key string, tokens int64) (bool, error) {
	// 检查键是否在白名单中
	if p.whitelist[key] {
		p.logger.Info("key in whitelist, bypassing rate limit", "key", key)
		return true, nil
	}
	return false, nil
}
