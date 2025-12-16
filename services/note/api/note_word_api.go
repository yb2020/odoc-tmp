package api

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	pb "github.com/yb2020/odoc/proto/gen/go/note"
	"github.com/yb2020/odoc/services/note/bean"
	noteInterface "github.com/yb2020/odoc/services/note/interfaces"
	"github.com/yb2020/odoc/services/note/service"
)

// NoteWordAPI 笔记单词API处理器
type NoteWordAPI struct {
	service               noteInterface.INoteWordService
	noteWordConfigService *service.NoteWordConfigService
	logger                logging.Logger
	tracer                opentracing.Tracer
}

// NewNoteWordAPI 创建笔记单词API处理器
func NewNoteWordAPI(service noteInterface.INoteWordService, noteWordConfigService *service.NoteWordConfigService, logger logging.Logger, tracer opentracing.Tracer) *NoteWordAPI {
	return &NoteWordAPI{
		service:               service,
		noteWordConfigService: noteWordConfigService,
		logger:                logger,
		tracer:                tracer,
	}
}

// @api_path: /api/note/paperNote/word/getByNoteId
// @method: POST
// @content-type: application/json
// @summary: 根据noteId获取生词列表
func (api *NoteWordAPI) GetNoteWordsByNoteId(c *gin.Context) {
	req := &pb.GetNoteWordsByNoteIdRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "Failed to parse request parameters", "error", err.Error())
		response.ErrorNoData(c, "Failed to parse request parameters")
		return
	}

	// 转换请求参数
	noteWordQuery := &bean.NoteWordQuery{
		NoteId:      req.NoteId,
		CurrentPage: int(req.CurrentPage),
		PageSize:    int(req.PageSize),
		MinLoadedId: req.MinLoadedId,
	}

	// 调用服务层方法
	noteWordsInfo, err := api.service.GetNoteWordsByNoteQuery(c.Request.Context(), noteWordQuery)
	if err != nil {
		api.logger.Error("msg", "Failed to get note words by note id", "error", err.Error())
		response.ErrorNoData(c, "Failed to get note words by note id")
		return
	}
	response.Success(c, "success", noteWordsInfo)
}

// @api_path: /api/note/paperNote/word/config
// @method: POST
// @content-type: application/json
// @summary: 保存或更新生词配置
func (api *NoteWordAPI) SaveOrUpdateNoteWordConfig(c *gin.Context) {
	req := &pb.ChangeWordConfigRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "Failed to parse request parameters", "error", err.Error())
		response.ErrorNoData(c, "Failed to parse request parameters")
		return
	}

	// 调用服务层方法
	_, err := api.service.SaveOrUpdateNoteWordConfig(c.Request.Context(), req.NoteId, req.Color, req.DisplayMode)
	if err != nil {
		api.logger.Error("msg", "Failed to save or update note word config", "error", err.Error())
		response.ErrorNoData(c, "Failed to save or update note word config")
		return
	}

	response.SuccessNoData(c, "success")
}

// @api_path: /api/note/paperNote/word/delete
// @method: POST
// @content-type: application/json
// @summary: 删除笔记单词
func (api *NoteWordAPI) DeleteNoteWord(c *gin.Context) {
	req := &pb.DeleteNoteWordRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析删除笔记单词请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析删除笔记单词请求失败")
		return
	}

	// 直接使用请求的上下文，它已经在中间件中被更新
	_, err := api.service.DeleteNoteWordById(c.Request.Context(), req.WordId)
	if err != nil {
		api.logger.Error("msg", "删除笔记单词失败", "error", err.Error())
		response.ErrorNoData(c, "删除笔记单词失败")
		return
	}

	response.SuccessNoData(c, "success")
}

// @api_path: /api/note/paperNote/word/update
// @method: POST
// @content-type: application/json
// @summary: 更新笔记单词
func (api *NoteWordAPI) UpdateNoteWord(c *gin.Context) {
	req := &pb.UpdateNoteWordRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析更新笔记单词请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析更新笔记单词请求失败")
		return
	}

	// 更新数据
	success, err := api.service.UpdateNoteWordTargetContent(c.Request.Context(), req)
	if err != nil {
		api.logger.Error("msg", "更新笔记单词失败", "error", err.Error())
		response.ErrorNoData(c, "更新笔记单词失败")
		return
	}

	resp := &pb.UpdateNoteWordResponse{
		Success: success,
	}

	response.Success(c, "success", resp)
}

// @api_path: /api/note/paperNote/word/save
// @method: POST
// @content-type: application/json
// @summary: 保存笔记单词
func (api *NoteWordAPI) SaveNoteWord(c *gin.Context) {
	req := &pb.SaveNoteWordRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析保存笔记单词请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析保存笔记单词请求失败")
		return
	}

	id, err := api.service.CreateNoteWord(c.Request.Context(), req)
	if err != nil {
		api.logger.Error("msg", "保存笔记单词失败", "error", err.Error())
		response.ErrorNoData(c, "保存笔记单词失败")
		return
	}

	resp := &pb.SaveNoteWordResponse{
		Id: id,
	}

	response.Success(c, "success", resp)
}

// // --管理端接口--------------------------------------------------------------------------------------//

// // @api_path: /api/note/noteManage/word/getList
// // @method: POST
// // @content-type: application/json
// // @summary: 笔记管理TAB——单词
// func (api *NoteWordAPI) GetUserWordList(c *gin.Context) {
// 	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "NoteWordAPI.GetUserWordList")
// 	defer span.Finish()

// 	req := &pb.GetWordListReq{}
// 	if err := transport.BindProto(c, req); err != nil {
// 		api.logger.Error("msg", "解析获取笔记总结列表请求失败", "error", err.Error())
// 		response.ErrorNoData(c, "解析获取笔记总结列表请求失败")
// 		return
// 	}

// 	userId, _ := userContext.GetUserID(ctx)

// 	resp, err := api.service.GetListByUserId(ctx, userId)
// 	if err != nil {
// 		api.logger.Error("msg", "获取笔记总结列表失败", "error", err.Error())
// 		response.ErrorNoData(c, "获取笔记总结列表失败")
// 	}

// 	response.Success(c, "success", resp)
// }

// // @api_path: /api/note/noteManage/word/getListByFolderId
// // @method: POST
// // @content-type: application/json
// // @summary: 笔记管理TAB——根据folderId获取单词列表
// func (api *NoteWordAPI) GetWordListByFolderIdReq(c *gin.Context) {
// 	span, _ := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "NoteWordAPI.GetWordListByFolderIdReq")
// 	defer span.Finish()

// 	req := &pb.GetWordListByFolderIdReq{}
// 	if err := transport.BindProto(c, req); err != nil {
// 		api.logger.Error("msg", "解析获取笔记单词列表请求失败", "error", err.Error())
// 		response.ErrorNoData(c, "解析获取笔记单词列表请求失败")
// 		return
// 	}

// 	// userId, _ := userContext.GetUserID(ctx)
// 	// resp, err := api.service.GetUserWordListByFolderId(ctx, req)
// 	// if err != nil {
// 	// 	api.logger.Error("msg", "根据folderId获取总结列表失败", "error", err.Error())
// 	// 	response.ErrorNoData(c, "根据folderId获取总结列表失败")
// 	// 	return
// 	// }
// 	var resp = &pb.GetWordListByFolderIdResponse{}

// 	response.Success(c, "success", resp)
// }

// // @api_path: /api/note/noteManage/extract/getList
// // @method: POST
// // @content-type: application/json
// // @summary: 笔记管理TAB——摘录
// func (api *NoteWordAPI) GetExtractListReq(c *gin.Context) {
// 	span, _ := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "NoteWordAPI.GetExtractListReq")
// 	defer span.Finish()

// 	req := &pb.GetExtractListReq{}
// 	if err := transport.BindProto(c, req); err != nil {
// 		api.logger.Error("msg", "解析获取笔记单词列表请求失败", "error", err.Error())
// 		response.ErrorNoData(c, "解析获取笔记单词列表请求失败")
// 		return
// 	}

// 	// userId, _ := userContext.GetUserID(ctx)
// 	// resp, err := api.service.GetUserWordListByFolderId(ctx, req)
// 	// if err != nil {
// 	// 	api.logger.Error("msg", "根据folderId获取总结列表失败", "error", err.Error())
// 	// 	response.ErrorNoData(c, "根据folderId获取总结列表失败")
// 	// 	return
// 	// }
// 	var resp = &pb.GetExtractListResponse{}

// 	response.Success(c, "success", resp)
// }
