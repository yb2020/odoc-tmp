package context

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

// SafeGoroutine 安全地在 goroutine 中执行函数，自动传递标准上下文
func SafeGoroutine(ctx context.Context, fn func(ctx context.Context)) {
	go func(ctx context.Context) {
		defer func() {
			if r := recover(); r != nil {
				// 记录 panic 信息和堆栈跟踪
				stack := debug.Stack()
				// 这里应该使用项目的日志系统记录错误
				// 例如: logger.Error("goroutine panic", "error", r, "stack", string(stack))
				// 简单起见，这里使用标准输出
				println("Goroutine panic:", r)
				println(string(stack))
			}
		}()
		fn(ctx)
	}(ctx)
}

// GinSafeGoroutine 安全地在 goroutine 中执行函数，自动从 Gin 上下文中提取用户上下文
func GinSafeGoroutine(c *gin.Context, fn func(ctx context.Context, uc *UserContext)) {
	// 获取标准上下文
	ctx := c.Request.Context()

	// 从 Gin 上下文中获取用户上下文
	var uc *UserContext
	if ucValue, exists := c.Get("userContext"); exists && ucValue != nil {
		if ucTyped, ok := ucValue.(*UserContext); ok {
			// 创建用户上下文的副本
			uc = ucTyped.Clone()
		}
	} else {
		// 如果没有找到用户上下文，尝试从 Gin 上下文中提取用户信息
		uc = FromGinContext(c)
	}

	// 启动 goroutine
	go func(ctx context.Context, uc *UserContext) {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				// 这里应该使用项目的日志系统记录错误
				println("Goroutine panic:", r)
				println(string(stack))
			}
		}()
		fn(ctx, uc)
	}(ctx, uc)
}

// WithTimeout 创建一个带超时的上下文，并在超时后自动取消
func WithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, timeout)
}

// WithCancel 创建一个可取消的上下文
func WithCancel(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(ctx)
}

// BackgroundContext 创建一个带用户上下文的后台上下文
func BackgroundContext(uc *UserContext) context.Context {
	return uc.ToContext(context.Background())
}

// RunWithUserContext 使用用户上下文运行函数，适用于非HTTP请求的场景
func RunWithUserContext(uc *UserContext, fn func(ctx context.Context)) {
	ctx := BackgroundContext(uc)
	fn(ctx)
}

// RunAsyncWithUserContext 使用用户上下文异步运行函数，适用于非HTTP请求的场景
func RunAsyncWithUserContext(uc *UserContext, fn func(ctx context.Context)) {
	ctx := BackgroundContext(uc)
	SafeGoroutine(ctx, fn)
}
