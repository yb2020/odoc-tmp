package openai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/yb2020/odoc/pkg/http_client"
	"github.com/yb2020/odoc/pkg/logging"
)

// ProviderType 表示不同的 LLM 提供商
type ProviderType string

const (
	ProviderOpenAI   ProviderType = "openai"
	ProviderDeepSeek ProviderType = "deepseek"
	ProviderStepFun  ProviderType = "stepfun"
	ProviderKimi     ProviderType = "kimi"
	// 可以添加更多提供商
)

// Role 表示消息角色
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleFunction  Role = "function"
)

// Message 表示对话消息
type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

// ChatCompletionRequest 表示聊天补全请求
type ChatCompletionRequest struct {
	Model            string     `json:"model"`
	Messages         []Message  `json:"messages"`
	Temperature      float64    `json:"temperature,omitempty"`
	TopP             float64    `json:"top_p,omitempty"`
	MaxTokens        int        `json:"max_tokens,omitempty"`
	Stream           bool       `json:"stream,omitempty"`
	PresencePenalty  float64    `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64    `json:"frequency_penalty,omitempty"`
	User             string     `json:"user,omitempty"`
	Stop             []string   `json:"stop,omitempty"`
	Functions        []Function `json:"functions,omitempty"`
}

// Function 表示函数调用定义
type Function struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Parameters  any    `json:"parameters,omitempty"`
}

// ChatCompletionChoice 表示聊天补全选择
type ChatCompletionChoice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage 表示令牌使用情况
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletionResponse 表示聊天补全响应
type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   Usage                  `json:"usage"`
}

// ChatCompletionStreamChoice 表示流式聊天补全选择
type ChatCompletionStreamChoice struct {
	Index        int     `json:"index"`
	Delta        Message `json:"delta"`
	FinishReason string  `json:"finish_reason"`
}

// ChatCompletionStreamResponse 表示流式聊天补全响应
type ChatCompletionStreamResponse struct {
	ID      string                       `json:"id"`
	Object  string                       `json:"object"`
	Created int64                        `json:"created"`
	Model   string                       `json:"model"`
	Choices []ChatCompletionStreamChoice `json:"choices"`
}

// StreamHandler 表示流式响应的处理函数
type StreamHandler func(response ChatCompletionStreamResponse) error

// ClientOption 表示客户端选项
type ClientOption func(*ClientConfig)

// ClientConfig 表示客户端配置
type ClientConfig struct {
	APIKey     string
	BaseURL    string
	Logger     logging.Logger
	HTTPClient http_client.HttpClient
	Timeout    time.Duration
}

// WithLogger 设置日志记录器
func WithLogger(logger logging.Logger) ClientOption {
	return func(c *ClientConfig) {
		c.Logger = logger
	}
}

// WithHTTPClient 设置 HTTP 客户端
func WithHTTPClient(client http_client.HttpClient) ClientOption {
	return func(c *ClientConfig) {
		c.HTTPClient = client
	}
}

// WithAPIKey 设置 API 密钥
func WithAPIKey(apiKey string) ClientOption {
	return func(c *ClientConfig) {
		c.APIKey = apiKey
	}
}

// WithBaseURL 设置基础 URL
func WithBaseURL(baseURL string) ClientOption {
	return func(c *ClientConfig) {
		c.BaseURL = baseURL
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *ClientConfig) {
		c.Timeout = timeout
	}
}

// LLMClient 表示 LLM 客户端接口
type LLMClient interface {
	// CreateChatCompletion 同步调用 LLM 模型
	CreateChatCompletion(ctx context.Context, request ChatCompletionRequest) (ChatCompletionResponse, error)

	// CreateChatCompletionStream 流式调用 LLM 模型
	CreateChatCompletionStream(ctx context.Context, request ChatCompletionRequest, handler StreamHandler) error

	// CollectStreamResponse 收集流式响应并返回完整结果
	CollectStreamResponse(ctx context.Context, request ChatCompletionRequest) (ChatCompletionResponse, error)
}

// Client 表示通用 LLM 客户端
type Client struct {
	provider ProviderType
	config   ClientConfig
	client   LLMClient
}

// NewClient 创建新的 LLM 客户端
func NewClient(provider ProviderType, options ...ClientOption) (*Client, error) {
	config := ClientConfig{
		Timeout: 30 * time.Second,
	}

	for _, option := range options {
		option(&config)
	}

	if config.APIKey == "" {
		return nil, errors.New("API key is required")
	}

	if config.HTTPClient == nil {
		if config.Logger == nil {
			return nil, errors.New("either HTTPClient or Logger must be provided")
		}
		config.HTTPClient = http_client.NewHttpClient(
			config.Logger,
			http_client.WithTimeout(config.Timeout),
		)
	}

	var client LLMClient
	var err error

	switch provider {
	case ProviderOpenAI:
		client, err = newOpenAIClient(config)
	case ProviderDeepSeek:
		client, err = newDeepSeekClient(config)
	case ProviderStepFun:
		client, err = newStepFunClient(config)
	case ProviderKimi:
		client, err = newKimiClient(config)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	if err != nil {
		return nil, err
	}

	return &Client{
		provider: provider,
		config:   config,
		client:   client,
	}, nil
}

// CreateChatCompletion 同步调用 LLM 模型
func (c *Client) CreateChatCompletion(ctx context.Context, request ChatCompletionRequest) (ChatCompletionResponse, error) {
	return c.client.CreateChatCompletion(ctx, request)
}

// CreateChatCompletionStream 流式调用 LLM 模型
func (c *Client) CreateChatCompletionStream(ctx context.Context, request ChatCompletionRequest, handler StreamHandler) error {
	return c.client.CreateChatCompletionStream(ctx, request, handler)
}

// CollectStreamResponse 收集流式响应并返回完整结果
func (c *Client) CollectStreamResponse(ctx context.Context, request ChatCompletionRequest) (ChatCompletionResponse, error) {
	return c.client.CollectStreamResponse(ctx, request)
}

// OpenAIClient 表示 OpenAI 客户端
type OpenAIClient struct {
	config ClientConfig
	client http_client.HttpClient
}

func newOpenAIClient(config ClientConfig) (LLMClient, error) {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.openai.com/v1"
	}

	return &OpenAIClient{
		config: config,
		client: config.HTTPClient,
	}, nil
}

func (c *OpenAIClient) CreateChatCompletion(ctx context.Context, request ChatCompletionRequest) (ChatCompletionResponse, error) {
	// 将请求体序列化为JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("marshal request: %w", err)
	}

	// 设置请求头
	headers := map[string]string{
		"Content-Type":  "application/json",
		"api-key":       c.config.APIKey,
		"Authorization": c.config.APIKey,
	}

	// 发送请求
	responseData, err := c.client.Post(c.config.BaseURL, requestBody, headers)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("send request: %w", err)
	}

	// 解析响应
	var response ChatCompletionResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("decode response: %w", err)
	}

	return response, nil
}

func (c *OpenAIClient) CreateChatCompletionStream(ctx context.Context, request ChatCompletionRequest, handler StreamHandler) error {
	// 确保请求是流式的
	request.Stream = true

	// 将请求体序列化为JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	// 设置请求头
	headers := map[string]string{
		"Content-Type":  "application/json",
		"api-key":       c.config.APIKey,
		"Authorization": c.config.APIKey,
		"Accept":        "text/event-stream",
	}

	// 创建一个缓冲区用于存储未完整处理的数据
	var buffer string

	// 发送流式请求
	return c.client.PostStream(c.config.BaseURL, requestBody, headers, func(chunk []byte) error {
		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// 将新接收的数据添加到缓冲区
		buffer += string(chunk)

		// 按行处理缓冲区
		for {
			// 查找换行符
			idx := strings.Index(buffer, "\n")
			if idx == -1 {
				break // 没有完整的行，等待更多数据
			}

			// 提取一行
			line := buffer[:idx]
			buffer = buffer[idx+1:] // 移除已处理的行

			// 去除前后空白
			line = strings.TrimSpace(line)

			// 跳过空行和注释
			if line == "" || strings.HasPrefix(line, ":") {
				continue
			}

			// 检查是否为数据行
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			// 提取数据部分
			data := strings.TrimPrefix(line, "data: ")

			// 检查是否为流结束标记
			if data == "[DONE]" {
				return nil
			}

			// 解析JSON数据
			var response ChatCompletionStreamResponse
			if err := json.Unmarshal([]byte(data), &response); err != nil {
				return fmt.Errorf("unmarshal response: %w", err)
			}

			// 调用处理函数
			if err := handler(response); err != nil {
				return err
			}
		}

		return nil
	})
}

func (c *OpenAIClient) CollectStreamResponse(ctx context.Context, request ChatCompletionRequest) (ChatCompletionResponse, error) {
	var fullResponse ChatCompletionResponse
	var mu sync.Mutex
	var content string

	// 强制设置为流式请求
	streamRequest := request
	streamRequest.Stream = true

	err := c.CreateChatCompletionStream(ctx, streamRequest, func(response ChatCompletionStreamResponse) error {
		mu.Lock()
		defer mu.Unlock()

		if len(response.Choices) > 0 {
			content += response.Choices[0].Delta.Content

			// 保存第一个响应的元数据
			if fullResponse.ID == "" {
				fullResponse.ID = response.ID
				fullResponse.Object = response.Object
				fullResponse.Created = response.Created
				fullResponse.Model = response.Model
			}

			// 检查是否完成
			if response.Choices[0].FinishReason != "" {
				fullResponse.Choices = []ChatCompletionChoice{
					{
						Index: 0,
						Message: Message{
							Role:    RoleAssistant,
							Content: content,
						},
						FinishReason: response.Choices[0].FinishReason,
					},
				}
			}
		}

		return nil
	})

	if err != nil {
		return ChatCompletionResponse{}, err
	}

	return fullResponse, nil
}

// DeepSeekClient 表示 DeepSeek 客户端
type DeepSeekClient struct {
	config ClientConfig
	client http_client.HttpClient
}

func newDeepSeekClient(config ClientConfig) (LLMClient, error) {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.deepseek.com/v1"
	}

	return &DeepSeekClient{
		config: config,
		client: config.HTTPClient,
	}, nil
}

func (c *DeepSeekClient) CreateChatCompletion(ctx context.Context, request ChatCompletionRequest) (ChatCompletionResponse, error) {
	// 简化实现，实际项目中需要根据DeepSeek API文档进行完整实现
	return ChatCompletionResponse{}, errors.New("DeepSeek API not fully implemented")
}

func (c *DeepSeekClient) CreateChatCompletionStream(ctx context.Context, request ChatCompletionRequest, handler StreamHandler) error {
	// 确保请求是流式的
	request.Stream = true

	// 将请求体序列化为JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	// 设置请求头 - DeepSeek 使用相同的 Bearer 认证方式
	headers := map[string]string{
		"Content-Type":  "application/json",
		"api-key":       c.config.APIKey,
		"Authorization": c.config.APIKey,
		"Accept":        "text/event-stream",
	}

	// 创建一个缓冲区用于存储未完整处理的数据
	var buffer string

	// 发送流式请求
	return c.client.PostStream(c.config.BaseURL, requestBody, headers, func(chunk []byte) error {
		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// 将新接收的数据添加到缓冲区
		buffer += string(chunk)

		// 按行处理缓冲区
		for {
			// 查找换行符
			idx := strings.Index(buffer, "\n")
			if idx == -1 {
				break // 没有完整的行，等待更多数据
			}

			// 提取一行
			line := buffer[:idx]
			buffer = buffer[idx+1:] // 移除已处理的行

			// 去除前后空白
			line = strings.TrimSpace(line)

			// 跳过空行和注释
			if line == "" || strings.HasPrefix(line, ":") {
				continue
			}

			// 检查是否为数据行
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			// 提取数据部分
			data := strings.TrimPrefix(line, "data: ")

			// 检查是否为流结束标记
			if data == "[DONE]" {
				return nil
			}

			// 解析JSON数据
			var response ChatCompletionStreamResponse
			if err := json.Unmarshal([]byte(data), &response); err != nil {
				return fmt.Errorf("unmarshal response: %w", err)
			}

			// 调用处理函数
			if err := handler(response); err != nil {
				return err
			}
		}

		return nil
	})
}

func (c *DeepSeekClient) CollectStreamResponse(ctx context.Context, request ChatCompletionRequest) (ChatCompletionResponse, error) {
	var fullResponse ChatCompletionResponse
	var fullContent strings.Builder
	var id string
	var model string
	var created int64

	// 调用流式API并收集响应
	err := c.CreateChatCompletionStream(ctx, request, func(response ChatCompletionStreamResponse) error {
		// 记录第一个响应的基本信息
		if id == "" && response.ID != "" {
			id = response.ID
		}
		if model == "" && response.Model != "" {
			model = response.Model
		}
		if created == 0 && response.Created != 0 {
			created = response.Created
		}

		// 收集内容
		if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
			fullContent.WriteString(response.Choices[0].Delta.Content)
		}

		return nil
	})

	if err != nil {
		return ChatCompletionResponse{}, err
	}

	// 构建完整响应
	fullResponse.ID = id
	fullResponse.Model = model
	fullResponse.Created = created
	fullResponse.Choices = []ChatCompletionChoice{
		{
			Message: Message{
				Role:    RoleAssistant,
				Content: fullContent.String(),
			},
			FinishReason: "stop",
		},
	}

	return fullResponse, nil
}

// StepFunClient 表示 StepFun 客户端
type StepFunClient struct {
	config ClientConfig
	client http_client.HttpClient
}

func newStepFunClient(config ClientConfig) (LLMClient, error) {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.stepfun.com/v1"
	}

	return &StepFunClient{
		config: config,
		client: config.HTTPClient,
	}, nil
}

func (c *StepFunClient) CreateChatCompletion(ctx context.Context, request ChatCompletionRequest) (ChatCompletionResponse, error) {
	// 简化实现，实际项目中需要根据StepFun API文档进行完整实现
	return ChatCompletionResponse{}, errors.New("StepFun API not fully implemented")
}

func (c *StepFunClient) CreateChatCompletionStream(ctx context.Context, request ChatCompletionRequest, handler StreamHandler) error {
	// 简化实现，实际项目中需要根据StepFun API文档进行完整实现
	return errors.New("StepFun streaming not fully implemented")
}

func (c *StepFunClient) CollectStreamResponse(ctx context.Context, request ChatCompletionRequest) (ChatCompletionResponse, error) {
	// 简化实现，实际项目中需要根据StepFun API文档进行完整实现
	return ChatCompletionResponse{}, errors.New("StepFun streaming not fully implemented")
}

// KimiClient 表示 Kimi 客户端
type KimiClient struct {
	config ClientConfig
	client http_client.HttpClient
}

func newKimiClient(config ClientConfig) (LLMClient, error) {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.kimi.ai/v1"
	}

	return &KimiClient{
		config: config,
		client: config.HTTPClient,
	}, nil
}

func (c *KimiClient) CreateChatCompletion(ctx context.Context, request ChatCompletionRequest) (ChatCompletionResponse, error) {
	// 简化实现，实际项目中需要根据Kimi API文档进行完整实现
	return ChatCompletionResponse{}, errors.New("Kimi API not fully implemented")
}

func (c *KimiClient) CreateChatCompletionStream(ctx context.Context, request ChatCompletionRequest, handler StreamHandler) error {
	// 简化实现，实际项目中需要根据Kimi API文档进行完整实现
	return errors.New("Kimi streaming not fully implemented")
}

func (c *KimiClient) CollectStreamResponse(ctx context.Context, request ChatCompletionRequest) (ChatCompletionResponse, error) {
	// 简化实现，实际项目中需要根据Kimi API文档进行完整实现
	return ChatCompletionResponse{}, errors.New("Kimi streaming not fully implemented")
}
