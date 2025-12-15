package middleware

import (
	"context"
)

// ContextKey is a type for context keys
type ContextKey string

const (
	// ContextKeyRequestID is the key for the request ID in the context
	ContextKeyRequestID ContextKey = "request_id"
	// ContextKeyUserID is the key for the user ID in the context
	ContextKeyUserID ContextKey = "user_id"
	// ContextKeyMethod is the key for the method name in the context
	ContextKeyMethod ContextKey = "method"
	// ContextKeyTraceID is the key for the trace ID in the context
	ContextKeyTraceID ContextKey = "trace_id"
)

// WithRequestID adds a request ID to the context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ContextKeyRequestID, requestID)
}

// GetRequestID gets the request ID from the context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(ContextKeyRequestID).(string); ok {
		return requestID
	}
	return ""
}

// WithUserID adds a user ID to the context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ContextKeyUserID, userID)
}

// GetUserID gets the user ID from the context
func GetUserID(ctx context.Context) (string, bool) {
	if userID, ok := ctx.Value(ContextKeyUserID).(string); ok {
		return userID, true
	}
	return "", false
}

// WithMethod adds a method name to the context
func WithMethod(ctx context.Context, method string) context.Context {
	return context.WithValue(ctx, ContextKeyMethod, method)
}

// GetMethod gets the method name from the context
func GetMethod(ctx context.Context) string {
	if method, ok := ctx.Value(ContextKeyMethod).(string); ok {
		return method
	}
	return ""
}

// WithTraceID adds a trace ID to the context
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, ContextKeyTraceID, traceID)
}

// GetTraceID gets the trace ID from the context
func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(ContextKeyTraceID).(string); ok {
		return traceID
	}
	return ""
}
