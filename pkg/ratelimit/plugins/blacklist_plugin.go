package plugins

import (
	"context"
	"fmt"

	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/ratelimit"
)

func init() {
	// 注册黑名单插件工厂
	ratelimit.RegisterPluginFactory("blacklist", NewBlacklistPlugin)
}

// BlacklistPlugin 黑名单插件
type BlacklistPlugin struct {
	blacklist map[string]bool
	logger    logging.Logger
}

// NewBlacklistPlugin 创建黑名单插件
func NewBlacklistPlugin(params map[string]interface{}, logger logging.Logger) (ratelimit.RateLimiterPlugin, error) {
	// 从参数中获取黑名单
	blacklistParam, ok := params["blacklist"]
	if !ok {
		return nil, fmt.Errorf("缺少必要参数: blacklist")
	}

	// 转换参数类型
	blacklistSlice, ok := blacklistParam.([]interface{})
	if !ok {
		return nil, fmt.Errorf("blacklist 参数类型错误，应为字符串数组")
	}

	// 构建黑名单映射
	blacklist := make(map[string]bool)
	for _, item := range blacklistSlice {
		if key, ok := item.(string); ok {
			blacklist[key] = true
		}
	}

	return &BlacklistPlugin{
		blacklist: blacklist,
		logger:    logger,
	}, nil
}

// Name 插件名称
func (p *BlacklistPlugin) Name() string {
	return "blacklist_plugin"
}

// DoFilter 过滤方法
func (p *BlacklistPlugin) DoFilter(ctx context.Context, key string, tokens int64) (bool, error) {
	// 检查键是否在黑名单中
	if p.blacklist[key] {
		p.logger.Info("key in blacklist, enforcing rate limit", "key", key)
		// 返回 false 表示不绕过限流，但这里我们不返回错误
		return false, nil
	}
	// 不在黑名单中，允许继续处理
	return false, nil
}
