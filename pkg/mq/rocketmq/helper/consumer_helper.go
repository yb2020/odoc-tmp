package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/distlock"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	mqi "github.com/yb2020/odoc/pkg/mq/interface"
	"github.com/yb2020/odoc/pkg/mq/rocketmq/consumer"
)

// MessageHandlerHelper 提供分布式锁功能的消息处理器
type MessageHandlerHelper struct {
	lockTemplate *distlock.LockTemplate
	logger       logging.Logger
	tracer       opentracing.Tracer
}

type HandlerLockOptions struct {
	// 锁key
	LockKey string
	// 锁过期时间
	LockExpiry *time.Duration
	// 锁超时时间
	LockTimeout *time.Duration
}

// NewMessageHandlerHelper 创建新的带锁消息处理器
func NewMessageHandlerHelper(
	lockTemplate *distlock.LockTemplate,
	logger logging.Logger,
	tracer opentracing.Tracer,
) *MessageHandlerHelper {
	return &MessageHandlerHelper{
		lockTemplate: lockTemplate,
		logger:       logger,
		tracer:       tracer,
	}
}

// withDistributedLock 使用分布式锁包装消息处理函数
func (h *MessageHandlerHelper) withDistributedLock(
	isStartLock bool,
	lockOptions HandlerLockOptions,
	handler func(context.Context, mqi.Message) error,
) func(context.Context, mqi.Message) error {
	return func(ctx context.Context, msg mqi.Message) error {
		lockKey := fmt.Sprintf("mq:biz:lock:%s:%s", msg.GetTopic(), lockOptions.LockKey)
		if isStartLock {
			lockKey = fmt.Sprintf("mq:public:lock:%s:%s:%s", msg.GetTopic(), lockOptions.LockKey, msg.GetMessageId())
		}
		distLock := &distlock.LockInfo{}

		// 获取锁过期时间，如果是整数值则自动转换为毫秒单位
		var lockExpiry time.Duration
		if lockOptions.LockExpiry == nil {
			lockExpiry = 30 * time.Second
		} else if *lockOptions.LockExpiry < time.Millisecond {
			// 如果值小于1毫秒，认为它是一个整数值，需要转换为毫秒
			lockExpiry = *lockOptions.LockExpiry * time.Millisecond
		} else {
			lockExpiry = *lockOptions.LockExpiry
		}

		// 获取锁超时时间，如果是整数值则自动转换为毫秒单位
		var lockTimeout time.Duration
		if lockOptions.LockTimeout == nil {
			lockTimeout = 5 * time.Second
		} else if *lockOptions.LockTimeout < time.Millisecond {
			// 如果值小于1毫秒，认为它是一个整数值，需要转换为毫秒
			lockTimeout = *lockOptions.LockTimeout * time.Millisecond
		} else {
			lockTimeout = *lockOptions.LockTimeout
		}

		// 尝试获取锁，设置适当的超时和重试参数
		distLock = h.lockTemplate.Lock(
			lockKey,
			lockExpiry,
			lockTimeout,
		)

		// 如果获取锁失败，表示消息正在被其他实例处理,返回错误,需要mq重试
		if distLock == nil {
			h.logger.Info("消息正在被其他实例处理，跳过",
				"topic", msg.GetTopic(),
				"messageId", msg.GetMessageId())
			// 返回错误表示消息处理失败，需要mq重试
			return errors.Biz("mq.UploadCallbackService error, failed to get lock!")
		}

		// 确保锁会被释放
		defer h.lockTemplate.ReleaseLock(distLock)

		// 使用span跟踪处理过程
		var span opentracing.Span
		var newCtx context.Context

		if h.tracer != nil {
			span, newCtx = opentracing.StartSpanFromContextWithTracer(
				ctx,
				h.tracer,
				fmt.Sprintf("%s.handleMessage", lockKey),
			)
			defer span.Finish()
		} else {
			newCtx = ctx
		}

		// 调用原始处理函数
		return handler(newCtx, msg)
	}
}

func (h *MessageHandlerHelper) HandlerMessageWithLock(ctx context.Context,
	msg mqi.Message,
	lockOptions HandlerLockOptions,
	handler func(context.Context, mqi.Message) error) error {
	return h.withDistributedLock(false, lockOptions, handler)(ctx, msg)
}

// StartConsumerWithLock 启动带有分布式锁的消费者
func (h *MessageHandlerHelper) StartConsumerWithLock(
	consumerClient *consumer.RocketMQConsumer,
	lockOptions HandlerLockOptions,
	handler func(context.Context, mqi.Message) error,
) (mqi.Consumer, error) {
	return h.startConsumer(consumerClient, &lockOptions, handler)
}

// StartConsumer 启动不带分布式锁的消费者
func (h *MessageHandlerHelper) StartConsumer(
	consumerClient *consumer.RocketMQConsumer,
	handler func(context.Context, mqi.Message) error,
) (mqi.Consumer, error) {
	return h.startConsumer(consumerClient, nil, handler)
}

func (h *MessageHandlerHelper) startConsumer(
	consumerClient *consumer.RocketMQConsumer,
	lockOptions *HandlerLockOptions,
	handler func(context.Context, mqi.Message) error,
) (mqi.Consumer, error) {

	// 使用传入的消费者实例，不再尝试创建新的消费者
	if consumerClient == nil {
		h.logger.Error("CRITICAL: Consumer client is nil, cannot start consumer")
		return nil, fmt.Errorf("consumer client is nil, cannot start consumer")
	}

	// 启动消费者
	if err := consumerClient.Start(); err != nil {
		h.logger.Error("Failed to start consumer",
			"topic", consumerClient.GetTopic(),
			"consumerGroup", consumerClient.GetConsumerGroup(),
			"error", err)
		return nil, fmt.Errorf("failed to start consumer for topic:%s, consumerGroup:%s: %w",
			consumerClient.GetTopic(), consumerClient.GetConsumerGroup(), err)
	}

	if lockOptions == nil {
		// 不带锁订阅主题
		// 在V5版本中，Subscribe方法内部会启动消息接收循环
		// 消息处理函数会在接收到消息时被调用
		err := consumerClient.Subscribe(handler)
		if err != nil {
			h.logger.Error("Failed to subscribe to topic",
				"topic", consumerClient.GetTopic(),
				"consumerGroup", consumerClient.GetConsumerGroup(),
				"error", err)
			return nil, fmt.Errorf("failed to subscribe to topic for %s: %w", consumerClient.GetTopic(), err)
		}

		h.logger.Info("Successfully subscribed to topic without distributed lock",
			"topic", consumerClient.GetTopic(),
			"consumerGroup", consumerClient.GetConsumerGroup())
		return consumerClient, nil
	}

	// 带锁订阅主题
	h.logger.Info("Subscribing to topic with distributed lock",
		"topic", consumerClient.GetTopic(),
		"consumerGroup", consumerClient.GetConsumerGroup(),
		"lockKey", lockOptions.LockKey)

	// 创建带分布式锁的处理函数
	// 在V5版本中，我们需要确保消息处理失败时不确认消息，让消息重新投递
	lockHandler := h.withDistributedLock(
		true,
		*lockOptions,
		handler,
	)

	// 订阅主题
	err := consumerClient.Subscribe(lockHandler)
	if err != nil {
		h.logger.Error("Failed to subscribe to topic with distributed lock",
			"topic", consumerClient.GetTopic(),
			"consumerGroup", consumerClient.GetConsumerGroup(),
			"lockKey", lockOptions.LockKey,
			"error", err)
		return nil, fmt.Errorf("failed to subscribe to topic for %s: %w", consumerClient.GetTopic(), err)
	}

	h.logger.Info("Successfully subscribed to topic with distributed lock",
		"topic", consumerClient.GetTopic(),
		"consumerGroup", consumerClient.GetConsumerGroup(),
		"lockKey", lockOptions.LockKey)

	return consumerClient, nil
}
