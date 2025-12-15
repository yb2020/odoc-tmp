package callback

import "github.com/yb2020/odoc/external/dify/proto"

// ChatStreamCallback 聊天流式回调接口
type ChatStreamCallback interface {
	// OnStart 处理开始事件
	OnStart(event proto.ChatMessageEvent) error
	// OnFinish 处理完成事件
	OnFinish(event proto.ChatMessageEvent) error
	// OnEnd 处理结束事件
	OnEnd() error
	// OnMessage 处理消息事件
	OnMessage(event proto.ChatMessageEvent) error
	// OnError 处理错误事件
	OnError(event proto.ChatMessageEvent) error
	// OnWorkflowStarted 处理 workflow 开始事件
	OnWorkflowStarted(event proto.ChatMessageEvent) error
	// OnWorkflowFinished 处理 workflow 结束事件
	OnWorkflowFinished(event proto.ChatMessageEvent) error
	// OnNodeStarted 处理 node 开始事件
	OnNodeStarted(event proto.ChatMessageEvent) error
	// OnNodeFinished 处理 node 结束事件
	OnNodeFinished(event proto.ChatMessageEvent) error
	// OnPing 处理 ping 事件
	OnPing(event proto.ChatMessageEvent) error
}

// CompletionStreamCallback 完成流式回调接口
type CompletionStreamCallback interface {
	OnError(error string) error

	OnMessage(string) error

	// OnStart 处理开始事件
	OnStart(taskID string) error

	// OnFinish 处理完成事件
	OnFinish(messageID string) error
}

// ChatflowStreamCallback 聊天流程回调接口
type ChatflowStreamCallback interface {
	// 处理消息事件
	OnMessage(event proto.ChatMessageEvent) error
	// 处理错误事件
	OnError(error string) error

	// OnStart 处理开始事件
	OnStart(taskID string) error

	// OnFinish 处理完成事件
	OnFinish(messageID string, conversationID string) error

	// OnAgentThought 处理代理思考事件
	OnAgentThought(agentName string, thought string) error

	// OnAgentAction 处理代理动作事件
	OnAgentAction(agentName string, action string, actionInput map[string]interface{}) error

	// OnAgentActionResult 处理代理动作结果事件
	OnAgentActionResult(agentName string, actionResult string) error

	// OnEnd 处理结束事件
	OnEnd() error
}

// WorkflowStreamCallback 工作流回调接口
type WorkflowStreamCallback interface {
	// OnMessage 处理消息事件
	OnMessage(message string) error
	// OnError 处理错误事件
	OnError(error string) error

	// OnStart 处理开始事件
	OnStart(taskID string) error

	// OnNodeStart 处理节点开始事件
	OnNodeStart(nodeID string, nodeName string, nodeType string) error

	// OnNodeFinish 处理节点完成事件
	OnNodeFinish(nodeID string, nodeName string, nodeType string, outputs map[string]interface{}) error

	// OnNodeError 处理节点错误事件
	OnNodeError(nodeID string, nodeName string, nodeType string, error string) error

	// OnFinish 处理完成事件
	OnFinish(outputs map[string]interface{}) error

	OnEnd() error
}
