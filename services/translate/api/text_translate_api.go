package api

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/config"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/errors"
	pkgi18n "github.com/yb2020/odoc/pkg/i18n"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/ratelimit"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	"github.com/yb2020/odoc/pkg/utils"
	pb "github.com/yb2020/odoc/proto/gen/go/membership"
	"github.com/yb2020/odoc/proto/gen/go/translate"
	membershipService "github.com/yb2020/odoc/services/membership/interfaces"
	"github.com/yb2020/odoc/services/translate/model"
	"github.com/yb2020/odoc/services/translate/service"
)

// 常量定义
const (
	TranslateChannelYoudao = "youdao"
	TranslateChannelBaidu  = "baidu"
	TranslateChannelGoogle = "google"
	TranslateChannelAI     = "ai"
)

var (
	// 单词模式正则表达式 - 匹配整个内容是否为单个单词
	wordPattern = regexp.MustCompile(`^\s*([A-Za-z]+)\s*$`)
)

// TextTranslateAPI 文本翻译API处理器
type TextTranslateAPI struct {
	textTranslateService     *service.TextTranslateService
	glossaryService          *service.GlossaryService
	wordPronunciationService *service.WordPronunciationService
	rateLimitService         *ratelimit.RateLimiterService
	membershipService        membershipService.IMembershipService
	config                   *config.Config
	logger                   logging.Logger
	tracer                   opentracing.Tracer
	translateLimiter         ratelimit.RateLimiter
}

// NewTextTranslateAPI 创建文本翻译API处理器
func NewTextTranslateAPI(
	config *config.Config,
	logger logging.Logger,
	tracer opentracing.Tracer,
	textTranslateService *service.TextTranslateService,
	glossaryService *service.GlossaryService,
	wordPronunciationService *service.WordPronunciationService,
	rateLimitService *ratelimit.RateLimiterService,
	membershipService membershipService.IMembershipService,
) *TextTranslateAPI {

	var limiter ratelimit.RateLimiter
	limiter, err := rateLimitService.CreateLimiter(ratelimit.LimiterConfig{
		Type:       ratelimit.CounterLimiterType,
		KeyPrefix:  "translate",
		MaxRate:    10,
		TimeUnit:   ratelimit.Second,
		Dimension:  ratelimit.User,
		ExpireTime: 5,
	})
	// 检查限流
	if err != nil {
		logger.Error("Failed to create rate limiter", "error", err)
	}

	return &TextTranslateAPI{
		textTranslateService:     textTranslateService,
		glossaryService:          glossaryService,
		wordPronunciationService: wordPronunciationService,
		rateLimitService:         rateLimitService,
		membershipService:        membershipService,
		config:                   config,
		logger:                   logger,
		tracer:                   tracer,
		translateLimiter:         limiter,
	}
}

// Translate 中文翻译服务
func (api *TextTranslateAPI) Translate(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "TextTranslateAPI.Translate")
	defer span.Finish()

	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())

	// 检查限流
	// Check rate limit
	if api.translateLimiter != nil {
		// Create a key specific to this user
		key := ratelimit.GetLimiterKey("translate", ratelimit.User, userId)
		result, err := api.translateLimiter.Allow(c.Request.Context(), key, 2)

		if err != nil {
			api.logger.Error("Rate limit check failed", "error", err)
		} else if !result.Allowed {
			c.Error(errors.System(errors.ErrorTypeRateLimit, "You're making requests too quickly. Please wait a moment.", nil))
			return
		}
	}

	// 使用 Proto 绑定器解析请求体
	req := &translate.TextTranslateRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("msg", "解析翻译请求失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 验证请求
	if len(req.Text) > api.config.Translate.Text.MaxLength {
		c.Error(errors.System(errors.ErrorTypeInvalidArgument, "translate content exceeds maximum length", nil))
		return
	}

	// 处理内容
	content := req.Text
	var glossaryResult *model.GlossaryTranslateModel

	// 处理术语库
	if req.UseGlossary != nil && *req.UseGlossary {
		var err error
		glossaryResult, err = api.glossaryService.ReplaceOriginalText(ctx, userId, content)
		if err != nil {
			api.logger.Warn("msg", "处理术语库失败", "error", err.Error())
		} else if glossaryResult != nil && glossaryResult.UserSelectedTextAfterReplace != "" {
			content = glossaryResult.UserSelectedTextAfterReplace
		}
	}

	// 检查是否是单词模式
	translateResp := &translate.TranslateResponse{}
	if matches := wordPattern.FindStringSubmatch(content); len(matches) > 1 {
		// 使用CreditFunTranslate包装单词翻译逻辑
		err := api.membershipService.CreditFunTranslate(ctx, pb.CreditServiceType_CREDIT_SERVICE_TYPE_TRANSLATE_WORD, 0, func(xctx context.Context, sessionId string) error {
			// 如果是单词，使用单词发音服务
			wordContent := strings.ToLower(matches[1])
			if req.PdfId == nil {
				*req.PdfId = "0"
			}

			resp, err := api.wordPronunciationService.GetTranslateResp(xctx, service.TranslateSourceYoudao, wordContent, *req.PdfId)
			if err != nil {
				return err
			}

			if resp != nil {
				// 处理术语库
				if req.UseGlossary != nil && *req.UseGlossary && glossaryResult != nil {
					api.glossaryService.DealTranslationText(xctx, resp, glossaryResult)
					if glossaryResult.OriginalTranslationAfterReplace != "" {
						resp.TargetContent = []string{glossaryResult.OriginalTranslationAfterReplace}
					}
				}
				translateResp = resp
			}
			return nil
		}, true)
		if err != nil {
			c.Error(err)
			return
		}

		// 在积分扣除成功后再返回响应
		if translateResp != nil {
			response.Success(c, "success", translateResp)
		}
		return
	}

	// 如果不是单词或单词翻译失败，使用普通翻译
	// 确定翻译类型
	var creditServiceType pb.CreditServiceType
	if req.Channel == translate.TranslateChannel_AI {
		creditServiceType = pb.CreditServiceType_CREDIT_SERVICE_TYPE_TRANSLATE_AI
	} else {
		creditServiceType = pb.CreditServiceType_CREDIT_SERVICE_TYPE_TRANSLATE_WORD // 默认单词翻译
	}

	// 使用CreditFunTranslate包装翻译逻辑
	err := api.membershipService.CreditFunTranslate(ctx, creditServiceType, 0, func(xctx context.Context, sessionId string) error {
		// 移除特殊字符（如果没有使用术语库）
		if req.UseGlossary == nil || !*req.UseGlossary {
			content = removeSpecialWords(content, api.config.Translate.Text.Special.Words.NeedReplace.List)
			if content == "" {
				api.logger.Error("msg", "翻译内容包含过多特殊字符", "content", req.Text)
				return errors.System(errors.ErrorTypeBiz, "translate content contains too many special characters", nil)
			}
		}

		// 计算MD5，检查缓存
		md5Hash := utils.MD5(content)
		cachedLog, err := api.textTranslateService.FindByMD5(xctx, md5Hash)
		if err == nil && cachedLog != nil && cachedLog.Channel == req.Channel.String() {
			// 处理术语库
			if req.UseGlossary != nil && *req.UseGlossary && glossaryResult != nil {
				// 这里应该处理缓存结果中的术语库替换
				api.glossaryService.DealTranslationTextTargetContent(xctx, cachedLog.TargetContent, glossaryResult)
				if glossaryResult.OriginalTranslationAfterReplace != "" {
					cachedLog.TargetContent = glossaryResult.OriginalTranslationAfterReplace
				}
			}

			resp, err := api.buildTranslateResponse(cachedLog)
			if err != nil {
				api.logger.Error("msg", "构建缓存翻译响应失败", "error", err.Error())
				return err
			}
			translateResp = resp
			return nil
		}

		// 处理术语库（如果前面没有处理过）
		if req.UseGlossary != nil && *req.UseGlossary && glossaryResult == nil {
			api.logger.Info("msg", "处理术语库", "content", content)
			glossaryResult, err = api.glossaryService.ReplaceOriginalText(xctx, userId, content)
			if err != nil {
				api.logger.Error("msg", "处理术语库失败", "error", err.Error())
			} else if glossaryResult != nil && glossaryResult.UserSelectedTextAfterReplace != "" {
				content = glossaryResult.UserSelectedTextAfterReplace
			}
			api.logger.Info("msg", "术语库处理后内容", "content", content)
		}

		// 调用翻译服务
		translateLog, err := api.textTranslateService.TranslateWithPdfID(xctx, req.Channel,
			utils.GetEnumFromPtr(req.SourceLanguage, pkgi18n.GetDefaultTranslateSourceLanguage()),
			utils.GetEnumFromPtr(req.TargetLanguage, pkgi18n.GetDefaultTranslateTargetLanguage()),
			content,
			*req.PdfId,
			utils.GetBoolPtrValue(req.UseGlossary, false))
		if err != nil {
			api.logger.Error("msg", "调用翻译服务失败", "error", err.Error(), "content", content)
			return err
		}

		// 处理术语库
		if req.UseGlossary != nil && *req.UseGlossary && glossaryResult != nil {
			// 这里应该处理翻译结果中的术语库替换
			api.glossaryService.DealTranslationTextTargetContent(xctx, translateLog.TargetContent, glossaryResult)
			if glossaryResult != nil && glossaryResult.OriginalTranslationAfterReplace != "" {
				translateLog.TargetContent = glossaryResult.OriginalTranslationAfterReplace
			}
		}

		// 构建响应
		resp, err := api.buildTranslateResponse(translateLog)
		if err != nil {
			api.logger.Error("msg", "构建翻译响应失败", "error", err.Error())
			return err
		}

		translateResp = resp
		return nil
	}, true)
	if err != nil {
		c.Error(err)
		return
	}

	// 在积分扣除成功后再返回响应
	if translateResp != nil {
		response.Success(c, "success", translateResp)
	}
}

// buildTranslateResponse 构建翻译响应
func (api *TextTranslateAPI) buildTranslateResponse(translateLog *model.TextTranslateLog) (*translate.TranslateResponse, error) {
	if translateLog == nil {
		return nil, errors.System(errors.ErrorTypeInternal, "translate log is empty", nil)
	}

	api.logger.Info("msg", "构建翻译响应", "translateLog", translateLog)
	api.logger.Info("msg", "构建翻译响应", "translateLog.TargetContent", translateLog.TargetContent)

	// 解析目标内容
	var targetContent []string
	if err := json.Unmarshal([]byte(translateLog.TargetContent), &targetContent); err != nil {
		return nil, errors.Wrap(err, "failed to parse target content")
	}

	// // 解析网络内容
	var networkContent [][]string
	if translateLog.NetworkContent != "" {
		if err := json.Unmarshal([]byte(translateLog.NetworkContent), &networkContent); err != nil {
			api.logger.Warn("msg", "failed to parse network content", "error", err.Error())
			// 继续处理，不返回错误
		}
	}

	// 构建网络内容响应
	var networkContentResp []*translate.TargetResp
	for _, network := range networkContent {
		if len(network) > 0 {
			networkContentResp = append(networkContentResp, &translate.TargetResp{
				TargetContent: network,
			})
		}
	}

	// 构建响应
	resp := &translate.TranslateResponse{
		TargetContent: targetContent,
		RequestId:     &translateLog.RequestId,
		TargetResp:    networkContentResp,
	}

	return resp, nil
}
func (api *TextTranslateAPI) buildFullTextBlockTranslateResponse(translateLog *model.TextTranslateLog) (*translate.FullTextBlockTranslateResp, error) {
	if translateLog == nil {
		return nil, errors.Biz("translate log is empty")
	}

	// 解析目标内容
	var targetContent []string
	if err := json.Unmarshal([]byte(translateLog.TargetContent), &targetContent); err != nil {
		return nil, errors.Biz("failed to parse target content")
	}

	// 创建响应对象
	resp := &translate.FullTextBlockTranslateResp{
		Engine:  translateLog.Channel,
		Content: targetContent[0],
		Cached:  true,
	}
	return resp, nil
}

// removeSpecialWords 移除特殊字符
func removeSpecialWords(content string, specialWords []string) string {
	if content == "" {
		return ""
	}

	result := content
	for _, word := range specialWords {
		result = strings.ReplaceAll(result, word, "")
	}

	// 如果移除特殊字符后内容为空，返回空字符串
	if strings.TrimSpace(result) == "" {
		return ""
	}

	return result
}

// GetTranslateTabs 获取翻译标签
func (api *TextTranslateAPI) GetTranslateTabs(c *gin.Context) {
	response.Success(c, "success", &translate.GetTranslateTabsResponse{
		Tabs: []string{
			translate.TranslateChannel_GOOGLE.String(),
			translate.TranslateChannel_YOUDAO.String(),
			translate.TranslateChannel_AI.String()},
	})
}

// TranslateInternal 内部翻译
func (api *TextTranslateAPI) TranslateInternal(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "TextTranslateAPI.TranslateInternal")
	defer span.Finish()

	// 解析请求
	req := &translate.FullTextBlockTranslateReq{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("解析翻译请求失败", "error", err.Error())
		c.Error(err)
		return
	}
	// 计算MD5，检查缓存
	md5Hash := utils.MD5(req.Content)
	cachedLog, err := api.textTranslateService.FindByMD5(ctx, md5Hash)
	if err != nil {
		api.logger.Error("msg", "获取缓存翻译失败", "error", err.Error())
		c.Error(err)
		return
	}
	if cachedLog != nil {
		// 构建响应
		resp, err := api.buildFullTextBlockTranslateResponse(cachedLog)
		if err != nil {
			api.logger.Error("msg", "构建翻译响应失败", "error", err.Error())
			c.Error(err)
			return
		}
		response.Success(c, "success", resp)
		return
	}

	// 调用翻译服务
	translateLog, err := api.textTranslateService.TranslateInternal(ctx,
		req.Content,
		utils.GetEnumFromPtr(req.SourceLanguage, pkgi18n.GetDefaultTranslateSourceLanguage()),
		utils.GetEnumFromPtr(req.TargetLanguage, pkgi18n.GetDefaultTranslateTargetLanguage()),
		*req.PaperId,
		*req.PdfId,
		*req.PageId,
		*req.BlockId)
	if err != nil {
		api.logger.Error("msg", "调用翻译服务失败", "error", err.Error())
		c.Error(err)
		return
	}

	resp, err := api.buildFullTextBlockTranslateResponse(translateLog)
	if err != nil {
		api.logger.Error("msg", "构建翻译响应失败", "error", err.Error())
		c.Error(err)
		return
	}
	response.Success(c, "success", resp)
}

// AiTranslate 流式翻译
func (api *TextTranslateAPI) AiTranslate(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "TextTranslateAPI.AiTranslate")
	defer span.Finish()

	// 获取用户ID
	userID, exists := userContext.GetUserID(ctx)
	if !exists {
		api.logger.Warn("用户ID不存在")
		response.SSEError(c, "user info is empty")
		return
	}

	// 解析请求
	req := &translate.AiTranslateRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("解析翻译请求失败", "error", err.Error())
		response.SSEError(c, "translate request error")
		return
	}

	// 验证请求
	if len(req.Text) > api.config.Translate.Text.MaxLength {
		response.SSEError(c, "translate text too long")
		return
	}

	startTime := time.Now()

	// 设置SSE头信息
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Flush()

	// 检查是否为单词翻译模式
	wordPattern := regexp.MustCompile(`^\s*([A-Za-z]+)\s*$`)
	isWordMode := wordPattern.MatchString(req.Text)

	// 处理术语表
	var glossaryBean *model.GlossaryTranslateModel
	var glossaryList []*translate.GlossaryInfo
	if req.UseGlossary != nil && *req.UseGlossary {
		// 获取术语表匹配结果
		var err error
		glossaryBean, err = api.glossaryService.ReplaceOriginalText(ctx, userID, req.Text)
		if err != nil {
			api.logger.Warn("获取术语表匹配失败", "error", err.Error())
		}

		// 构建术语表列表
		if glossaryBean != nil && len(glossaryBean.RelInfos) > 0 {
			for _, relInfo := range glossaryBean.RelInfos {
				glossaryInfo := &translate.GlossaryInfo{
					OriginalText:    relInfo.OriginalText,
					TranslationText: relInfo.Translation,
				}
				glossaryList = append(glossaryList, glossaryInfo)
			}
		}
	}

	// 这里参照其他的方法，加入MD5校验，直接返回
	md5Hash := utils.MD5(req.Text)
	if md5Hash == "" {
		response.SSEError(c, "md5 hash is empty")
		return
	}

	// 检查缓存
	cachedLog, err := api.textTranslateService.FindByMD5(ctx, md5Hash)
	if err == nil && cachedLog != nil {
		// 处理术语库
		if req.UseGlossary != nil && *req.UseGlossary && glossaryBean != nil {
			api.glossaryService.DealTranslationTextTargetContent(ctx, cachedLog.TargetContent, glossaryBean)
			if glossaryBean.OriginalTranslationAfterReplace != "" {
				cachedLog.TargetContent = glossaryBean.OriginalTranslationAfterReplace
			}
		}
		// 构建新的响应结构
		translateResp, _ := api.buildTranslateResponse(cachedLog)
		translateResp.GlossaryList = glossaryList

		response.SSESuccess(c, "success", translateResp)
		return
	}

	if isWordMode {
		// 使用CreditFunTranslate包装单词翻译逻辑
		err := api.membershipService.CreditFunTranslate(ctx, pb.CreditServiceType_CREDIT_SERVICE_TYPE_TRANSLATE_WORD, 0, func(xctx context.Context, sessionId string) error {
			// 处理单词发音
			wordContent := req.Text
			translateResp, err := api.wordPronunciationService.GetTranslateResp(xctx, service.TranslateSourceYoudao, wordContent, *req.PdfId)
			if err != nil {
				return err
			}

			if translateResp != nil {
				// 处理术语库
				if req.UseGlossary != nil && *req.UseGlossary && glossaryBean != nil {
					api.glossaryService.DealTranslationText(xctx, translateResp, glossaryBean)
					if glossaryBean.OriginalTranslationAfterReplace != "" {
						translateResp.TargetContent = []string{glossaryBean.OriginalTranslationAfterReplace}
					}
				}
				// NOTICE: 这里返回术语库，前端显示有bug，暂不返回
				// translateResp.GlossaryList = glossaryList

				response.SSESuccess(c, "success", translateResp)
			}
			return nil
		}, true)
		if err != nil {
			response.SSEError(c, "translate service error")
			return
		}
		return
	}

	// 使用CreditFunTranslate包装AI翻译逻辑
	creditErr := api.membershipService.CreditFunTranslate(ctx, pb.CreditServiceType_CREDIT_SERVICE_TYPE_TRANSLATE_AI, 0, func(xctx context.Context, sessionId string) error {
		// 创建通道用于接收翻译结果
		errorChan := make(chan error)
		completeChan := make(chan struct{})

		// 获取用户上下文
		uc := userContext.GetUserContext(xctx)
		// 在goroutine中处理翻译
		userContext.RunAsyncWithUserContext(uc, func(bgCtx context.Context) {
			// 定义处理流式响应的回调函数
			resultCallback := func(content string) error {
				// 构建流式响应
				sseResp := &translate.TranslateResponse{
					SseData: &content,
				}
				response.SSESuccess(c, "success", sseResp)
				c.Writer.Flush()
				return nil
			}

			// 定义完成回调
			completeCallback := func(translatedText string) error {
				// 发送完整的翻译结果
				finalResp := &translate.TranslateResponse{
					TargetContent: []string{translatedText},
					GlossaryList:  glossaryList,
				}

				response.SSESuccess(c, "success", finalResp)
				c.Writer.Flush()
				return nil
			}

			// 调用服务层的流式翻译方法
			streamErr := api.textTranslateService.StreamTranslateWithCallback(
				bgCtx,
				req.Text,
				utils.GetEnumFromPtr(req.SourceLanguage, translate.TranslateLanguage_EN_US),
				utils.GetEnumFromPtr(req.TargetLanguage, translate.TranslateLanguage_ZH_CN),
				userID,
				*req.PdfId,
				utils.GetBoolPtrValue(req.UseGlossary, false),
				startTime,
				resultCallback,
				completeCallback,
			)

			if streamErr != nil {
				errorChan <- streamErr
				return
			}

			close(completeChan)
		})

		// 监听客户端连接断开
		clientGone := c.Request.Context().Done()

		// 发送翻译结果
		for {
			select {
			case err := <-errorChan:
				api.logger.Error("流式翻译失败", "error", err.Error())
				return err
			case <-completeChan:
				return nil
			case <-clientGone:
				// 客户端断开连接
				return nil
			}
		}
	}, true)
	if creditErr != nil {
		response.SSEError(c, "translate service error")
		return
	}
}
