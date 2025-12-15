package api

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	pb "github.com/yb2020/odoc-proto/gen/go/membership"
	"github.com/yb2020/odoc-proto/gen/go/translate"
	"github.com/yb2020/odoc/config"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/ratelimit"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	"github.com/yb2020/odoc/services/membership/interfaces"
	"github.com/yb2020/odoc/services/translate/service"
)

// OCRTextTranslateAPI OCR翻译API处理器
type OCRTextTranslateAPI struct {
	ocrTranslateService *service.OCRTranslateService
	membershipService   interfaces.IMembershipService
	config              *config.Config
	logger              logging.Logger
	tracer              opentracing.Tracer
	extractTextLimiter  ratelimit.RateLimiter
}

// NewOCRTextTranslateAPI 创建OCR翻译API处理器
func NewOCRTextTranslateAPI(
	config *config.Config,
	logger logging.Logger,
	tracer opentracing.Tracer,
	ocrTranslateService *service.OCRTranslateService,
	rateLimitService *ratelimit.RateLimiterService,
	membershipService interfaces.IMembershipService,
) *OCRTextTranslateAPI {
	var limiter ratelimit.RateLimiter
	if rateLimitService != nil {
		var err error
		limiter, err = rateLimitService.CreateLimiter(ratelimit.LimiterConfig{
			Type:       ratelimit.CounterLimiterType,
			KeyPrefix:  "ocr_extract_text",
			MaxRate:    2,
			TimeUnit:   ratelimit.Second,
			Dimension:  ratelimit.User,
			ExpireTime: 60,
		})
		// 检查限流
		if err != nil {
			logger.Error("Failed to create rate limiter", "error", err)
		}
	}

	return &OCRTextTranslateAPI{
		ocrTranslateService: ocrTranslateService,
		membershipService:   membershipService,
		config:              config,
		logger:              logger,
		tracer:              tracer,
		extractTextLimiter:  limiter,
	}
}

// ExtractText OCR提取文本
func (api *OCRTextTranslateAPI) ExtractText(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "OCRTextTranslateAPI.ExtractText")
	defer span.Finish()

	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())

	// 检查限流
	if api.extractTextLimiter != nil {
		// 创建特定于此用户的键
		key := ratelimit.GetLimiterKey("ocr_extract_text", ratelimit.User, fmt.Sprintf("%d", userId))
		result, err := api.extractTextLimiter.Allow(c.Request.Context(), key, 1)

		if err != nil {
			api.logger.Error("Rate limit check failed", "error", err)
		} else if !result.Allowed {
			c.Error(errors.System(errors.ErrorTypeBiz, "your request is too frequent, please try again later", nil))
			return
		}
	}

	// 解析请求
	req := &translate.OcrExtractTextRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("msg", "解析OCR提取文本请求失败", "error", err.Error())
		c.Error(err)
		return
	}
	resp, err := api.ocrTranslateService.ExtractText(ctx, req)
	if err != nil {
		api.logger.Error("OCR提取文本失败", "error", err.Error())
		c.Error(err)
		return
	}
	response.Success(c, "success", resp)
	return
}

// OCRTranslate OCR翻译
func (api *OCRTextTranslateAPI) OCRTranslate(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "OCRTranslateAPI.OCRTranslate")
	defer span.Finish()

	// 解析请求
	request := &translate.OcrTranslateRequest{}
	if err := transport.BindProto(c, request); err != nil {
		api.logger.Warn("msg", "解析OCR翻译请求失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 查看sourcelanguage和targetlanguage，没有的话默认中文和英文
	if request.SourceLanguage == nil {
		enUS := translate.TranslateLanguage_EN_US
		request.SourceLanguage = &enUS
	}
	if request.TargetLanguage == nil {
		zhCN := translate.TranslateLanguage_ZH_CN
		request.TargetLanguage = &zhCN
	}

	var finalResp *translate.TranslateResponse
	// 使用CreditFunTranslate包装OCR翻译逻辑
	err := api.membershipService.CreditFunTranslate(ctx, pb.CreditServiceType_CREDIT_SERVICE_TYPE_TRANSLATE_OCR, 0, func(xctx context.Context, sessionId string) error {
		// 调用服务
		resp, err := api.ocrTranslateService.OCRTranslate(xctx, request)
		if err != nil {
			api.logger.Error("OCR翻译失败", "error", err.Error())
			return err
		}
		finalResp = resp
		return nil
	}, true)
	if err != nil {
		c.Error(err)
		return
	}

	// 在积分扣除成功后再返回响应
	if finalResp != nil {
		response.Success(c, "success", finalResp)
	}
}

// TranslateCorrect 翻译纠错服务
func (api *TextTranslateAPI) TranslateCorrect(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "TextTranslateAPI.TranslateCorrect")
	defer span.Finish()

	// 解析请求
	req := &translate.TranslateCorrectInfoReq{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("解析翻译纠错请求失败", "error", err.Error())
		response.Error(c, "请求参数错误", nil)
		return
	}

	// 根据不同渠道处理纠错
	var err error
	switch req.Channel {
	case translate.TranslateChannel_GOOGLE, translate.TranslateChannel_YOUDAO, translate.TranslateChannel_AI:
		err = api.textTranslateService.TranslateCorrect(ctx, req.RequestId, req.Channel)
	default:
		api.logger.Warn("不支持的翻译渠道", "channel", req.Channel)
		c.Error(err)
		return
	}

	if err != nil {
		api.logger.Error("处理翻译纠错失败", "error", err.Error(), "requestId", req.RequestId, "channel", req.Channel)
		c.Error(err)
		return
	}

	response.SuccessNoData(c, "thanks for your feedback")
}
