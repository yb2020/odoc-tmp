// pkg/eventbus/eventbus.go
package eventbus

import (
	"context"
	"sync"

	appContext "github.com/yb2020/odoc/pkg/context"
)

type EventType string

type Event struct {
	Type EventType
	Data any
}

type EventHandler func(ctx context.Context, event Event)

type EventBus struct {
	handlers map[EventType][]EventHandler
	mu       sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[EventType][]EventHandler),
	}
}

func (b *EventBus) Subscribe(eventType EventType, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.handlers == nil {
		b.handlers = make(map[EventType][]EventHandler)
	}

	if _, exists := b.handlers[eventType]; !exists {
		b.handlers[eventType] = make([]EventHandler, 0)
	}

	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// PublishSync 同步发布事件，直接在当前 goroutine 中处理
func (b *EventBus) publishSync(ctx context.Context, event Event) {
	b.mu.RLock()
	handlers, exists := b.handlers[event.Type]
	b.mu.RUnlock()

	if !exists {
		return
	}

	for _, handler := range handlers {
		// 同步处理，直接使用原始上下文
		handler(ctx, event)
	}
}

// PublishAsync 异步发布事件，在新的 goroutine 中处理
// 会从原始上下文中提取用户上下文并在新的 goroutine 中使用
func (b *EventBus) publishAsync(ctx context.Context, event Event) {
	b.mu.RLock()
	handlers, exists := b.handlers[event.Type]
	b.mu.RUnlock()

	if !exists {
		return
	}

	// 从原始上下文中获取用户上下文
	uc := appContext.GetUserContext(ctx)

	for _, handler := range handlers {
		// 使用安全的 goroutine 处理，保留用户上下文
		if uc != nil {
			// 如果有用户上下文，使用 RunAsyncWithUserContext
			h := handler // 创建一个副本以避免闭包问题
			appContext.RunAsyncWithUserContext(uc, func(eventCtx context.Context) {
				h(eventCtx, event)
			})
		} else {
			// 如果没有用户上下文，使用 SafeGoroutine
			h := handler // 创建一个副本以避免闭包问题
			appContext.SafeGoroutine(context.Background(), func(eventCtx context.Context) {
				h(eventCtx, event)
			})
		}
	}
}

// Publish 为了向后兼容，默认使用异步方式发布事件
func (b *EventBus) Publish(ctx context.Context, event Event, async bool) {
	if async {
		b.publishAsync(ctx, event)
	} else {
		b.publishSync(ctx, event)
	}
}
