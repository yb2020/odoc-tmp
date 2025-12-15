package proto

import "time"

// SimpleResponse 简单响应结构
type SimpleResponse struct {
	Result string `json:"result"`
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Error       string `json:"error"`
	ErrorCode   string `json:"error_code"`
	Description string `json:"description,omitempty"`
}

// DifyConfig Dify客户端配置
type DifyConfig struct {
	BaseURL string
	APIKey  string
	Timeout time.Duration
}

// NewDifyConfig 创建新的Dify配置
func NewDifyConfig(baseURL, apiKey string) *DifyConfig {
	return &DifyConfig{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Timeout: 30 * time.Second, // 默认超时时间
	}
}

// WithTimeout 设置超时时间
func (c *DifyConfig) WithTimeout(timeout time.Duration) *DifyConfig {
	c.Timeout = timeout
	return c
}
