package tracing

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// FromGinContext 从 gin.Context 中获取或创建 span
// 如果 gin.Context 中已经有 span，则创建一个子 span
// 如果没有，则创建一个新的 span
func FromGinContext(c *gin.Context, operationName string, tracer opentracing.Tracer) (context.Context, opentracing.Span) {
	// 从请求上下文中获取 span
	parentSpan := opentracing.SpanFromContext(c.Request.Context())
	
	var span opentracing.Span
	// 如果已经有 span，则创建一个子 span
	if parentSpan != nil {
		span = tracer.StartSpan(
			operationName,
			opentracing.ChildOf(parentSpan.Context()),
		)
	} else {
		// 如果没有 span，则创建一个新的 span
		span = tracer.StartSpan(operationName)
	}
	
	// 添加一些有用的标签
	span.SetTag("operation", operationName)
	
	// 获取请求 ID（如果有）
	if requestID := c.GetHeader("X-Request-ID"); requestID != "" {
		span.SetTag("request_id", requestID)
	}
	
	// 创建新的上下文
	ctx := opentracing.ContextWithSpan(c.Request.Context(), span)
	return ctx, span
}

// StartSpanFromContext 从已有上下文中创建子 span
func StartSpanFromContext(ctx context.Context, operationName string, tracer opentracing.Tracer) (context.Context, opentracing.Span) {
	var span opentracing.Span
	
	// 从上下文中获取父 span
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan != nil {
		// 创建子 span
		span = tracer.StartSpan(
			operationName,
			opentracing.ChildOf(parentSpan.Context()),
		)
	} else {
		// 创建新 span
		span = tracer.StartSpan(operationName)
	}
	
	// 添加标签
	span.SetTag("operation", operationName)
	
	// 创建新的上下文
	newCtx := opentracing.ContextWithSpan(ctx, span)
	return newCtx, span
}

// FinishSpan 完成 span 并记录错误（如果有）
func FinishSpan(span opentracing.Span, err error) {
	if err != nil {
		ext.Error.Set(span, true)
		span.SetTag("error.message", err.Error())
	}
	span.Finish()
}