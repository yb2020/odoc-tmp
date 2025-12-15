package ratelimit

import (
	"context"
	"strconv"
	"time"

	"github.com/yb2020/odoc/internal/database"
	"github.com/yb2020/odoc/pkg/logging"
)

// TokenBucketLimiter 令牌桶限流器
type TokenBucketLimiter struct {
	*BaseLimiter
	redis      database.RedisClient
	config     LimiterConfig
	timeWindow time.Duration
}

// NewTokenBucketLimiter 创建令牌桶限流器
func NewTokenBucketLimiter(redis database.RedisClient, logger logging.Logger, config LimiterConfig) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		BaseLimiter: NewBaseLimiter(logger),
		redis:       redis,
		config:      config,
		timeWindow:  GetTimeUnitDuration(config.TimeUnit),
	}
}

// Allow 判断请求是否允许通过
func (l *TokenBucketLimiter) Allow(ctx context.Context, key string, tokens int64) (*RateLimitResult, error) {
	// 运行插件
	bypass, err := l.RunPlugins(ctx, key, tokens)
	if err != nil {
		l.logger.Error("run plugins error", "error", err)
	}
	
	// 如果插件返回 true，表示请求可以绕过限流
	if bypass {
		return &RateLimitResult{
			Allowed:    true,
			Remaining:  l.config.MaxRate,
			RetryAfter: 0,
		}, nil
	}

	// 构建限流键
	limiterKey := GetLimiterKey(l.config.KeyPrefix, l.config.Dimension, key)
	
	// 获取当前时间戳（毫秒）
	now := time.Now().UnixNano() / int64(time.Millisecond)
	
	// 执行 Lua 脚本
	script := `
	local key = KEYS[1]
	local timestamp_key = KEYS[1] .. ":ts"
	local rate = tonumber(ARGV[1])
	local capacity = tonumber(ARGV[2])
	local now = tonumber(ARGV[3])
	local requested = tonumber(ARGV[4])
	local expire = tonumber(ARGV[5])
	
	local last_tokens = tonumber(redis.call("get", key) or capacity)
	local last_refreshed = tonumber(redis.call("get", timestamp_key) or now)
	
	local delta = math.max(0, now - last_refreshed)
	local filled_tokens = math.min(capacity, last_tokens + (delta * rate / 1000))
	
	local allowed = filled_tokens >= requested
	local new_tokens = filled_tokens
	
	if allowed then
		new_tokens = filled_tokens - requested
	end
	
	redis.call("setex", key, expire, new_tokens)
	redis.call("setex", timestamp_key, expire, now)
	
	return {allowed and 200 or 401, new_tokens}
	`
	
	// 计算每毫秒的令牌生成速率
	ratePerMs := float64(l.config.MaxRate) / float64(l.timeWindow/time.Millisecond)
	
	// 执行脚本
	result, err := l.redis.Do(ctx, "EVAL", script, 2, limiterKey, limiterKey+":ts", 
		ratePerMs, l.config.MaxRate, now, tokens, l.config.ExpireTime).Result()
	
	if err != nil {
		l.logger.Error("token bucket limiter error", "error", err)
		// 出错时默认放行
		return &RateLimitResult{
			Allowed:    true,
			Remaining:  l.config.MaxRate,
			RetryAfter: 0,
		}, nil
	}
	
	// 解析结果
	resultArray := result.([]interface{})
	statusCode, _ := strconv.ParseInt(resultArray[0].(string), 10, 64)
	remaining, _ := strconv.ParseInt(resultArray[1].(string), 10, 64)
	
	// 计算重试时间
	var retryAfter int64 = 0
	if statusCode == 401 {
		// 计算需要等待的时间（毫秒）
		tokensNeeded := tokens - remaining
		retryAfter = int64(float64(tokensNeeded) / ratePerMs)
	}
	
	return &RateLimitResult{
		Allowed:    statusCode == 200,
		Remaining:  remaining,
		RetryAfter: retryAfter,
	}, nil
}
