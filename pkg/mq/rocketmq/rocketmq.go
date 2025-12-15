package rocketmq

import (
	"os"
	"strconv"

	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/logging"
)

// InitRocketMQEnvironment 初始化RocketMQ环境
func InitRocketMQEnvironment(config *config.Config, logger logging.Logger) {
	// 设置系统属性
	// 设置环境变量 - RocketMQ V5 客户端使用不同的环境变量
	// 设置日志级别，默认为ERROR，减少冗余日志
	os.Setenv("rocketmq.log.level", config.RocketMQ.Client.LogLevel)

	// 设置请求超时时间，默认30秒
	os.Setenv("rocketmq.client.timeout.request", strconv.Itoa(config.RocketMQ.Client.RequestTimeout))

	// 设置重试次数，默认3次
	os.Setenv("rocketmq.client.retry.times", strconv.Itoa(config.RocketMQ.Client.RetryTimes))

	// 设置重试间隔，默认3秒
	os.Setenv("rocketmq.client.retry.interval", strconv.Itoa(config.RocketMQ.Client.RetryInterval))

	// 设置最大重试间隔，默认15秒
	os.Setenv("rocketmq.client.retry.max-interval", strconv.Itoa(config.RocketMQ.Client.RetryMaxInterval))

	// 设置最大消息大小，默认4MB
	os.Setenv("rocketmq.client.max.message.size", strconv.Itoa(config.RocketMQ.Client.MaxMessageSize))

	// 禁用统计信息日志
	os.Setenv("rocketmq.client.log.stats.enable", "false")

	// 设置V5客户端的日志配置
	// V5客户端中日志设置主要通过环境变量控制
	// 设置日志路径为空字符串，让日志输出到标准输出，便于查看错误
	// 由于我们已经设置了日志级别为ERROR，所以只会输出错误日志
	os.Setenv("rocketmq.log.path", "")

	// 设置日志输出到控制台，便于查看错误
	os.Setenv("mq.consoleAppender.enabled", "true")

	// 禁用文件日志输出
	os.Setenv("mq.fileAppender.enabled", "false")
}
