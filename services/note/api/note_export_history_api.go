package api

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	pb "github.com/yb2020/odoc/proto/gen/go/note"
	"github.com/yb2020/odoc/services/note/model"
	"github.com/yb2020/odoc/services/note/service"
)

// NoteExportHistoryAPI 笔记导出历史API处理器
type NoteExportHistoryAPI struct {
	service *service.NoteExportHistoryService
	logger  logging.Logger
	tracer  opentracing.Tracer
}

// NewNoteExportHistoryAPI 创建笔记导出历史API处理器
func NewNoteExportHistoryAPI(service *service.NoteExportHistoryService, logger logging.Logger, tracer opentracing.Tracer) *NoteExportHistoryAPI {
	return &NoteExportHistoryAPI{
		service: service,
		logger:  logger,
		tracer:  tracer,
	}
}

// CreateNoteExportHistory 创建笔记导出历史
func (api *NoteExportHistoryAPI) CreateNoteExportHistory(c *gin.Context) {
	req := &pb.CreateNoteExportHistoryRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析创建笔记导出历史请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析创建笔记导出历史请求失败")
		return
	}

	// 转换请求参数
	noteExportHistory := &model.NoteExportHistory{
		NoteId:  req.NoteId,
		Version: int64(req.Version),
	}

	// 直接使用请求的上下文，它已经在中间件中被更新
	id, err := api.service.CreateNoteExportHistory(c.Request.Context(), noteExportHistory)
	if err != nil {
		api.logger.Error("msg", "创建笔记导出历史失败", "error", err.Error())
		response.ErrorNoData(c, "创建笔记导出历史失败")
		return
	}

	resp := &pb.CreateNoteExportHistoryResponse{
		Id: id,
	}

	response.Success(c, "success", resp)
}

// GetNoteExportHistory 获取笔记导出历史
func (api *NoteExportHistoryAPI) GetNoteExportHistory(c *gin.Context) {
	req := &pb.GetNoteExportHistoryRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析获取笔记导出历史请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析获取笔记导出历史请求失败")
		return
	}

	// 直接使用请求的上下文，它已经在中间件中被更新
	noteExportHistory, err := api.service.GetNoteExportHistoryById(c.Request.Context(), req.Id)
	if err != nil {
		api.logger.Error("msg", "获取笔记导出历史失败", "error", err.Error())
		response.ErrorNoData(c, "获取笔记导出历史失败")
		return
	}

	resp := &pb.GetNoteExportHistoryResponse{
		NoteExportHistory: &pb.NoteExportHistory{
			Id:      noteExportHistory.Id,
			NoteId:  noteExportHistory.NoteId,
			Version: int64(noteExportHistory.Version),
		},
	}

	response.Success(c, "success", resp)
}

// UpdateNoteExportHistory 更新笔记导出历史
func (api *NoteExportHistoryAPI) UpdateNoteExportHistory(c *gin.Context) {
	req := &pb.UpdateNoteExportHistoryRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析更新笔记导出历史请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析更新笔记导出历史请求失败")
		return
	}

	// 先获取原始数据
	noteExportHistory, err := api.service.GetNoteExportHistoryById(c.Request.Context(), req.Id)
	if err != nil {
		api.logger.Error("msg", "获取笔记导出历史失败", "error", err.Error())
		response.ErrorNoData(c, "获取笔记导出历史失败")
		return
	}

	// 更新数据
	success, err := api.service.UpdateNoteExportHistory(c.Request.Context(), noteExportHistory)
	if err != nil {
		api.logger.Error("msg", "更新笔记导出历史失败", "error", err.Error())
		response.ErrorNoData(c, "更新笔记导出历史失败")
		return
	}

	resp := &pb.UpdateNoteExportHistoryResponse{
		Success: success,
	}

	response.Success(c, "success", resp)
}

// DeleteNoteExportHistory 删除笔记导出历史
func (api *NoteExportHistoryAPI) DeleteNoteExportHistory(c *gin.Context) {
	req := &pb.DeleteNoteExportHistoryRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析删除笔记导出历史请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析删除笔记导出历史请求失败")
		return
	}

	// 直接使用请求的上下文，它已经在中间件中被更新
	_, err := api.service.DeleteNoteExportHistoryById(c.Request.Context(), req.Id)
	if err != nil {
		api.logger.Error("msg", "删除笔记导出历史失败", "error", err.Error())
		response.ErrorNoData(c, "删除笔记导出历史失败")
		return
	}

	resp := &pb.DeleteNoteExportHistoryResponse{
		Success: true,
	}

	response.Success(c, "success", resp)
}

// GetNoteExportHistoriesByNoteId 根据笔记ID获取笔记导出历史列表
func (api *NoteExportHistoryAPI) GetNoteExportHistoriesByNoteId(c *gin.Context) {
	req := &pb.GetNoteExportHistoriesByNoteIdRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析获取笔记导出历史列表请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析获取笔记导出历史列表请求失败")
		return
	}

	// 直接使用请求的上下文，它已经在中间件中被更新
	noteExportHistories, err := api.service.GetNoteExportHistoriesByNoteId(c.Request.Context(), req.NoteId)
	if err != nil {
		api.logger.Error("msg", "获取笔记导出历史列表失败", "error", err.Error())
		response.ErrorNoData(c, "获取笔记导出历史列表失败")
		return
	}

	// 转换为响应格式
	pbNoteExportHistories := make([]*pb.NoteExportHistory, 0, len(noteExportHistories))
	for _, noteExportHistory := range noteExportHistories {
		pbNoteExportHistories = append(pbNoteExportHistories, &pb.NoteExportHistory{
			Id:      noteExportHistory.Id,
			NoteId:  noteExportHistory.NoteId,
			Version: noteExportHistory.Version,
		})
	}

	resp := &pb.GetNoteExportHistoriesByNoteIdResponse{
		NoteExportHistories: pbNoteExportHistories,
	}

	response.Success(c, "success", resp)
}

// GetLatestNoteExportHistoryByNoteId 根据笔记ID获取最新的笔记导出历史
func (api *NoteExportHistoryAPI) GetLatestNoteExportHistoryByNoteId(c *gin.Context) {
	req := &pb.GetLatestNoteExportHistoryByNoteIdRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析获取最新笔记导出历史请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析获取最新笔记导出历史请求失败")
		return
	}

	// 直接使用请求的上下文，它已经在中间件中被更新
	noteExportHistory, err := api.service.GetLatestNoteExportHistoryByNoteId(c.Request.Context(), req.NoteId)
	if err != nil {
		api.logger.Error("msg", "获取最新笔记导出历史失败", "error", err.Error())
		response.ErrorNoData(c, "获取最新笔记导出历史失败")
		return
	}

	resp := &pb.GetLatestNoteExportHistoryByNoteIdResponse{
		NoteExportHistory: &pb.NoteExportHistory{
			Id:      noteExportHistory.Id,
			NoteId:  noteExportHistory.NoteId,
			Version: noteExportHistory.Version,
		},
	}

	response.Success(c, "success", resp)
}
