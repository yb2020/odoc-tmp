package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yb2020/odoc/internal/database"
	"github.com/yb2020/odoc/pkg/logging"
)

// RedisCache 实现基于Redis的缓存
type RedisCache struct {
	client     database.RedisClient
	logger     logging.Logger
	prefix     string
	expiration time.Duration
}

// NewRedisCache 创建一个新的Redis缓存
func NewRedisCache(client database.RedisClient, logger logging.Logger, expiration time.Duration, prefix string) *RedisCache {
	if expiration == 0 {
		expiration = 60 * time.Minute
	}

	return &RedisCache{
		client:     client,
		logger:     logger,
		prefix:     prefix,
		expiration: expiration,
	}
}

// formatKey 格式化缓存键
func (c *RedisCache) formatKey(key string) string {
	return fmt.Sprintf("%s:%s", c.prefix, key)
}

// generateKey 根据函数名和参数生成缓存键
func (c *RedisCache) generateKey(funcName string, args ...interface{}) string {
	key := funcName

	// 将参数添加到键中
	if len(args) > 0 {
		for _, arg := range args {
			// 对于基本类型，直接添加到键中
			key += fmt.Sprintf(":%v", arg)
		}
	}

	return c.formatKey(key)
}

// Get 从缓存中获取值
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	formattedKey := c.formatKey(key)

	// 从Redis获取值
	val, err := c.client.Get(ctx, formattedKey).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil // 键不存在，但不是错误
		}
		return false, err
	}

	// 反序列化JSON
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return false, err
	}

	return true, nil
}

func (c *RedisCache) GetNotBizPrefix(ctx context.Context, key string, dest interface{}) (bool, error) {
	// 从Redis获取值
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil // 键不存在，但不是错误
		}
		return false, err
	}
	// 反序列化JSON
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return false, err
	}

	return true, nil
}

// Set 将值存入缓存
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	formattedKey := c.formatKey(key)

	// 如果没有指定过期时间，使用默认值
	if expiration == 0 {
		expiration = c.expiration
	}

	// 序列化为JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// 存入Redis
	return c.client.Set(ctx, formattedKey, jsonData, expiration).Err()
}

func (c *RedisCache) SetNotBizPrefix(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// 如果没有指定过期时间，使用默认值
	if expiration == 0 {
		expiration = c.expiration
	}
	// 序列化为JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	// 存入Redis
	return c.client.Set(ctx, key, jsonData, expiration).Err()
}

// Delete 从缓存中删除值
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	formattedKey := c.formatKey(key)
	return c.client.Del(ctx, formattedKey).Err()
}

// DeletePattern 从缓存中删除匹配模式的所有键
func (c *RedisCache) DeletePattern(ctx context.Context, pattern string) error {
	formattedPattern := c.formatKey(pattern)

	// 使用KEYS命令查找匹配的键（在生产环境中应谨慎使用）
	keys, err := c.client.Do(ctx, "KEYS", formattedPattern).Result()
	if err != nil {
		return err
	}

	// 将结果转换为字符串切片
	stringKeys := make([]string, 0)
	if keyList, ok := keys.([]interface{}); ok {
		for _, key := range keyList {
			if strKey, ok := key.(string); ok {
				stringKeys = append(stringKeys, strKey)
			}
		}
	}

	// 如果没有匹配的键，直接返回
	if len(stringKeys) == 0 {
		return nil
	}

	// 删除所有匹配的键
	return c.client.Del(ctx, stringKeys...).Err()
}

// Wrap 包装一个函数，为其添加缓存功能
// 注意：这个方法仅适用于返回单个结果和错误的函数
func (c *RedisCache) Wrap(funcName string, fn interface{}) interface{} {
	// 获取函数的反射值
	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()

	// 检查函数类型
	if fnType.Kind() != reflect.Func {
		panic("Wrap: 参数必须是函数")
	}

	// 检查函数返回值
	if fnType.NumOut() != 2 {
		panic("Wrap: 函数必须有两个返回值")
	}

	// 创建一个新的函数
	wrappedFn := reflect.MakeFunc(fnType, func(args []reflect.Value) []reflect.Value {
		// 将参数转换为interface{}切片
		var interfaceArgs []interface{}
		for _, arg := range args {
			interfaceArgs = append(interfaceArgs, arg.Interface())
		}

		// 生成缓存键
		cacheKey := c.generateKey(funcName, interfaceArgs...)

		// 尝试从缓存获取结果
		resultType := fnType.Out(0)
		resultPtr := reflect.New(resultType)

		found, err := c.Get(context.Background(), cacheKey, resultPtr.Interface())
		if err != nil {
			c.logger.Error("msg", "从缓存获取结果失败", "error", err)
		} else if found {
			// 从缓存获取成功，返回缓存的结果
			return []reflect.Value{
				resultPtr.Elem(),
				reflect.Zero(fnType.Out(1)), // 返回nil错误
			}
		}

		// 调用原始函数
		results := fnValue.Call(args)

		// 检查错误
		if !results[1].IsNil() {
			// 如果有错误，直接返回原始结果
			return results
		}

		// 将结果存入缓存
		if err := c.Set(context.Background(), cacheKey, results[0].Interface(), c.expiration); err != nil {
			c.logger.Error("msg", "将结果存入缓存失败", "error", err)
		}

		// 返回原始结果
		return results
	})

	return wrappedFn.Interface()
}

// ClearCache 清除指定函数的缓存
func (c *RedisCache) ClearCache(ctx context.Context, funcName string, args ...interface{}) error {
	cacheKey := c.generateKey(funcName, args...)
	return c.Delete(ctx, cacheKey)
}

// ClearCachePattern 清除匹配模式的所有缓存
func (c *RedisCache) ClearCachePattern(ctx context.Context, pattern string) error {
	return c.DeletePattern(ctx, pattern)
}
