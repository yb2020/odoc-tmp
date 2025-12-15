package consumer

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/logging"
	mqi "github.com/yb2020/odoc/pkg/mq/interface"
	"github.com/yb2020/odoc/pkg/mq/rocketmq"

	rmq "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
)

// RocketMQConsumerV5 RocketMQ V5消费者
type RocketMQConsumer struct {
	config                   *config.Config
	logger                   logging.Logger
	consumer                 rmq.SimpleConsumer
	consumerConfig           *MQConsumerOptions
	exitChan                 chan struct{}
	wg                       sync.WaitGroup
	running                  bool
	mutex                    sync.Mutex
	handleError              func(context.Context, mqi.Message, error) (bool, error)
	handleMaxRetriesExceeded func(context.Context, mqi.Message, error) error
}

// MQConsumerOptions 消费者配置
type MQConsumerOptions struct {
	Topic             string
	ConsumerGroup     string
	ConsumerTag       string
	BatchSize         int  // 批量拉取消息数量
	MaxReconsumeTimes int  // 消费失败重试次数
	ConsumeTimeout    int  // 默认消费超时时间（毫秒）
	ManualAck         bool // 是否手动提交偏移量
	MaxAwaitTime      int  // 消息接收等待时间（毫秒）
}

// NewRocketMQConsumerV5 创建新的RocketMQ V5消费者实例
func NewRocketMQConsumer(config *config.Config, logger logging.Logger, consumerConfig MQConsumerOptions) (*RocketMQConsumer, error) {
	if !config.RocketMQ.Enabled {
		return nil, errors.New("RocketMQ is not enabled")
	}

	return &RocketMQConsumer{
		config:         config,
		logger:         logger,
		consumerConfig: &consumerConfig,
		exitChan:       make(chan struct{}),
		running:        false,
	}, nil
}

// --- 新增方法 ---
// SetErrorHandler 注册一个自定义的错误处理函数
func (c *RocketMQConsumer) SetErrorHandler(handler func(context.Context, mqi.Message, error) (bool, error)) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.handleError = handler
}

// SetMaxRetriesExceededHandler 注册一个超过最大重试次数后的处理函数
func (c *RocketMQConsumer) SetMaxRetriesExceededHandler(handler func(context.Context, mqi.Message, error) error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.handleMaxRetriesExceeded = handler
}

// Start 启动消费者
func (c *RocketMQConsumer) Start() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.running {
		return nil
	}
	// 设置订阅表达式
	subscriptionExpressions := map[string]*rmq.FilterExpression{
		c.consumerConfig.Topic: rmq.NewFilterExpression(c.consumerConfig.ConsumerTag),
	}

	// 设置消息接收等待时间
	awaitDuration := time.Duration(c.consumerConfig.MaxAwaitTime) * time.Millisecond

	// 创建消费者配置
	rmqConfig := &rmq.Config{
		Endpoint:      c.config.RocketMQ.GrpcAddress,
		ConsumerGroup: c.consumerConfig.ConsumerGroup,
	}

	// 如果配置了AccessKey和SecretKey，则添加认证信息
	if c.config.RocketMQ.AccessKey != "" && c.config.RocketMQ.SecretKey != "" {
		rmqConfig.Credentials = &credentials.SessionCredentials{
			AccessKey:    c.config.RocketMQ.AccessKey,
			AccessSecret: c.config.RocketMQ.SecretKey,
		}
	}

	// 创建消费者选项
	opts := []rmq.SimpleConsumerOption{
		rmq.WithAwaitDuration(awaitDuration),
		rmq.WithSubscriptionExpressions(subscriptionExpressions),
	}

	var err error
	c.consumer, err = rmq.NewSimpleConsumer(rmqConfig, opts...)
	if err != nil {
		c.logger.Error("CRITICAL: Failed to create RocketMQ V5 consumer",
			"error", err,
			"group", c.consumerConfig.ConsumerGroup,
			"endpoint", c.config.RocketMQ.GrpcAddress)
		return fmt.Errorf("new RocketMQ V5 consumer failed: %w", err)
	}

	// 启动消费者
	err = c.consumer.Start()
	if err != nil {
		c.logger.Error("CRITICAL: Failed to start RocketMQ V5 consumer",
			"error", err,
			"group", c.consumerConfig.ConsumerGroup,
			"topic", c.consumerConfig.Topic)
		return fmt.Errorf("start RocketMQ V5 consumer failed: %w", err)
	}
	c.running = true
	c.exitChan = make(chan struct{})

	c.logger.Info("RocketMQ V5 consumer started successfully",
		"group", c.consumerConfig.ConsumerGroup,
		"topic", c.consumerConfig.Topic,
		"tag", c.consumerConfig.ConsumerTag)
	return nil
}

// Shutdown 关闭消费者
func (c *RocketMQConsumer) Shutdown() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.running {
		return nil
	}

	if c.consumer == nil {
		return nil
	}

	// 关闭退出通道，通知所有接收循环退出
	close(c.exitChan)

	// 等待所有接收循环退出
	c.wg.Wait()

	// 优雅关闭消费者
	c.logger.Info("Gracefully stopping RocketMQ V5 consumer", "group", c.consumerConfig.ConsumerGroup)
	c.consumer.GracefulStop()
	c.logger.Info("RocketMQ V5 consumer shutdown successfully", "group", c.consumerConfig.ConsumerGroup)

	c.running = false
	return nil
}

// Subscribe 订阅主题
func (c *RocketMQConsumer) Subscribe(handler func(context.Context, mqi.Message) error) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.consumer == nil {
		c.logger.Error("CRITICAL: Cannot subscribe - consumer is not initialized",
			"topic", c.consumerConfig.Topic,
			"consumerGroup", c.consumerConfig.ConsumerGroup)
		return errors.New("consumer is not initialized, call Start() first")
	}

	if !c.running {
		c.logger.Error("CRITICAL: Cannot subscribe - consumer is not running",
			"topic", c.consumerConfig.Topic,
			"consumerGroup", c.consumerConfig.ConsumerGroup)
		return errors.New("consumer is not running, call Start() first")
	}

	// 启动消息接收循环
	c.wg.Add(1)
	go c.receiveMessages(handler)

	return nil
}

// receiveMessages 接收消息循环
func (c *RocketMQConsumer) receiveMessages(handler func(context.Context, mqi.Message) error) {
	defer c.wg.Done()

	// 设置消息接收参数 - 批量拉取消息数量
	maxMessageNum := int32(c.consumerConfig.BatchSize)
	if maxMessageNum <= 0 {
		maxMessageNum = 16 // 默认值
	}

	// 消费超时时间 - 用于业务处理逻辑
	consumeTimeout := c.consumerConfig.ConsumeTimeout // 默认超时时间（毫秒）

	// 确保消费超时时间有效
	if consumeTimeout <= 0 {
		consumeTimeout = 60000 // 默认设置为60秒
	}

	// 消息可见性时间 - 设置和超时时间保持一致
	invisibleDuration := time.Duration(consumeTimeout) * time.Millisecond
	// invisibleDuration := time.Duration(10000) * time.Millisecond

	for {
		select {
		case <-c.exitChan:
			return
		default:
			// 接收消息
			messages, err := c.consumer.Receive(context.Background(), maxMessageNum, invisibleDuration)

			if err != nil {

				time.Sleep(time.Duration(c.consumerConfig.MaxAwaitTime) * time.Millisecond) // 失败后等待一段时间再重试
				continue
			}
			// 处理接收到的消息
			if len(messages) == 0 {
				continue
			}

			for _, msg := range messages {
				// 转换为通用消息格
				message := c.fromRocketMQMessage(msg)

				// 创建带超时的上下文
				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(consumeTimeout)*time.Millisecond)

				// 调用处理函数
				err := handler(ctx, message)

				// 无论成功失败都取消上下文
				cancel()

				if err != nil {
					// 获取当前消息的重试次数
					deliveryAttempt := msg.GetDeliveryAttempt()

					if errors.Is(err, context.DeadlineExceeded) {
						c.logger.Error("消息处理超时",
							"messageId", msg.GetMessageId(),
							"timeout", consumeTimeout,
							"topic", c.consumerConfig.Topic,
							"deliveryAttempt", deliveryAttempt)
					} else {
						c.logger.Error("消息处理失败",
							"messageId", msg.GetMessageId(),
							"error", err,
							"topic", c.consumerConfig.Topic,
							"deliveryAttempt", deliveryAttempt)
					}
					// 如果注册了自定义错误处理器，则调用它
					if c.handleError != nil {
						// 使用一个新的后台上下文来执行错误处理，以防上面的ctx已经超时
						hasAutoAck, handleErr := c.handleError(context.Background(), message, err)
						if handleErr != nil {
							c.logger.Error("自定义错误处理失败",
								"messageId", msg.GetMessageId(),
								"error", handleErr,
								"topic", c.consumerConfig.Topic,
								"deliveryAttempt", deliveryAttempt)
						}
						if hasAutoAck {
							// 自动确认消息
							if err := c.consumer.Ack(context.Background(), msg); err != nil {
								c.logger.Error("确认消息失败",
									"messageId", msg.GetMessageId(),
									"error", err)
							}
						}
					}
					// 检查重试次数是否超过最大重试次数
					if c.consumerConfig.MaxReconsumeTimes > 0 && int(deliveryAttempt) >= c.consumerConfig.MaxReconsumeTimes {
						// 超过最大重试次数，确认消息以避免消息堆积
						if c.handleMaxRetriesExceeded != nil {
							if err := c.handleMaxRetriesExceeded(context.Background(), message, err); err != nil {
								c.logger.Error("超限处理器执行失败",
									"messageId", msg.GetMessageId(),
									"error", err,
									"topic", c.consumerConfig.Topic,
									"deliveryAttempt", deliveryAttempt)
							}
						}

						if err := c.consumer.Ack(context.Background(), msg); err != nil {
							c.logger.Error("确认超过重试次数的消息失败",
								"messageId", msg.GetMessageId(),
								"error", err)
						}
					} else {
						// 未超过最大重试次数，不确认消息，让它重新投递
					}
					continue
				}

				// 根据ManualAck配置决定是否手动确认消息
				// 如果ManualAck为false，则由框架自动确认
				// 如果ManualAck为true，则需要消息处理函数自己实现确认逻辑
				if !c.consumerConfig.ManualAck {
					// 确认消息
					if err := c.consumer.Ack(context.Background(), msg); err != nil {
						c.logger.Error("确认消息失败",
							"messageId", msg.GetMessageId(),
							"error", err)
					}
				}
			}

			// 如果没有收到消息，等待一段时间再尝试
			if len(messages) == 0 {
				time.Sleep(time.Duration(c.consumerConfig.MaxAwaitTime) * time.Millisecond)
			}
		}
	}
}

// fromRocketMQMessage 将RocketMQ V5消息转换为通用消息
func (c *RocketMQConsumer) fromRocketMQMessage(rmqMsg *rmq.MessageView) mqi.Message {
	// 在V5中，很多方法返回指针，需要处理空指针情况
	var tag string
	if rmqMsg.GetTag() != nil {
		tag = *rmqMsg.GetTag()
	}

	// 获取消息键（可能为空）
	var keys string
	if rmqMsg.GetKeys() != nil && len(rmqMsg.GetKeys()) > 0 {
		keys = strings.Join(rmqMsg.GetKeys(), ",")
	}

	// 获取所有属性
	properties := rmqMsg.GetProperties()

	// 从属性中获取userId和businessKey
	userId := properties["userId"]
	businessKey := properties["businessKey"]

	return &rocketmq.Message{
		Topic:       rmqMsg.GetTopic(),
		Tags:        tag,
		Keys:        keys,
		UserId:      userId,
		MessageId:   rmqMsg.GetMessageId(),
		BusinessKey: businessKey,
		Body:        rmqMsg.GetBody(),
		Properties:  properties,
	}
}

// GetTopic 获取主题
func (c *RocketMQConsumer) GetTopic() string {
	return c.consumerConfig.Topic
}

// GetConsumerGroup 获取消费者组
func (c *RocketMQConsumer) GetConsumerGroup() string {
	return c.consumerConfig.ConsumerGroup
}
