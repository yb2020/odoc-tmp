package rocketmq

import (
	v5 "github.com/apache/rocketmq-clients/golang/v5"
	mq_interface "github.com/yb2020/odoc/pkg/mq/interface"
)

// Message 实现 mq_interface.Message 接口
type Message struct {
	Topic       string
	Body        []byte
	Tags        string
	Keys        string
	UserId      string
	BusinessKey string
	MessageId   string
	Properties  map[string]string
}

// GetTopic 获取消息主题
func (m *Message) GetTopic() string {
	return m.Topic
}

// GetBody 获取消息体
func (m *Message) GetBody() []byte {
	return m.Body
}

// GetTags 获取消息标签
func (m *Message) GetTags() string {
	return m.Tags
}

// GetUserId 获取消息用户ID
func (m *Message) GetUserId() string {
	return m.UserId
}

// GetProperties 获取消息属性
func (m *Message) GetProperties() map[string]string {
	return m.Properties
}

// GetMessageId 获取消息ID
func (m *Message) GetMessageId() string {
	return m.MessageId
}

// GetKeys 获取消息键
func (m *Message) GetKeys() string {
	return m.Keys
}

// GetBusinessKey 获取消息业务键
func (m *Message) GetBusinessKey() string {
	return m.BusinessKey
}

// FromRocketMQMessage 从RocketMQ V5消息转换为通用消息
func FromRocketMQMessage(msg *v5.MessageView) mq_interface.Message {
	// 获取所有属性
	properties := msg.GetProperties()
	
	// 获取标签，注意V5中GetTag返回的是*string
	var tag string
	if msg.GetTag() != nil {
		tag = *msg.GetTag()
	}

	// 获取消息ID
	messageId := msg.GetMessageId()

	// 获取消息键，注意V5中GetKeys返回的是[]string
	var key string
	keys := msg.GetKeys()
	if len(keys) > 0 {
		key = keys[0]
	}

	// 创建并返回消息
	return &Message{
		Topic:       msg.GetTopic(),
		Body:        msg.GetBody(),
		Tags:        tag,
		Keys:        key,
		UserId:      properties["userId"],      // 从属性中获取UserId（如果有）
		BusinessKey: properties["businessKey"], // 从属性中获取BusinessKey（如果有）
		MessageId:   messageId,
		Properties:  properties,
	}
}

// ToRocketMQMessage 将通用消息转换为RocketMQ V5消息
func ToRocketMQMessage(message mq_interface.Message) *v5.Message {
	// 创建V5消息
	msg := &v5.Message{
		Topic: message.GetTopic(),
		Body:  message.GetBody(),
	}

	// 设置标签
	if message.GetTags() != "" {
		msg.SetTag(message.GetTags())
	}

	// 设置键（注意V5中SetKeys接受单个字符串，而不是字符串切片）
	if message.GetKeys() != "" {
		msg.SetKeys(message.GetKeys())
	}

	// 设置用户ID作为属性（如果有）
	if message.GetUserId() != "" {
		msg.AddProperty("userId", message.GetUserId())
	}

	// 设置业务键作为属性（如果有）
	if message.GetBusinessKey() != "" {
		msg.AddProperty("businessKey", message.GetBusinessKey())
	}

	// 设置其他属性（如果有）
	if props := message.GetProperties(); props != nil {
		for k, v := range props {
			msg.AddProperty(k, v)
		}
	}

	return msg
}
