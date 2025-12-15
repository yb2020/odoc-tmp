package proto

// ConversationMessagesResponse 对话消息响应
type ConversationMessagesResponse struct {
	Data    []Message `json:"data"`
	HasMore bool      `json:"has_more"`
	Limit   int32     `json:"limit"`
}

// ConversationsResponse 对话列表响应
type ConversationsResponse struct {
	Data []Conversation `json:"data"`
}

// Conversation 对话
type Conversation struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Inputs    Config `json:"inputs"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
}

// Message 消息
type Message struct {
	ID                string               `json:"id"`
	ConversationID    string               `json:"conversation_id"`
	Inputs            Config               `json:"inputs"`
	Query             string               `json:"query"`
	Answer            string               `json:"answer"`
	MessageFiles      []MessageFile        `json:"message_files"`
	Feedback          Feedback             `json:"feedback"`
	CreatedAt         int64                `json:"created_at"`
	RetriverResources []RetrieverResources `json:"retriever_resources"`
}

type MessageFile struct {
	Id string `json:"id"`
	// (string) 文件类型，image 图片
	Type string `json:"type"`
	// (string) 预览图片地址
	Url string `json:"url"`
	// (string) 文件归属方，user 或 assistant
	BelongsTo string `json:"belongs_to"`
}

// Feedback 反馈
type Feedback struct {
	Rating   string `json:"rating"`
	Content  string `json:"content"`
	From     string `json:"from"`
	FromUser string `json:"from_user"`
}
