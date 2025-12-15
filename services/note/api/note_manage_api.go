package api

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	pb "github.com/yb2020/odoc-proto/gen/go/note"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	noteInterface "github.com/yb2020/odoc/services/note/interfaces"
	"github.com/yb2020/odoc/services/note/service"
	pdfInterfaces "github.com/yb2020/odoc/services/pdf/interfaces"
)

// NoteManageAPI 笔记管理API处理器
type NoteManageAPI struct {
	logger         logging.Logger
	tracer         opentracing.Tracer
	wordService    noteInterface.INoteWordService
	summaryService *service.NoteSummaryService
	pdfService     pdfInterfaces.IPaperPdfService
}

// NewNoteManageAPI 创建笔记管理API处理器
func NewNoteManageAPI(logger logging.Logger, tracer opentracing.Tracer, wordService noteInterface.INoteWordService, summaryService *service.NoteSummaryService, pdfService pdfInterfaces.IPaperPdfService) *NoteManageAPI {
	return &NoteManageAPI{
		logger:         logger,
		tracer:         tracer,
		wordService:    wordService,
		summaryService: summaryService,
		pdfService:     pdfService,
	}
}

// SetPaperPdfService 设置论文PDF服务
func (s *NoteManageAPI) SetPaperPdfService(pdfService pdfInterfaces.IPaperPdfService) error {
	if pdfService == nil {
		return errors.Biz("pdfService cannot be nil")
	}
	s.pdfService = pdfService
	return nil
}

// --管理端接口--------------------------------------------------------------------------------------//

// @api_path: /api/note/noteManage/word/getList
// @method: POST
// @content-type: application/json
// @summary: 笔记管理TAB——单词
func (api *NoteManageAPI) GetUserWordList(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "NoteManageAPI.GetUserWordList")
	defer span.Finish()

	req := &pb.GetWordListReq{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析获取笔记总结列表请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析获取笔记总结列表请求失败")
		return
	}

	userId, _ := userContext.GetUserID(ctx)

	resp, err := api.wordService.GetListByUserId(ctx, userId)
	if err != nil {
		api.logger.Error("msg", "获取笔记总结列表失败", "error", err.Error())
		response.ErrorNoData(c, "获取笔记总结列表失败")
	}

	response.Success(c, "success", resp)
}

// @api_path: /api/note/noteManage/word/getListByFolderId
// @method: POST
// @content-type: application/json
// @summary: 笔记管理TAB——根据folderId获取单词列表
func (api *NoteManageAPI) GetWordListByFolderIdReq(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "NoteManageAPI.GetWordListByFolderIdReq")
	defer span.Finish()

	req := &pb.GetWordListByFolderIdReq{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析获取笔记单词列表请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析获取笔记单词列表请求失败")
		return
	}

	if req.FolderId == nil {
		var defaultFolderId string = "0"
		req.FolderId = &defaultFolderId
	}
	if req.CurrentPage == 0 || req.CurrentPage < 0 {
		req.CurrentPage = 1
	}
	if req.PageSize == 0 || req.PageSize < 0 {
		req.PageSize = 10
	}

	userId, _ := userContext.GetUserID(ctx)
	resp, err := api.wordService.GetUserWordListByFolderId(ctx, userId, req)
	if err != nil {
		api.logger.Error("msg", "根据folderId获取总结列表失败", "error", err.Error())
		response.ErrorNoData(c, "根据folderId获取总结列表失败")
		return
	}

	response.Success(c, "success", resp)
}

// @api_path: /api/note/noteManage/summary/getList
// @method: POST
// @content-type: application/json
// @summary: 笔记管理TAB——总结
func (api *NoteManageAPI) GetUserSummaryList(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "NoteManageAPI.GetUserSummaryList")
	defer span.Finish()

	req := &pb.GetSummaryListReq{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析获取笔记总结列表请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析获取笔记总结列表请求失败")
		return
	}

	userId, _ := userContext.GetUserID(ctx)

	resp, err := api.summaryService.GetSummaryListByUserId(ctx, userId)
	if err != nil {
		api.logger.Error("msg", "获取笔记总结列表失败", "error", err.Error())
		response.ErrorNoData(c, "获取笔记总结列表失败")
	}

	response.Success(c, "success", resp)
}

// @api_path: /api/note/noteManage/summary/getListByFolderId
// @method: POST
// @content-type: application/json
// @summary: 笔记管理TAB——根据folderId获取总结列表
func (api *NoteManageAPI) GetUserSummaryListByFolderId(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "NoteManageAPI.GetUserSummaryListByFolderId")
	defer span.Finish()

	req := &pb.GetSummaryListByFolderIdReq{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析获取笔记总结列表请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析获取笔记总结列表请求失败")
		return
	}

	if req.FolderId == nil {
		var defaultFolderId string = "0"
		req.FolderId = &defaultFolderId
	}
	if req.CurrentPage == 0 || req.CurrentPage < 0 {
		req.CurrentPage = 1
	}
	if req.PageSize == 0 || req.PageSize < 0 {
		req.PageSize = 10
	}

	userId, _ := userContext.GetUserID(ctx)
	resp, err := api.summaryService.GetUserSummaryListByFolderId(ctx, userId, req)
	if err != nil {
		api.logger.Error("msg", "根据folderId获取总结列表失败", "error", err.Error())
		response.ErrorNoData(c, "根据folderId获取总结列表失败")
		return
	}

	response.Success(c, "success", resp)
}

// @api_path: /api/note/noteManage/extract/getList
// @method: POST
// @content-type: application/json
// @summary: 笔记管理TAB——摘录
func (api *NoteManageAPI) GetExtractListReq(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "NoteManageAPI.GetExtractListReq")
	defer span.Finish()

	req := &pb.GetExtractListReq{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析获取笔记单词列表请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析获取笔记单词列表请求失败")
		return
	}

	userId, _ := userContext.GetUserID(ctx)

	resp, err := api.summaryService.GetExtractListByUserId(ctx, userId)
	if err != nil {
		api.logger.Error("msg", "获取笔记总结列表失败", "error", err.Error())
		response.ErrorNoData(c, "获取笔记总结列表失败")
	}
	//var resp = &pb.GetExtractListResponse{}

	response.Success(c, "success", resp)
}

// @api_path: /api/note/noteManage/extract/getMarkTagListByFolderId
// @method: POST
// @content-type: application/json
// @summary: 笔记管理TAB——根据folderId获取标签列表
func (api *NoteManageAPI) GetMarkTagListByFolderIdReq(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "NoteManageAPI.GetMarkTagListByFolderIdReq")
	defer span.Finish()

	req := &pb.GetMarkTagListByFolderIdReq{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析获取笔记总结列表请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析获取笔记总结列表请求失败")
		return
	}

	if req.FolderId == "" {
		req.FolderId = "0"
	}

	userId, _ := userContext.GetUserID(ctx)
	markTagInfos, err := api.pdfService.GetMarkTagInfosByFolderId(ctx, userId, req.FolderId)
	if err != nil {
		api.logger.Error("msg", "根据folderId获取总结列表失败", "error", err.Error())
		response.ErrorNoData(c, "根据folderId获取总结列表失败")
		return
	}
	var resp = &pb.GetMarkTagListByFolderIdResponse{
		MarkTagList: markTagInfos,
	}

	response.Success(c, "success", resp)
}
