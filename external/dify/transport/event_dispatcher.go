package transport

import (
	"encoding/json"
	"fmt"

	"github.com/yb2020/odoc/external/dify/callback"
	"github.com/yb2020/odoc/external/dify/proto"
	"github.com/yb2020/odoc/pkg/logging"
)

// EventDispatcher 事件分发器
// 负责将事件数据分发到相应的回调函数
type EventDispatcher struct {
	logger logging.Logger
}

// NewEventDispatcher 创建一个新的事件分发器
func NewEventDispatcher(logger logging.Logger) *EventDispatcher {
	return &EventDispatcher{
		logger: logger,
	}
}

// DispatchChatFlowEvent 分发工作流事件
// 根据事件类型将事件数据分发到相应的回调函数
func (d *EventDispatcher) DispatchChatFlowEvent(data []byte, eventType string, callback callback.ChatflowStreamCallback) error {
	switch eventType {
	case string(proto.EventMessage):
		return d.handleChatflowMessage(data, callback)
	case string(proto.EventError):
		return d.handleChatflowError(data, callback)
	case string(proto.EventAgentThought):
		return d.handleAgentThought(data, callback)
	case string(proto.EventAgentAction):
		return d.handleAgentAction(data, callback)
	case string(proto.EventAgentActionResult):
		return d.handleAgentActionResult(data, callback)
	case string(proto.EventStart):
		return d.handleChatflowStart(data, callback)
	case string(proto.EventFinish):
		return d.handleChatflowFinish(data, callback)
	case string(proto.EventEnd):
		return callback.OnEnd()
	case string(proto.EventPing):
		// 忽略ping事件
		return nil
	default:
		d.logger.Warn("msg", "Received unknown event type", "eventType", eventType)
		return nil
	}
}

// handleChatflowMessage 处理工作流编排对话消息事件
func (d *EventDispatcher) handleChatflowMessage(data []byte, callback callback.ChatflowStreamCallback) error {
	var event proto.ChatMessageEvent
	if err := UnmarshalEvent(data, &event); err != nil {
		d.logger.Error("msg", "Failed to unmarshal message event", "error", err)
		return err
	}

	return callback.OnMessage(event)
}

// handleChatflowError 处理工作流编排对话错误事件
func (d *EventDispatcher) handleChatflowError(data []byte, callback callback.ChatflowStreamCallback) error {
	var event proto.ChatErrorEvent
	if err := UnmarshalEvent(data, &event); err != nil {
		d.logger.Error("msg", "Failed to unmarshal error event", "error", err)
		return err
	}

	return callback.OnError(event.Data.Error)
}

// handleAgentThought 处理代理思考事件
func (d *EventDispatcher) handleAgentThought(data []byte, callback callback.ChatflowStreamCallback) error {
	var event proto.AgentThoughtEvent
	if err := UnmarshalEvent(data, &event); err != nil {
		d.logger.Error("msg", "Failed to unmarshal agent thought event", "error", err)
		return err
	}

	return callback.OnAgentThought(event.Data.AgentName, event.Data.Thought)
}

// handleAgentAction 处理代理动作事件
func (d *EventDispatcher) handleAgentAction(data []byte, callback callback.ChatflowStreamCallback) error {
	var event proto.AgentActionEvent
	if err := UnmarshalEvent(data, &event); err != nil {
		d.logger.Error("msg", "Failed to unmarshal agent action event", "error", err)
		return err
	}

	return callback.OnAgentAction(event.Data.AgentName, event.Data.Action, event.Data.ActionInput)
}

// handleAgentActionResult 处理代理动作结果事件
func (d *EventDispatcher) handleAgentActionResult(data []byte, callback callback.ChatflowStreamCallback) error {
	var event proto.AgentActionResultEvent
	if err := UnmarshalEvent(data, &event); err != nil {
		d.logger.Error("msg", "Failed to unmarshal agent action result event", "error", err)
		return err
	}

	return callback.OnAgentActionResult(event.Data.AgentName, event.Data.ActionResult)
}

// handleChatflowStart 处理工作流编排对话开始事件
func (d *EventDispatcher) handleChatflowStart(data []byte, callback callback.ChatflowStreamCallback) error {
	var event proto.ChatStartEvent
	if err := UnmarshalEvent(data, &event); err != nil {
		d.logger.Error("msg", "Failed to unmarshal start event", "error", err)
		return err
	}

	return callback.OnStart(event.Data.TaskID)
}

// handleChatflowFinish 处理工作流编排对话完成事件
func (d *EventDispatcher) handleChatflowFinish(data []byte, callback callback.ChatflowStreamCallback) error {
	var event proto.ChatFinishEvent
	if err := UnmarshalEvent(data, &event); err != nil {
		d.logger.Error("msg", "Failed to unmarshal finish event", "error", err)
		return err
	}

	return callback.OnFinish(event.Data.MessageID, event.Data.ConversationID)
}

// ProcessStreamLine 处理流式响应行
// 解析事件类型并调用事件处理函数
func (d *EventDispatcher) ProcessStreamLine(line string, callback callback.ChatflowStreamCallback) error {
	// 忽略空行和ping事件
	if line == "" {
		return nil
	}

	// 解析事件数据
	var event map[string]interface{}
	if err := json.Unmarshal([]byte(line), &event); err != nil {
		d.logger.Error("msg", "Failed to unmarshal event line", "error", err, "line", line)
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	// 获取事件类型
	eventTypeVal, ok := event["event"]
	if !ok {
		d.logger.Warn("msg", "Event type not found in event data", "data", line)
		return nil
	}

	eventType, ok := eventTypeVal.(string)
	if !ok {
		d.logger.Warn("msg", "Event type is not a string", "eventType", eventTypeVal)
		return nil
	}

	// 分发事件
	return d.DispatchChatFlowEvent([]byte(line), eventType, callback)
}
