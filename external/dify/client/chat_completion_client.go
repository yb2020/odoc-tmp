package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/external/dify/callback"
	"github.com/yb2020/odoc/external/dify/constant"
	"github.com/yb2020/odoc/external/dify/proto"
	"github.com/yb2020/odoc/external/dify/transport"
	"github.com/yb2020/odoc/pkg/logging"
)

// DifyCompletionClient Dify文本生成型应用客户端接口
// 包含文本生成型应用相关的功能
type DifyCompletionClient interface {
	// SendCompletionMessage 发送文本生成请求（阻塞模式）
	SendCompletionMessage(ctx context.Context, request *proto.CompletionRequest) (*proto.CompletionResponse, error)

	// SendCompletionMessageStream 发送文本生成请求（流式模式）
	SendCompletionMessageStream(ctx context.Context, request *proto.CompletionRequest, callback callback.CompletionStreamCallback) error

	// StopCompletion 停止文本生成
	StopCompletion(ctx context.Context, taskID string, user string) (*proto.SimpleResponse, error)

	// TextToAudio 文字转语音
	TextToAudio(ctx context.Context, messageID string, text string, user string) ([]byte, error)
}

// DifyChatCompletionClientImpl Dify聊天客户端实现
type DifyChatCompletionClientImpl struct {
	DifyBaseClientImpl
	logger logging.Logger
	tracer opentracing.Tracer
}

// DifyCompletionClientImpl Dify文本生成型应用客户端实现
type DifyCompletionClientImpl struct {
	DifyBaseClientImpl
	logger logging.Logger
	tracer opentracing.Tracer
}

// NewDifyChatCompletionClient 创建新的Dify聊天客户端
func NewDifyChatCompletionClient(logger logging.Logger, tracer opentracing.Tracer, apiBaseUrl, apiKey string) DifyCompletionClient {
	return &DifyCompletionClientImpl{
		DifyBaseClientImpl: DifyBaseClientImpl{
			apiKey:     apiKey,
			apiBaseUrl: apiBaseUrl,
			httpClient: transport.NewDefaultHttpClient(logger, 5, 120),
			tracer:     tracer,
		},
		logger: logger,
		tracer: tracer,
	}
}

// NewDifyChatCompletionClientWithHttpClient 使用自定义HTTP客户端创建新的Dify聊天客户端
func NewDifyChatCompletionClientWithHttpClient(logger logging.Logger, tracer opentracing.Tracer, apiBaseUrl, apiKey string, httpClient transport.HttpClient) DifyCompletionClient {
	return &DifyCompletionClientImpl{
		DifyBaseClientImpl: DifyBaseClientImpl{
			apiKey:     apiKey,
			apiBaseUrl: apiBaseUrl,
			httpClient: httpClient,
			tracer:     tracer,
		},
		logger: logger,
		tracer: tracer,
	}
}

// NewDifyCompletionClient 创建新的Dify文本生成型应用客户端
func NewDifyCompletionClient(logger logging.Logger, tracer opentracing.Tracer, apiKey, apiBaseUrl string) DifyCompletionClient {
	return &DifyCompletionClientImpl{
		DifyBaseClientImpl: DifyBaseClientImpl{
			apiKey:     apiKey,
			apiBaseUrl: apiBaseUrl,
			httpClient: transport.NewDefaultHttpClient(logger, 5, 120),
			tracer:     tracer,
		},
		logger: logger,
		tracer: tracer,
	}
}

// NewDifyCompletionClientWithHttpClient 使用自定义HTTP客户端创建新的Dify文本生成型应用客户端
func NewDifyCompletionClientWithHttpClient(logger logging.Logger, tracer opentracing.Tracer, apiKey, apiBaseUrl string, httpClient transport.HttpClient) DifyCompletionClient {
	return &DifyCompletionClientImpl{
		DifyBaseClientImpl: DifyBaseClientImpl{
			apiKey:     apiKey,
			apiBaseUrl: apiBaseUrl,
			httpClient: httpClient,
			tracer:     tracer,
		},
		logger: logger,
		tracer: tracer,
	}
}

// SendChatMessage 发送聊天消息
func (c *DifyChatCompletionClientImpl) SendChatMessage(ctx context.Context, message *proto.ChatMessage) (*proto.ChatMessageResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyChatCompletionClientImpl.SendChatMessage")
	defer span.Finish()

	if message == nil {
		c.logger.Error("消息不能为空")
		return nil, fmt.Errorf("消息不能为空")
	}

	// 构建URL
	url := transport.BuildUrl(c.apiBaseUrl, constant.CHAT_MESSAGES_PATH)

	// 将消息转换为JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + c.apiKey
	headers["Content-Type"] = transport.ContentTypeJson

	// 发送请求
	resp, err := c.httpClient.Post(ctx, url, bytes.NewBuffer(jsonData), transport.ContentTypeJson, headers)
	if err != nil {
		return nil, err
	}

	// 处理响应
	respBody, err := transport.HandleResponse(resp)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var response proto.ChatMessageResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// SendCompletionMessage 发送文本生成请求（阻塞模式）
func (c *DifyCompletionClientImpl) SendCompletionMessage(ctx context.Context, request *proto.CompletionRequest) (*proto.CompletionResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DifyCompletionClient.SendCompletionMessage")
	defer span.Finish()

	if request == nil {
		return nil, errors.New("request cannot be nil")
	}

	// 设置阻塞响应模式
	request.ResponseMode = "blocking"

	// 构建URL
	url := fmt.Sprintf("%s%s", c.apiBaseUrl, constant.COMPLETION_MESSAGES_PATH)

	// 将请求转换为JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 设置请求头
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + c.apiKey
	headers["Content-Type"] = transport.ContentTypeJson

	// 发送请求
	resp, err := c.httpClient.Post(ctx, url, bytes.NewBuffer(jsonData), transport.ContentTypeJson, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to send completion message: %w", err)
	}

	// 处理响应
	respBody, err := transport.HandleResponse(resp)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var response proto.CompletionResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// SendCompletionMessageStream 发送文本生成请求（流式模式）
func (c *DifyCompletionClientImpl) SendCompletionMessageStream(ctx context.Context, request *proto.CompletionRequest, callback callback.CompletionStreamCallback) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DifyCompletionClient.SendCompletionMessageStream")
	defer span.Finish()

	if request == nil {
		return errors.New("request cannot be nil")
	}

	if callback == nil {
		return errors.New("callback cannot be nil")
	}

	// 设置流式响应模式
	request.ResponseMode = "streaming"

	// 构建URL
	url := fmt.Sprintf("%s%s", c.apiBaseUrl, constant.COMPLETION_MESSAGES_PATH)

	// 将请求转换为JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// 设置请求头
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + c.apiKey
	headers["Content-Type"] = transport.ContentTypeJson
	headers["Accept"] = "text/event-stream"

	// 发送请求
	resp, err := c.httpClient.Post(ctx, url, bytes.NewBuffer(jsonData), transport.ContentTypeJson, headers)
	if err != nil {
		return fmt.Errorf("failed to send stream completion message: %w", err)
	}

	// 检查响应状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		errorCode, errorMessage, _ := transport.ExtractErrorInfo(string(respBody))
		return fmt.Errorf("HTTP error %d: %s (%s)", resp.StatusCode, errorMessage, errorCode)
	}

	// 处理流式响应
	go func() {
		defer resp.Body.Close()
		reader := bufio.NewReader(resp.Body)

		var messageID string

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					callback.OnError(err.Error())
				}
				break
			}

			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			// 提取数据部分
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				callback.OnFinish(messageID)
				break
			}

			// 解析事件类型和数据
			var event proto.BaseEvent
			if err := json.Unmarshal([]byte(data), &event); err != nil {
				callback.OnError(err.Error())
				continue
			}

			switch event.Event {
			case proto.EventMessage:
				callback.OnMessage(event.Data)
			case proto.EventError:
				callback.OnError(event.Data)
			case proto.EventMessageEnd:
				var messageEndData struct {
					TaskID    string `json:"task_id"`
					MessageID string `json:"message_id"`
				}
				if err := json.Unmarshal([]byte(event.Data), &messageEndData); err == nil {
					messageID = messageEndData.MessageID
				}
			}
		}
	}()

	return nil
}

// StopCompletion 停止文本生成
func (c *DifyCompletionClientImpl) StopCompletion(ctx context.Context, taskID string, user string) (*proto.SimpleResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DifyCompletionClient.StopCompletion")
	defer span.Finish()

	if taskID == "" {
		return nil, errors.New("taskID cannot be empty")
	}

	if user == "" {
		return nil, errors.New("user cannot be empty")
	}

	// 构建URL
	url := fmt.Sprintf("%s%s/%s%s", c.apiBaseUrl, constant.COMPLETION_MESSAGES_PATH, taskID, constant.STOP_PATH)

	// 构建请求体
	request := map[string]string{
		"user": user,
	}

	// 将请求转换为JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 设置请求头
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + c.apiKey
	headers["Content-Type"] = transport.ContentTypeJson

	// 发送请求
	resp, err := c.httpClient.Post(ctx, url, bytes.NewBuffer(jsonData), transport.ContentTypeJson, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to stop completion: %w", err)
	}

	// 处理响应
	respBody, err := transport.HandleResponse(resp)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var response proto.SimpleResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// TextToAudio 文字转语音
func (c *DifyCompletionClientImpl) TextToAudio(ctx context.Context, messageID string, text string, user string) ([]byte, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DifyCompletionClient.TextToAudio")
	defer span.Finish()

	if text == "" && messageID == "" {
		return nil, errors.New("either text or messageID must be provided")
	}

	if user == "" {
		return nil, errors.New("user cannot be empty")
	}

	// 构建URL
	url := fmt.Sprintf("%s%s", c.apiBaseUrl, constant.TEXT_TO_AUDIO_PATH)

	// 构建请求体
	request := map[string]string{
		"user": user,
	}

	if messageID != "" {
		request["message_id"] = messageID
	} else {
		request["text"] = text
	}

	// 将请求转换为JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 设置请求头
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + c.apiKey
	headers["Content-Type"] = transport.ContentTypeJson
	headers["Accept"] = "audio/mpeg"

	// 发送请求
	resp, err := c.httpClient.Post(ctx, url, bytes.NewBuffer(jsonData), transport.ContentTypeJson, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// 检查响应状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		errorCode, errorMessage, _ := transport.ExtractErrorInfo(string(respBody))
		return nil, fmt.Errorf("HTTP error %d: %s (%s)", resp.StatusCode, errorMessage, errorCode)
	}

	// 读取音频数据
	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio data: %w", err)
	}
	defer resp.Body.Close()

	return audioData, nil
}
