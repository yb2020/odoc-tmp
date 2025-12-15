package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"net/http"
	"strings"
	"time"

	timeUtil "github.com/jinzhu/now"
	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc-proto/gen/go/translate"
	"github.com/yb2020/odoc/config"
	ocrExternalApi "github.com/yb2020/odoc/external/ocr/api"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/ratelimit"
	"github.com/yb2020/odoc/pkg/utils"
	"github.com/yb2020/odoc/services/translate/dao"
	"github.com/yb2020/odoc/services/translate/model"
)

// 常量定义
const (
	OCRTranslateStatusFailed  = 0
	OCRTranslateStatusSuccess = 1
)

// OCRExtractTextRequest OCR提取文本请求
type OCRExtractTextRequest struct {
	Images []string `json:"images"`
}

// OCRExtractTextResponse OCR提取文本响应
type OCRExtractTextResponse struct {
	Results [][]OCRExtractTextResponseItem `json:"results"`
}

// OCRExtractTextResponseItem OCR提取文本响应项
type OCRExtractTextResponseItem struct {
	Text string `json:"text"`
}

// OCRTranslateService OCR翻译服务
type OCRTranslateService struct {
	ocrTranslateDAO      dao.OCRTranslateDAO
	textTranslateService *TextTranslateService
	glossaryService      *GlossaryService
	rateLimiterService   *ratelimit.RateLimiterService
	imageOCRApiService   *ocrExternalApi.ImageOCRApiService
	config               *config.Config
	logger               logging.Logger
	tracer               opentracing.Tracer
	httpClient           *http.Client
}

// NewOCRTranslateService 创建OCR翻译服务
func NewOCRTranslateService(
	config *config.Config,
	logger logging.Logger,
	tracer opentracing.Tracer,
	ocrTranslateDAO dao.OCRTranslateDAO,
	textTranslateService *TextTranslateService,
	glossaryService *GlossaryService,
	rateLimiterService *ratelimit.RateLimiterService,
	imageOCRApiService *ocrExternalApi.ImageOCRApiService,
) *OCRTranslateService {
	// 创建带超时的HTTP客户端
	httpClient := &http.Client{
		Timeout: 25 * time.Second,
	}

	return &OCRTranslateService{
		ocrTranslateDAO:      ocrTranslateDAO,
		textTranslateService: textTranslateService,
		glossaryService:      glossaryService,
		rateLimiterService:   rateLimiterService,
		imageOCRApiService:   imageOCRApiService,
		config:               config,
		logger:               logger,
		tracer:               tracer,
		httpClient:           httpClient,
	}
}

// ExtractText 从图片中提取文本
func (s *OCRTranslateService) ExtractText(ctx context.Context, req *translate.OcrExtractTextRequest) (*translate.OcrExtractTextResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OCRTranslateService.ExtractText")
	defer span.Finish()

	// 获取用户ID
	userID, ok := userContext.GetUserID(ctx)
	if !ok {
		return nil, errors.Biz("please login first")
	}
	imgFileByte, err := utils.HandleImageToBytes(req.PicBase64)
	if err != nil {
		s.logger.Error("处理图片失败", "error", err)
		return nil, errors.Biz("handle image failed")
	}
	// 调用OCR服务提取文本
	ocrText, err := s.imageOCRApiService.ExtractImgToText(ctx, imgFileByte)
	if err != nil {
		s.logger.Error("OCR提取文本失败", "error", err)
		return nil, errors.Biz("ocr extract text failed")
	}

	// 替换换行符为空格
	ocrText = strings.ReplaceAll(ocrText, "\n", " ")
	s.logger.Info("OCR提取文本成功", "text", ocrText)

	// 保存OCR翻译记录
	if ocrText != "" {
		id := idgen.GenerateUUID()
		log := &model.OCRTranslateLog{
			UserId:      userID,
			RequestId:   fmt.Sprintf("%d", id),
			ImageBase64: req.PicBase64,
			OCRText:     ocrText,
			Status:      OCRTranslateStatusSuccess,
			CostMS:      0, // 这里可以记录处理时间
		}

		err = s.ocrTranslateDAO.Save(ctx, log)
		if err != nil {
			s.logger.Error("保存OCR翻译记录失败", "error", err)
			// 继续处理，不返回错误
		}
	}

	// 构建响应
	response := &translate.OcrExtractTextResponse{
		Text: ocrText,
	}

	return response, nil
}

// OCRTranslate OCR翻译
func (s *OCRTranslateService) OCRTranslate(ctx context.Context, req *translate.OcrTranslateRequest) (*translate.TranslateResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OCRTranslateService.OCRTranslate")
	defer span.Finish()

	// 获取用户ID
	userID, ok := userContext.GetUserID(ctx)
	if !ok {
		return nil, errors.Biz("please login first")
	}
	// 检查用户OCR翻译次数限制
	imgFileByte, err := utils.HandleImageToBytes(req.PicBase64)
	if err != nil {
		s.logger.Error("处理图片失败", "error", err)
		return nil, errors.Biz("handle image failed")
	}
	// 调用OCR服务提取文本
	ocrText, err := s.imageOCRApiService.ExtractImgToText(ctx, imgFileByte)
	if err != nil {
		s.logger.Error("OCR提取文本失败", "error", err)
		return nil, errors.Biz("ocr extract text failed")
	}

	s.logger.Info("OCR提取文本成功", "text", ocrText)

	if ocrText == "" {
		return nil, errors.Biz("ocr identify failed or image too large, please try again later")
	}

	// 保存OCR翻译记录
	id := idgen.GenerateUUID()
	log := &model.OCRTranslateLog{
		UserId:      userID,
		RequestId:   fmt.Sprintf("%d", id),
		ImageBase64: req.PicBase64,
		OCRText:     ocrText,
		Status:      OCRTranslateStatusSuccess,
		CostMS:      0, // 这里可以记录处理时间
	}

	if err := s.ocrTranslateDAO.Save(ctx, log); err != nil {
		s.logger.Error("保存OCR翻译记录失败", "error", err)
		// 继续处理，不返回错误
	}

	// 如果不是AI翻译渠道，则调用文本翻译服务
	if req.Channel != translate.TranslateChannel_AI {
		// 设置源语言和目标语言
		sourceLanguage := translate.TranslateLanguage_EN_US
		targetLanguage := translate.TranslateLanguage_ZH_CN

		if req.SourceLanguage != nil {
			sourceLanguage = *req.SourceLanguage
		}

		if req.TargetLanguage != nil {
			targetLanguage = *req.TargetLanguage
		}

		// 调用文本翻译服务
		translateLog, err := s.textTranslateService.TranslateWithPdfID(
			ctx,
			req.Channel,
			sourceLanguage,
			targetLanguage,
			ocrText,
			*req.PdfId,
			req.UseGlossary != nil && *req.UseGlossary,
		)

		if err != nil {
			s.logger.Error("文本翻译失败", "error", err)
			return nil, errors.Biz("text translate failed")
		}

		// 构建翻译响应
		resp, err := s.buildTranslateResponse(translateLog)
		if err != nil {
			return nil, err
		}

		// 添加OCR提取的文本
		resp.OcrExtractText = &ocrText

		return resp, nil
	}

	// 对于AI翻译渠道，只返回OCR提取的文本
	resp := &translate.TranslateResponse{
		OcrExtractText: &ocrText,
	}

	return resp, nil
}

// GetRemainingCount 获取用户剩余OCR翻译次数
func (s *OCRTranslateService) GetRemainingCount(ctx context.Context) (int, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OCRTranslateService.GetRemainingCount")
	defer span.Finish()

	// 获取用户ID
	userID, ok := userContext.GetUserID(ctx)
	if !ok {
		return 0, errors.Biz("please login first")
	}

	// 获取用户本周OCR翻译记录
	startOfWeek := timeUtil.BeginningOfWeek()
	endOfWeek := timeUtil.EndOfWeek()

	logs, err := s.ocrTranslateDAO.GetUserOCRHistoryBetween(ctx, userID, startOfWeek, endOfWeek)
	if err != nil {
		s.logger.Error("获取用户OCR翻译历史失败", "error", err)
		return 0, errors.Biz("get user ocr translate history failed")
	}

	// TODO: 实现VIP用户检查逻辑
	maxLimit := 10 // 默认限制，实际应从配置或VIP服务获取
	remainCount := maxLimit - len(logs)
	if remainCount < 0 {
		remainCount = 0
	}

	return remainCount, nil
}

// 内部辅助方法

// decodeBase64ToImage 将base64编码的图片转换为image.Image
func (s *OCRTranslateService) decodeBase64ToImage(base64Str string) (image.Image, error) {
	// 分离base64头部和数据部分
	parts := strings.Split(base64Str, ",")
	var base64Data string
	if len(parts) > 1 {
		base64Data = parts[1]
	} else {
		base64Data = parts[0]
	}

	// 解码base64数据
	imgData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, err
	}

	// 读取图片
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}

	return img, nil
}

// buildTranslateResponse 构建翻译响应
func (s *OCRTranslateService) buildTranslateResponse(translateLog *model.TextTranslateLog) (*translate.TranslateResponse, error) {
	if translateLog == nil {
		return nil, errors.System(errors.ErrorTypeInternal, "translate log is empty", nil)
	}

	// 解析目标内容
	var targetContentList []string
	if translateLog.TargetContent != "" {
		targetContentList = []string{translateLog.TargetContent}
	}

	// 构建响应
	resp := &translate.TranslateResponse{
		TargetContent: targetContentList,
		RequestId:     &translateLog.RequestId,
	}

	return resp, nil
}
