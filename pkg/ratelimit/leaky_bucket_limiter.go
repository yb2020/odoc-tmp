package ratelimit

import (
	"context"
	"strconv"
	"time"

	"github.com/yb2020/odoc/internal/database"
	"github.com/yb2020/odoc/pkg/logging"
)

// LeakyBucketLimiter 漏桶限流器
type LeakyBucketLimiter struct {
	*BaseLimiter
	redis      database.RedisClient
	config     LimiterConfig
	timeWindow time.Duration
}

// NewLeakyBucketLimiter 创建漏桶限流器
func NewLeakyBucketLimiter(redis database.RedisClient, logger logging.Logger, config LimiterConfig) *LeakyBucketLimiter {
	return &LeakyBucketLimiter{
		BaseLimiter: NewBaseLimiter(logger),
		redis:       redis,
		config:      config,
		timeWindow:  GetTimeUnitDuration(config.TimeUnit),
	}
}

// Allow 判断请求是否允许通过
func (l *LeakyBucketLimiter) Allow(ctx context.Context, key string, tokens int64) (*RateLimitResult, error) {
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
	local last_time_key = KEYS[1] .. ":last"
	local queue_key = KEYS[1] .. ":queue"
	local rate = tonumber(ARGV[1])
	local capacity = tonumber(ARGV[2])
	local now = tonumber(ARGV[3])
	local requested = tonumber(ARGV[4])
	local expire = tonumber(ARGV[5])
	
	-- 获取上次漏水时间和当前队列大小
	local last_time = tonumber(redis.call("get", last_time_key) or now)
	local queue_size = tonumber(redis.call("get", queue_key) or "0")
	
	-- 计算从上次漏水到现在漏出的水量
	local elapsed = now - last_time
	local leaked = math.min(queue_size, math.floor(elapsed * rate / 1000))
	
	-- 计算新的队列大小
	local new_queue_size = queue_size - leaked
	
	-- 判断是否可以加入新的请求
	if new_queue_size + requested <= capacity then
		-- 更新队列大小和最后漏水时间
		redis.call("set", queue_key, new_queue_size + requested, "EX", expire)
		redis.call("set", last_time_key, now, "EX", expire)
		
		-- 计算剩余容量
		local remaining = capacity - (new_queue_size + requested)
		return {200, remaining}
	else
		-- 计算剩余容量
		local remaining = capacity - new_queue_size
		
		-- 更新队列大小和最后漏水时间
		redis.call("set", queue_key, new_queue_size, "EX", expire)
		redis.call("set", last_time_key, now, "EX", expire)
		
		return {401, remaining}
	end
	`

	// 计算每毫秒的漏水速率
	ratePerMs := float64(l.config.MaxRate) / float64(l.timeWindow/time.Millisecond)

	// 执行脚本
	result, err := l.redis.Do(ctx, "EVAL", script, 3, limiterKey, limiterKey+":last", limiterKey+":queue",
		ratePerMs, l.config.MaxRate, now, tokens, l.config.ExpireTime).Result()

	if err != nil {
		l.logger.Error("leaky bucket limiter error", "error", err)
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
		if tokensNeeded > 0 {
			retryAfter = int64(float64(tokensNeeded) / ratePerMs)
		}
	}

	return &RateLimitResult{
		Allowed:    statusCode == 200,
		Remaining:  remaining,
		RetryAfter: retryAfter,
	}, nil
}
