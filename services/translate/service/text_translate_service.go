package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc-proto/gen/go/translate"
	"github.com/yb2020/odoc/config"
	LLM "github.com/yb2020/odoc/external/LLM"
	externalTranslateApi "github.com/yb2020/odoc/external/translate"
	"github.com/yb2020/odoc/internal/biz"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/http_client"
	pkgi18n "github.com/yb2020/odoc/pkg/i18n"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	baseModel "github.com/yb2020/odoc/pkg/model"
	"github.com/yb2020/odoc/pkg/utils"
	"github.com/yb2020/odoc/services/translate/dao"
	"github.com/yb2020/odoc/services/translate/model"
)

// 常量定义
const (
	TranslateStatusFailed  = 0
	TranslateStatusSuccess = 1
	FeedBackResult         = 1 // 已纠错
	SplitLabel             = ":"
)

// TranslateResult 翻译结果
type TranslateResult struct {
	TargetList  []string    `json:"targetList"`
	NetworkList [][]string  `json:"networkList"`
	ErrorCode   int         `json:"errorCode"`
	Translation []string    `json:"translation"`
	Web         []WebResult `json:"web,omitempty"`
}

// WebResult 网络翻译结果
type WebResult struct {
	Key   string   `json:"key"`
	Value []string `json:"value"`
}

// TextTranslateService 文本翻译服务
type TextTranslateService struct {
	textTranslateDAO dao.TextTranslateDAO
	config           *config.Config
	logger           logging.Logger
	tracer           opentracing.Tracer
	llmClient        *LLM.Client
	httpClient       http_client.HttpClient
}

// NewTextTranslateService 创建文本翻译服务
func NewTextTranslateService(
	config *config.Config,
	logger logging.Logger,
	tracer opentracing.Tracer,
	httpClient http_client.HttpClient,
	textTranslateDAO dao.TextTranslateDAO,
) *TextTranslateService {
	var llmClient *LLM.Client

	service := &TextTranslateService{
		textTranslateDAO: textTranslateDAO,
		config:           config,
		logger:           logger,
		tracer:           tracer,
		httpClient:       httpClient,
	}
	switch config.LLM.UseChannel {
	case "deepseek":
		// 创建DeepSeek LLM客户端配置
		clientConfig := LLM.ClientConfig{
			APIKey:     config.LLM.Channel.Deepseek.APIKey,
			BaseURL:    config.LLM.Channel.Deepseek.URL,
			Logger:     logger,
			HTTPClient: httpClient,
			Timeout:    200 * time.Second,
		}

		// 创建DeepSeek LLM客户端
		llmClient, _ = LLM.NewClient(LLM.ProviderDeepSeek,
			LLM.WithAPIKey(clientConfig.APIKey),
			LLM.WithBaseURL(clientConfig.BaseURL),
			LLM.WithHTTPClient(clientConfig.HTTPClient),
			LLM.WithLogger(clientConfig.Logger),
			LLM.WithTimeout(clientConfig.Timeout),
		)
		// if err != nil {

		// 	return nil, fmt.Errorf("创建DeepSeek客户端失败: %w", err)
		// }
		service.llmClient = llmClient
	case "gpt4o_mini":
		// 创建GPT4o-mini LLM客户端配置
		gpt4oMiniClientConfig := LLM.ClientConfig{
			APIKey:     config.LLM.Channel.Gpt4oMini.APIKey,
			BaseURL:    config.LLM.Channel.Gpt4oMini.URL,
			Logger:     logger,
			HTTPClient: httpClient,
			Timeout:    200 * time.Second,
		}

		// 创建GPT4o-mini LLM客户端
		llmClient, _ = LLM.NewClient(LLM.ProviderOpenAI,
			LLM.WithAPIKey(gpt4oMiniClientConfig.APIKey),
			LLM.WithBaseURL(gpt4oMiniClientConfig.BaseURL),
			LLM.WithHTTPClient(gpt4oMiniClientConfig.HTTPClient),
			LLM.WithLogger(gpt4oMiniClientConfig.Logger),
			LLM.WithTimeout(gpt4oMiniClientConfig.Timeout),
		)
		// if err != nil {
		// return nil, fmt.Errorf("创建GPT4o-mini客户端失败: %w", err)
		// }
		service.llmClient = llmClient
	}

	return service
}

// FindByMD5 根据MD5查找翻译记录
func (s *TextTranslateService) FindByMD5(ctx context.Context, md5Hash string) (*model.TextTranslateLog, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "TextTranslateService.FindByMD5")
	defer span.Finish()

	return s.textTranslateDAO.FindByMD5(ctx, md5Hash)
}

// RemoveById 根据ID删除翻译记录
func (s *TextTranslateService) RemoveById(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "TextTranslateService.RemoveById")
	defer span.Finish()

	return s.textTranslateDAO.DeleteById(ctx, id)
}

// Translate 翻译文本（指定源语言和目标语言）
func (s *TextTranslateService) Translate(ctx context.Context, translateChannel translate.TranslateChannel, sourceContent, sourceLanguage, targetLanguage string, useGlossary bool) (*model.TextTranslateLog, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "TextTranslateService.Translate")
	defer span.Finish()

	return s.TranslateWithPdfID(ctx, translateChannel, GetTranslateLanguage(sourceLanguage),
		GetTranslateLanguage(targetLanguage), sourceContent, "0", useGlossary)
}

// TranslateWithPdfID 带PDF ID的翻译
func (s *TextTranslateService) TranslateWithPdfID(ctx context.Context,
	translateChannel translate.TranslateChannel,
	sourceLanguage, targetLanguage translate.TranslateLanguage,
	sourceContent string,
	pdfID string,
	useGlossary bool) (*model.TextTranslateLog, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "TextTranslateService.TranslateWithPdfID")
	defer span.Finish()

	startTime := time.Now()

	// 计算内容的MD5哈希值
	md5Hash := utils.MD5(sourceContent)

	// 获取用户ID
	userId, _ := userContext.GetUserID(ctx)

	// 检查缓存
	cachedLog, err := s.textTranslateDAO.FindByMD5(ctx, md5Hash)
	if err == nil && cachedLog != nil && cachedLog.Channel == translateChannel.String() {
		s.logger.Info("从缓存获取翻译结果", "md5", md5Hash)
		return cachedLog, nil
	}

	// 调用翻译API
	translateResult, err := s.callTranslateAPI(ctx, translateChannel, sourceContent, sourceLanguage.String(), targetLanguage.String())
	if err != nil {
		s.logger.Error("调用翻译API失败", "error", err)

		// 记录失败的翻译日志
		id := idgen.GenerateUUID()
		log := &model.TextTranslateLog{
			BaseModel: baseModel.BaseModel{
				Id: id,
			},
			UserId:         userId,
			PdfId:          pdfID,
			Channel:        translateChannel.String(),
			RequestId:      generateRequestID(id, md5Hash),
			SourceLanguage: sourceLanguage.String(),
			TargetLanguage: targetLanguage.String(),
			SourceContent:  sourceContent,
			Md5Hash:        md5Hash,
			UseGlossary:    useGlossary,
			Status:         TranslateStatusFailed,
			CostMs:         time.Since(startTime).Milliseconds(),
		}

		// 异步保存日志
		uc := userContext.GetUserContext(ctx)
		userContext.RunAsyncWithUserContext(uc, func(bgCtx context.Context) {
			s.SaveTranslateLog(bgCtx, log)
		})

		return nil, errors.Biz("translate failed")
	}

	// 解析翻译结果
	targetContent, _, err := s.parseTranslateResult(translateResult)
	if err != nil {
		s.logger.Error("解析翻译结果失败", "error", err)
		return nil, errors.Biz("translate failed")
	}

	// 创建翻译日志
	id := idgen.GenerateUUID()
	log := &model.TextTranslateLog{
		BaseModel: baseModel.BaseModel{
			Id: id,
		},
		UserId:         userId,
		PdfId:          pdfID,
		Channel:        translateChannel.String(),
		RequestId:      generateRequestID(id, md5Hash),
		SourceLanguage: sourceLanguage.String(),
		TargetLanguage: targetLanguage.String(),
		SourceContent:  sourceContent,
		TargetContent:  targetContent,
		Md5Hash:        md5Hash,
		UseGlossary:    useGlossary,
		Status:         TranslateStatusSuccess,
		CostMs:         time.Since(startTime).Milliseconds(),
	}
	uc := userContext.GetUserContext(ctx)
	// 异步保存日志
	userContext.RunAsyncWithUserContext(uc, func(bgCtx context.Context) {
		s.SaveTranslateLog(bgCtx, log)
	})

	return log, nil
}

// GetTranslateLanguage 将语言字符串转换为TranslateLanguage枚举值
func GetTranslateLanguage(lang string) translate.TranslateLanguage {
	// 使用转换层标准化语言格式
	normalizedLang := pkgi18n.GlobalConverter.NormalizeToRFC5646(lang)
	return pkgi18n.GetTranslateLanguageEnum(normalizedLang)
}

// TranslateText 简单翻译文本，只返回翻译结果字符串
func (s *TextTranslateService) TranslateText(ctx context.Context,
	translateChannel translate.TranslateChannel,
	sourceLanguage, targetLanguage translate.TranslateLanguage,
	sourceContent string) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "TextTranslateService.TranslateText")
	defer span.Finish()

	// 调用翻译API
	translateResult, err := s.callTranslateAPI(ctx, translateChannel, sourceContent, sourceLanguage.String(), targetLanguage.String())
	if err != nil {
		s.logger.Error("调用翻译API失败", "error", err)
		return "", errors.Biz("translate failed")
	}

	// 解析翻译结果
	var result TranslateResult
	if err := json.Unmarshal([]byte(translateResult), &result); err != nil {
		s.logger.Error("解析翻译结果失败", "error", err)
		return "", errors.Biz("translate failed")
	}

	if result.ErrorCode != 0 {
		return "", errors.Biz(fmt.Sprintf("translate failed: %d", result.ErrorCode))
	}

	if len(result.Translation) == 0 {
		return "", nil
	}

	targetJSON, err := json.Marshal(result.Translation)
	if err != nil {
		s.logger.Error("解析翻译结果失败", "error", err)
		return "", errors.Biz("translate failed")
	}

	return string(targetJSON), nil
}

// TranslateCorrect 标记翻译结果为不正确
func (s *TextTranslateService) TranslateCorrect(ctx context.Context, requestID string, channel translate.TranslateChannel) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "TextTranslateService.TranslateCorrect")
	defer span.Finish()

	// 查找翻译记录
	log, err := s.textTranslateDAO.FindByRequestIDAndChannel(ctx, requestID, channel)
	if err != nil {
		s.logger.Error("查找翻译记录失败", "requestID", requestID, "error", err)
		return errors.Biz("translate failed")
	}

	if log == nil {
		s.logger.Warn("未找到翻译记录", "requestID", requestID)
		return errors.Biz("feed back failed")
	}

	// 更新状态为已纠错
	log.FeedBackResult = FeedBackResult
	err = s.textTranslateDAO.ModifyExcludeNull(ctx, log)
	if err != nil {
		s.logger.Error("更新翻译记录状态失败", "error", err)
		return errors.Biz("feed back failed")
	}

	return nil
}

// callTranslateAPI 调用翻译API
func (s *TextTranslateService) callTranslateAPI(ctx context.Context, channel translate.TranslateChannel, sourceContent, sourceLanguage, targetLanguage string) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "TextTranslateService.callTranslateAPI")
	defer span.Finish()

	switch channel {
	case translate.TranslateChannel_GOOGLE:
		googleFreeClient := externalTranslateApi.NewGoogleFreeTranslateClient(*s.config, s.httpClient, s.logger)
		return googleFreeClient.Translate(sourceContent, sourceLanguage, targetLanguage)
	case translate.TranslateChannel_YOUDAO:
		youdaoClient := externalTranslateApi.NewYoudaoTranslateClient(*s.config, s.httpClient, s.logger)
		return youdaoClient.Translate(sourceContent, sourceLanguage, targetLanguage)
	}

	return "", errors.Biz("translate failed")
}

// parseTranslateResult 解析翻译结果
func (s *TextTranslateService) parseTranslateResult(translateResult string) (string, string, error) {
	var result TranslateResult
	if err := json.Unmarshal([]byte(translateResult), &result); err != nil {
		return "", "", err
	}

	// 解析目标翻译
	var targetList []string
	if len(result.Translation) > 0 {
		targetList = result.Translation
	}
	targetJSON, err := json.Marshal(targetList)
	if err != nil {
		return "", "", err
	}

	// 解析网络翻译
	var networkList [][]string
	if len(result.Web) > 0 {
		for _, web := range result.Web {
			networkList = append(networkList, web.Value)
		}
	}
	networkJSON, err := json.Marshal(networkList)
	if err != nil {
		return "", "", err
	}

	return string(targetJSON), string(networkJSON), nil
}

// saveTranslateLog 保存翻译日志
func (s *TextTranslateService) SaveTranslateLog(ctx context.Context, log *model.TextTranslateLog) {
	if log == nil || log.TargetContent == "" || len(log.TargetContent) == 0 { // 空值不保存
		return
	}
	// 保存到数据库
	if err := s.textTranslateDAO.SaveExcludeNull(ctx, log); err != nil {
		s.logger.Error("保存翻译日志失败", "error", err)
	}
}

// generateRequestID 生成请求ID
func generateRequestID(id string, md5Hash string) string {
	requestID := fmt.Sprintf("%d%s%s", id, SplitLabel, md5Hash)
	return base64.StdEncoding.EncodeToString([]byte(requestID))
}

// StreamTranslateWithLLM 使用LLM进行流式翻译
func (s *TextTranslateService) StreamTranslateWithLLM(ctx context.Context, text string, sourceLanguage, targetLanguage translate.TranslateLanguage, handler func(content string) error) (string, error) {
	// 构建提示
	var prompt string
	prompt = fmt.Sprintf("请将以下文本从%s翻译成%s，保持原文的格式和语气，不要有任何解释或额外内容：：\n\n%s",
		sourceLanguage.String(), targetLanguage.String(), text)

	// 创建聊天请求
	chatRequest := LLM.ChatCompletionRequest{
		Messages: []LLM.Message{
			{
				Role:    LLM.RoleUser,
				Content: prompt,
			},
		},
		Stream: true,
		TopP:   0.1,
		// Temperature: 0.7,
		MaxTokens: 1500,
	}

	// 收集完整的翻译结果
	var fullContent strings.Builder

	// 使用流式API
	err := s.llmClient.CreateChatCompletionStream(ctx, chatRequest, func(response LLM.ChatCompletionStreamResponse) error {
		if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
			content := response.Choices[0].Delta.Content
			fullContent.WriteString(content)

			// 调用回调函数处理内容
			if handler != nil {
				if err := handler(content); err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("调用流式API失败: %w", err)
	}

	return fullContent.String(), nil
}

// StreamTranslateWithCallback 使用LLM进行流式翻译，并通过回调函数处理结果
func (s *TextTranslateService) StreamTranslateWithCallback(
	ctx context.Context,
	text string,
	sourceLanguage translate.TranslateLanguage,
	targetLanguage translate.TranslateLanguage,
	userID string,
	pdfID string,
	userGlossary bool,
	startTime time.Time,
	resultCallback func(content string) error,
	completeCallback func(translatedText string) error,
) error {
	// 使用服务层的StreamTranslateWithLLM方法进行流式翻译
	var fullContent strings.Builder

	// 定义处理流式响应的回调函数
	handler := func(content string) error {
		if resultCallback != nil {
			return resultCallback(content)
		}
		return nil
	}

	// 调用服务层的流式翻译方法
	translatedText, err := s.StreamTranslateWithLLM(
		ctx,
		text,
		sourceLanguage,
		targetLanguage,
		handler,
	)

	if err != nil {
		return err
	}

	fullContent.WriteString(translatedText)

	// 调用完成回调
	if completeCallback != nil {
		if err := completeCallback(fullContent.String()); err != nil {
			return err
		}
	}
	id := idgen.GenerateUUID()
	md5Hash := utils.MD5(text)
	uc := userContext.GetUserContext(ctx)
	targetContent := []string{fullContent.String()}
	resultJSON, err := json.Marshal(targetContent)
	if err != nil {
		s.logger.Error("序列化翻译结果失败", "error", err)
		return errors.System(errors.ErrorTypeExternalInterface, fmt.Sprintf("serialize target content failed: %s", err.Error()), err)
	}

	// 保存翻译日志，使用新的上下文但保留用户信息
	userContext.RunAsyncWithUserContext(uc, func(bgCtx context.Context) {
		// 创建翻译日志
		translateLog := &model.TextTranslateLog{
			BaseModel: baseModel.BaseModel{
				Id: id,
			},
			UserId:         userID,
			SourceContent:  text,
			RequestId:      generateRequestID(id, md5Hash),
			TargetContent:  string(resultJSON),
			SourceLanguage: sourceLanguage.String(),
			TargetLanguage: targetLanguage.String(),
			Channel:        translate.TranslateChannel_AI.String(),
			PdfId:          pdfID,
			Status:         TranslateStatusSuccess,
			UseGlossary:    userGlossary,
			Md5Hash:        md5Hash,
			CostMs:         time.Since(startTime).Milliseconds(),
		}

		// 保存日志
		s.SaveTranslateLog(bgCtx, translateLog)
	})

	return nil
}

// TranslateInternal 内部翻译功能，尝试多个翻译渠道
func (s *TextTranslateService) TranslateInternal(ctx context.Context, content string,
	sourceLanguage, targetLanguage translate.TranslateLanguage,
	paperId, pdfId, pageId string, blockId string) (*model.TextTranslateLog, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "TextTranslateService.TranslateInternal")
	defer span.Finish()

	startTime := time.Now()
	id := idgen.GenerateUUID()

	targetContent := ""
	var err error
	var translateChannel translate.TranslateChannel

	// 尝试按照配置的顺序使用不同的翻译渠道
	for _, channelName := range s.config.Translate.FullTextTranslate.InternalOrderChannel {
		translateChannel = getTranslateChannelFromString(channelName)
		s.logger.Info("尝试使用翻译渠道", "channel", translateChannel.String())
		targetContent, err = s.TranslateText(ctx, translateChannel, sourceLanguage, targetLanguage, content)
		if err == nil {
			// 翻译成功，跳出循环
			s.logger.Info("翻译成功", "channel", translateChannel.String())
			break
		}
		s.logger.Warn("翻译失败，尝试下一个渠道", "channel", translateChannel.String(), "error", err)
	}

	// 如果所有渠道都失败了
	if err != nil {
		return nil, errors.BizWithStatus(biz.Translate_StatusAllFailed, "translate failed, all channels failed")
	}

	// 创建并保存翻译日志
	translateLog := s.insertServiceInvokeLog(ctx, id, pdfId, paperId, pageId, blockId, translateChannel,
		sourceLanguage, targetLanguage, content, TranslateStatusSuccess, targetContent, startTime)

	return translateLog, nil
}

// getTranslateChannelFromString 将字符串转换为TranslateChannel枚举
func getTranslateChannelFromString(channelName string) translate.TranslateChannel {
	switch strings.ToUpper(channelName) {
	case translate.TranslateChannel_GOOGLE.String():
		return translate.TranslateChannel_GOOGLE
	case translate.TranslateChannel_YOUDAO.String():
		return translate.TranslateChannel_YOUDAO
	case translate.TranslateChannel_AI.String():
		return translate.TranslateChannel_AI
	default:
		// 默认返回GOOGLE作为后备选项
		return translate.TranslateChannel_GOOGLE
	}
}

// insertServiceInvokeLog 创建并保存翻译日志
func (s *TextTranslateService) insertServiceInvokeLog(ctx context.Context, id, pdfId, paperId, pageId string, blockId string, channelType translate.TranslateChannel,
	sourceLanguage, targetLanguage translate.TranslateLanguage, sourceContent string, status int, targetContent string, startTime time.Time) *model.TextTranslateLog {
	md5Hash := utils.MD5(sourceContent)

	// 创建翻译日志
	translateLog := &model.TextTranslateLog{
		BaseModel: baseModel.BaseModel{
			Id: id,
		},
		UserId:         "-1",
		PdfId:          pdfId,
		Channel:        channelType.String(),
		RequestId:      generateRequestID(id, md5Hash),
		SourceLanguage: sourceLanguage.String(),
		TargetLanguage: targetLanguage.String(),
		SourceContent:  sourceContent,
		TargetContent:  targetContent,
		Md5Hash:        md5Hash,
		UseGlossary:    false,
		Status:         status,
		CostMs:         time.Since(startTime).Milliseconds(),
		ExtParams:      fmt.Sprintf("{\"paperId\":%d,\"pageId\":%d,\"blockId\":\"%s\"}", paperId, pageId, blockId),
	}

	uc := userContext.GetUserContext(ctx)
	// 异步保存日志
	userContext.RunAsyncWithUserContext(uc, func(bgCtx context.Context) {
		s.SaveTranslateLog(bgCtx, translateLog)
	})

	return translateLog
}
