package cache

import (
	"context"
	"time"

	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/internal/database"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/memory"
	"github.com/yb2020/odoc/pkg/redis"
)

// 缓存类型常量
const (
	CacheTypeRedis  = "redis"
	CacheTypeMemory = "memory"
)

// Cache 定义缓存接口
type Cache interface {
	// Get 从缓存中获取值
	Get(ctx context.Context, key string, dest interface{}) (bool, error)
	// GetNotBizPrefix 从缓存中获取值，不带业务前缀
	GetNotBizPrefix(ctx context.Context, key string, dest interface{}) (bool, error)
	// Set 将值存入缓存
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	// SetNotBizPrefix 将值存入缓存，不带业务前缀
	SetNotBizPrefix(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Delete 从缓存中删除值
	Delete(ctx context.Context, key string) error

	// DeletePattern 从缓存中删除匹配模式的所有键
	DeletePattern(ctx context.Context, pattern string) error

	// Wrap 包装一个函数，为其添加缓存功能
	Wrap(funcName string, fn interface{}) interface{}

	// ClearCache 清除指定函数的缓存
	ClearCache(ctx context.Context, funcName string, args ...interface{}) error

	// ClearCachePattern 清除匹配模式的所有缓存
	ClearCachePattern(ctx context.Context, pattern string) error
}

// NewCache 根据配置创建缓存实例
// 通过 config.Cache.Type 自动选择缓存类型（redis 或 memory）
func NewCache(logger logging.Logger, expiration time.Duration, prefix string) Cache {
	cfg := config.GetConfig()
	if expiration == 0 {
		expiration = time.Duration(cfg.Cache.Expiration) * time.Minute
	}

	switch cfg.Cache.Type {
	case CacheTypeMemory:
		logger.Info("msg", "使用内存缓存", "prefix", prefix, "expiration", expiration)
		return memory.NewMemoryCache(logger, expiration, prefix)
	case CacheTypeRedis:
		fallthrough
	default:
		logger.Info("msg", "使用Redis缓存", "prefix", prefix, "expiration", expiration)
		return redis.NewRedisCache(database.GetRedisClient(), logger, expiration, prefix)
	}
}
