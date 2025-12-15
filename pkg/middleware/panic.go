package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-kit/kit/endpoint"
	"github.com/yb2020/odoc/pkg/logging"
)

// ErrorReporter 定义了错误报告接口
type ErrorReporter interface {
	ReportError(ctx context.Context, err error, stack []byte) error
}

// DefaultErrorReporter 是默认的错误报告实现
type DefaultErrorReporter struct {
	logger logging.Logger
	tracer opentracing.Tracer
}

// NewDefaultErrorReporter 创建一个新的默认错误报告器
func NewDefaultErrorReporter(logger logging.Logger, tracer opentracing.Tracer) *DefaultErrorReporter {
	return &DefaultErrorReporter{
		logger: logger,
		tracer: tracer,
	}
}

// ReportError 报告错误到日志和追踪系统
func (r *DefaultErrorReporter) ReportError(ctx context.Context, err error, stack []byte) error {
	// 记录到日志
	r.logger.Error(
		"msg", "Panic recovered",
		"error", err,
		"stack", string(stack),
	)

	// 记录到追踪系统
	if r.tracer != nil {
		span := opentracing.SpanFromContext(ctx)
		if span != nil {
			ext.Error.Set(span, true)
			span.LogKV(
				"event", "error",
				"error.kind", "panic",
				"error.object", err,
				"error.stack", string(stack),
			)
		}
	}

	return err
}

// PanicHandlingMiddleware 创建一个中间件，用于捕获 panic 并记录堆栈信息到 Jaeger span
func PanicHandlingMiddleware(reporter ErrorReporter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func() {
				if r := recover(); r != nil {
					// 获取堆栈信息
					stackTrace := debug.Stack()

					// 转换 panic 值为错误
					var recoverErr error
					switch t := r.(type) {
					case string:
						recoverErr = fmt.Errorf("panic: %s", t)
					case error:
						recoverErr = t
					default:
						recoverErr = fmt.Errorf("panic: %v", r)
					}

					// 报告错误
					reporter.ReportError(ctx, recoverErr, stackTrace)

					// 返回错误
					err = recoverErr
				}
			}()

			// 调用下一个中间件或端点
			return next(ctx, request)
		}
	}
}

// HTTPPanicHandlingMiddleware 创建一个 HTTP 中间件，用于捕获 HTTP 处理过程中的 panic
// 并将堆栈信息记录到 Jaeger span
func HTTPPanicHandlingMiddleware(reporter ErrorReporter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 捕获 panic
			defer func() {
				if err := recover(); err != nil {
					// 获取堆栈信息
					stackTrace := debug.Stack()

					// 转换 panic 值为错误
					var recoverErr error
					switch t := err.(type) {
					case string:
						recoverErr = fmt.Errorf("panic: %s", t)
					case error:
						recoverErr = t
					default:
						recoverErr = fmt.Errorf("panic: %v", err)
					}

					// 报告错误
					reporter.ReportError(r.Context(), recoverErr, stackTrace)

					// 返回 500 错误
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()

			// 处理请求
			next.ServeHTTP(w, r)
		})
	}
}

// GinPanicHandlingMiddleware 创建一个 Gin 中间件，用于捕获 Gin 处理过程中的 panic
func GinPanicHandlingMiddleware(reporter ErrorReporter) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取堆栈信息
				stackTrace := debug.Stack()

				// 转换 panic 值为错误
				var recoverErr error
				switch t := err.(type) {
				case string:
					recoverErr = fmt.Errorf("panic: %s", t)
				case error:
					recoverErr = t
				default:
					recoverErr = fmt.Errorf("panic: %v", err)
				}

				// 报告错误
				reporter.ReportError(c.Request.Context(), recoverErr, stackTrace)

				// 中止请求处理并返回 500 错误
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		// 处理请求
		c.Next()
	}
}

// GRPCPanicHandlingUnaryInterceptor 创建一个 gRPC 一元拦截器，用于捕获 gRPC 处理过程中的 panic
func GRPCPanicHandlingUnaryInterceptor(reporter ErrorReporter) grpc.UnaryServerInterceptor {
	// 自定义恢复函数
	recoveryFunc := func(p interface{}) error {
		// 获取堆栈信息
		stackTrace := debug.Stack()

		// 转换 panic 值为错误
		var recoverErr error
		switch t := p.(type) {
		case string:
			recoverErr = fmt.Errorf("panic: %s", t)
		case error:
			recoverErr = t
		default:
			recoverErr = fmt.Errorf("panic: %v", p)
		}

		// 报告错误（注意：这里没有上下文，所以传入 context.Background()）
		reporter.ReportError(context.Background(), recoverErr, stackTrace)

		// 返回 gRPC 错误
		return status.Errorf(codes.Internal, "panic: %v", p)
	}

	// 使用 grpc_recovery 包创建拦截器
	return grpc_recovery.UnaryServerInterceptor(
		grpc_recovery.WithRecoveryHandler(recoveryFunc),
	)
}

// GRPCPanicHandlingStreamInterceptor 创建一个 gRPC 流拦截器，用于捕获 gRPC 流处理过程中的 panic
func GRPCPanicHandlingStreamInterceptor(reporter ErrorReporter) grpc.StreamServerInterceptor {
	// 自定义恢复函数
	recoveryFunc := func(p interface{}) error {
		// 获取堆栈信息
		stackTrace := debug.Stack()

		// 转换 panic 值为错误
		var recoverErr error
		switch t := p.(type) {
		case string:
			recoverErr = fmt.Errorf("panic: %s", t)
		case error:
			recoverErr = t
		default:
			recoverErr = fmt.Errorf("panic: %v", p)
		}

		// 报告错误（注意：这里没有上下文，所以传入 context.Background()）
		reporter.ReportError(context.Background(), recoverErr, stackTrace)

		// 返回 gRPC 错误
		return status.Errorf(codes.Internal, "panic: %v", p)
	}

	// 使用 grpc_recovery 包创建拦截器
	return grpc_recovery.StreamServerInterceptor(
		grpc_recovery.WithRecoveryHandler(recoveryFunc),
	)
}

// CreateGRPCServerOptions 创建包含错误处理的 gRPC 服务器选项
func CreateGRPCServerOptions(reporter ErrorReporter, extraUnaryInterceptors []grpc.UnaryServerInterceptor, extraStreamInterceptors []grpc.StreamServerInterceptor) []grpc.ServerOption {
	// 创建 panic 处理拦截器
	panicUnaryInterceptor := GRPCPanicHandlingUnaryInterceptor(reporter)
	panicStreamInterceptor := GRPCPanicHandlingStreamInterceptor(reporter)

	// 合并拦截器
	unaryInterceptors := append([]grpc.UnaryServerInterceptor{panicUnaryInterceptor}, extraUnaryInterceptors...)
	streamInterceptors := append([]grpc.StreamServerInterceptor{panicStreamInterceptor}, extraStreamInterceptors...)

	// 创建服务器选项
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryInterceptors...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamInterceptors...)),
	}
}
