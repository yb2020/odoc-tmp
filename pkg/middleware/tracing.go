package middleware

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/yb2020/odoc/pkg/logging"
)

// TracingMiddleware returns an endpoint middleware that starts a new span
// for each request and adds the span context to the request context.
func TracingMiddleware(tracer opentracing.Tracer, operationName string, logger logging.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			// Try to extract parent span context from the incoming context
			var parentSpanCtx opentracing.SpanContext
			if parent := opentracing.SpanFromContext(ctx); parent != nil {
				parentSpanCtx = parent.Context()
			}

			// Start a new span
			span := tracer.StartSpan(
				operationName,
				opentracing.ChildOf(parentSpanCtx),
				opentracing.Tag{Key: "component", Value: "endpoint"},
			)
			defer span.Finish()

			// Add method name to span
			if method := GetMethod(ctx); method != "" {
				span.SetTag("method", method)
			}

			// Add request ID to span
			if requestID := GetRequestID(ctx); requestID != "" {
				span.SetTag("request_id", requestID)
			}

			// Add trace ID to context
			traceID := ""
			if jaegerCtx, ok := span.Context().(jaeger.SpanContext); ok {
				traceID = jaegerCtx.TraceID().String()
			}
			ctx = WithTraceID(ctx, traceID)

			// Add span to context
			ctx = opentracing.ContextWithSpan(ctx, span)

			// Execute the next endpoint
			resp, err := next(ctx, request)

			// Record error if any
			if err != nil {
				ext.Error.Set(span, true)
				span.SetTag("error.message", err.Error())
				logger.Error(
					"msg", "Request failed",
					"operation", operationName,
					"error", err.Error(),
				)
			}

			return resp, err
		}
	}
}
