package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	gocache "github.com/patrickmn/go-cache"
	"github.com/yb2020/odoc/pkg/logging"
)

// MemoryCache 基于 go-cache 的本地缓存实现
type MemoryCache struct {
	cache      *gocache.Cache
	logger     logging.Logger
	prefix     string
	expiration time.Duration
}

// NewMemoryCache 创建一个新的内存缓存
// expiration: 默认过期时间
// prefix: 缓存键前缀（业务隔离）
func NewMemoryCache(logger logging.Logger, expiration time.Duration, prefix string) *MemoryCache {
	if expiration == 0 {
		expiration = 60 * time.Minute
	}

	// 创建 go-cache 实例
	// 第一个参数：默认过期时间
	// 第二个参数：清理间隔（每 10 分钟清理一次过期项）
	c := gocache.New(expiration, 10*time.Minute)

	return &MemoryCache{
		cache:      c,
		logger:     logger,
		prefix:     prefix,
		expiration: expiration,
	}
}

// Close 关闭缓存（go-cache 会自动停止清理 goroutine）
func (c *MemoryCache) Close() {
	// go-cache 的清理 goroutine 会在 GC 时自动停止
	// 如果需要立即清空，可以调用 Flush
	c.cache.Flush()
}

// formatKey 格式化缓存键
func (c *MemoryCache) formatKey(key string) string {
	return fmt.Sprintf("%s:%s", c.prefix, key)
}

// generateKey 根据函数名和参数生成缓存键
func (c *MemoryCache) generateKey(funcName string, args ...interface{}) string {
	key := funcName
	if len(args) > 0 {
		for _, arg := range args {
			key += fmt.Sprintf(":%v", arg)
		}
	}
	return c.formatKey(key)
}

// Get 从缓存中获取值
func (c *MemoryCache) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	formattedKey := c.formatKey(key)
	return c.getInternal(formattedKey, dest)
}

// GetNotBizPrefix 从缓存中获取值，不带业务前缀
func (c *MemoryCache) GetNotBizPrefix(ctx context.Context, key string, dest interface{}) (bool, error) {
	return c.getInternal(key, dest)
}

// getInternal 内部获取方法
func (c *MemoryCache) getInternal(key string, dest interface{}) (bool, error) {
	data, found := c.cache.Get(key)
	if !found {
		return false, nil
	}

	// 数据以 []byte 形式存储，需要反序列化
	jsonData, ok := data.([]byte)
	if !ok {
		return false, fmt.Errorf("缓存数据类型错误")
	}

	if err := json.Unmarshal(jsonData, dest); err != nil {
		return false, err
	}

	return true, nil
}

// Set 将值存入缓存
func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	formattedKey := c.formatKey(key)
	return c.setInternal(formattedKey, value, expiration)
}

// SetNotBizPrefix 将值存入缓存，不带业务前缀
func (c *MemoryCache) SetNotBizPrefix(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.setInternal(key, value, expiration)
}

// setInternal 内部设置方法
func (c *MemoryCache) setInternal(key string, value interface{}, expiration time.Duration) error {
	if expiration == 0 {
		expiration = c.expiration
	}

	// 序列化为 JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	c.cache.Set(key, jsonData, expiration)
	return nil
}

// Delete 从缓存中删除值
func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	formattedKey := c.formatKey(key)
	c.cache.Delete(formattedKey)
	return nil
}

// DeletePattern 从缓存中删除匹配模式的所有键
func (c *MemoryCache) DeletePattern(ctx context.Context, pattern string) error {
	formattedPattern := c.formatKey(pattern)
	// 将通配符 * 转换为前缀匹配
	prefix := strings.TrimSuffix(formattedPattern, "*")

	// 遍历所有缓存项，删除匹配的
	items := c.cache.Items()
	for key := range items {
		if strings.HasPrefix(key, prefix) {
			c.cache.Delete(key)
		}
	}

	return nil
}

// Wrap 包装一个函数，为其添加缓存功能
func (c *MemoryCache) Wrap(funcName string, fn interface{}) interface{} {
	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()

	if fnType.Kind() != reflect.Func {
		panic("Wrap: 参数必须是函数")
	}

	if fnType.NumOut() != 2 {
		panic("Wrap: 函数必须有两个返回值")
	}

	wrappedFn := reflect.MakeFunc(fnType, func(args []reflect.Value) []reflect.Value {
		var interfaceArgs []interface{}
		for _, arg := range args {
			interfaceArgs = append(interfaceArgs, arg.Interface())
		}

		cacheKey := c.generateKey(funcName, interfaceArgs...)

		resultType := fnType.Out(0)
		resultPtr := reflect.New(resultType)

		found, err := c.Get(context.Background(), cacheKey, resultPtr.Interface())
		if err != nil {
			c.logger.Error("msg", "从缓存获取结果失败", "error", err)
		} else if found {
			return []reflect.Value{
				resultPtr.Elem(),
				reflect.Zero(fnType.Out(1)),
			}
		}

		results := fnValue.Call(args)

		if !results[1].IsNil() {
			return results
		}

		if err := c.Set(context.Background(), cacheKey, results[0].Interface(), c.expiration); err != nil {
			c.logger.Error("msg", "将结果存入缓存失败", "error", err)
		}

		return results
	})

	return wrappedFn.Interface()
}

// ClearCache 清除指定函数的缓存
func (c *MemoryCache) ClearCache(ctx context.Context, funcName string, args ...interface{}) error {
	cacheKey := c.generateKey(funcName, args...)
	return c.Delete(ctx, cacheKey)
}

// ClearCachePattern 清除匹配模式的所有缓存
func (c *MemoryCache) ClearCachePattern(ctx context.Context, pattern string) error {
	return c.DeletePattern(ctx, pattern)
}

// ItemCount 返回缓存中的项目数量（用于调试）
func (c *MemoryCache) ItemCount() int {
	return c.cache.ItemCount()
}

// Flush 清空所有缓存
func (c *MemoryCache) Flush() {
	c.cache.Flush()
}
