package api

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	pb "github.com/yb2020/odoc-proto/gen/go/note"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	docService "github.com/yb2020/odoc/services/doc/service"
	noteInterface "github.com/yb2020/odoc/services/note/interfaces"
	"github.com/yb2020/odoc/services/note/model"
	"github.com/yb2020/odoc/services/note/service"
)

// NoteSummaryAPI 笔记摘要API处理器
type NoteSummaryAPI struct {
	service          *service.NoteSummaryService
	userDocService   *docService.UserDocService
	logger           logging.Logger
	tracer           opentracing.Tracer
	paperNoteService noteInterface.IPaperNoteService
}

// NewNoteSummaryAPI 创建笔记摘要API处理器
func NewNoteSummaryAPI(service *service.NoteSummaryService, userDocService *docService.UserDocService, logger logging.Logger, tracer opentracing.Tracer, paperNoteService noteInterface.IPaperNoteService) *NoteSummaryAPI {
	return &NoteSummaryAPI{
		service:          service,
		userDocService:   userDocService,
		logger:           logger,
		tracer:           tracer,
		paperNoteService: paperNoteService,
	}
}

// @api_path: /api/note/paperNote/summary/getByNoteId
// @method: GET
// @content-type: application/json
// @summary: 获取笔记摘要
func (api *NoteSummaryAPI) GetNoteSummaryByNoteId(c *gin.Context) {
	req := &pb.GetNoteSummaryByNoteIdRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析获取笔记摘要请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析获取笔记摘要请求失败")
		return
	}

	userDoc, err := api.userDocService.GetUserDocByNoteId(c.Request.Context(), req.NoteId)
	if err != nil {
		api.logger.Error("msg", "获取用户文档失败", "error", err.Error())
		response.ErrorNoData(c, "获取用户文档失败")
		return
	}
	// 在构建响应之前添加
	docName := ""
	if userDoc != nil {
		docName = userDoc.DocName
	}

	// 直接使用请求的上下文，它已经在中间件中被更新
	noteSummary, err := api.service.GetNoteSummaryByNoteId(c.Request.Context(), req.NoteId)
	if err != nil {
		api.logger.Error("msg", "获取笔记摘要失败", "error", err.Error())
		response.ErrorNoData(c, "获取笔记摘要失败")
		return
	}

	resp := &pb.GetNoteSummaryByNoteIdResponse{
		DocName: docName,
	}

	if noteSummary != nil {
		resp.Content = noteSummary.Content
		resp.ModifyDate = uint64(noteSummary.UpdatedAt.UnixMilli())
	}

	response.Success(c, "success", resp)
}

/**
 * @api_path: /api/note/paperNote/summary/saveOrUpdate
 * @method: POST
 * @content-type: application/json
 * @summary: 添加/更新总结
 */
func (api *NoteSummaryAPI) SaveOrUpdateSummary(c *gin.Context) {
	req := &pb.SaveOrUpdateSummaryReq{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析更新笔记摘要请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析更新笔记摘要请求失败")
		return
	}
	// 笔记摘要不能为空判断
	if !(req.Content != "") {
		response.ErrorNoData(c, "笔记摘要不能为空")
		return
	}

	// 实现添加/更新总结
	userId, _ := userContext.GetUserID(c.Request.Context())
	paperNote, err := api.paperNoteService.GetPaperNoteById(c.Request.Context(), req.NoteId)
	if err != nil {
		api.logger.Error("msg", "获取论文笔记失败", "error", err.Error())
		response.ErrorNoData(c, "获取论文笔记失败")
		return
	}
	// 检查笔记是否存在以及是否属于自己
	if userId == "" || paperNote == nil || paperNote.CreatorId != userId {
		api.logger.Error("msg", "笔记不存在或者不属于您本人")
		response.ErrorNoData(c, "笔记不存在或者不属于您本人")
		return
	}

	// 先获取原始数据
	noteSummary, err := api.service.GetNoteSummaryByNoteId(c.Request.Context(), req.NoteId)
	if err != nil {
		api.logger.Error("msg", "获取笔记摘要失败", "error", err.Error())
		response.ErrorNoData(c, "获取笔记摘要失败")
		return
	}
	// 如果不存在，创建新摘要
	if noteSummary == nil {
		noteSummary = &model.NoteSummary{
			NoteId:  req.NoteId,
			Content: req.Content,
			UserId:  userId,
		}
		_, err = api.service.CreateNoteSummary(c.Request.Context(), noteSummary)
		if err != nil {
			api.logger.Error("msg", "创建笔记摘要失败", "error", err.Error())
			response.ErrorNoData(c, "创建笔记摘要失败")
			return
		}
	} else {
		noteSummary.Content = req.Content
		noteSummary.UserId = userId
		// 更新数据
		_, err := api.service.UpdateNoteSummary(c.Request.Context(), noteSummary)
		if err != nil {
			api.logger.Error("msg", "更新笔记摘要失败", "error", err.Error())
			response.ErrorNoData(c, "更新笔记摘要失败")
			return
		}
	}
	response.SuccessNoData(c, "success")
}

// // --管理端接口--------------------------------------------------------------------------------------//
// // @api_path: /api/note/noteManage/summary/getList
// // @method: POST
// // @content-type: application/json
// // @summary: 笔记管理TAB——总结
// func (api *NoteSummaryAPI) GetUserSummaryList(c *gin.Context) {
// 	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "NoteSummaryAPI.GetNoteSummaryByNoteId")
// 	defer span.Finish()

// 	req := &pb.GetSummaryListReq{}
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

// // @api_path: /api/note/noteManage/summary/getListByFolderId
// // @method: POST
// // @content-type: application/json
// // @summary: 笔记管理TAB——根据folderId获取总结列表
// func (api *NoteSummaryAPI) GetUserSummaryListByFolderId(c *gin.Context) {
// 	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "NoteSummaryAPI.GetUserSummaryListByFolderId")
// 	defer span.Finish()

// 	req := &pb.GetSummaryListByFolderIdReq{}
// 	if err := transport.BindProto(c, req); err != nil {
// 		api.logger.Error("msg", "解析获取笔记总结列表请求失败", "error", err.Error())
// 		response.ErrorNoData(c, "解析获取笔记总结列表请求失败")
// 		return
// 	}

// 	// userId, _ := userContext.GetUserID(ctx)
// 	resp, err := api.service.GetUserSummaryListByFolderId(ctx, req)
// 	if err != nil {
// 		api.logger.Error("msg", "根据folderId获取总结列表失败", "error", err.Error())
// 		response.ErrorNoData(c, "根据folderId获取总结列表失败")
// 		return
// 	}
// 	//var resp = &pb.GetSummaryListByFolderIdResponse{}

// 	response.Success(c, "success", resp)
// }
