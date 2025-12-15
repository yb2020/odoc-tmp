package ratelimit

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/yb2020/odoc/internal/database"
	"github.com/yb2020/odoc/pkg/logging"
)

// CounterLimiter 计数器限流器
type CounterLimiter struct {
	*BaseLimiter
	redis      database.RedisClient
	config     LimiterConfig
	timeWindow time.Duration
}

// NewCounterLimiter 创建计数器限流器
func NewCounterLimiter(redis database.RedisClient, logger logging.Logger, config LimiterConfig) *CounterLimiter {
	return &CounterLimiter{
		BaseLimiter: NewBaseLimiter(logger),
		redis:       redis,
		config:      config,
		timeWindow:  GetTimeUnitDuration(config.TimeUnit),
	}
}

// Allow 判断请求是否允许通过
func (l *CounterLimiter) Allow(ctx context.Context, key string, tokens int64) (*RateLimitResult, error) {
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

	// 使用传入的key作为限流键，不再重新构建
	// 如果key不包含前缀和维度，则构建完整的key
	var limiterKey string
	if strings.Contains(key, ":") {
		// 已经是完整的key，直接使用
		limiterKey = key
	} else {
		// 构建限流键
		limiterKey = GetLimiterKey(l.config.KeyPrefix, l.config.Dimension, key)
	}

	// 执行 Lua 脚本
	script := `
	local key = KEYS[1]
	local max_rate = tonumber(ARGV[1])
	local tokens = tonumber(ARGV[2])
	local window = tonumber(ARGV[3])
	local current = tonumber(redis.call('get', key) or "0")
	local new_current = current + tokens
	
	if new_current <= max_rate then
		redis.call('set', key, new_current, 'EX', window)
		return {200, max_rate - new_current}
	else
		return {401, max_rate - current}
	end
	`

	// 执行脚本
	result, err := l.redis.Do(ctx, "EVAL", script, 1, limiterKey, l.config.MaxRate, tokens, l.config.ExpireTime).Result()
	if err != nil {
		l.logger.Error("counter limiter error", "error", err)
		// 出错时默认放行
		return &RateLimitResult{
			Allowed:    true,
			Remaining:  l.config.MaxRate,
			RetryAfter: 0,
		}, nil
	}

	// 解析结果
	resultArray := result.([]interface{})

	// 安全地解析状态码和剩余次数
	var statusCode, remaining int64

	// 处理不同类型的返回值
	if sc, ok := resultArray[0].(int64); ok {
		statusCode = sc
	} else if scStr, ok := resultArray[0].(string); ok {
		statusCode, _ = strconv.ParseInt(scStr, 10, 64)
	} else {
		// 默认允许通过
		return &RateLimitResult{
			Allowed:    true,
			Remaining:  l.config.MaxRate,
			RetryAfter: 0,
		}, nil
	}

	if rem, ok := resultArray[1].(int64); ok {
		remaining = rem
	} else if remStr, ok := resultArray[1].(string); ok {
		remaining, _ = strconv.ParseInt(remStr, 10, 64)
	} else {
		remaining = 0
	}

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

// 计数器限流脚本
const counterLimiterScript = `
local seconds = "s"
local minute = "m"
local hour = "h"
local day = "d"

local success = 200
local badRequest = 400
local notEnoughToken = 401

--interval 自用户上次取到令牌至现在，经历的毫秒数
--maxRate 最大速率
--timeUnit 最大速率对应的时间单位
--计算在特定时间间隔内需要添加多少个令牌
local function getAddTokenCount(interval, maxRate, timeUnit)
    local intervalUnit = 1
    if seconds == timeUnit then
        intervalUnit = 1000 * intervalUnit
    elseif minute == timeUnit then
        intervalUnit = 60 * 1000 * intervalUnit
    elseif hour == timeUnit then
        intervalUnit = 60 * 60 * 1000 * intervalUnit
    elseif day == timeUnit then
        intervalUnit = 24 * 60 * 60 * 1000 * intervalUnit
    end
    local intervalForAddToken = intervalUnit / maxRate;
    return math.floor(interval / intervalForAddToken);
end

local function invalidExpireTimeAndTimeUnit(timeUnit, expireTimeBySeconds)
    if seconds == timeUnit then
        if expireTimeBySeconds < 1 then
            return true
        end
    elseif minute == timeUnit then
        if expireTimeBySeconds < 60 then
            return true
        end
    elseif hour == timeUnit then
        if expireTimeBySeconds < (60 * 60) then
            return true
        end
    elseif day == timeUnit then
        if expireTimeBySeconds < (24 * 60 * 60) then
            return true
        end
    end
    return false
end

--key 唯一标识
--currentTimeMillis 当前时间戳
--maxRate 限制速率
--timeUnit 限制速率的时间单位 支持 s(seconds) m(minute) h(hour) d(day)
--requestTokenNumber  请求的令牌数
--expireTimeBySeconds 过期时间
--返回table, 第1个值是状态码，第2个值是桶内当前的令牌数
local function tryAcquire(key, currentTimeMillis, maxRate, requestTokenNumber, timeUnit, expireTimeBySeconds)
    currentTimeMillis = tonumber(currentTimeMillis)
    maxRate = tonumber(maxRate)
    requestTokenNumber = tonumber(requestTokenNumber)
    expireTimeBySeconds = tonumber(expireTimeBySeconds)
    local rateLimitMetaData = redis.call("hmget", key, "lastAccess", "currentTokenCount")
    local lastAccess = tonumber(rateLimitMetaData[1])
    local currentTokenCount = tonumber(rateLimitMetaData[2])
    if (type(lastAccess) == "boolean" or lastAccess == nil) then
        lastAccess = 0
        currentTokenCount = 0
        redis.call("hmset", key, "lastAccess", lastAccess, "currentTokenCount", currentTokenCount)
        redis.call("expire", key, expireTimeBySeconds)
    end
    --参数校验
    if currentTimeMillis < lastAccess or invalidExpireTimeAndTimeUnit(timeUnit, expireTimeBySeconds) then
        redis.call("del", key)
        return { badRequest, currentTokenCount }
    end
    local interval = currentTimeMillis - lastAccess
    local addTokenNumberOfInterval = getAddTokenCount(interval, maxRate, timeUnit)
    if addTokenNumberOfInterval > 0 then
        currentTokenCount = currentTokenCount + addTokenNumberOfInterval
    end

    if currentTokenCount >= maxRate then
        currentTokenCount = maxRate
    end

    if currentTokenCount > 0 and currentTokenCount >= requestTokenNumber then
        currentTokenCount = currentTokenCount - requestTokenNumber
        lastAccess = currentTimeMillis
        redis.call("hmset", key, "lastAccess", lastAccess, "currentTokenCount", currentTokenCount)
        return { success, currentTokenCount }
    end

    return { notEnoughToken, currentTokenCount }
end

local key = KEYS[1]
local currentTimeMillis = ARGV[1]
local maxRate = ARGV[2]
local requestTokenNumber = ARGV[3]
local timeUnit = ARGV[4]
local expireTimeBySeconds = ARGV[5]

return tryAcquire(key, currentTimeMillis, maxRate, requestTokenNumber, timeUnit, expireTimeBySeconds)
`
