package proto

// EventType 事件类型
type EventType string

const (
	// 通用事件
	// EventMessage 消息事件 - LLM 返回文本块事件
	EventMessage EventType = "message"
	// EventMessageEnd 消息结束事件
	EventMessageEnd EventType = "message_end"
	// EventMessageReplace 消息内容替换事件
	EventMessageReplace EventType = "message_replace"
	// EventTTSMessage TTS 音频流事件
	EventTTSMessage EventType = "tts_message"
	// EventTTSMessageEnd TTS 音频流结束事件
	EventTTSMessageEnd EventType = "tts_message_end"
	// EventError 错误事件
	EventError EventType = "error"
	// EventPing ping事件 - 心跳事件
	EventPing EventType = "ping"

	// Agent 相关事件
	// EventAgentMessage Agent模式下返回文本块事件
	EventAgentMessage EventType = "agent_message"
	// EventAgentThought Agent模式下思考步骤事件
	EventAgentThought EventType = "agent_thought"
	// EventAgentAction 代理动作事件
	EventAgentAction EventType = "agent_action"
	// EventAgentActionResult 代理动作结果事件
	EventAgentActionResult EventType = "agent_action_result"
	// EventMessageFile 文件事件
	EventMessageFile EventType = "message_file"

	// Workflow 相关事件
	// EventWorkflowStarted workflow 开始执行
	EventWorkflowStarted EventType = "workflow_started"
	// EventNodeStarted node 开始执行
	EventNodeStarted EventType = "node_started"
	// EventNodeFinished node 执行结束
	EventNodeFinished EventType = "node_finished"
	// EventWorkflowFinished workflow 执行结束
	EventWorkflowFinished EventType = "workflow_finished"
	// EventWorkflowTextChunk workflow llm模型输入结果
	EventWorkflowTextChunk EventType = "text_chunk"

	// 兼容旧版本的事件类型
	// EventStart 开始事件
	EventStart EventType = "start"
	// EventFinish 完成事件
	EventFinish EventType = "finish"
	// EventEnd 结束事件
	EventEnd EventType = "end"
)

// BaseEvent 基础事件
type BaseEvent struct {
	Event EventType `json:"event"`
	Data  string    `json:"data"`
}

// ChatMessageEvent 聊天消息事件
type ChatMessageEvent struct {
	Event          EventType     `json:"event"`
	TaskID         string        `json:"task_id"`
	WorkflowRunID  string        `json:"workflow_run_id"`
	ConversationID string        `json:"conversation_id"`
	Answer         string        `json:"answer"`
	MessageID      string        `json:"message_id"`
	Metadata       ChunkMetadata `json:"metadata"`
	Audio          string        `json:"audio"` // 当是tts_message_end时，会返回音频
	CreateAt       int64         `json:"create_at"`
	Data           struct {
		ID                string                 `json:"id"`
		WorkflowID        string                 `json:"workflow_id"`
		NodeID            string                 `json:"node_id"`
		NodeType          string                 `json:"node_type"`
		Index             int                    `json:"index"`
		PredecessorNodeID string                 `json:"predecessor_node_id"`
		Title             string                 `json:"title"`
		Inputs            map[string]interface{} `json:"inputs"`
		Outputs           map[string]interface{} `json:"outputs"`
		Status            string                 `json:"status"` //running / succeeded / failed / stopped
		ElapsedTime       float64                `json:"elapsed_time"`
		TotalTokens       int                    `json:"total_tokens"`
		TotalPrice        string                 `json:"total_price"`
		TotalSteps        int                    `json:"total_steps"`
		Currency          string                 `json:"currency"`
		ExecutionMetadata map[string]interface{} `json:"execution_metadata"`
		SequenceNumber    int                    `json:"sequence_number"`
		CreatedAt         int64                  `json:"created_at"`
		FinishedAt        int64                  `json:"finished_at"`
		Error             string                 `json:"error"`
	} `json:"data"`
}

// 当是message_end时，会返回这些信息
type ChunkMetadata struct {
	Usage              usage                `json:"usage"`
	RetrieverResources []RetrieverResources `json:"retriever_resources"`
}

type usage struct {
	PromptTokens        int     `json:"prompt_tokens"`
	PromptUnitPrice     string  `json:"prompt_unit_price"`
	PromptPriceUnit     string  `json:"prompt_price_unit"`
	PromptPrice         string  `json:"prompt_price"`
	CompletionTokens    int     `json:"completion_tokens"`
	CompletionUnitPrice string  `json:"completion_unit_price"`
	CompletionPriceUnit string  `json:"completion_price_unit"`
	CompletionPrice     string  `json:"completion_price"`
	TotalTokens         int     `json:"total_tokens"`
	TotalPrice          string  `json:"total_price"`
	Currency            string  `json:"currency"`
	Latency             float64 `json:"latency"`
}

type RetrieverResources struct {
	Position     int     `json:"position"`
	DatasetID    string  `json:"dataset_id"`
	DatasetName  string  `json:"dataset_name"`
	DocumentID   string  `json:"document_id"`
	DocumentName string  `json:"document_name"`
	SegmentID    string  `json:"segment_id"`
	Score        float64 `json:"score"`
	Content      string  `json:"content"`
}

// ChatErrorEvent 聊天错误事件
type ChatErrorEvent struct {
	Event EventType `json:"event"`
	Data  struct {
		Error string `json:"error"`
	} `json:"data"`
}

// AgentThoughtEvent 代理思考事件
type AgentThoughtEvent struct {
	Event EventType `json:"event"`
	Data  struct {
		AgentName string `json:"agent_name"`
		Thought   string `json:"thought"`
	} `json:"data"`
}

// AgentActionEvent 代理动作事件
type AgentActionEvent struct {
	Event EventType `json:"event"`
	Data  struct {
		AgentName   string                 `json:"agent_name"`
		Action      string                 `json:"action"`
		ActionInput map[string]interface{} `json:"action_input"`
	} `json:"data"`
}

// AgentActionResultEvent 代理动作结果事件
type AgentActionResultEvent struct {
	Event EventType `json:"event"`
	Data  struct {
		AgentName    string `json:"agent_name"`
		ActionResult string `json:"action_result"`
	} `json:"data"`
}

// ChatStartEvent 聊天开始事件
type ChatStartEvent struct {
	Event EventType `json:"event"`
	Data  struct {
		TaskID string `json:"task_id"`
	} `json:"data"`
}

// ChatFinishEvent 聊天完成事件
type ChatFinishEvent struct {
	Event EventType `json:"event"`
	Data  struct {
		MessageID      string `json:"message_id"`
		ConversationID string `json:"conversation_id"`
	} `json:"data"`
}

// MessageReplaceEvent 消息内容替换事件
type MessageReplaceEvent struct {
	Event EventType `json:"event"`
	Data  struct {
		Message string `json:"message"`
	} `json:"data"`
}

// TTSMessageEvent TTS 音频流事件
type TTSMessageEvent struct {
	Event EventType `json:"event"`
	Data  struct {
		Audio string `json:"audio"`
	} `json:"data"`
}

// TTSMessageEndEvent TTS 音频流结束事件
type TTSMessageEndEvent struct {
	Event EventType `json:"event"`
	Data  struct {
		TaskID string `json:"task_id"`
	} `json:"data"`
}

// MessageFileEvent 文件事件
type MessageFileEvent struct {
	Event EventType `json:"event"`
	Data  struct {
		FileID   string `json:"file_id"`
		FileName string `json:"file_name"`
		FileURL  string `json:"file_url"`
	} `json:"data"`
}

// WorkflowStartedEvent workflow 开始执行事件
type WorkflowStartedEvent struct {
	Event EventType `json:"event"`
	Data  struct {
		WorkflowID string `json:"workflow_id"`
		TaskID     string `json:"task_id"`
	} `json:"data"`
}

// NodeStartedEvent node 开始执行事件
type NodeStartedEvent struct {
	Event EventType `json:"event"`
	Data  struct {
		NodeID   string `json:"node_id"`
		NodeName string `json:"node_name"`
	} `json:"data"`
}

// NodeFinishedEvent node 执行结束事件
type NodeFinishedEvent struct {
	Event EventType `json:"event"`
	Data  struct {
		NodeID   string `json:"node_id"`
		NodeName string `json:"node_name"`
		Output   string `json:"output"`
	} `json:"data"`
}

// WorkflowFinishedEvent workflow 执行结束事件
type WorkflowFinishedEvent struct {
	Event EventType `json:"event"`
	Data  struct {
		WorkflowID string `json:"workflow_id"`
		TaskID     string `json:"task_id"`
	} `json:"data"`
}

// WorkflowTextChunkEvent workflow llm模型输入结果事件
type WorkflowTextChunkEvent struct {
	Event EventType `json:"event"`
	Data  struct {
		Text string `json:"text"`
	} `json:"data"`
}
