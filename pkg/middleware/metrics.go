package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/metrics"
)

// MetricsMiddleware returns an endpoint middleware that records metrics for each request
func MetricsMiddleware(metrics *metrics.Metrics) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			method := GetMethod(ctx)
			if method == "" {
				method = "unknown"
			}

			endpoint := "unknown"
			if e, ok := ctx.Value("endpoint").(string); ok {
				endpoint = e
			}

			// Start timer
			startTime := time.Now()

			// Execute the next endpoint
			resp, err := next(ctx, request)

			// Record latency
			latency := time.Since(startTime).Seconds()
			metrics.ObserveRequestLatency(method, endpoint, latency)

			// Record status and errors
			status := "success"
			if err != nil {
				status = "error"

				// Determine error type
				errorType := "internal"
				if appErr, ok := err.(*errors.SystemError); ok {
					errorType = string(appErr.Type)
				}

				// Record error count
				metrics.IncrementErrorCount(method, endpoint, errorType)
			}

			// Record request count
			metrics.IncrementRequestCount(method, endpoint, status)

			return resp, err
		}
	}
}

// HTTPMetricsMiddleware is a middleware that records metrics for HTTP requests
func HTTPMetricsMiddleware(metrics *metrics.Metrics) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			method := r.Method
			endpoint := r.URL.Path

			// Start timer
			startTime := time.Now()

			// Create a response wrapper to capture the status code
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Execute the next handler
			next.ServeHTTP(rw, r)

			// Record latency
			latency := time.Since(startTime).Seconds()
			metrics.ObserveRequestLatency(method, endpoint, latency)

			// Record request count
			metrics.IncrementRequestCount(method, endpoint, strconv.Itoa(rw.statusCode))

			// Record error count if status code >= 400
			if rw.statusCode >= 400 {
				errorType := "http_" + strconv.Itoa(rw.statusCode/100) + "xx"
				metrics.IncrementErrorCount(method, endpoint, errorType)
			}
		})
	}
}

// responseWriter is a wrapper around http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code and calls the underlying ResponseWriter's WriteHeader
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures the status code if not already set and calls the underlying ResponseWriter's Write
func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	return rw.ResponseWriter.Write(b)
}
