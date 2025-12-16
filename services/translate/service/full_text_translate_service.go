package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/cache"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/distlock"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/http_client"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/metrics"
	"github.com/yb2020/odoc/pkg/utils"
	pb "github.com/yb2020/odoc/proto/gen/go/translate"
	docService "github.com/yb2020/odoc/services/doc/service"
	membershipInterface "github.com/yb2020/odoc/services/membership/interfaces"
	noteInterface "github.com/yb2020/odoc/services/note/interfaces"
	ossConstant "github.com/yb2020/odoc/services/oss/constant"
	ossService "github.com/yb2020/odoc/services/oss/service"
	pdfService "github.com/yb2020/odoc/services/pdf/service"
	"github.com/yb2020/odoc/services/translate/dao"
	"github.com/yb2020/odoc/services/translate/dto"
	"github.com/yb2020/odoc/services/translate/model"
)

// FullTextTranslateService 全文翻译自研服务
type FullTextTranslateService struct {
	fullTextTranslateDAO dao.FullTextTranslateDAO
	membershipService    membershipInterface.IMembershipService
	config               *config.Config
	logger               logging.Logger
	tracer               opentracing.Tracer
	httpClient           http_client.HttpClient
	cacheClient          cache.Cache
	metrics              *metrics.Metrics
	lockTemplate         *distlock.LockTemplate
	paperNoteService     noteInterface.IPaperNoteService
	paperPdfService      *pdfService.PaperPdfService
	userDocService       *docService.UserDocService
	ossService           ossService.OssServiceInterface
}

// NewFullTextTranslateService 创建全文翻译自研服务
func NewFullTextTranslateService(
	cfg *config.Config,
	logger logging.Logger,
	tracer opentracing.Tracer,
	cache cache.Cache,
	httpClient http_client.HttpClient,
	metrics *metrics.Metrics,
	translateHistoryDAO dao.FullTextTranslateDAO,
	lockTemplate *distlock.LockTemplate,
	paperNoteService noteInterface.IPaperNoteService,
	paperPdfService *pdfService.PaperPdfService,
	userDocService *docService.UserDocService,
	ossService ossService.OssServiceInterface,
	membershipService membershipInterface.IMembershipService,
) *FullTextTranslateService {
	// 从配置中读取参数
	// 注意：这里假设配置结构中已有这些字段，实际使用时需要确保配置结构匹配
	return &FullTextTranslateService{
		fullTextTranslateDAO: translateHistoryDAO,
		config:               cfg,
		logger:               logger,
		tracer:               opentracing.GlobalTracer(),
		httpClient:           httpClient,
		cacheClient:          cache,
		metrics:              metrics,
		lockTemplate:         lockTemplate,
		paperNoteService:     paperNoteService,
		paperPdfService:      paperPdfService,
		userDocService:       userDocService,
		ossService:           ossService,
		membershipService:    membershipService,
	}
}

const fullTextTranslateProgressCacheKey = "full_text_translate_progress_%s"

// const fullTextTranslateProgressCacheKey = "apirun.%s"

// GetRightInfo 获取全文翻译权益信息
func (s *FullTextTranslateService) GetRightInfo(ctx context.Context, userId string) (*pb.GetFullTextTranslateRightInfoResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FullTextTranslateService.GetRightInfo")
	defer span.Finish()

	response := &pb.GetFullTextTranslateRightInfoResponse{}

	// 根据当前语言环境设置规则描述
	// if s.i18n.IsEn() {
	// 	response.RuleDesc = s.config.Translate.FullTextTranslate.EnTemplate
	// } else {
	// 	response.RuleDesc = s.config.Translate.FullTextTranslate.RuleDescTemplate
	// }

	// 设置剩余数量，如果小于0则设置为-1
	// if remainTicketCount < 0 {
	// 	response.RemainCount = -1
	// } else {
	// 	response.RemainCount = int32(remainTicketCount)
	// }

	//先mock权益
	remainCount := int32(-1)
	response.RuleDesc = "xxxxxx规则"
	response.RemainCount = &remainCount

	return response, nil
}

// Translate 翻译文档
func (s *FullTextTranslateService) Translate(ctx context.Context, req *pb.FullTextTranslateRequest, sessionId string) (*pb.FullTextTranslateResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FullTextTranslateService.Translate")
	defer span.Finish()

	// 获取用户ID
	userId := *req.UserId
	if userId == "" {
		// 从上下文中获取用户ID
		userId, _ = userContext.GetUserID(ctx)
	}

	// 获取PDF ID
	pdfId := req.PdfId

	// 检查用户是否有正在翻译的记录
	translatingHistories, err := s.fullTextTranslateDAO.FindByPdfIdAndUserIdAndStatusWithinDays(ctx, pdfId, userId,
		pb.FullTranslateFlowStatus_TRANSLATING, s.config.Translate.FullTextTranslate.HistoryVisibleDays)
	if err != nil {
		s.logger.Error("查询正在翻译的历史记录失败", "error", err.Error(), "pdfId", pdfId, "userId", userId)
		return nil, errors.Biz("query translating history failed")
	}
	if len(translatingHistories) > 0 {
		return nil, errors.Biz("translation is in progress, please refresh the page and wait for the translation result")
	}

	pdfInfo, err := s.paperPdfService.GetById(ctx, pdfId)
	if err != nil {
		s.logger.Error("获取PDF信息失败", "error", err.Error(), "pdfId", pdfId)
		return nil, errors.Biz("get pdf info failed")
	}
	if pdfInfo == nil {
		return nil, errors.Biz("pdf info is nil")
	}
	userDoc, err := s.userDocService.GetByUserIdAndPdfId(ctx, userId, pdfId)
	if err != nil {
		return nil, errors.Biz("get doc name failed")
	}
	fileSHA256 := pdfInfo.FileSHA256
	if userDoc == nil {
		return nil, errors.Biz("doc name is nil")
	}
	docName := ""
	if userDoc.PaperTitle == "" {
		docName = userDoc.DocName
	} else {
		docName = userDoc.PaperTitle
	}

	// 检查用户是否有翻译完成的记录
	finishedHistories, err := s.fullTextTranslateDAO.FindByPdfIdAndUserIdAndStatusWithinDays(
		ctx, pdfId, userId, pb.FullTranslateFlowStatus_TRANSLATE_FINISHED, s.config.Translate.FullTextTranslate.HistoryVisibleDays)
	if err != nil {
		s.logger.Error("查询已完成翻译的历史记录失败", "error", err.Error(), "pdfId", pdfId, "userId", userId)
		return nil, errors.System(errors.ErrorTypeDatabase, "query finished history failed", err)
	}
	if len(finishedHistories) > 0 {
		return nil, errors.Biz("translation is completed, please refresh the page and get the translation result")
	}

	// 获取已完成的翻译记录
	finishedHistory, err := s.fullTextTranslateDAO.FindLatestByFileSHA256AndStatus(ctx, fileSHA256, pb.FullTranslateFlowStatus_TRANSLATE_FINISHED)
	if err != nil {
		s.logger.Error("查询最新已完成翻译的历史记录失败", "error", err.Error(), "fileSHA256", fileSHA256)
		return nil, errors.Biz("query latest finished history failed")
	}

	// 分布式锁，防止重复翻译
	lockKey := fmt.Sprintf("full_text_translate_lock_%d", pdfId)
	lockInfo := s.lockTemplate.LockWithRetry(
		lockKey,
		30000,
		50000,
		100,
	)

	defer s.lockTemplate.ReleaseLock(lockInfo)

	if lockInfo == nil {
		s.logger.Error("获取分布式锁失败", "error", err, "pdfId", pdfId)
		return nil, errors.Biz("you are translating this pdf, please wait")
	}

	// 如果已有完成的翻译记录，直接使用
	if finishedHistory != nil {
		s.logger.Info("PDF已有翻译完成的记录", "pdfId", pdfId)

		// 创建新的翻译历史记录
		history := &model.FullTextTranslate{
			UserId:           userId,
			FileSHA256:       fileSHA256,
			DocName:          docName,
			SourcePdfId:      pdfId,
			Status:           pb.FullTranslateFlowStatus_TRANSLATE_FINISHED,
			SourceLanguage:   finishedHistory.SourceLanguage,
			TargetLanguage:   finishedHistory.TargetLanguage,
			Alignment:        finishedHistory.Alignment,
			TargetBucketName: finishedHistory.TargetBucketName,
			TargetObjectKey:  finishedHistory.TargetObjectKey,
			SessionId:        finishedHistory.SessionId,
		}

		err = s.fullTextTranslateDAO.Save(ctx, history)
		if err != nil {
			s.logger.Error("保存翻译历史记录失败", "error", err.Error())
			return nil, errors.Biz("保存翻译历史记录失败")
		}

		return &pb.FullTextTranslateResponse{
			ErrorCode: "0",
		}, nil
	}

	// 检查是否有正在翻译的记录（其他用户）
	translatingHistory, err := s.fullTextTranslateDAO.FindLatestByFileSHA256AndStatus(ctx, fileSHA256, pb.FullTranslateFlowStatus_TRANSLATING)
	if err != nil {
		s.logger.Error("查询最新正在翻译的历史记录失败", "error", err.Error(), "fileSHA256", fileSHA256)
		return nil, errors.Biz("query latest translating history failed")
	}

	if translatingHistory != nil {
		s.logger.Info("PDF正在被翻译", "pdfId", pdfId)

		// 创建新的翻译历史记录
		history := &model.FullTextTranslate{
			UserId:         userId,
			SourcePdfId:    pdfId,
			DocName:        docName,
			FileSHA256:     fileSHA256,
			Status:         pb.FullTranslateFlowStatus_TRANSLATING,
			SourceLanguage: pb.TranslateLanguage_EN_US.String(),
			TargetLanguage: pb.TranslateLanguage_ZH_CN.String(),
			FlowNumber:     translatingHistory.FlowNumber,
			SessionId:      sessionId,
		}

		err = s.fullTextTranslateDAO.Save(ctx, history)
		if err != nil {
			s.logger.Error("保存翻译历史记录失败", "error", err.Error())
			return nil, errors.Biz("保存翻译历史记录失败")
		}

		return &pb.FullTextTranslateResponse{
			ErrorCode: "0",
		}, nil
	}

	// 检查文件语言
	// if s.langCheckSwitch {
	// 	langResult, err := s.checkFileLang(ctx, pdfUrl)
	// 	if err != nil {
	// 		s.logger.Error("检查文件语言失败", "error", err.Error(), "url", pdfUrl)
	// 		return nil, errors.Biz("文档语言检测服务异常,请稍后重试")
	// 	}

	// 	if langResult.Language != s.allowedLang {
	// 		return nil, errors.Biz("当前PDF不是英文文档,暂不支持全文翻译")
	// 	}
	// }

	// 检查文件大小和页数
	// 这里需要根据实际情况实现
	// ...

	// 获取PDF信息  这里不需要提前获取url进行参数传递，每次用的时候在获取最新的地址就好了
	// var pdfUrl string
	// // 否则，通过 pdfId 获取 URL
	// urlPtr, err := s.paperPdfService.GetPdfUrlById(ctx, pdfId, 60)
	// if err != nil {
	// 	s.logger.Error("获取PDF文件地址失败", "error", err)
	// 	return nil, errors.System(errors.ErrorTypeBiz, "get pdf url failed", err)
	// }
	// if urlPtr == nil {
	// 	return nil, errors.System(errors.ErrorTypeBiz, "get pdf url failed: url is nil", nil)
	// }
	// pdfUrl = *urlPtr

	// // 检查 pdfUrl 是否为空
	// if pdfUrl == "" {
	// 	return nil, errors.System(errors.ErrorTypeBiz, "pdf url is empty", nil)
	// }

	// 创建新的翻译历史记录
	history := &model.FullTextTranslate{
		UserId:         userId,
		SourcePdfId:    pdfId,
		DocName:        docName,
		FileSHA256:     fileSHA256,
		Status:         pb.FullTranslateFlowStatus_TRANSLATING,
		SourceLanguage: pb.TranslateLanguage_EN_US.String(),
		TargetLanguage: pb.TranslateLanguage_ZH_CN.String(),
		FlowNumber:     idgen.GenerateUUID(),
		SessionId:      sessionId,
	}

	err = s.fullTextTranslateDAO.Save(ctx, history)
	if err != nil {
		s.logger.Error("保存翻译历史记录失败", "error", err.Error())
		return nil, errors.Biz("保存翻译历史记录失败")
	}

	// 异步调用翻译服务
	// 使用 UserSafeGoroutine 确保即使请求上下文结束，异步操作也能继续执行
	uc := userContext.GetUserContext(ctx)
	userContext.RunAsyncWithUserContext(uc, func(newCtx context.Context) {
		s.callTranslateService(newCtx, pdfId, history)
	})

	return &pb.FullTextTranslateResponse{
		ErrorCode: "0",
	}, nil
}

// callTranslateService 调用翻译服务
func (s *FullTextTranslateService) callTranslateService(ctx context.Context, pdfId string, history *model.FullTextTranslate) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FullTextTranslateService.callTranslateService")
	defer span.Finish()

	startTime := time.Now()
	result, err := s.doCallTranslateService(ctx, pdfId, &s.config.Translate.FullTextTranslate.UseWatermark)
	s.logger.Info("调用翻译服务完成", "cost", time.Since(startTime), "pdfId", pdfId)

	if err != nil || result == nil || result.Token == "" || s.config.Translate.FullTextTranslate.MockFailSwitch {
		errMsg := "翻译服务错误"
		if err != nil {
			errMsg = err.Error()
		} else if result != nil && result.Message != "" {
			errMsg = result.Message
		}

		s.handleTranslateFail(ctx, history, errMsg)
		return
	}

	// 保存流水号
	flowNumber := result.Token
	history.FlowNumber = flowNumber

	err = s.fullTextTranslateDAO.ModifyExcludeNull(ctx, history)
	if err != nil {
		s.logger.Error("更新翻译历史记录失败", "error", err.Error(), "historyId", history.Id)
	}

	// 设置Redis缓存，用于跟踪进度
	progressCacheKey := fmt.Sprintf(fullTextTranslateProgressCacheKey, flowNumber)
	err = s.cacheClient.Set(ctx, progressCacheKey, 0, time.Duration(s.config.Translate.FullTextTranslate.TranslateTimeOut)*time.Second)
	if err != nil {
		s.logger.Error("设置翻译进度缓存失败", "error", err.Error(), "flowNumber", flowNumber)
	}
}

// doCallTranslateService 实际调用翻译服务
func (s *FullTextTranslateService) doCallTranslateService(ctx context.Context, pdfId string, useWatermark *bool) (*dto.FullTranslateUploadResult, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FullTextTranslateService.doCallTranslateService")
	defer span.Finish()

	if useWatermark == nil {
		useWatermark = &s.config.Translate.FullTextTranslate.UseWatermark
	}
	pdf, err := s.paperPdfService.GetById(ctx, pdfId)
	if err != nil {
		s.logger.Error("获取论文PDF失败", "error", err)
		return nil, errors.Biz("pdf.paper_pdf.errors.get_failed")
	}
	if pdf == nil {
		return nil, errors.Biz("pdf.paper_pdf.errors.not_found")
	}
	// 构建请求URL
	requestURL := fmt.Sprintf("%s%s?pdf_id=%d&file_sha256=%s&watermark=%t",
		s.config.Translate.FullTextTranslate.TranslateBaseURL,
		s.config.Translate.FullTextTranslate.TranslateServiceURI,
		pdfId,
		pdf.FileSHA256,
		*useWatermark)

	// 设置请求头
	headers := map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   s.config.OAuth2.AppID,
	}

	// 使用已有的HTTP客户端组件发送请求
	body, err := s.httpClient.Post(requestURL, nil, headers)
	if err != nil {
		s.logger.Error("发送HTTP请求失败", "error", err.Error(), "url", requestURL)
		return nil, err
	}

	// 解析响应
	var result dto.FullTranslateUploadResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		s.logger.Error("解析HTTP响应失败", "error", err.Error(), "body", string(body))
		return nil, err
	}

	s.logger.Info("调用翻译服务响应", "result", result)
	return &result, nil
}

// GetTranslateStatus 获取翻译状态
func (s *FullTextTranslateService) GetTranslateStatus(ctx context.Context, req *pb.GetTranslateStatusReq) (*pb.GetTranslateStatusResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FullTextTranslateService.GetTranslateStatus")
	defer span.Finish()

	// 获取用户ID
	var userId string
	if req.UserId != nil && *req.UserId != "" {
		userId = *req.UserId
	} else {
		// 从上下文中获取用户ID
		userId, _ = userContext.GetUserID(ctx)
	}

	// 获取PDF ID
	pdfId := req.PdfId
	if pdfId == "" {
		return nil, errors.Biz("PDF ID can't be empty")
	}

	// 获取用户的翻译历史记录
	// 这里需要根据实际情况调整，例如考虑VIP用户可以查看更长时间的历史记录
	histories, err := s.fullTextTranslateDAO.FindByPdfIdAndUserIdWithinDays(ctx, pdfId, userId,
		s.config.Translate.FullTextTranslate.HistoryVisibleDays)
	if err != nil {
		s.logger.Error("查询翻译历史记录失败", "error", err.Error(), "pdfId", pdfId, "userId", userId)
		return nil, errors.Biz("query history failed")
	}

	if len(histories) == 0 {
		emptyString := ""
		return &pb.GetTranslateStatusResponse{
			Alignment:          &emptyString,
			Status:             pb.FullTranslateFlowStatus_WITHOUT_TRANSLATE_HISTORY,
			TranslationFileUrl: &emptyString,
			ProgressPercent:    &emptyString,
		}, nil
	}

	// 获取最新的翻译历史记录
	history := histories[0]

	// 如果创建时间太近，直接返回翻译中状态
	if time.Since(history.CreatedAt) < time.Duration(s.config.Translate.FullTextTranslate.DuplicateDelaySeconds)*time.Second {
		if history.Status == pb.FullTranslateFlowStatus_TRANSLATE_FINISHED {
			return s.returnResponse(ctx, history)
		}
		return &pb.GetTranslateStatusResponse{
			Status: history.Status,
		}, nil
	}

	// 如果已经翻译完成，直接返回结果
	if history.Status == pb.FullTranslateFlowStatus_TRANSLATE_FINISHED {
		return s.returnResponse(ctx, history)
	}

	// 如果翻译失败，直接返回结果
	if history.Status == pb.FullTranslateFlowStatus_TRANSLATE_FAIL {
		return &pb.GetTranslateStatusResponse{
			Status: pb.FullTranslateFlowStatus_TRANSLATE_FAIL,
		}, nil
	}

	// 查询是否有其他用户的翻译完成记录
	finishedHistory, err := s.fullTextTranslateDAO.FindLatestByFileSHA256AndStatus(ctx, history.FileSHA256, pb.FullTranslateFlowStatus_TRANSLATE_FINISHED)
	if err != nil {
		s.logger.Error("查询最新已完成翻译的历史记录失败", "error", err.Error(), "pdfId", pdfId)
		return nil, errors.Biz("query latest finished history failed")
	}

	if finishedHistory != nil {
		// 更新用户的记录状态
		history.Status = finishedHistory.Status
		history.TargetObjectKey = finishedHistory.TargetObjectKey
		history.Alignment = finishedHistory.Alignment
		err = s.fullTextTranslateDAO.Save(ctx, history)
		if err != nil {
			s.logger.Error("更新翻译历史记录失败", "error", err.Error(), "historyId", history.Id)
		}

		return s.returnResponse(ctx, history)
	}

	// 获取流水号
	flowNumber := history.FlowNumber

	if flowNumber == "" || flowNumber == "0" {
		return &pb.GetTranslateStatusResponse{
			ProgressPercent: stringPtr("0"),
			Status:          pb.FullTranslateFlowStatus_WITHOUT_TRANSLATE_HISTORY,
		}, nil
	}

	// 从Redis获取进度
	progressCacheKey := fmt.Sprintf(fullTextTranslateProgressCacheKey, flowNumber)
	var progressValue int
	found, err := s.cacheClient.Get(ctx, progressCacheKey, &progressValue)
	if err != nil || !found {
		s.logger.Error("获取翻译进度失败", "error", err, "flowNumber", flowNumber)
		s.handleTranslateFail(ctx, history, "自研服务超时")
		return &pb.GetTranslateStatusResponse{
			Status: pb.FullTranslateFlowStatus_WITHOUT_TRANSLATE_HISTORY,
		}, nil
	}

	return &pb.GetTranslateStatusResponse{
		ProgressPercent: stringPtr(fmt.Sprintf("%v", progressValue)),
		Status:          pb.FullTranslateFlowStatus_TRANSLATING,
	}, nil
}

func (s *FullTextTranslateService) returnResponse(ctx context.Context, history *model.FullTextTranslate) (*pb.GetTranslateStatusResponse, error) {
	if history.Status == pb.FullTranslateFlowStatus_TRANSLATE_FINISHED {
		duration := 120
		// 获取PDF文件地址
		fileUrl, err := s.ossService.GetFileTemporaryURL(ctx, ossConstant.BucketTypeToEnum(s.config, history.TargetBucketName), history.TargetObjectKey, utils.GetIntPtrValue(&duration, 60))
		if err != nil {
			return nil, err
		}
		return &pb.GetTranslateStatusResponse{
			Status:             pb.FullTranslateFlowStatus_TRANSLATE_FINISHED,
			TranslationFileUrl: &fileUrl,
			Alignment:          &history.Alignment,
		}, nil
	}
	return &pb.GetTranslateStatusResponse{
		Status: pb.FullTranslateFlowStatus_WITHOUT_TRANSLATE_HISTORY,
	}, nil
}

// UpdateProgress 更新翻译进度
func (s *FullTextTranslateService) UpdateProgress(ctx context.Context, result *dto.FullTranslateProgressResult, updateType int) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FullTextTranslateService.UpdateProgress")
	defer span.Finish()

	if result.Progress >= 100 {
		// 查询历史记录
		historyList, err := s.fullTextTranslateDAO.FindByFlowNumber(ctx, result.Token)
		if err != nil {
			s.logger.Error("查询翻译历史记录失败", "error", err.Error())
			return errors.Biz("query history failed")
		}

		if historyList == nil {
			s.logger.Error("未找到对应的翻译历史记录", "token", result.Token)
			return errors.Biz("query history failed")
		}

		for _, history := range historyList {
			//过滤非翻译中的数据
			if history.Status != pb.FullTranslateFlowStatus_TRANSLATING {
				continue
			}
			if updateType == 1 {
				// 翻译完成
				startTime := history.CreatedAt
				costTime := time.Since(startTime).Milliseconds()
				// 记录指标
				s.metrics.ObserveRequestLatency("POST", "full_text_translate", float64(costTime)/1000.0)
				if len(result.Alignment) > 0 {
					history.Alignment = string(result.Alignment)
				}
				err = s.fullTextTranslateDAO.ModifyExcludeNull(ctx, &history)
				if err != nil {
					s.logger.Error("更新翻译历史记录失败", "error", err.Error(), "historyId", history.Id)
					return errors.Biz("update history failed")
				}
			} else {
				history.Status = pb.FullTranslateFlowStatus_TRANSLATE_FINISHED
				history.TargetBucketName = result.Data.BucketName
				history.TargetObjectKey = result.Data.ObjectKey
				err = s.fullTextTranslateDAO.ModifyExcludeNull(ctx, &history)
				if err != nil {
					s.logger.Error("更新翻译历史记录失败", "error", err.Error(), "historyId", history.Id)
					return errors.Biz("update history failed")
				}
				//设置用户id
				uc := userContext.NewUserContext().SetUserID(history.UserId)
				ctx = uc.ToContext(ctx)

				// 根据session确认消息
				err = s.membershipService.ConfirmCreditFun(ctx, history.SessionId)
				if err != nil {
					s.logger.Error("确认翻译失败", "error", err.Error(), "sessionId", history.SessionId)
					return errors.Biz("confirm translation failed")
				}
			}

		}

	} else if result.Progress < 0 {
		// 查询历史记录
		historyList, err := s.fullTextTranslateDAO.FindByFlowNumber(ctx, result.Token)
		if err != nil {
			s.logger.Error("查询翻译历史记录失败", "error", err.Error())
			return errors.Biz("query history failed")
		}

		if historyList == nil {
			s.logger.Error("未找到对应的翻译历史记录", "token", result.Token)
			return errors.Biz("query history failed")
		}
		// 翻译失败
		s.metrics.IncrementErrorCount("POST", "full_text_translate", "translation_failed")
		for _, history := range historyList {
			//过滤非翻译中的数据
			if history.Status != pb.FullTranslateFlowStatus_TRANSLATING {
				continue
			}
			s.handleTranslateFail(ctx, &history, result.Message)
		}
		return nil
	} else {
		// 更新进度
		progressCacheKey := fmt.Sprintf(fullTextTranslateProgressCacheKey, result.Token)
		err := s.cacheClient.Set(ctx, progressCacheKey, result.Progress, time.Duration(s.config.Translate.FullTextTranslate.TranslateTimeOut)*time.Second)
		if err != nil {
			s.logger.Error("更新翻译进度失败", "error", err.Error(), "token", result.Token, "progress", result.Progress)
			return errors.Biz("update progress failed")
		}
	}

	return nil
}

// ReTranslate 重新翻译
func (s *FullTextTranslateService) ReTranslate(ctx context.Context, noteId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FullTextTranslateService.ReTranslate")
	defer span.Finish()

	// 获取笔记信息
	note, err := s.paperNoteService.GetPaperNoteById(ctx, noteId)
	if err != nil {
		s.logger.Error("查询笔记失败", "error", err.Error(), "noteId", noteId)
		return "0", errors.Biz("query note failed")
	}

	if note == nil {
		return "0", errors.Biz("not found note")
	}

	// 获取翻译历史记录
	histories, err := s.fullTextTranslateDAO.FindByPdfId(ctx, "1")
	if err != nil {
		s.logger.Error("查询翻译历史记录失败", "error", err.Error(), "pdfId", "1")
		return "0", errors.Biz("query history failed")
	}

	if len(histories) == 0 {
		return "0", errors.Biz("not found history")
	}

	// 获取源文件URL和错误文件URL
	sourceFileUrl := histories[0].SourcePdfId
	errorFileUrl := histories[0].TargetObjectKey

	// 创建修复记录
	fix := &model.FullTextTranslateFix{
		NoteId:            noteId,
		ErrorFileObjectId: errorFileUrl,
		Progress:          0,
	}

	// err = s.fullTextTranslateFixDAO.Save(ctx, fix)
	// if err != nil {
	// 	s.logger.Error("保存翻译修复记录失败", "error", err.Error())
	// 	return 0, errors.Biz("保存翻译修复记录失败")
	// }

	// 异步调用重新翻译服务
	uc := userContext.GetUserContext(ctx)
	userContext.RunAsyncWithUserContext(uc, func(newCtx context.Context) {
		s.reTranslate(newCtx, sourceFileUrl, fix.Id)
	})

	return fix.Id, nil
}

// reTranslate 调用重新翻译服务
func (s *FullTextTranslateService) reTranslate(ctx context.Context, url string, fixId string) (*dto.FullTranslateUploadResult, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FullTextTranslateService.reTranslate")
	defer span.Finish()

	// 构建请求URL
	requestURL := fmt.Sprintf("%s?pdf_url=%s&fixId=%d", s.config.Translate.FullTextTranslate.TranslateServiceFixURI, url, fixId)

	// 设置请求头
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	// 使用已有的HTTP客户端组件发送请求
	body, err := s.httpClient.Get(requestURL, headers)
	if err != nil {
		s.logger.Error("发送HTTP请求失败", "error", err.Error(), "url", requestURL)
		return nil, err
	}

	// 解析响应
	var result dto.FullTranslateUploadResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		s.logger.Error("解析HTTP响应失败", "error", err.Error(), "body", string(body))
		return nil, err
	}

	s.logger.Info("调用重新翻译服务响应", "result", result)
	return &result, nil
}

// FixTranslateResult 修复翻译结果
func (s *FullTextTranslateService) FixTranslateResult(ctx context.Context, result *dto.FullTranslateProgressResult) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FullTextTranslateService.FixTranslateResult")
	defer span.Finish()

	fixId := result.FixId
	if fixId == 0 {
		return errors.Biz("修复ID不能为空")
	}

	// 获取修复记录
	// fix, err := s.fullTextTranslateFixDAO.FindById(ctx, fixId)
	// if err != nil {
	// 	s.logger.Error("查询翻译修复记录失败", "error", err.Error(), "fixId", fixId)
	// 	return errors.Biz("查询翻译修复记录失败")
	// }

	// if fix == nil {
	// 	s.logger.Error("未找到翻译修复记录", "fixId", fixId)
	// 	return errors.Biz("未找到翻译修复记录")
	// }

	// if result.Progress >= 100 {
	// 	// 翻译完成
	// 	fix.Progress = 100
	// 	fix.TranslationFileUrl = result.Data.StaticDomain + "/" + result.Data.FileFullName

	// 	err = s.fullTextTranslateFixDAO.Save(ctx, fix)
	// 	if err != nil {
	// 		s.logger.Error("更新翻译修复记录失败", "error", err.Error(), "fixId", fixId)
	// 		return errors.Biz("更新翻译修复记录失败")
	// 	}

	// 	// 更新所有相关的翻译历史记录
	// 	if fix.NoteId > 0 {
	// 		// note, err := s.paperNoteDAO.FindById(ctx, fix.NoteId)
	// 		// if err != nil {
	// 		// 	s.logger.Error("查询笔记失败", "error", err.Error(), "noteId", fix.NoteId)
	// 		// } else if note != nil {
	// 		// 	histories, err := s.fullTextTranslateHistoryDAO.FindByPdfId(ctx, note.PdfId)
	// 		// 	if err != nil {
	// 		// 		s.logger.Error("查询翻译历史记录失败", "error", err.Error(), "pdfId", note.PdfId)
	// 		// 	} else {
	// 		// 		for _, history := range histories {
	// 		// 			history.TargetContent = fix.TranslationFileUrl
	// 		// 			err = s.fullTextTranslateHistoryDAO.Save(ctx, history)
	// 		// 			if err != nil {
	// 		// 				s.logger.Error("更新翻译历史记录失败", "error", err.Error(), "historyId", history.Id)
	// 		// 			}
	// 		// 		}
	// 		// 	}
	// 		// }
	// 	}
	// } else {
	// 	// 更新进度
	// 	fix.Progress = result.Progress
	// 	err = s.fullTextTranslateFixDAO.Save(ctx, fix)
	// 	if err != nil {
	// 		s.logger.Error("更新翻译修复记录失败", "error", err.Error(), "fixId", fixId)
	// 		return errors.Biz("更新翻译修复记录失败")
	// 	}
	// }

	return nil
}

// GetReTranslateResult 获取重新翻译结果
// func (s *FullTextTranslateService) GetReTranslateResult(ctx context.Context, req *translate.ReTranslateResultRequest) (*translate.PageInfo, error) {
// 	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FullTextTranslateService.GetReTranslateResult")
// 	defer span.Finish()

// 	// 查询修复记录
// 	// 这里需要根据实际情况实现分页查询
// 	// ...

// 	return &translate.PageInfo{}, nil
// }

// handleTranslateFail 处理翻译失败
func (s *FullTextTranslateService) handleTranslateFail(ctx context.Context, history *model.FullTextTranslate, message string) {
	history.Status = 0

	// 更新扩展参数
	extParams := map[string]interface{}{
		"message": message,
	}
	extParamsBytes, _ := json.Marshal(extParams)
	history.Message = string(extParamsBytes)
	history.Status = pb.FullTranslateFlowStatus_TRANSLATE_FAIL

	err := s.fullTextTranslateDAO.ModifyExcludeNull(ctx, history)
	if err != nil {
		s.logger.Error("更新翻译历史记录失败", "error", err.Error(), "historyId", history.Id)
	}

	//设置用户id
	uc := userContext.NewUserContext().SetUserID(history.UserId)
	ctx = uc.ToContext(ctx)
	// VIP回退权益
	err = s.membershipService.RetrieveCreditFun(ctx, history.SessionId)
	if err != nil {
		s.logger.Error("确认翻译失败", "error", err.Error(), "sessionId", history.SessionId)
	}

	// 发送告警通知
	s.logger.Error("翻译失败", "message", message, "pdfId", 1)

	// 异步发送告警
	userContext.RunAsyncWithUserContext(uc, func(newCtx context.Context) {
		// 这里需要根据实际情况调用告警服务
		// ...
	})
}

// 辅助函数：创建字符串指针
func stringPtr(s string) *string {
	return &s
}

// 获取用户的全文翻译历史
func (s *FullTextTranslateService) GetHistoryList(ctx context.Context, userId string) (*pb.GetHistoryListResponse, error) {
	fullTextTranslateHistories, err := s.fullTextTranslateDAO.FindByUserIdAndIntervalDays(ctx, userId, s.config.Translate.FullTextTranslate.HistoryVisibleDays)
	if err != nil {
		s.logger.Error("获取全文翻译历史失败", "error", err.Error())
		return nil, errors.System(errors.ErrorTypeDatabase, "get history list failed", err)
	}

	var results []*pb.HistoryInfo

	for _, history := range fullTextTranslateHistories {
		historyInfo := &pb.HistoryInfo{
			DocName:       history.DocName,
			PdfId:         &history.SourcePdfId,
			Status:        history.Status,
			TranslateTime: uint64(history.CreatedAt.Unix()),
		}
		results = append(results, historyInfo)
	}
	return &pb.GetHistoryListResponse{HistoryList: results}, nil
}
