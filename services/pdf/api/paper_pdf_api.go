package api

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	pb "github.com/yb2020/odoc/proto/gen/go/pdf"
	"github.com/yb2020/odoc/services/pdf/model"
	"github.com/yb2020/odoc/services/pdf/service"
)

// PaperPdfAPI 论文PDF API处理器
type PaperPdfAPI struct {
	paperPdfService         *service.PaperPdfService
	logger                  logging.Logger
	tracer                  opentracing.Tracer
	pdfReaderSettingService *service.PdfReaderSettingService
}

// NewPaperPdfAPI 创建论文PDF API处理器
func NewPaperPdfAPI(
	paperPdfService *service.PaperPdfService,
	logger logging.Logger,
	tracer opentracing.Tracer,
	pdfReaderSettingService *service.PdfReaderSettingService,
) *PaperPdfAPI {
	return &PaperPdfAPI{
		paperPdfService:         paperPdfService,
		logger:                  logger,
		tracer:                  tracer,
		pdfReaderSettingService: pdfReaderSettingService,
	}
}

// @api /api/pdf/getPdfStatusInfo/v2
// @method post
// @summary 获取PDF状态信息
func (api *PaperPdfAPI) GetPdfStatusInfo(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "PaperPdfAPI.GetPdfStatusInfo")
	defer span.Finish()

	var req pb.GetPdfStatusInfoRequest
	if err := transport.BindProto(c, &req); err != nil {
		api.logger.Error("msg", "解析请求参数失败", "error", err.Error())
		response.ErrorNoData(c, "解析请求参数失败")
		return
	}

	pdfStatusInfo, err := api.paperPdfService.GetPdfStatusInfo(ctx, req.PdfId, req.NoteId)
	if err != nil {
		api.logger.Error("msg", "获取PDF状态信息失败", "error", err.Error())
		response.ErrorNoData(c, "获取PDF状态信息失败")
		return
	}

	response.Success(c, "success", pdfStatusInfo)
}

// @api /api/pdf/pdfReader/getSetting 获取PDF阅读设置
// @method post
// @summary 获取PDF阅读设置
func (api *PaperPdfAPI) GetPdfReaderSetting(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "PaperPdfAPI.GetPdfReaderSetting")
	defer span.Finish()

	// 获取用户ID
	userId, _ := userContext.GetUserID(ctx)
	api.logger.Info("msg", "获取PDF状态信息", "userId", userId)

	var req pb.GetPdfReaderSettingRequest
	if err := transport.BindProto(c, &req); err != nil {
		api.logger.Error("msg", "解析请求参数失败", "error", err.Error())
		response.ErrorNoData(c, "解析请求参数失败")
		return
	}

	// TODO: 需要实现REDIS缓存

	//获取PDF阅读设置
	setting, err := api.pdfReaderSettingService.GetByUserId(ctx, userId, int(req.ClientType))
	if err != nil {
		api.logger.Error("msg", "获取PDF阅读设置失败", "error", err.Error())
		response.ErrorNoData(c, "获取PDF阅读设置失败")
		return
	}
	var settingStr string
	if setting != nil {
		settingStr = setting.Setting
	}

	settingResponse := &pb.GetPdfReaderSettingResponse{
		Setting: settingStr,
	}

	response.Success(c, "success", settingResponse)

}

// @api {post} /api/pdf/pdfReader/recordSetting 记录PDF阅读设置
// @method post
// @summary 记录PDF阅读设置
func (api *PaperPdfAPI) RecordPdfReaderSetting(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "PaperPdfAPI.RecordPdfReaderSetting")
	defer span.Finish()

	// 获取用户ID
	userId, _ := userContext.GetUserID(ctx)
	api.logger.Info("msg", "获取PDF状态信息", "userId", userId)

	var req pb.RecordPdfReaderSettingRequest
	if err := transport.BindProto(c, &req); err != nil {
		api.logger.Error("msg", "解析请求参数失败", "error", err.Error())
		response.ErrorNoData(c, "解析请求参数失败")
		return
	}

	// TODO: 需要实现REDIS缓存

	//获取PDF阅读设置
	setting, err := api.pdfReaderSettingService.GetByUserId(ctx, userId, int(req.ClientType))
	if err != nil {
		api.logger.Error("msg", "获取PDF阅读设置失败", "error", err.Error())
		response.ErrorNoData(c, "获取PDF阅读设置失败")
		return
	}

	if setting != nil {
		setting.Setting = req.Setting
		_, err = api.pdfReaderSettingService.UpdatePdfReaderSetting(ctx, setting)
		if err != nil {
			api.logger.Error("msg", "更新PDF阅读设置失败", "error", err.Error())
			response.ErrorNoData(c, "更新PDF阅读设置失败")
			return
		}
	} else {
		setting = &model.PdfReaderSetting{
			ClientType: int(req.ClientType),
			Setting:    req.Setting,
		}
		_, err = api.pdfReaderSettingService.CreatePdfReaderSetting(ctx, setting)
		if err != nil {
			api.logger.Error("msg", "创建PDF阅读设置失败", "error", err.Error())
			response.ErrorNoData(c, "创建PDF阅读设置失败")
			return
		}
	}

	response.SuccessNoData(c, "success")

}
