package api

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	pb "github.com/yb2020/odoc-proto/gen/go/pdf"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	userDocService "github.com/yb2020/odoc/services/doc/service"
	noteInterfaces "github.com/yb2020/odoc/services/note/interfaces"
	paperNotesService "github.com/yb2020/odoc/services/note/service"
	paperService "github.com/yb2020/odoc/services/paper/service"
	"github.com/yb2020/odoc/services/pdf/service"
	userService "github.com/yb2020/odoc/services/user/service"
)

// PdfParseAPI 论文PDF解析API处理器
type PdfParseAPI struct {
	paperPdfService         *service.PaperPdfService
	logger                  logging.Logger
	tracer                  opentracing.Tracer
	userService             *userService.UserService
	paperService            *paperService.PaperService
	userDocService          *userDocService.UserDocService
	paperNoteService        noteInterfaces.IPaperNoteService
	paperNoteAccessService  *paperNotesService.PaperNoteAccessService
	pdfReaderSettingService *service.PdfReaderSettingService
	pdfParseService         *service.PdfParseService
}

// NewPdfParseAPI 创建PDF解析API处理器
func NewPdfParseAPI(
	paperPdfService *service.PaperPdfService,
	logger logging.Logger,
	tracer opentracing.Tracer,
	paperService *paperService.PaperService,
	userService *userService.UserService,
	userDocService *userDocService.UserDocService,
	paperNoteService noteInterfaces.IPaperNoteService,
	paperNoteAccessService *paperNotesService.PaperNoteAccessService,
	pdfReaderSettingService *service.PdfReaderSettingService,
	pdfParseService *service.PdfParseService,
) *PdfParseAPI {
	return &PdfParseAPI{
		paperPdfService:         paperPdfService,
		logger:                  logger,
		tracer:                  tracer,
		paperService:            paperService,
		userService:             userService,
		userDocService:          userDocService,
		paperNoteService:        paperNoteService,
		paperNoteAccessService:  paperNoteAccessService,
		pdfReaderSettingService: pdfReaderSettingService,
		pdfParseService:         pdfParseService,
	}
}

// GetReferenceMarkers 获取参考文献标记
// url: /api/pdf/parser/getReferenceMarkers
// method: POST
// @api {post} /api/pdf/parser/getReferenceMarkers 获取参考文献标记
func (api *PdfParseAPI) GetReferenceMarkers(c *gin.Context) {
	span, _ := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "PaperPdfParseAPI.GetReferenceMarkers")
	defer span.Finish()

	var req pb.GetReferenceMarkersRequest
	if err := transport.BindProto(c, &req); err != nil {
		api.logger.Error("msg", "parse request parameters failed", "error", err.Error())
		response.ErrorNoData(c, "parse request parameters failed")
		return
	}
	res, err := api.pdfParseService.GetReferenceMarkers(c.Request.Context(), &req)
	if err != nil {
		api.logger.Error("msg", "get reference markers failed", "error", err.Error())
		response.ErrorNoData(c, "get reference markers failed")
		return
	}

	response.Success(c, "success", res)
}

// GetFiguresAndTables 获取图表和表格信息
// url: /api/pdf/parser/getFiguresAndTables
// method: POST
// @api {post} /api/pdf/parser/getFiguresAndTables 获取图表和表格信息
func (api *PdfParseAPI) GetFiguresAndTables(c *gin.Context) {
	span, _ := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "PaperPdfParseAPI.GetFiguresAndTables")
	defer span.Finish()

	var req pb.GetFiguresAndTablesListRequest
	if err := transport.BindProto(c, &req); err != nil {
		api.logger.Error("msg", "parse request parameters failed", "error", err.Error())
		response.ErrorNoData(c, "parse request parameters failed")
		return
	}
	res, err := api.pdfParseService.GetFiguresAndTables(c.Request.Context(), &req)
	if err != nil {
		api.logger.Error("msg", "get figures and tables failed", "error", err.Error())
		response.ErrorNoData(c, "get figures and tables failed")
		return
	}
	response.Success(c, "success", res)
}

// @api_path: /api/pdf/parser/getReference
// @method: POST
// @content-type: application/json
// @summary: 提取参考信息
func (api *PdfParseAPI) GetReference(c *gin.Context) {
	span, _ := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "PaperPdfParseAPI.GetReference")
	defer span.Finish()

	var req pb.GetReferenceRequest
	if err := transport.BindProto(c, &req); err != nil {
		response.ErrorNoData(c, "parse request parameters failed")
		return
	}
	res, err := api.pdfParseService.GetReference(c.Request.Context(), &req)
	if err != nil {
		response.ErrorNoData(c, "get reference failed")
		return
	}
	response.Success(c, "success", res)
}

// @api_path: /api/pdf/parser/getCatalogue
// @method: POST
// @content-type: application/json
// @summary: 提取目录信息
func (api *PdfParseAPI) GetCatalogue(c *gin.Context) {
	span, _ := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "PaperPdfParseAPI.GetCatalogue")
	defer span.Finish()

	var req pb.GetCatalogueRequest
	if err := transport.BindProto(c, &req); err != nil {
		response.ErrorNoData(c, "parse request parameters failed")
		return
	}

	res, err := api.pdfParseService.GetCatalogue(c.Request.Context(), &req)
	if err != nil {
		response.ErrorNoData(c, "get catalogue failed")
		return
	}
	response.Success(c, "success", res)
}

// 重新解析
func (api *PdfParseAPI) ReParse(c *gin.Context) {
	span, _ := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "PaperPdfParseAPI.ReParse")
	defer span.Finish()

	var req pb.GetCatalogueRequest
	if err := transport.BindProto(c, &req); err != nil {
		response.ErrorNoData(c, "parse request parameters failed")
		return
	}
	if err := api.pdfParseService.ReParse(c.Request.Context(), &req); err != nil {
		response.ErrorNoData(c, "re parse failed")
		return
	}
	response.Success(c, "success", nil)
}
