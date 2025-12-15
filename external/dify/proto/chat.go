package proto

import "time"

// ChatMessage 聊天消息请求
type ChatMessage struct {
	Query          string                 `json:"query,omitempty"`
	User           string                 `json:"user"`
	Inputs         map[string]interface{} `json:"inputs,omitempty"`
	ResponseMode   string                 `json:"response_mode,omitempty"` // blocking 或 streaming
	ConversationID string                 `json:"conversation_id,omitempty"`
	Files          []*File                `json:"files,omitempty"`
}

type File struct {
	// 文件类型。document: TXT,MD,PDF等; image: JPG,PNG等; audio: MP3,WAV等; video: MP4,MOV等; custom: 其他。
	Type string `json:"type"`
	// 文件传输方式。Available options: remote_url, local_file
	TransferMethod string `json:"transfer_method"`
	//图片地址 (当 transfer_method 为 remote_url 时)。
	Url string `json:"url"`
	// 上传文件ID (当 transfer_method 为 local_file 时)。
	UploadFileId string `json:"upload_file_id"`
}

// NewChatMessage 创建新的聊天消息请求
func NewChatMessage(query, user string) *ChatMessage {
	return &ChatMessage{
		Query: query,
		User:  user,
	}
}

// WithInputs 设置输入参数
func (m *ChatMessage) WithInputs(inputs map[string]interface{}) *ChatMessage {
	m.Inputs = inputs
	return m
}

// WithResponseMode 设置响应模式
func (m *ChatMessage) WithResponseMode(responseMode string) *ChatMessage {
	m.ResponseMode = responseMode
	return m
}

// WithConversationID 设置对话ID
func (m *ChatMessage) WithConversationID(conversationID string) *ChatMessage {
	m.ConversationID = conversationID
	return m
}

// WithFiles 设置文件列表
func (m *ChatMessage) WithFiles(files []*File) *ChatMessage {
	m.Files = files
	return m
}

// ChatMessageResponse 聊天消息响应
type ChatMessageResponse struct {
	ID             string                 `json:"id"`
	Answer         string                 `json:"answer"`
	ConversationID string                 `json:"conversation_id"`
	CreatedAt      time.Time              `json:"created_at"`
	TaskID         string                 `json:"task_id,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	Documents      []ChatDocument         `json:"documents,omitempty"`
}

// ChatDocument 文档信息
type ChatDocument struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Source    string `json:"source"`
	SourceURL string `json:"source_url,omitempty"`
}

// AppParametersResponse 应用参数响应
type AppParametersResponse struct {
	Inputs []InputParameter `json:"inputs"`
}

// InputParameter 输入参数
type InputParameter struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Options     []any  `json:"options,omitempty"`
	Default     any    `json:"default,omitempty"`
	Description string `json:"description,omitempty"`
}

// ConversationResponse 对话响应
type ConversationResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Inputs    any       `json:"inputs"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MessagesResponse 消息列表响应
type MessagesResponse struct {
	Data    []MessageResponse `json:"data"`
	HasMore bool              `json:"has_more"`
	Limit   int               `json:"limit"`
}

// MessageResponse 消息响应
type MessageResponse struct {
	ID             string                 `json:"id"`
	ConversationID string                 `json:"conversation_id"`
	Query          string                 `json:"query"`
	Answer         string                 `json:"answer"`
	Feedback       *FeedbackResponse      `json:"feedback,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// FeedbackResponse 反馈响应
type FeedbackResponse struct {
	Rating  string `json:"rating"`
	Content string `json:"content,omitempty"`
}

// dify返回的推荐问题响应
type SuggestedQuestionsResult struct {
	Result string   `json:"result"`
	Data   []string `json:"data"`
}

// SuggestedQuestionsResponse 推荐问题响应
type SuggestedQuestionsResponse struct {
	Questions []string `json:"questions"`
}

// AudioToTextResponse 语音转文本响应
type AudioToTextResponse struct {
	Text string `json:"text"`
}

// TextToAudioResponse 文本转语音响应
type TextToAudioResponse struct {
	AudioURL string `json:"audio_url"`
}

// AppMetaResponse 应用元数据响应
type AppMetaResponse struct {
	AppMode string `json:"app_mode"`
}
