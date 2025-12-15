package tracing

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/middleware"
)

// NewGinTracingMiddleware 创建一个 Gin 中间件，用于处理分布式跟踪
func NewGinTracingMiddleware(tracer opentracing.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头中提取跟踪信息
		spanCtx, err := tracer.Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request.Header),
		)

		// 创建一个新的 span
		var span opentracing.Span
		if err != nil {
			// 请求头中没有跟踪上下文，创建一个新的
			span = tracer.StartSpan("HTTP " + c.Request.Method + " " + c.FullPath())
		} else {
			// 继续跟踪
			span = tracer.StartSpan(
				"HTTP "+c.Request.Method+" "+c.FullPath(),
				opentracing.ChildOf(spanCtx),
			)
		}
		defer span.Finish()

		// 设置 span 标签
		ext.HTTPMethod.Set(span, c.Request.Method)
		ext.HTTPUrl.Set(span, c.Request.URL.String())
		ext.Component.Set(span, "HTTP")

		// 生成请求 ID（如果不存在）
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
			c.Header("X-Request-ID", requestID)
		}

		// 添加请求 ID 到 span
		span.SetTag("request_id", requestID)

		// 将 span 添加到上下文
		ctx := c.Request.Context()
		ctx = middleware.WithRequestID(ctx, requestID)
		ctx = opentracing.ContextWithSpan(ctx, span)
		c.Request = c.Request.WithContext(ctx)

		// 处理请求
		c.Next()

		// 添加响应状态码
		ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
	}
}

// NewEndpointTracingMiddleware 创建一个 endpoint 中间件，用于处理分布式跟踪
func NewEndpointTracingMiddleware(tracer opentracing.Tracer, operationName string, logger logging.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			// 尝试从上下文中提取父 span 上下文
			var parentSpanCtx opentracing.SpanContext
			if parent := opentracing.SpanFromContext(ctx); parent != nil {
				parentSpanCtx = parent.Context()
			}

			// 创建一个新的 span
			span := tracer.StartSpan(
				operationName,
				opentracing.ChildOf(parentSpanCtx),
				opentracing.Tag{Key: "component", Value: "endpoint"},
			)
			defer span.Finish()

			// 添加方法名到 span
			if method := middleware.GetMethod(ctx); method != "" {
				span.SetTag("method", method)
			}

			// 添加请求 ID 到 span
			if requestID := middleware.GetRequestID(ctx); requestID != "" {
				span.SetTag("request_id", requestID)
			}

			// 添加跟踪 ID 到上下文
			traceID := ""
			if jaegerCtx, ok := span.Context().(jaeger.SpanContext); ok {
				traceID = jaegerCtx.TraceID().String()
			}
			ctx = middleware.WithTraceID(ctx, traceID)

			// 将 span 添加到上下文
			ctx = opentracing.ContextWithSpan(ctx, span)

			// 执行下一个 endpoint
			resp, err := next(ctx, request)

			// 记录错误（如果有）
			if err != nil {
				ext.Error.Set(span, true)
				span.SetTag("error.message", err.Error())
				if logger != nil {
					logger.Error(
						"msg", "Request failed",
						"operation", operationName,
						"error", err.Error(),
					)
				}
			}

			return resp, err
		}
	}
}

// generateRequestID 生成一个唯一的请求 ID
func generateRequestID() string {
	// 使用时间戳和随机数生成唯一 ID
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("req-%d-%d", time.Now().UnixNano(), rand.Intn(1000000))
}
