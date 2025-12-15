package mq_interface

import (
	"context"
)

// Message 消息接口
type Message interface {
	// GetTopic 获取消息主题
	GetTopic() string
	// GetTags 获取消息标签
	GetTags() string
	// GetKeys 获取消息键
	GetKeys() string
	// GetUserId 获取消息用户ID
	GetUserId() string
	// GetBusinessKey 获取消息业务键，从一定程度上解决消费消息的幂等性问题，不保证100%，需要100%的情况下，仍然需要自己在业务流上做唯一性判断
	GetBusinessKey() string
	// GetMessageId 获取消息ID
	GetMessageId() string
	// GetBody 获取消息体
	GetBody() []byte
	// GetProperties 获取消息属性
	GetProperties() map[string]string
}

// Producer 生产者接口
type Producer interface {
	// Start 启动生产者
	Start() error
	// Shutdown 关闭生产者
	Shutdown() error
	// SendSync 同步发送消息
	SendSync(ctx context.Context, msg Message) (string, error)
	// SendAsync 异步发送消息
	SendAsync(ctx context.Context, msg Message, callback func(context.Context, string, error)) error
	// SendOneWay 单向发送消息（不关心结果）
	SendOneWay(ctx context.Context, msg Message) error
}

// Consumer 消费者接口
type Consumer interface {
	// Start 启动消费者
	Start() error
	// Shutdown 关闭消费者
	Shutdown() error
	// Subscribe 订阅主题
	Subscribe(handler func(context.Context, Message) error) error
}

// MessageQueue 消息队列接口
type MessageQueue interface {
	// NewProducer 创建生产者
	NewProducer(group string) (Producer, error)
	// NewConsumer 创建消费者
	NewConsumer(group string) (Consumer, error)
}
