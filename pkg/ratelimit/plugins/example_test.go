package plugins

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/yb2020/odoc/internal/database"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/ratelimit"
)

// MockRedisClient 是一个模拟的 Redis 客户端
type MockRedisClient struct {
	mock.Mock
	database.RedisClient
}

// Do 模拟 Redis 命令执行
func (m *MockRedisClient) Do(ctx context.Context, args ...interface{}) *redis.Cmd {
	mockArgs := m.Called(ctx, args)
	return mockArgs.Get(0).(*redis.Cmd)
}

// 创建模拟的Redis命令结果
func createMockRedisCmd(val interface{}, err error) *redis.Cmd {
	cmd := redis.NewCmd(context.Background())
	cmd.SetVal(val)
	cmd.SetErr(err)
	return cmd
}

// 测试插件系统
func TestPluginSystem(t *testing.T) {
	// 创建模拟的 Redis 客户端和日志记录器
	mockRedis := &MockRedisClient{}
	logger := logging.NewLogger("debug", "console")

	// 创建限流服务
	service := ratelimit.NewRateLimiterService(mockRedis, logger)

	// 创建带有插件配置的限流器配置
	pluginConfigs := []ratelimit.PluginConfig{
		{
			Type: "whitelist",
			Params: map[string]interface{}{
				"whitelist": []interface{}{"user1", "user2"},
			},
		},
		{
			Type: "blacklist",
			Params: map[string]interface{}{
				"blacklist": []interface{}{"blocked_user"},
			},
		},
		{
			Type: "time_based",
			Params: map[string]interface{}{
				"time_ranges": []interface{}{
					map[string]interface{}{
						"start": "22:00",
						"end":   "06:00",
					},
				},
			},
		},
	}

	config := ratelimit.LimiterConfig{
		Type:          ratelimit.TokenBucketLimiterType,
		KeyPrefix:     "api",
		MaxRate:       100,
		TimeUnit:      ratelimit.Minute,
		Dimension:     ratelimit.User,
		ExpireTime:    3600,
		PluginConfigs: pluginConfigs,
	}

	// 创建限流器
	limiter, err := service.CreateLimiter(config)
	if err != nil {
		t.Fatalf("创建限流器失败: %v", err)
	}

	// 测试白名单用户
	// 白名单插件会绕过限流检查，所以不会调用 Redis
	result, err := limiter.Allow(context.Background(), "user1", 1)
	if err != nil {
		t.Errorf("白名单用户测试失败: %v", err)
	}
	if !result.Allowed {
		t.Errorf("白名单用户应该被允许通过")
	}

	// 测试普通用户
	normalUserCmd := createMockRedisCmd(mockRedisResult(200, 99), nil)
	mockRedis.On("Do", mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
		if len(args) < 4 || args[0] != "EVAL" {
			return false
		}
		// 检查键名是否包含 normal_user
		key, ok := args[3].(string)
		return ok && strings.Contains(key, "normal_user")
	})).Return(normalUserCmd).Once()
	
	result, err = limiter.Allow(context.Background(), "normal_user", 1)
	if err != nil {
		t.Errorf("普通用户测试失败: %v", err)
	}
	if !result.Allowed {
		t.Errorf("普通用户应该被允许通过")
	}

	// 测试黑名单用户
	blacklistCmd := createMockRedisCmd(mockRedisResult(401, 0), nil)
	mockRedis.On("Do", mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
		if len(args) < 4 || args[0] != "EVAL" {
			return false
		}
		// 检查键名是否包含 blocked_user
		key, ok := args[3].(string)
		return ok && strings.Contains(key, "blocked_user")
	})).Return(blacklistCmd).Once()
	
	result, err = limiter.Allow(context.Background(), "blocked_user", 1)
	if err != nil {
		t.Errorf("黑名单用户测试失败: %v", err)
	}
	if result.Allowed {
		t.Errorf("黑名单用户不应该被允许通过")
	}
	
	// 验证所有预期的调用都已发生
	mockRedis.AssertExpectations(t)
}

// 模拟 Redis 结果
func mockRedisResult(status, remaining int) interface{} {
	return []interface{}{
		fmt.Sprintf("%d", status),
		fmt.Sprintf("%d", remaining),
	}
}
