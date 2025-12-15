package producer

import (
	"context"
	"errors"
	"strings"
	"sync"

	// V2 版本导入（已注释）
	// rocketmq "github.com/apache/rocketmq-client-go/core"
	// "github.com/apache/rocketmq-client-go/v2"
	// "github.com/apache/rocketmq-client-go/v2/primitive"
	// "github.com/apache/rocketmq-client-go/v2/producer"

	// V5 版本导入
	rmq "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/logging"
	mqi "github.com/yb2020/odoc/pkg/mq/interface"
	customRocketMQ "github.com/yb2020/odoc/pkg/mq/rocketmq"
)

// RocketMQProducerV5 实现 mq_interface.Producer 接口（V5版本）
type RocketMQProducer struct {
	producer rmq.Producer
	config   *config.Config
	logger   logging.Logger
	started  bool
	mutex    sync.Mutex
}

// ErrRocketMQNotEnabled RocketMQ未启用错误
var ErrRocketMQNotEnabled = errors.New("RocketMQ is not enabled")

// NewRocketMQProducerV5 创建新的RocketMQ V5版本生产者
func NewRocketMQProducer(config *config.Config, logger logging.Logger) (*RocketMQProducer, error) {
	if !config.RocketMQ.Enabled {
		return nil, ErrRocketMQNotEnabled
	}

	logger.Info("creating RocketMQ V5 producer", "grpcAddress", config.RocketMQ.GrpcAddress)
	// 解析 grpcAddress 地址，格式为 ip:port
	parts := strings.Split(config.RocketMQ.GrpcAddress, ":")
	if len(parts) != 2 {
		return nil, errors.New("invalid grpcAddress format, should be ip:port")
	}

	// 创建生产者配置
	producerConfig := &rmq.Config{
		Endpoint: config.RocketMQ.GrpcAddress,
	}

	// 如果配置了AccessKey和SecretKey，则添加认证信息
	if config.RocketMQ.AccessKey != "" && config.RocketMQ.SecretKey != "" {
		producerConfig.Credentials = &credentials.SessionCredentials{
			AccessKey:    config.RocketMQ.AccessKey,
			AccessSecret: config.RocketMQ.SecretKey, // 注意：V5 版本使用 AccessSecret 而不是 SecretKey
		}
	}

	// 重置日志配置以应用新的设置
	rmq.ResetLogger()

	// 收集所有可能使用的主题
	topics := []string{
		config.RocketMQ.Topic.UploadCallback.Name,
		config.RocketMQ.Topic.ParsePdfHeader.Name,
		config.RocketMQ.Topic.ParsePdfText.Name,
		config.RocketMQ.Topic.FullTextTranslateUploadHandler.Name,
		config.RocketMQ.Topic.FullTextTranslateProgress.Name,
	}

	// 检查事件主题是否存在
	if config.RocketMQ.Topic.Event.Doc2DifyIntegrationEvent.Name != "" {
		topics = append(topics, config.RocketMQ.Topic.Event.Doc2DifyIntegrationEvent.Name)
	}

	// 创建生产者
	// 注意：RocketMQ V5 客户端使用具体的选项函数，而不是 Option 类型
	// 使用WithTopics选项预先指定所有可能使用的主题，避免发送消息时查询路由信息导致的空指针异常
	p, err := rmq.NewProducer(
		producerConfig,
		rmq.WithTopics(topics...),
	)
	if err != nil {
		logger.Error("failed to create RocketMQ V5 producer", "error", err, "topics", topics)
		return nil, err
	}

	logger.Info("RocketMQ V5 producer created successfully", "topics", topics)
	return &RocketMQProducer{
		producer: p,
		config:   config,
		logger:   logger,
		started:  false,
		mutex:    sync.Mutex{},
	}, nil
}

// // NewRocketMQProducer 创建新的RocketMQ生产者
// func NewRocketMQProducer(config *config.Config, logger logging.Logger) (*RocketMQProducer, error) {
// 	if !config.RocketMQ.Enabled {
// 		return nil, ErrRocketMQNotEnabled
// 	}

// 	logger.Info("creating RocketMQ producer", "nameServer", config.RocketMQ.NameServer)

// 	// 创建生产者选项
// 	opts := []producer.Option{
// 		// producer.WithGroupName(group),
// 		// 使用WithNameServer而不是WithNameServerDomain，以正确处理IP:端口格式
// 		producer.WithNameServer([]string{config.RocketMQ.NameServer}),
// 		producer.WithRetry(3),                         // 默认重试3次
// 		producer.WithSendMsgTimeout(10 * time.Second), // 增加超时时间到10秒
// 		producer.WithCreateTopicKey("TBW102"),         // 添加自动创建主题的选项
// 	}

// 	// 如果配置了AccessKey和SecretKey，则添加认证信息
// 	if config.RocketMQ.AccessKey != "" && config.RocketMQ.SecretKey != "" {
// 		opts = append(opts, producer.WithCredentials(primitive.Credentials{
// 			AccessKey: config.RocketMQ.AccessKey,
// 			SecretKey: config.RocketMQ.SecretKey,
// 		}))
// 	}

// 	p, err := rocketmq.NewProducer(opts...)
// 	if err != nil {
// 		logger.Error("failed to create RocketMQ producer", "error", err)
// 		return nil, err
// 	}

// 	logger.Info("RocketMQ producer created successfully")
// 	return &RocketMQProducer{
// 		producer: p,
// 		config:   config,
// 		logger:   logger,
// 		started:  false,
// 	}, nil
// }

// // Start 启动生产者
// func (p *RocketMQProducer) Start() error {
// 	p.logger.Info("starting RocketMQ producer")

// 	if p.started {
// 		return nil
// 	}

// 	if p.producer == nil {
// 		return errors.New("producer not initialized")
// 	}

// 	err := p.producer.Start()
// 	if err != nil {
// 		p.logger.Error("failed to start RocketMQ producer", "error", err)
// 		return err
// 	}

// 	p.started = true
// 	p.logger.Info("RocketMQ producer started successfully")
// 	return nil
// }

// // Shutdown 关闭生产者
// func (p *RocketMQProducer) Shutdown() error {
// 	if p.producer == nil {
// 		return nil
// 	}

// 	err := p.producer.Shutdown()
// 	if err != nil {
// 		p.logger.Error("failed to shutdown RocketMQ producer", "error", err)
// 		return err
// 	}

// 	p.started = false
// 	p.logger.Info("RocketMQ producer shutdown successfully")
// 	return nil
// }

func (p *RocketMQProducer) Start() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.logger.Info("starting RocketMQ V5 producer")

	if p.started {
		return nil
	}

	if p.producer == nil {
		return errors.New("producer not initialized")
	}

	// 添加实际的Start调用
	err := p.producer.Start()
	if err != nil {
		p.logger.Error("failed to start RocketMQ V5 producer", "error", err)
		return err
	}

	p.started = true
	p.logger.Info("RocketMQ V5 producer started successfully")
	return nil
}

// Shutdown 关闭生产者 (V5版本)
func (p *RocketMQProducer) Shutdown() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.producer == nil {
		return nil
	}
	// V5 版本使用 GracefulStop 代替 Shutdown
	err := p.producer.GracefulStop()
	if err != nil {
		p.logger.Error("failed to shutdown RocketMQ V5 producer", "error", err)
		return err
	}

	p.started = false
	p.logger.Info("RocketMQ V5 producer shutdown successfully")
	return nil
}

// // SendSync 同步发送消息
// func (p *RocketMQProducer) SendSync(ctx context.Context, message mqi.Message) (string, error) {
// 	if !p.started {
// 		return "", errors.New("producer not started")
// 	}

// 	msg := ToRocketMQMessage(message)
// 	result, err := p.producer.SendSync(ctx, msg)
// 	if err != nil {
// 		p.logger.Error("failed to send message synchronously",
// 			"topic", message.GetTopic(),
// 			"tags", message.GetTags(),
// 			"keys", message.GetKeys(),
// 			"error", err)
// 		return "", err
// 	}

// 	p.logger.Debug("message sent synchronously",
// 		"topic", message.GetTopic(),
// 		"tags", message.GetTags(),
// 		"keys", message.GetKeys(),
// 		"messageId", result.MsgID)
// 	return result.MsgID, nil
// }

// SendSync 同步发送消息 (V5版本)
func (p *RocketMQProducer) SendSync(ctx context.Context, message mqi.Message) (string, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.started {
		return "", errors.New("producer not started")
	}

	msg := customRocketMQ.ToRocketMQMessage(message)

	// V5 版本的 Send 方法返回的是一个切片，可能包含多个结果
	results, err := p.producer.Send(ctx, msg)
	if err != nil {
		p.logger.Error("failed to send message synchronously",
			"topic", message.GetTopic(),
			"tags", message.GetTags(),
			"keys", message.GetKeys(),
			"error", err)
		return "", err
	}

	// 检查是否有结果返回
	if len(results) == 0 {
		p.logger.Error("no receipt returned from RocketMQ V5 producer")
		return "", errors.New("no receipt returned")
	}

	// 如果返回多个结果，记录该情况，但仍然取第一个
	if len(results) > 1 {
		p.logger.Warn("multiple receipts returned from RocketMQ V5 producer",
			"count", len(results),
			"topic", message.GetTopic(),
			"tags", message.GetTags(),
			"keys", message.GetKeys())
	}

	// 取第一个结果的消息 ID
	// 根据 RocketMQ V5 客户端的行为，单条消息发送时应该只返回一个结果
	messageId := results[0].MessageID

	return messageId, nil
}

// // SendAsync 异步发送消息
// func (p *RocketMQProducer) SendAsync(ctx context.Context, message mqi.Message, callback func(context.Context, string, error)) error {
// 	if !p.started {
// 		return errors.New("producer not started")
// 	}

// 	msg := ToRocketMQMessage(message)
// 	err := p.producer.SendAsync(ctx, func(ctx context.Context, result *primitive.SendResult, err error) {
// 		if err != nil {
// 			p.logger.Error("async message send failed",
// 				"topic", message.GetTopic(),
// 				"tags", message.GetTags(),
// 				"keys", message.GetKeys(),
// 				"error", err)
// 			if callback != nil {
// 				callback(ctx, "", err)
// 			}
// 			return
// 		}

// 		p.logger.Debug("message sent asynchronously",
// 			"topic", message.GetTopic(),
// 			"tags", message.GetTags(),
// 			"keys", message.GetKeys(),
// 			"messageId", result.MsgID)
// 		if callback != nil {
// 			callback(ctx, result.MsgID, nil)
// 		}
// 	}, msg)

// 	if err != nil {
// 		p.logger.Error("failed to send message asynchronously",
// 			"topic", message.GetTopic(),
// 			"tags", message.GetTags(),
// 			"keys", message.GetKeys(),
// 			"error", err)
// 		return err
// 	}

// 	return nil
// }

// SendAsync 异步发送消息 (V5版本)
// 注意：V5 版本没有直接的异步发送 API，我们使用 goroutine 包装同步发送实现异步发送
func (p *RocketMQProducer) SendAsync(ctx context.Context, message mqi.Message, callback func(context.Context, string, error)) error {
	if !p.started {
		return errors.New("producer not started")
	}

	// 创建一个新的上下文，以避免上下文过期
	ctxCopy := context.Background()
	if deadline, ok := ctx.Deadline(); ok {
		var cancel context.CancelFunc
		ctxCopy, cancel = context.WithDeadline(context.Background(), deadline)
		defer cancel()
	}

	// 使用 goroutine 实现异步发送
	go func(ctx context.Context, msg mqi.Message) {
		// 调用同步发送方法
		msgId, err := p.SendSync(ctx, msg)
		if err != nil {
			p.logger.Error("async message send failed",
				"topic", message.GetTopic(),
				"tags", message.GetTags(),
				"keys", message.GetKeys(),
				"error", err)
			if callback != nil {
				callback(ctx, "", err)
			}
			return
		}

		p.logger.Debug("message sent asynchronously",
			"topic", message.GetTopic(),
			"tags", message.GetTags(),
			"keys", message.GetKeys(),
			"messageId", msgId)
		if callback != nil {
			callback(ctx, msgId, nil)
		}
	}(ctxCopy, message)

	return nil
}

// // SendOneWay 单向发送消息（不关心结果）
// func (p *RocketMQProducer) SendOneWay(ctx context.Context, message mqi.Message) error {
// 	if !p.started {
// 		return errors.New("producer not started")
// 	}

// 	msg := ToRocketMQMessage(message)
// 	err := p.producer.SendOneWay(ctx, msg)
// 	if err != nil {
// 		p.logger.Error("failed to send message one-way",
// 			"topic", message.GetTopic(),
// 			"tags", message.GetTags(),
// 			"keys", message.GetKeys(),
// 			"error", err)
// 		return err
// 	}

// 	p.logger.Debug("message sent one-way",
// 		"topic", message.GetTopic(),
// 		"tags", message.GetTags(),
// 		"keys", message.GetKeys())
// 	return nil
// }

// SendOneWay 单向发送消息（不关心结果） (V5版本)
// 注意：V5 版本没有直接的单向发送 API，我们使用同步发送并忽略结果来实现
func (p *RocketMQProducer) SendOneWay(ctx context.Context, message mqi.Message) error {
	if !p.started {
		return errors.New("producer not started")
	}

	// 创建一个新的上下文，以避免上下文过期
	ctxCopy := context.Background()
	if deadline, ok := ctx.Deadline(); ok {
		var cancel context.CancelFunc
		ctxCopy, cancel = context.WithDeadline(context.Background(), deadline)
		defer cancel()
	}

	// 使用 goroutine 实现单向发送，不关心结果
	go func(ctx context.Context, msg mqi.Message) {
		// 调用同步发送方法，但忽略结果
		_, err := p.SendSync(ctx, msg)
		if err != nil {
			p.logger.Error("failed to send message one-way",
				"topic", message.GetTopic(),
				"tags", message.GetTags(),
				"keys", message.GetKeys(),
				"error", err)
			// 单向发送模式下忽略错误
			return
		}

		p.logger.Debug("message sent one-way",
			"topic", message.GetTopic(),
			"tags", message.GetTags(),
			"keys", message.GetKeys())
	}(ctxCopy, message)

	return nil
}

// // ToRocketMQMessage 将通用消息转换为RocketMQ消息
// func ToRocketMQMessage(message mqi.Message) *primitive.Message {
// 	msg := primitive.NewMessage(message.GetTopic(), message.GetBody())

// 	// 设置标签
// 	if message.GetTags() != "" {
// 		msg.WithTag(message.GetTags())
// 	}

// 	// 设置键
// 	if message.GetKeys() != "" {
// 		msg.WithKeys([]string{message.GetKeys()})
// 	}

// 	// 添加特殊属性，将语言标记为JAVA，解决服务器无法识别GO语言代码的问题
// 	msg.WithProperty("LANGUAGE", "JAVA")

// 	// 添加UserId作为属性
// 	if message.GetUserId() != "" {
// 		msg.WithProperty("userId", message.GetUserId())
// 	}

// 	// 添加BusinessKey作为属性
// 	if message.GetBusinessKey() != "" {
// 		msg.WithProperty("businessKey", message.GetBusinessKey())
// 	}

// 	// 设置属性
// 	if message.GetProperties() != nil {
// 		for k, v := range message.GetProperties() {
// 			msg.WithProperty(k, v)
// 		}
// 	}

// 	return msg
// }
