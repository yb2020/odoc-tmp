package proto

// AppInfoResponse 应用信息响应
type AppInfoResponse struct {
	AppID                 string                `json:"app_id"`
	ModelConfig           ModelConfig           `json:"model_config"`
	UserInputForm         []UserInputFormItem   `json:"user_input_form"`
	ConversationHistories []ConversationHistory `json:"conversation_histories"`
}

// ModelConfig 模型配置
type ModelConfig struct {
	Provider    string `json:"provider"`
	Model       string `json:"model"`
	Configs     Config `json:"configs"`
	Parameters  Config `json:"parameters"`
	PromptType  string `json:"prompt_type"`
	Credentials Config `json:"credentials"`
}

// Config 配置
type Config struct {
	Temperature      float64 `json:"temperature"`
	TopP             float64 `json:"top_p"`
	PresencePenalty  float64 `json:"presence_penalty"`
	FrequencyPenalty float64 `json:"frequency_penalty"`
	MaxTokens        int     `json:"max_tokens"`
}

// UserInputFormItem 用户输入表单项
type UserInputFormItem struct {
	VariableID string `json:"variable_id"`
	Label      string `json:"label"`
	Required   bool   `json:"required"`
	Default    string `json:"default"`
	Type       string `json:"type"`
	Options    []struct {
		Label string `json:"label"`
		Value string `json:"value"`
	} `json:"options"`
}

// ConversationHistory 对话历史
type ConversationHistory struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Inputs    Config `json:"inputs"`
	CreatedAt string `json:"created_at"`
}
