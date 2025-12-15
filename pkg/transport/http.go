package transport

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/metrics"
	"github.com/yb2020/odoc/pkg/middleware"
	"github.com/yb2020/odoc/pkg/tracing"
)

// GinHandlerFunc 是自定义的 Gin 处理函数类型
type GinHandlerFunc func(*gin.Context, opentracing.Tracer, logging.Logger)

// NewHTTPHandler creates a new HTTP handler for the service using Gin
func NewHTTPHandler(
	logger logging.Logger,
	tracer opentracing.Tracer,
	metrics *metrics.Metrics,
	customHandlers ...func(*gin.Engine, opentracing.Tracer, logging.Logger),
) http.Handler {
	// 设置为发布模式
	gin.SetMode(gin.ReleaseMode)

	// 创建 Gin 引擎
	r := gin.New()

	// 添加 CORS 中间件（必须在其他中间件之前）
	// 使用自定义中间件，完全放开 CORS 限制（开发环境）
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			// 显式列出所有允许的请求头（浏览器对 * 支持不一致）
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With, Access-Control-Request-Token, X-Custom-Handle-Error, Accept-Language, X-Traceid-Header, Cache-Control, Pragma")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type, X-Request-Id")
			c.Header("Access-Control-Max-Age", "43200")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// 创建错误报告器
	errorReporter := middleware.NewDefaultErrorReporter(logger, tracer)

	// 添加自定义 panic 处理中间件（替代 gin.Recovery()）
	r.Use(middleware.GinPanicHandlingMiddleware(errorReporter))

	// 添加跟踪中间件
	r.Use(tracing.NewGinTracingMiddleware(tracer))

	// 添加日志中间件
	r.Use(loggerMiddleware(logger))

	errorHandler := errors.NewErrorHandler(logger)

	// 添加业务错误处理中间件, 只能有一个
	r.Use(middleware.ErrorHandler(errorHandler))

	// 如果有指标中间件，添加它
	if metrics != nil {
		r.Use(metricsMiddleware(metrics))
	}

	// 添加用户代理中间件
	r.Use(middleware.UserAgentMiddleware())

	// 指标端点
	if metrics != nil {
		r.GET("/metrics", gin.WrapH(metrics.Handler))
	}

	// 健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 应用自定义处理器
	for _, handler := range customHandlers {
		handler(r, tracer, logger)
	}

	return r
}

// loggerMiddleware 创建一个 Gin 中间件，用于请求日志记录
func loggerMiddleware(logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()

		// 处理请求
		c.Next()

		// 请求处理完成后记录日志
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		raw := c.Request.URL.RawQuery

		// 构建完整路径
		if raw != "" {
			path = path + "?" + raw
		}

		logger.Info("request",
			"status", statusCode,
			"method", method,
			"path", path,
			"client_ip", clientIP,
			"latency", latency,
		)
	}
}

// metricsMiddleware 创建一个 Gin 中间件，用于指标收集
func metricsMiddleware(metrics *metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()

		// 处理请求
		c.Next()

		// 记录指标
		if metrics != nil {
			statusCode := c.Writer.Status()
			method := c.Request.Method
			endpoint := c.FullPath()
			latency := time.Since(start).Seconds()

			// 增加请求计数
			metrics.IncrementRequestCount(method, endpoint, http.StatusText(statusCode))

			// 观察请求延迟
			metrics.ObserveRequestLatency(method, endpoint, latency)
		}
	}
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	// 使用时间戳生成唯一请求 ID
	return "req-" + time.Now().Format("20060102150405.000000000")
}

// DefaultTestHandlers 提供默认的测试处理器
func DefaultTestHandlers(r *gin.Engine, tracer opentracing.Tracer, logger logging.Logger) {
	// 健康检查端点 - 返回服务状态
	r.GET("/health-check", func(c *gin.Context) {
		// 创建一个新的 span
		span := tracer.StartSpan("health-check-endpoint")
		defer span.Finish()

		// 将 span 添加到上下文
		ctx := c.Request.Context()
		ctx = opentracing.ContextWithSpan(ctx, span)
		c.Request = c.Request.WithContext(ctx)

		// 记录请求信息
		logger.Info("msg", "收到健康检查请求")
		span.LogKV("event", "health-check-request")

		// 返回响应
		c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"version": "1.0.0",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Hello World 端点 - 简单测试接口
	r.GET("/hello", func(c *gin.Context) {
		// 创建一个新的 span
		span := tracer.StartSpan("hello-endpoint")
		defer span.Finish()

		// 将 span 添加到上下文
		ctx := c.Request.Context()
		ctx = opentracing.ContextWithSpan(ctx, span)
		c.Request = c.Request.WithContext(ctx)

		// 记录请求信息
		logger.Info("msg", "收到 hello 请求")
		span.LogKV("event", "hello-request")

		// 返回响应
		c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!"})
	})

	// 日志测试端点 - 测试不同级别的日志
	r.GET("/log-test", func(c *gin.Context) {
		// 创建一个新的 span
		span := tracer.StartSpan("log-test-endpoint")
		defer span.Finish()

		// 将 span 添加到上下文
		ctx := c.Request.Context()
		ctx = opentracing.ContextWithSpan(ctx, span)
		c.Request = c.Request.WithContext(ctx)

		// 记录请求信息
		logger.Info("msg", "收到日志测试请求")
		span.LogKV("event", "log-test-request")

		// 记录不同级别的日志
		logger.Debug("msg", "这是一条调试日志", "level", "debug")
		logger.Info("msg", "这是一条信息日志", "level", "info")
		logger.Warn("msg", "这是一条警告日志", "level", "warn")
		logger.Error("msg", "这是一条错误日志", "level", "error")

		// 在 span 中记录日志
		span.LogKV(
			"event", "log-test",
			"message", "在 span 中记录了不同级别的日志",
		)

		// 返回响应
		c.JSON(http.StatusOK, map[string]string{
			"message": "日志测试完成",
			"status":  "success",
		})
	})

	// 跟踪测试端点 - 测试分布式跟踪
	r.GET("/trace-test", func(c *gin.Context) {
		// 解析查询参数
		depthStr := c.Query("depth")
		depth := 3 // 默认值
		if depthStr != "" {
			if d, err := strconv.Atoi(depthStr); err == nil && d > 0 {
				depth = d
			}
		}

		// 创建一个新的 span
		span := tracer.StartSpan("trace-test-endpoint")
		defer span.Finish()

		// 将 span 添加到上下文
		ctx := c.Request.Context()
		ctx = opentracing.ContextWithSpan(ctx, span)
		c.Request = c.Request.WithContext(ctx)

		// 记录请求信息
		logger.Info("msg", "收到跟踪测试请求", "depth", depth)
		span.LogKV("event", "trace-test-request", "depth", depth)

		// 递归创建嵌套 span
		result := createNestedSpans(ctx, tracer, logger, depth, 1)

		// 返回响应
		c.JSON(http.StatusOK, map[string]interface{}{
			"message": "跟踪测试完成",
			"depth":   depth,
			"result":  result,
		})
	})

	// Panic 测试端点 - 直接触发 panic
	r.GET("/panic-test", func(c *gin.Context) {
		// 创建一个新的 span
		span := tracer.StartSpan("panic-test-endpoint")
		defer span.Finish()

		// 将 span 添加到上下文
		ctx := c.Request.Context()
		ctx = opentracing.ContextWithSpan(ctx, span)
		c.Request = c.Request.WithContext(ctx)

		// 记录请求信息
		logger.Info("msg", "收到 panic-test 请求")
		span.LogKV("event", "panic-test-request")

		// 记录即将触发 panic 的信息
		span.LogKV(
			"event", "about-to-panic",
			"message", "即将触发一个模拟的 panic",
		)
		logger.Info("msg", "即将触发 panic")

		// 模拟一个除以零的错误
		a := 1
		b := 0
		// 这会触发一个 panic
		_ = a / b
		// 这行代码不会执行
		fmt.Println("This code will not be executed")
	})
}

// createNestedSpans 递归创建嵌套的 span
func createNestedSpans(ctx context.Context, tracer opentracing.Tracer, logger logging.Logger, maxDepth, currentDepth int) map[string]interface{} {
	// 从上下文中获取父 span
	parentSpan := opentracing.SpanFromContext(ctx)

	// 创建一个新的子 span
	childSpan := tracer.StartSpan(
		fmt.Sprintf("nested-span-depth-%d", currentDepth),
		opentracing.ChildOf(parentSpan.Context()),
	)
	defer childSpan.Finish()

	// 记录当前深度
	childSpan.SetTag("depth", currentDepth)
	childSpan.LogKV("event", "nested-span-created", "depth", currentDepth)

	// 模拟一些处理时间
	time.Sleep(time.Duration(10*currentDepth) * time.Millisecond)

	// 记录一些信息
	logger.Info("msg", fmt.Sprintf("创建了深度为 %d 的嵌套 span", currentDepth))

	// 构建结果
	result := map[string]interface{}{
		"depth":     currentDepth,
		"timestamp": time.Now().Format(time.RFC3339Nano),
	}

	// 如果还没有达到最大深度，继续递归
	if currentDepth < maxDepth {
		// 创建一个新的上下文，包含当前 span
		newCtx := opentracing.ContextWithSpan(ctx, childSpan)

		// 递归创建下一级 span
		childResult := createNestedSpans(newCtx, tracer, logger, maxDepth, currentDepth+1)

		// 将子结果添加到当前结果
		result["child"] = childResult
	}

	return result
}
