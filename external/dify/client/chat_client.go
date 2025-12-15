package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/external/dify/callback"
	"github.com/yb2020/odoc/external/dify/constant"
	"github.com/yb2020/odoc/external/dify/proto"
	"github.com/yb2020/odoc/external/dify/transport"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/utils"
)

// DifyChatClient Dify聊天客户端接口
type DifyChatClient interface {
	// SendChatMessage 发送聊天消息
	SendChatMessage(ctx context.Context, message *proto.ChatMessage) (*proto.ChatMessageResponse, error)

	// SendStreamChatMessage 发送流式聊天消息
	SendStreamChatMessage(ctx context.Context, message *proto.ChatMessage, callback callback.ChatStreamCallback) error

	// StopChatMessage 停止聊天消息
	StopChatMessage(ctx context.Context, taskID string, user string) (*proto.SimpleResponse, error)

	// FeedbackMessage 消息反馈（
	FeedbackMessage(ctx context.Context, messageID string, rating string, user string, content string) (*proto.SimpleResponse, error)

	// GetSuggestedQuestions 获取下一轮建议问题列表
	GetSuggestedQuestions(ctx context.Context, messageID string, user string) (*proto.SuggestedQuestionsResponse, error)

	// GetConversationMessages 获取对话消息
	GetConversationMessages(ctx context.Context, conversationId, user, firstId string, limit int) (*proto.ConversationMessagesResponse, error)

	// GetConversations 获取对话列表
	GetConversations(ctx context.Context, user string, firstID string, limit int) (*proto.ConversationsResponse, error)

	// RenameConversation 重命名对话
	RenameConversation(ctx context.Context, conversationID string, name string) (*proto.SimpleResponse, error)

	// DeleteConversation 删除对话
	DeleteConversation(ctx context.Context, conversationID string) (*proto.SimpleResponse, error)

	// AudioToText 语音转文字
	AudioToText(ctx context.Context, filePath string, user string) (*proto.AudioToTextResponse, error)

	// TextToAudio 文字转语音
	TextToAudio(ctx context.Context, messageID string, text string, user string) ([]byte, error)

	// GetAppMeta 获取应用元数据
	GetAppMeta(ctx context.Context) (*proto.AppMetaResponse, error)

	// GetAnnotations 获取标注列表
	GetAnnotations(ctx context.Context, page int, limit int) (*proto.AnnotationListResponse, error)

	// SaveAnnotation 创建标注
	SaveAnnotation(ctx context.Context, question string, answer string) (*proto.Annotation, error)

	// UpdateAnnotation 更新标注
	UpdateAnnotation(ctx context.Context, annotationID string, question string, answer string) (*proto.Annotation, error)

	// DeleteAnnotation 删除标注
	DeleteAnnotation(ctx context.Context, annotationID string) (*proto.SimpleResponse, error)

	// AnnotationReply 标注回复初始设置
	AnnotationReply(ctx context.Context, action string, embeddingProviderName string, embeddingModelName string, scoreThreshold int) (*proto.AnnotationReply, error)

	// GetAnnotationReply 查询标注回复初始设置任务状态
	GetAnnotationReply(ctx context.Context, action string, jobID string) (*proto.AnnotationReply, error)

	// SendCompletionMessage 发送完成消息
	SendCompletionMessage(ctx context.Context, request *proto.CompletionRequest) (*proto.CompletionResponse, error)

	// UploadFile 上传文件
	UploadFile(ctx context.Context, request *proto.FileUploadRequest) (*proto.FileUploadResponse, error)
}

// DifyChatClientImpl Dify聊天客户端实现
type DifyChatClientImpl struct {
	DifyBaseClientImpl
	logger logging.Logger
	tracer opentracing.Tracer
}

// NewDifyChatClient 创建新的Dify聊天客户端
func NewDifyChatClient(logger logging.Logger, tracer opentracing.Tracer, apiBaseUrl, apiKey string, timeout int, responseHeaderTimeout int) DifyChatClient {
	return &DifyChatClientImpl{
		DifyBaseClientImpl: DifyBaseClientImpl{
			apiKey:     apiKey,
			apiBaseUrl: apiBaseUrl,
			httpClient: transport.NewDefaultHttpClient(logger, timeout, responseHeaderTimeout),
			tracer:     tracer,
		},
		logger: logger,
		tracer: tracer,
	}
}

// NewDifyChatClientWithHttpClient 使用自定义HTTP客户端创建新的Dify聊天客户端
func NewDifyChatClientWithHttpClient(logger logging.Logger, tracer opentracing.Tracer, apiBaseUrl, apiKey string, httpClient transport.HttpClient) DifyChatClient {
	return &DifyChatClientImpl{
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
func (c *DifyChatClientImpl) SendChatMessage(ctx context.Context, message *proto.ChatMessage) (*proto.ChatMessageResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyChatClientImpl.SendChatMessage")
	defer span.Finish()

	if message == nil {
		c.logger.Error("消息不能为空")
		return nil, fmt.Errorf("消息不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.CHAT_MESSAGES_PATH)

	// 发送JSON请求并解析响应
	var response proto.ChatMessageResponse
	err := c.doJSONRequest(ctx, http.MethodPost, url, message, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// SendStreamChatMessage 发送流式聊天消息
func (c *DifyChatClientImpl) SendStreamChatMessage(ctx context.Context, message *proto.ChatMessage, callback callback.ChatStreamCallback) error {
	if message == nil {
		return fmt.Errorf("消息不能为空")
	}
	if callback == nil {
		return fmt.Errorf("回调不能为空")
	}

	// 设置流式响应模式
	message.ResponseMode = "streaming"

	// 构建URL
	url := transport.BuildUrl(c.apiBaseUrl, constant.CHAT_MESSAGES_PATH)
	// 设置请求头
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + c.apiKey
	headers["Content-Type"] = transport.ContentTypeJson
	headers["Accept"] = "text/event-stream"

	// 创建一个缓冲区用于存储未完整处理的数据
	buffer := ""
	// 发送请求
	return c.httpClient.PostStream(ctx, url, message, headers, func(chunk []byte) error {
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
			rowJsonData := strings.TrimPrefix(line, "data: ")

			// 检查是否为流结束标记
			if rowJsonData == "[DONE]" {
				return nil
			}

			// 解析事件类型和数据
			event := proto.ChatMessageEvent{}
			if err := utils.ParseJsonByObj(rowJsonData, &event); err != nil {
				event.Data.Error = err.Error()
				callback.OnError(event)
				return nil
			}
			switch event.Event {
			case proto.EventWorkflowStarted:
				callback.OnWorkflowStarted(event)
				callback.OnStart(event) //直接套用
			case proto.EventWorkflowFinished:
				callback.OnWorkflowFinished(event)
				callback.OnEnd() //直接套用
			case proto.EventNodeStarted:
				callback.OnNodeStarted(event)
			case proto.EventNodeFinished:
				callback.OnNodeFinished(event)
			case proto.EventMessage:
				callback.OnMessage(event)
			case proto.EventError:
				callback.OnError(event)
			case proto.EventMessageEnd:
				callback.OnFinish(event)
			case proto.EventPing:
				callback.OnPing(event)
			}
		}
		return nil
	})
}

// GetConversationMessages 获取对话消息
func (c *DifyChatClientImpl) GetConversationMessages(ctx context.Context, conversationId string, user string, firstId string, limit int) (*proto.ConversationMessagesResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyChatClientImpl.GetConversationMessages")
	defer span.Finish()

	if conversationId == "" {
		return nil, fmt.Errorf("对话ID不能为空")
	}

	// 准备查询参数
	queryParams := map[string]string{}
	if firstId != "" {
		queryParams["first_id"] = firstId
	}
	if user != "" {
		queryParams["user"] = user
	}
	if limit > 0 {
		queryParams["limit"] = fmt.Sprintf("%d", limit)
	}

	// 创建响应对象
	response := &proto.ConversationMessagesResponse{}

	// 发送请求并解析响应
	path := fmt.Sprintf("/conversations/%s/messages", conversationId)
	err := c.getWithParams(ctx, path, queryParams, response)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation messages: %w", err)
	}

	return response, nil
}

// GetConversations 获取对话列表
func (c *DifyChatClientImpl) GetConversations(ctx context.Context, user string, firstID string, limit int) (*proto.ConversationsResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyChatClientImpl.GetConversations")
	defer span.Finish()

	if user == "" {
		return nil, fmt.Errorf("用户不能为空")
	}

	// 准备查询参数
	queryParams := map[string]string{
		"user": user,
	}

	if firstID != "" {
		queryParams["first_id"] = firstID
	}

	if limit > 0 {
		queryParams["limit"] = fmt.Sprintf("%d", limit)
	}

	// 发送请求并解析响应
	var response proto.ConversationsResponse
	err := c.getWithParams(ctx, constant.CONVERSATIONS_PATH, queryParams, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// RenameConversation 重命名对话
func (c *DifyChatClientImpl) RenameConversation(ctx context.Context, conversationID string, name string) (*proto.SimpleResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyChatClientImpl.RenameConversation")
	defer span.Finish()

	if conversationID == "" {
		return nil, fmt.Errorf("对话ID不能为空")
	}
	if name == "" {
		return nil, fmt.Errorf("名称不能为空")
	}

	// 构建请求体
	request := map[string]string{
		"name": name,
	}

	// 创建响应对象
	response := &proto.SimpleResponse{}

	// 发送请求并解析响应
	path := c.buildURL(constant.CONVERSATIONS_PATH, conversationID)
	err := c.doJSONRequest(ctx, http.MethodPut, path, request, response)
	if err != nil {
		return nil, fmt.Errorf("failed to rename conversation: %w", err)
	}

	return response, nil
}

// DeleteConversation 删除对话
func (c *DifyChatClientImpl) DeleteConversation(ctx context.Context, conversationID string) (*proto.SimpleResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyChatClientImpl.DeleteConversation")
	defer span.Finish()

	if conversationID == "" {
		return nil, fmt.Errorf("对话ID不能为空")
	}

	// 创建响应对象
	response := &proto.SimpleResponse{}

	// 发送请求并解析响应
	path := c.buildURL(constant.CONVERSATIONS_PATH, conversationID)
	err := c.doJSONRequest(ctx, http.MethodDelete, path, nil, response)
	if err != nil {
		return nil, fmt.Errorf("failed to delete conversation: %w", err)
	}

	return response, nil
}

// SendCompletionMessage 发送完成消息
func (c *DifyChatClientImpl) SendCompletionMessage(ctx context.Context, request *proto.CompletionRequest) (*proto.CompletionResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyChatClientImpl.SendCompletionMessage")
	defer span.Finish()

	if request == nil {
		return nil, fmt.Errorf("请求不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.COMPLETION_MESSAGES_PATH)

	// 发送JSON请求并解析响应
	var response proto.CompletionResponse
	err := c.doJSONRequest(ctx, http.MethodPost, url, request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// StopChatMessage 停止聊天消息
func (c *DifyChatClientImpl) StopChatMessage(ctx context.Context, taskID string, user string) (*proto.SimpleResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyChatClientImpl.StopChatMessage")
	defer span.Finish()

	if taskID == "" {
		return nil, errors.New("taskID cannot be empty")
	}

	if user == "" {
		return nil, errors.New("user cannot be empty")
	}

	// 构建请求体
	request := map[string]string{
		"user": user,
	}

	// 创建响应对象
	response := &proto.SimpleResponse{}

	// 发送请求并解析响应
	path := c.buildURL(constant.CHAT_MESSAGES_PATH, taskID, constant.STOP_PATH)
	err := c.doJSONRequest(ctx, http.MethodPost, path, request, response)
	if err != nil {
		return nil, fmt.Errorf("failed to stop chat message: %w", err)
	}

	return response, nil
}

// FeedbackMessage 消息反馈（点赞）
func (c *DifyChatClientImpl) FeedbackMessage(ctx context.Context, messageID string, rating string, user string, content string) (*proto.SimpleResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyChatClientImpl.FeedbackMessage")
	defer span.Finish()

	if messageID == "" {
		return nil, errors.New("messageID cannot be empty")
	}

	if rating == "" {
		return nil, errors.New("rating cannot be empty")
	}

	if user == "" {
		return nil, errors.New("user cannot be empty")
	}

	// 构建请求体
	request := map[string]string{
		"rating":  rating,
		"user":    user,
		"content": content,
	}

	// 创建响应对象
	response := &proto.SimpleResponse{}

	// 发送请求并解析响应
	path := c.buildURL(constant.CHAT_MESSAGES_PATH, messageID, constant.FEEDBACKS_PATH)
	err := c.doJSONRequest(ctx, http.MethodPost, path, request, response)
	if err != nil {
		return nil, fmt.Errorf("failed to send feedback: %w", err)
	}

	return response, nil
}

// GetSuggestedQuestions 获取下一轮建议问题列表
func (c *DifyChatClientImpl) GetSuggestedQuestions(ctx context.Context, messageID string, user string) (*proto.SuggestedQuestionsResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyChatClientImpl.GetSuggestedQuestions")
	defer span.Finish()

	if messageID == "" {
		return nil, errors.New("messageID cannot be empty")
	}

	if user == "" {
		return nil, errors.New("user cannot be empty")
	}

	// 构建路径
	path := fmt.Sprintf("%s/%s%s", constant.MESSAGES_PATH, messageID, constant.SUGGESTED_QUESTIONS_PATH)

	// 创建响应对象
	difyResponse := &proto.SuggestedQuestionsResult{}
	// 发送请求并解析响应
	err := c.getWithParams(ctx, path, map[string]string{
		"user": user,
	}, difyResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggested questions: %w", err)
	}
	response := &proto.SuggestedQuestionsResponse{
		Questions: difyResponse.Data,
	}
	return response, nil
}

// AudioToText 语音转文字
func (c *DifyChatClientImpl) AudioToText(ctx context.Context, filePath string, user string) (*proto.AudioToTextResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyChatClientImpl.AudioToText")
	defer span.Finish()

	if filePath == "" {
		return nil, errors.New("filePath cannot be empty")
	}

	if user == "" {
		return nil, errors.New("user cannot be empty")
	}

	// 构建URL
	url := c.buildURL(constant.AUDIO_TO_TEXT_PATH)

	// 准备表单数据
	formData := map[string]interface{}{
		"user": user,
	}

	// 准备文件参数
	fileParams := map[string]string{
		"file": filePath,
	}

	// 发送多部分表单请求并解析响应
	var response proto.AudioToTextResponse
	err := c.doMultipartRequest(ctx, http.MethodPost, url, formData, fileParams, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to convert audio to text: %w", err)
	}

	return &response, nil
}

// TextToAudio 文字转语音
func (c *DifyChatClientImpl) TextToAudio(ctx context.Context, messageID string, text string, user string) ([]byte, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyChatClientImpl.TextToAudio")
	defer span.Finish()

	if text == "" && messageID == "" {
		return nil, errors.New("either text or messageID must be provided")
	}

	if user == "" {
		return nil, errors.New("user cannot be empty")
	}

	// 构建URL
	url := c.buildURL(constant.TEXT_TO_AUDIO_PATH)

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

	// 设置特殊的请求头
	headers := map[string]string{
		"Accept": "audio/mpeg", // 指定返回音频格式
	}

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

// GetAppMeta 获取应用元数据
func (c *DifyChatClientImpl) GetAppMeta(ctx context.Context) (*proto.AppMetaResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyChatClientImpl.GetAppMeta")
	defer span.Finish()

	// 构建URL
	url := c.buildURL(constant.META_PATH)

	// 发送请求并解析响应
	var response proto.AppMetaResponse
	err := c.doJSONRequest(ctx, http.MethodGet, url, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get app meta: %w", err)
	}

	return &response, nil
}

// GetAnnotations 获取标注列表
func (c *DifyChatClientImpl) GetAnnotations(ctx context.Context, page int, limit int) (*proto.AnnotationListResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyChatClientImpl.GetAnnotations")
	defer span.Finish()

	// 准备查询参数
	queryParams := map[string]string{}
	if page > 0 {
		queryParams["page"] = fmt.Sprintf("%d", page)
	}
	if limit > 0 {
		queryParams["limit"] = fmt.Sprintf("%d", limit)
	}

	// 创建响应对象
	response := &proto.AnnotationListResponse{}

	// 发送请求并解析响应
	err := c.getWithParams(ctx, constant.APPS_ANNOTATIONS_PATH, queryParams, response)
	if err != nil {
		return nil, fmt.Errorf("failed to get annotations: %w", err)
	}

	return response, nil
}

// SaveAnnotation 创建标注
func (c *DifyChatClientImpl) SaveAnnotation(ctx context.Context, question string, answer string) (*proto.Annotation, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyChatClientImpl.SaveAnnotation")
	defer span.Finish()

	if question == "" {
		return nil, errors.New("question cannot be empty")
	}

	if answer == "" {
		return nil, errors.New("answer cannot be empty")
	}

	// 构建请求体
	request := map[string]string{
		"question": question,
		"answer":   answer,
	}

	// 创建响应对象
	response := &proto.Annotation{}

	// 发送请求并解析响应
	err := c.doJSONRequest(ctx, http.MethodPost, c.buildURL(constant.APPS_ANNOTATIONS_PATH), request, response)
	if err != nil {
		return nil, fmt.Errorf("failed to save annotation: %w", err)
	}

	return response, nil
}

// UpdateAnnotation 更新标注
func (c *DifyChatClientImpl) UpdateAnnotation(ctx context.Context, annotationID string, question string, answer string) (*proto.Annotation, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyChatClientImpl.UpdateAnnotation")
	defer span.Finish()

	if annotationID == "" {
		return nil, errors.New("annotationID cannot be empty")
	}

	if question == "" {
		return nil, errors.New("question cannot be empty")
	}

	if answer == "" {
		return nil, errors.New("answer cannot be empty")
	}

	// 构建请求体
	request := map[string]string{
		"question": question,
		"answer":   answer,
	}

	// 创建响应对象
	response := &proto.Annotation{}

	// 发送请求并解析响应
	path := c.buildURL(constant.APPS_ANNOTATIONS_PATH, annotationID)
	err := c.doJSONRequest(ctx, http.MethodPut, path, request, response)
	if err != nil {
		return nil, fmt.Errorf("failed to update annotation: %w", err)
	}

	return response, nil
}

// DeleteAnnotation 删除标注
func (c *DifyChatClientImpl) DeleteAnnotation(ctx context.Context, annotationID string) (*proto.SimpleResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyChatClientImpl.DeleteAnnotation")
	defer span.Finish()

	if annotationID == "" {
		return nil, errors.New("annotationID cannot be empty")
	}

	// 创建响应对象
	response := &proto.SimpleResponse{}

	// 发送请求并解析响应
	path := c.buildURL(constant.APPS_ANNOTATIONS_PATH, annotationID)
	err := c.doJSONRequest(ctx, http.MethodDelete, path, nil, response)
	if err != nil {
		return nil, fmt.Errorf("failed to delete annotation: %w", err)
	}

	return response, nil
}

// AnnotationReply 标注回复初始设置
func (c *DifyChatClientImpl) AnnotationReply(ctx context.Context, action string, embeddingProviderName string, embeddingModelName string, scoreThreshold int) (*proto.AnnotationReply, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyChatClientImpl.AnnotationReply")
	defer span.Finish()

	if action == "" {
		return nil, errors.New("action cannot be empty")
	}

	// 构建请求体
	request := map[string]interface{}{
		"action": action,
	}

	if embeddingProviderName != "" {
		request["embedding_provider_name"] = embeddingProviderName
	}

	if embeddingModelName != "" {
		request["embedding_model_name"] = embeddingModelName
	}

	if scoreThreshold > 0 {
		request["score_threshold"] = scoreThreshold
	}

	// 创建响应对象
	response := &proto.AnnotationReply{}

	// 发送请求并解析响应
	path := c.buildURL(constant.APPS_ANNOTATIONS_REPLY_PATH)
	err := c.doJSONRequest(ctx, http.MethodPost, path, request, response)
	if err != nil {
		return nil, fmt.Errorf("failed to set annotation reply: %w", err)
	}

	return response, nil
}

// GetAnnotationReply 查询标注回复初始设置任务状态
func (c *DifyChatClientImpl) GetAnnotationReply(ctx context.Context, action string, jobID string) (*proto.AnnotationReply, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyChatClientImpl.GetAnnotationReply")
	defer span.Finish()

	if action == "" {
		return nil, errors.New("action cannot be empty")
	}

	if jobID == "" {
		return nil, errors.New("jobID cannot be empty")
	}

	// 准备查询参数
	queryParams := map[string]string{
		"action": action,
		"job_id": jobID,
	}

	// 创建响应对象
	response := &proto.AnnotationReply{}

	// 发送请求并解析响应
	err := c.getWithParams(ctx, constant.APPS_ANNOTATIONS_REPLY_PATH, queryParams, response)
	if err != nil {
		return nil, fmt.Errorf("failed to get annotation reply status: %w", err)
	}

	return response, nil
}

// UploadFile 上传文件
func (c *DifyChatClientImpl) UploadFile(ctx context.Context, request *proto.FileUploadRequest) (*proto.FileUploadResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyChatClientImpl.UploadFile")
	defer span.Finish()

	if request == nil {
		return nil, errors.New("request cannot be nil")
	}

	if request.File == nil {
		return nil, errors.New("file cannot be nil")
	}

	// 构建URL
	url := c.buildURL(constant.FILES_UPLOAD_PATH)

	// 创建multipart表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加用户字段
	if err := writer.WriteField("user", request.User); err != nil {
		return nil, fmt.Errorf("failed to add user field: %w", err)
	}

	// 创建文件表单字段
	var part io.Writer
	var err error
	if request.MediaType != "" {
		// 如果指定了 MediaType，使用它
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", request.File.Filename))
		h.Set("Content-Type", request.MediaType)
		part, err = writer.CreatePart(h)
	} else {
		// 否则使用默认行为
		part, err = writer.CreateFormFile("file", request.File.Filename)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// 复制文件内容到表单字段
	file, err := request.File.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if _, err = io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// 关闭writer以完成multipart消息
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// 执行请求
	respBody, err := c.doRequest(ctx, http.MethodPost, url, body, writer.FormDataContentType(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// 解析响应
	var response proto.FileUploadResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}
