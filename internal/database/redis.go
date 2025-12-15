package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yb2020/odoc/pkg/logging"
)

// RedisConfig Redis 配置
type RedisConfig struct {
	Enabled         bool          `yaml:"enabled" json:"enabled"`
	Host            string        `yaml:"host" json:"host"`
	Port            int           `yaml:"port" json:"port"`
	Password        string        `yaml:"password" json:"password"`
	DB              int           `yaml:"db" json:"db"`
	PoolSize        int           `yaml:"poolSize" json:"poolSize"`
	MinIdleConns    int           `yaml:"minIdleConns" json:"minIdleConns"`
	DialTimeout     time.Duration `yaml:"dialTimeout" json:"dialTimeout"`
	ReadTimeout     time.Duration `yaml:"readTimeout" json:"readTimeout"`
	WriteTimeout    time.Duration `yaml:"writeTimeout" json:"writeTimeout"`
	MaxConnAge      time.Duration `yaml:"maxConnAge" json:"maxConnAge"`
	MaxRetries      int           `yaml:"maxRetries" json:"maxRetries"`
	MinRetryBackoff time.Duration `yaml:"minRetryBackoff" json:"minRetryBackoff"`
	MaxRetryBackoff time.Duration `yaml:"maxRetryBackoff" json:"maxRetryBackoff"`
}

// NewRedisClient 创建一个新的 Redis 客户端
func NewRedisClient(config RedisConfig, appLogger logging.Logger) (*redis.Client, error) {
	if !config.Enabled {
		appLogger.Info("msg", "Redis 客户端已禁用")
		return nil, nil
	}

	// 创建 Redis 客户端选项
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	}

	// 设置连接池参数
	if config.PoolSize > 0 {
		options.PoolSize = config.PoolSize
	}
	if config.MinIdleConns > 0 {
		options.MinIdleConns = config.MinIdleConns
	}
	if config.DialTimeout > 0 {
		options.DialTimeout = config.DialTimeout * time.Second
	}
	if config.ReadTimeout > 0 {
		options.ReadTimeout = config.ReadTimeout * time.Second
	}
	if config.WriteTimeout > 0 {
		options.WriteTimeout = config.WriteTimeout * time.Second
	}
	// 注意：在最新版本的go-redis中，MaxConnAge已更名为ConnMaxLifetime
	if config.MaxConnAge > 0 {
		options.ConnMaxLifetime = config.MaxConnAge * time.Second
	}
	if config.MaxRetries > 0 {
		options.MaxRetries = config.MaxRetries
	}
	if config.MinRetryBackoff > 0 {
		options.MinRetryBackoff = config.MinRetryBackoff * time.Millisecond
	}
	if config.MaxRetryBackoff > 0 {
		options.MaxRetryBackoff = config.MaxRetryBackoff * time.Millisecond
	}

	// 创建 Redis 客户端
	client := redis.NewClient(options)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("连接 Redis 服务器失败: %w", err)
	}

	appLogger.Info("msg", "Redis 客户端连接成功", "host", config.Host, "port", config.Port, "db", config.DB)

	return client, nil
}

// RedisClient Redis 客户端接口
type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	TTL(ctx context.Context, key string) *redis.DurationCmd
	Incr(ctx context.Context, key string) *redis.IntCmd
	HGet(ctx context.Context, key, field string) *redis.StringCmd
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd
	HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd
	HExists(ctx context.Context, key, field string) *redis.BoolCmd
	LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	LPop(ctx context.Context, key string) *redis.StringCmd
	RPop(ctx context.Context, key string) *redis.StringCmd
	LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd
	SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
	SMembers(ctx context.Context, key string) *redis.StringSliceCmd
	SRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
	SIsMember(ctx context.Context, key string, member interface{}) *redis.BoolCmd
	ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd
	ZRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd
	ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd
	ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
	ZScore(ctx context.Context, key, member string) *redis.FloatCmd
	Pipeline() redis.Pipeliner
	Ping(ctx context.Context) *redis.StatusCmd
	Do(ctx context.Context, args ...interface{}) *redis.Cmd
	SCard(ctx context.Context, key string) *redis.IntCmd
	Close() error
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd
}

// RedisClientWrapper 包装 redis.Client，实现 RedisClient 接口
type RedisClientWrapper struct {
	*redis.Client
}

// 确保 RedisClientWrapper 实现了 RedisClient 接口
var _ RedisClient = (*RedisClientWrapper)(nil)

// 全局Redis客户端实例
var globalRedisClient RedisClient

// InitRedisClient 初始化全局Redis客户端
func InitRedisClient(config RedisConfig, logger logging.Logger) error {
	client, err := NewRedisClient(config, logger)
	if err != nil {
		return err
	}

	if client != nil {
		// 测试连接
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Ping(ctx).Err(); err != nil {
			logger.Error("msg", "Redis连接测试失败", "error", err.Error())
			return err
		}

		logger.Info("msg", "Redis连接成功", "addr", fmt.Sprintf("%s:%d", config.Host, config.Port))
		// 将 redis.Client 包装为 RedisClientWrapper
		globalRedisClient = &RedisClientWrapper{Client: client}
	}

	return nil
}

// GetRedisClient 获取全局Redis客户端
func GetRedisClient() RedisClient {
	return globalRedisClient
}

func SCard(ctx context.Context, key string) *redis.IntCmd {
	return globalRedisClient.SCard(ctx, key)
}
