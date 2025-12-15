package middleware

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/yb2020/odoc/pkg/logging"
)

// LoggingMiddleware returns an endpoint middleware that logs the
// duration of each request and any error that may occur.
func LoggingMiddleware(logger logging.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				logger.Info(
					"method", getMethodName(ctx),
					"request", request,
					"error", err,
					"took", time.Since(begin),
				)
			}(time.Now())
			return next(ctx, request)
		}
	}
}

// getMethodName extracts the method name from the context
func getMethodName(ctx context.Context) string {
	if method, ok := ctx.Value(ContextKeyMethod).(string); ok {
		return method
	}
	return "unknown"
}
