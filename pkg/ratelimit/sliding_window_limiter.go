package ratelimit

import (
	"context"
	"strconv"
	"time"

	"github.com/yb2020/odoc/internal/database"
	"github.com/yb2020/odoc/pkg/logging"
)

// SlidingWindowLimiter 滑动窗口限流器
type SlidingWindowLimiter struct {
	*BaseLimiter
	redis      database.RedisClient
	config     LimiterConfig
	timeWindow time.Duration
}

// NewSlidingWindowLimiter 创建滑动窗口限流器
func NewSlidingWindowLimiter(redis database.RedisClient, logger logging.Logger, config LimiterConfig) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		BaseLimiter: NewBaseLimiter(logger),
		redis:       redis,
		config:      config,
		timeWindow:  GetTimeUnitDuration(config.TimeUnit),
	}
}

// Allow 判断请求是否允许通过
func (l *SlidingWindowLimiter) Allow(ctx context.Context, key string, tokens int64) (*RateLimitResult, error) {
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
	local current_window_key = KEYS[1] .. ":current"
	local previous_window_key = KEYS[1] .. ":previous"
	local max_rate = tonumber(ARGV[1])
	local window_size_ms = tonumber(ARGV[2])
	local current_time = tonumber(ARGV[3])
	local tokens = tonumber(ARGV[4])
	local expire_time = tonumber(ARGV[5])
	
	-- 计算当前窗口开始时间
	local window_start = math.floor(current_time / window_size_ms) * window_size_ms
	local previous_window_start = window_start - window_size_ms
	
	-- 获取当前窗口和前一个窗口的请求数
	local current_count = tonumber(redis.call("get", current_window_key) or "0")
	local previous_count = tonumber(redis.call("get", previous_window_key) or "0")
	
	-- 计算当前时间在窗口中的位置比例
	local position = (current_time - window_start) / window_size_ms
	
	-- 计算滑动窗口中的请求数
	local sliding_window_count = current_count + previous_count * (1 - position)
	
	-- 判断是否超过限制
	if sliding_window_count + tokens <= max_rate then
		-- 更新当前窗口计数
		redis.call("incrby", current_window_key, tokens)
		-- 设置过期时间
		redis.call("expire", current_window_key, expire_time)
		redis.call("expire", previous_window_key, expire_time)
		
		-- 如果到了新窗口，将当前窗口设为前一个窗口
		if redis.call("exists", previous_window_key) == 0 then
			redis.call("set", previous_window_key, current_count, "EX", expire_time)
		end
		
		return {200, max_rate - (sliding_window_count + tokens)}
	else
		return {401, max_rate - sliding_window_count}
	end
	`
	
	// 执行脚本
	result, err := l.redis.Do(ctx, "EVAL", script, 2, limiterKey, limiterKey+":current", limiterKey+":previous", 
		l.config.MaxRate, l.timeWindow.Milliseconds(), now, tokens, l.config.ExpireTime).Result()
	
	if err != nil {
		l.logger.Error("sliding window limiter error", "error", err)
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
		retryAfter = int64(float64(l.timeWindow) / float64(l.config.MaxRate) * 1000)
	}
	
	return &RateLimitResult{
		Allowed:    statusCode == 200,
		Remaining:  remaining,
		RetryAfter: retryAfter,
	}, nil
}
