package ratelimit

import (
	"context"
	"errors"
	"fmt"
	"plugin"
	"sync"
	"time"

	"github.com/yb2020/odoc/internal/database"
	"github.com/yb2020/odoc/pkg/logging"
)

// 限流策略类型
type LimiterType string

const (
	// 计数器限流
	CounterLimiterType LimiterType = "counter"
	// 滑动窗口限流
	SlidingWindowLimiterType LimiterType = "sliding_window"
	// 令牌桶限流
	TokenBucketLimiterType LimiterType = "token_bucket"
	// 漏桶限流
	LeakyBucketLimiterType LimiterType = "leaky_bucket"
)

// 时间单位
type TimeUnit string

const (
	Second TimeUnit = "s"
	Minute TimeUnit = "m"
	Hour   TimeUnit = "h"
	Day    TimeUnit = "d"
)

// 限流维度
type LimiterDimension string

const (
	// 全局限流
	Global LimiterDimension = "global"
	// 用户级别限流
	User LimiterDimension = "user"
	// IP级别限流
	IP LimiterDimension = "ip"
	// 接口级别限流
	API LimiterDimension = "api"
)

// 限流器配置
type LimiterConfig struct {
	// 限流器类型
	Type LimiterType `json:"type"`
	// 限流键前缀
	KeyPrefix string `json:"key_prefix"`
	// 最大速率
	MaxRate int64 `json:"max_rate"`
	// 时间单位
	TimeUnit TimeUnit `json:"time_unit"`
	// 限流维度
	Dimension LimiterDimension `json:"dimension"`
	// 过期时间（秒）
	ExpireTime int64 `json:"expire_time"`
	// 插件配置
	PluginConfigs []PluginConfig `json:"plugin_configs"`
	// 动态插件路径
	DynamicPluginPaths []string `json:"dynamic_plugin_paths"`
}

// PluginConfig 插件配置
type PluginConfig struct {
	// 插件类型
	Type string `json:"type"`
	// 插件参数
	Params map[string]interface{} `json:"params"`
}

// 限流结果
type RateLimitResult struct {
	// 是否允许通过
	Allowed bool `json:"allowed"`
	// 剩余令牌数
	Remaining int64 `json:"remaining"`
	// 重试等待时间（毫秒）
	RetryAfter int64 `json:"retry_after"`
}

// RateLimiterPlugin 限流器插件接口
type RateLimiterPlugin interface {
	// Name 插件名称
	Name() string
	// DoFilter 过滤方法，返回 true 表示不需要限制，返回 false 表示需要限制
	DoFilter(ctx context.Context, key string, tokens int64) (bool, error)
}

// PluginFactory 插件工厂函数类型
type PluginFactory func(params map[string]interface{}, logger logging.Logger) (RateLimiterPlugin, error)

// PluginRegistry 插件注册表
var (
	pluginRegistry = make(map[string]PluginFactory)
	registryMutex  = &sync.RWMutex{}
)

// RegisterPluginFactory 注册插件工厂
func RegisterPluginFactory(pluginType string, factory PluginFactory) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	pluginRegistry[pluginType] = factory
}

// CreatePlugin 创建插件实例
func CreatePlugin(config PluginConfig, logger logging.Logger) (RateLimiterPlugin, error) {
	registryMutex.RLock()
	factory, exists := pluginRegistry[config.Type]
	registryMutex.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("未注册的插件类型: %s", config.Type)
	}
	
	return factory(config.Params, logger)
}

// LoadDynamicPlugin 加载动态插件
func LoadDynamicPlugin(pluginPath string, params map[string]interface{}, logger logging.Logger) (RateLimiterPlugin, error) {
	// 加载插件
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("加载插件失败: %w", err)
	}
	
	// 查找 NewPlugin 符号
	newPluginSym, err := p.Lookup("NewPlugin")
	if err != nil {
		return nil, fmt.Errorf("查找 NewPlugin 符号失败: %w", err)
	}
	
	// 类型断言
	newPlugin, ok := newPluginSym.(func(map[string]interface{}, logging.Logger) (RateLimiterPlugin, error))
	if !ok {
		return nil, fmt.Errorf("NewPlugin 符号类型错误")
	}
	
	// 创建插件实例
	return newPlugin(params, logger)
}

// RateLimiter 限流器接口
type RateLimiter interface {
	// Allow 判断请求是否允许通过
	// key: 限流键
	// tokens: 请求的令牌数量
	Allow(ctx context.Context, key string, tokens int64) (*RateLimitResult, error)
	
	// RegisterPlugin 注册插件
	RegisterPlugin(plugin RateLimiterPlugin)
	
	// RemovePlugin 移除插件
	RemovePlugin(pluginName string)
	
	// GetPlugins 获取所有插件
	GetPlugins() []RateLimiterPlugin
}

// BaseLimiter 基础限流器，实现插件功能
type BaseLimiter struct {
	plugins map[string]RateLimiterPlugin
	logger  logging.Logger
}

// NewBaseLimiter 创建基础限流器
func NewBaseLimiter(logger logging.Logger) *BaseLimiter {
	return &BaseLimiter{
		plugins: make(map[string]RateLimiterPlugin),
		logger:  logger,
	}
}

// RegisterPlugin 注册插件
func (b *BaseLimiter) RegisterPlugin(plugin RateLimiterPlugin) {
	b.plugins[plugin.Name()] = plugin
	b.logger.Info("plugin registered", "name", plugin.Name())
}

// RemovePlugin 移除插件
func (b *BaseLimiter) RemovePlugin(pluginName string) {
	if _, exists := b.plugins[pluginName]; exists {
		delete(b.plugins, pluginName)
		b.logger.Info("plugin removed", "name", pluginName)
	}
}

// GetPlugins 获取所有插件
func (b *BaseLimiter) GetPlugins() []RateLimiterPlugin {
	plugins := make([]RateLimiterPlugin, 0, len(b.plugins))
	for _, plugin := range b.plugins {
		plugins = append(plugins, plugin)
	}
	return plugins
}

// RunPlugins 运行所有插件
func (b *BaseLimiter) RunPlugins(ctx context.Context, key string, tokens int64) (bool, error) {
	for name, plugin := range b.plugins {
		bypass, err := plugin.DoFilter(ctx, key, tokens)
		if err != nil {
			b.logger.Error("plugin error", "name", name, "error", err)
			continue
		}
		
		// 如果插件返回 true，表示请求可以绕过限流
		if bypass {
			b.logger.Info("request bypassed by plugin", "name", name, "key", key)
			return true, nil
		}
	}
	return false, nil
}

// RateLimiterService 限流服务
type RateLimiterService struct {
	redis  database.RedisClient
	logger logging.Logger
}

// NewRateLimiterService 创建限流服务
func NewRateLimiterService(redis database.RedisClient, logger logging.Logger) *RateLimiterService {
	return &RateLimiterService{
		redis:  redis,
		logger: logger,
	}
}

// CreateLimiter 创建限流器
func (s *RateLimiterService) CreateLimiter(config LimiterConfig) (RateLimiter, error) {
	var limiter RateLimiter
	var err error
	
	// 根据类型创建限流器
	switch config.Type {
	case CounterLimiterType:
		limiter = NewCounterLimiter(s.redis, s.logger, config)
	case SlidingWindowLimiterType:
		limiter = NewSlidingWindowLimiter(s.redis, s.logger, config)
	case TokenBucketLimiterType:
		limiter = NewTokenBucketLimiter(s.redis, s.logger, config)
	case LeakyBucketLimiterType:
		limiter = NewLeakyBucketLimiter(s.redis, s.logger, config)
	default:
		return nil, errors.New("不支持的限流器类型")
	}
	
	// 创建并注册内置插件
	for _, pluginConfig := range config.PluginConfigs {
		plugin, err := CreatePlugin(pluginConfig, s.logger)
		if err != nil {
			s.logger.Error("创建插件失败", "type", pluginConfig.Type, "error", err)
			continue
		}
		limiter.RegisterPlugin(plugin)
	}
	
	// 加载并注册动态插件
	for i, pluginPath := range config.DynamicPluginPaths {
		// 为动态插件创建默认参数
		params := make(map[string]interface{})
		if i < len(config.PluginConfigs) {
			params = config.PluginConfigs[i].Params
		}
		
		plugin, err := LoadDynamicPlugin(pluginPath, params, s.logger)
		if err != nil {
			s.logger.Error("加载动态插件失败", "path", pluginPath, "error", err)
			continue
		}
		limiter.RegisterPlugin(plugin)
	}
	
	return limiter, err
}

// GetLimiterKey 根据维度获取限流键
func GetLimiterKey(prefix string, dimension LimiterDimension, id string) string {
	return fmt.Sprintf("%s:%s:%s", prefix, dimension, id)
}

// GetTimeUnitDuration 获取时间单位对应的持续时间
func GetTimeUnitDuration(unit TimeUnit) time.Duration {
	switch unit {
	case Second:
		return time.Second
	case Minute:
		return time.Minute
	case Hour:
		return time.Hour
	case Day:
		return time.Hour * 24
	default:
		return time.Second
	}
}
