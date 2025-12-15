package ratelimit

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yb2020/odoc/internal/database"
	"github.com/yb2020/odoc/pkg/logging"
)

// 模拟 Redis 客户端
type MockRedisClient struct {
	mock.Mock
	database.RedisClient
}

func (m *MockRedisClient) Do(ctx context.Context, args ...interface{}) *redis.Cmd {
	mockArgs := m.Called(ctx, args)
	return mockArgs.Get(0).(*redis.Cmd)
}

// 模拟 Logger
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(keyvals ...interface{}) error {
	args := m.Called(keyvals)
	return args.Error(0)
}

func (m *MockLogger) Info(keyvals ...interface{}) error {
	args := m.Called(keyvals)
	return args.Error(0)
}

func (m *MockLogger) Warn(keyvals ...interface{}) error {
	args := m.Called(keyvals)
	return args.Error(0)
}

func (m *MockLogger) Error(keyvals ...interface{}) error {
	args := m.Called(keyvals)
	return args.Error(0)
}

func (m *MockLogger) With(keyvals ...interface{}) logging.Logger {
	args := m.Called(keyvals)
	return args.Get(0).(logging.Logger)
}

func TestGetTimeUnitDuration(t *testing.T) {
	tests := []struct {
		name     string
		unit     TimeUnit
		expected time.Duration
	}{
		{
			name:     "Second",
			unit:     Second,
			expected: time.Second,
		},
		{
			name:     "Minute",
			unit:     Minute,
			expected: time.Minute,
		},
		{
			name:     "Hour",
			unit:     Hour,
			expected: time.Hour,
		},
		{
			name:     "Day",
			unit:     Day,
			expected: time.Hour * 24,
		},
		{
			name:     "Invalid",
			unit:     TimeUnit("invalid"),
			expected: time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration := GetTimeUnitDuration(tt.unit)
			assert.Equal(t, tt.expected, duration)
		})
	}
}

func TestGetLimiterKey(t *testing.T) {
	tests := []struct {
		name      string
		prefix    string
		dimension LimiterDimension
		id        string
		expected  string
	}{
		{
			name:      "User dimension",
			prefix:    "app",
			dimension: User,
			id:        "123",
			expected:  "app:user:123",
		},
		{
			name:      "IP dimension",
			prefix:    "api",
			dimension: IP,
			id:        "192.168.1.1",
			expected:  "api:ip:192.168.1.1",
		},
		{
			name:      "API dimension",
			prefix:    "gateway",
			dimension: API,
			id:        "translate",
			expected:  "gateway:api:translate",
		},
		{
			name:      "Global dimension",
			prefix:    "global",
			dimension: Global,
			id:        "all",
			expected:  "global:global:all",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := GetLimiterKey(tt.prefix, tt.dimension, tt.id)
			assert.Equal(t, tt.expected, key)
		})
	}
}

func TestRateLimiterService_CreateLimiter(t *testing.T) {
	// 创建模拟对象
	mockRedis := new(MockRedisClient)
	mockLogger := new(MockLogger)

	// 设置 mockLogger 的预期行为
	mockLogger.On("With", mock.Anything).Return(mockLogger)
	mockLogger.On("Error", mock.Anything).Return(nil)
	mockLogger.On("Info", mock.Anything).Return(nil)

	// 创建限流服务
	service := NewRateLimiterService(mockRedis, mockLogger)

	tests := []struct {
		name        string
		config      LimiterConfig
		expectedErr bool
	}{
		{
			name: "Create CounterLimiter",
			config: LimiterConfig{
				Type:       CounterLimiterType,
				KeyPrefix:  "test",
				MaxRate:    100,
				TimeUnit:   Minute,
				Dimension:  User,
				ExpireTime: 3600,
			},
			expectedErr: false,
		},
		{
			name: "Create SlidingWindowLimiter",
			config: LimiterConfig{
				Type:       SlidingWindowLimiterType,
				KeyPrefix:  "test",
				MaxRate:    100,
				TimeUnit:   Minute,
				Dimension:  User,
				ExpireTime: 3600,
			},
			expectedErr: false,
		},
		{
			name: "Create TokenBucketLimiter",
			config: LimiterConfig{
				Type:       TokenBucketLimiterType,
				KeyPrefix:  "test",
				MaxRate:    100,
				TimeUnit:   Minute,
				Dimension:  User,
				ExpireTime: 3600,
			},
			expectedErr: false,
		},
		{
			name: "Create LeakyBucketLimiter",
			config: LimiterConfig{
				Type:       LeakyBucketLimiterType,
				KeyPrefix:  "test",
				MaxRate:    100,
				TimeUnit:   Minute,
				Dimension:  User,
				ExpireTime: 3600,
			},
			expectedErr: false,
		},
		{
			name: "Invalid Limiter Type",
			config: LimiterConfig{
				Type:       LimiterType("invalid"),
				KeyPrefix:  "test",
				MaxRate:    100,
				TimeUnit:   Minute,
				Dimension:  User,
				ExpireTime: 3600,
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiter, err := service.CreateLimiter(tt.config)
			if tt.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, limiter)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, limiter)
			}
		})
	}
}

// 创建模拟的Redis命令结果
func createMockRedisCmd(val interface{}, err error) *redis.Cmd {
	cmd := redis.NewCmd(context.Background())
	cmd.SetVal(val)
	cmd.SetErr(err)
	return cmd
}

func TestCounterLimiter_Allow(t *testing.T) {
	// 创建模拟对象
	mockRedis := new(MockRedisClient)
	mockLogger := new(MockLogger)

	// 配置限流器
	config := LimiterConfig{
		Type:       CounterLimiterType,
		KeyPrefix:  "test",
		MaxRate:    10,
		TimeUnit:   Second,
		Dimension:  User,
		ExpireTime: 60,
	}

	// 设置 mockLogger 的预期行为
	mockLogger.On("With", mock.Anything).Return(mockLogger)
	mockLogger.On("Error", mock.Anything).Return(nil)

	// 创建限流器
	limiter := NewCounterLimiter(mockRedis, mockLogger, config)

	// 测试用例
	tests := []struct {
		name           string
		key            string
		tokens         int64
		redisResult    interface{}
		redisError     error
		expectedResult *RateLimitResult
		expectError    bool
	}{
		{
			name:        "Allow request",
			key:         "user1",
			tokens:      1,
			redisResult: []interface{}{int64(200), int64(9)}, // 状态码200，剩余9个令牌
			redisError:  nil,
			expectedResult: &RateLimitResult{
				Allowed:    true,
				Remaining:  9,
				RetryAfter: 0,
			},
			expectError: false,
		},
		{
			name:        "Deny request",
			key:         "user2",
			tokens:      1,
			redisResult: []interface{}{int64(401), int64(0)}, // 状态码401，剩余0个令牌
			redisError:  nil,
			expectedResult: &RateLimitResult{
				Allowed:    false,
				Remaining:  0,
				RetryAfter: 100, // 1000ms/10 = 100ms
			},
			expectError: false,
		},
		{
			name:        "Redis error",
			key:         "user3",
			tokens:      1,
			redisResult: nil,
			redisError:  redis.ErrClosed,
			expectedResult: &RateLimitResult{
				Allowed:    true, // 默认放行
				Remaining:  10,
				RetryAfter: 0,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置模拟行为
			ctx := context.Background()
			mockRedis.On("Do", ctx, mock.Anything).Return(createMockRedisCmd(tt.redisResult, tt.redisError)).Once()

			// 调用被测函数
			result, err := limiter.Allow(ctx, tt.key, tt.tokens)

			// 验证结果
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.Allowed, result.Allowed)
				assert.Equal(t, tt.expectedResult.Remaining, result.Remaining)
				// 由于RetryAfter是动态计算的，只检查是否大于0
				if tt.expectedResult.RetryAfter > 0 {
					assert.Greater(t, result.RetryAfter, int64(0))
				} else {
					assert.Equal(t, tt.expectedResult.RetryAfter, result.RetryAfter)
				}
			}

			// 验证模拟对象的调用
			mockRedis.AssertExpectations(t)
		})
	}
}

func TestSlidingWindowLimiter_Allow(t *testing.T) {
	// 创建模拟对象
	mockRedis := new(MockRedisClient)
	mockLogger := new(MockLogger)

	// 配置限流器
	config := LimiterConfig{
		Type:       SlidingWindowLimiterType,
		KeyPrefix:  "test",
		MaxRate:    10,
		TimeUnit:   Second,
		Dimension:  User,
		ExpireTime: 60,
	}

	// 设置 mockLogger 的预期行为
	mockLogger.On("With", mock.Anything).Return(mockLogger)
	mockLogger.On("Error", mock.Anything).Return(nil)

	// 创建限流器
	limiter := NewSlidingWindowLimiter(mockRedis, mockLogger, config)

	// 测试用例
	tests := []struct {
		name           string
		key            string
		tokens         int64
		redisResult    interface{}
		redisError     error
		expectedResult *RateLimitResult
		expectError    bool
	}{
		{
			name:        "Allow request",
			key:         "user1",
			tokens:      1,
			redisResult: int64(9), // 剩余9个令牌
			redisError:  nil,
			expectedResult: &RateLimitResult{
				Allowed:    true,
				Remaining:  9,
				RetryAfter: 0,
			},
			expectError: false,
		},
		{
			name:        "Deny request",
			key:         "user2",
			tokens:      1,
			redisResult: int64(0), // 剩余0个令牌
			redisError:  nil,
			expectedResult: &RateLimitResult{
				Allowed:    false,
				Remaining:  0,
				RetryAfter: 100, // 1000ms/10 = 100ms
			},
			expectError: false,
		},
		{
			name:        "Redis error",
			key:         "user3",
			tokens:      1,
			redisResult: nil,
			redisError:  redis.ErrClosed,
			expectedResult: &RateLimitResult{
				Allowed:    true, // 默认放行
				Remaining:  10,
				RetryAfter: 0,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置模拟行为
			ctx := context.Background()
			mockRedis.On("Do", ctx, mock.Anything).Return(createMockRedisCmd(tt.redisResult, tt.redisError)).Once()

			// 调用被测函数
			result, err := limiter.Allow(ctx, tt.key, tt.tokens)

			// 验证结果
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.Allowed, result.Allowed)
				assert.Equal(t, tt.expectedResult.Remaining, result.Remaining)
				// 由于RetryAfter是动态计算的，只检查是否大于0
				if tt.expectedResult.RetryAfter > 0 {
					assert.Greater(t, result.RetryAfter, int64(0))
				} else {
					assert.Equal(t, tt.expectedResult.RetryAfter, result.RetryAfter)
				}
			}

			// 验证模拟对象的调用
			mockRedis.AssertExpectations(t)
		})
	}
}

func TestTokenBucketLimiter_Allow(t *testing.T) {
	// 创建模拟对象
	mockRedis := new(MockRedisClient)
	mockLogger := new(MockLogger)

	// 配置限流器
	config := LimiterConfig{
		Type:       TokenBucketLimiterType,
		KeyPrefix:  "test",
		MaxRate:    10,
		TimeUnit:   Second,
		Dimension:  User,
		ExpireTime: 60,
	}

	// 设置 mockLogger 的预期行为
	mockLogger.On("With", mock.Anything).Return(mockLogger)
	mockLogger.On("Error", mock.Anything).Return(nil)

	// 创建限流器
	limiter := NewTokenBucketLimiter(mockRedis, mockLogger, config)

	// 测试用例
	tests := []struct {
		name           string
		key            string
		tokens         int64
		redisResult    interface{}
		redisError     error
		expectedResult *RateLimitResult
		expectError    bool
	}{
		{
			name:        "Allow request",
			key:         "user1",
			tokens:      1,
			redisResult: []interface{}{int64(1), int64(9)}, // 允许通过，剩余9个令牌
			redisError:  nil,
			expectedResult: &RateLimitResult{
				Allowed:    true,
				Remaining:  9,
				RetryAfter: 0,
			},
			expectError: false,
		},
		{
			name:        "Deny request",
			key:         "user2",
			tokens:      1,
			redisResult: []interface{}{int64(0), int64(0)}, // 不允许通过，剩余0个令牌
			redisError:  nil,
			expectedResult: &RateLimitResult{
				Allowed:    false,
				Remaining:  0,
				RetryAfter: 100, // 需要等待100ms
			},
			expectError: false,
		},
		{
			name:        "Redis error",
			key:         "user3",
			tokens:      1,
			redisResult: nil,
			redisError:  redis.ErrClosed,
			expectedResult: &RateLimitResult{
				Allowed:    true, // 默认放行
				Remaining:  10,
				RetryAfter: 0,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置模拟行为
			ctx := context.Background()
			mockRedis.On("Do", ctx, mock.Anything).Return(createMockRedisCmd(tt.redisResult, tt.redisError)).Once()
			
			// 调用被测函数
			result, err := limiter.Allow(ctx, tt.key, tt.tokens)

			// 验证结果
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.Allowed, result.Allowed)
				assert.Equal(t, tt.expectedResult.Remaining, result.Remaining)
				// 由于RetryAfter是动态计算的，只检查是否大于0
				if tt.expectedResult.RetryAfter > 0 {
					assert.Greater(t, result.RetryAfter, int64(0))
				} else {
					assert.Equal(t, tt.expectedResult.RetryAfter, result.RetryAfter)
				}
			}

			// 验证模拟对象的调用
			mockRedis.AssertExpectations(t)
		})
	}
}

func TestLeakyBucketLimiter_Allow(t *testing.T) {
	// 创建模拟对象
	mockRedis := new(MockRedisClient)
	mockLogger := new(MockLogger)

	// 配置限流器
	config := LimiterConfig{
		Type:       LeakyBucketLimiterType,
		KeyPrefix:  "test",
		MaxRate:    10,
		TimeUnit:   Second,
		Dimension:  User,
		ExpireTime: 60,
	}

	// 设置 mockLogger 的预期行为
	mockLogger.On("With", mock.Anything).Return(mockLogger)
	mockLogger.On("Error", mock.Anything).Return(nil)

	// 创建限流器
	limiter := NewLeakyBucketLimiter(mockRedis, mockLogger, config)

	// 测试用例
	tests := []struct {
		name           string
		key            string
		tokens         int64
		redisResult    interface{}
		redisError     error
		expectedResult *RateLimitResult
		expectError    bool
	}{
		{
			name:        "Allow request",
			key:         "user1",
			tokens:      1,
			redisResult: []interface{}{int64(1), int64(0)}, // 允许通过，等待时间0
			redisError:  nil,
			expectedResult: &RateLimitResult{
				Allowed:    true,
				Remaining:  10, // 桶容量 - 等待时间
				RetryAfter: 0,
			},
			expectError: false,
		},
		{
			name:        "Deny request",
			key:         "user2",
			tokens:      1,
			redisResult: []interface{}{int64(0), int64(100)}, // 不允许通过，需要等待100ms
			redisError:  nil,
			expectedResult: &RateLimitResult{
				Allowed:    false,
				Remaining:  0,
				RetryAfter: 100,
			},
			expectError: false,
		},
		{
			name:        "Redis error",
			key:         "user3",
			tokens:      1,
			redisResult: nil,
			redisError:  redis.ErrClosed,
			expectedResult: &RateLimitResult{
				Allowed:    true, // 默认放行
				Remaining:  10,
				RetryAfter: 0,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置模拟行为
			ctx := context.Background()
			mockRedis.On("Do", ctx, mock.Anything).Return(createMockRedisCmd(tt.redisResult, tt.redisError)).Once()
			
			// 调用被测函数
			result, err := limiter.Allow(ctx, tt.key, tt.tokens)

			// 验证结果
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.Allowed, result.Allowed)
				// 由于剩余空间是动态计算的，只检查是否大于等于0
				assert.GreaterOrEqual(t, result.Remaining, int64(0))
				if tt.expectedResult.RetryAfter > 0 {
					assert.Equal(t, tt.expectedResult.RetryAfter, result.RetryAfter)
				} else {
					assert.Equal(t, tt.expectedResult.RetryAfter, result.RetryAfter)
				}
			}

			// 验证模拟对象的调用
			mockRedis.AssertExpectations(t)
		})
	}
}
