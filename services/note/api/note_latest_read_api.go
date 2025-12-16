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

// NoteLatestReadAPI 笔记最近阅读API处理器
type NoteLatestReadAPI struct {
	service *service.NoteLatestReadService
	logger  logging.Logger
	tracer  opentracing.Tracer
}

// NewNoteLatestReadAPI 创建笔记最近阅读API处理器
func NewNoteLatestReadAPI(service *service.NoteLatestReadService, logger logging.Logger, tracer opentracing.Tracer) *NoteLatestReadAPI {
	return &NoteLatestReadAPI{
		service: service,
		logger:  logger,
		tracer:  tracer,
	}
}

// CreateNoteLatestRead 创建笔记最近阅读
func (api *NoteLatestReadAPI) CreateNoteLatestRead(c *gin.Context) {
	req := &pb.CreateNoteLatestReadRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析创建笔记最近阅读请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析创建笔记最近阅读请求失败")
		return
	}

	// 转换请求参数
	noteLatestRead := &model.NoteLatestRead{
		NoteId: req.NoteId,
	}

	// 直接使用请求的上下文，它已经在中间件中被更新
	id, err := api.service.CreateNoteLatestRead(c.Request.Context(), noteLatestRead)
	if err != nil {
		api.logger.Error("msg", "创建笔记最近阅读失败", "error", err.Error())
		response.ErrorNoData(c, "创建笔记最近阅读失败")
		return
	}

	resp := &pb.CreateNoteLatestReadResponse{
		Id: id,
	}

	response.Success(c, "success", resp)
}

// GetNoteLatestRead 获取笔记最近阅读
func (api *NoteLatestReadAPI) GetNoteLatestRead(c *gin.Context) {
	req := &pb.GetNoteLatestReadRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析获取笔记最近阅读请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析获取笔记最近阅读请求失败")
		return
	}

	// 直接使用请求的上下文，它已经在中间件中被更新
	noteLatestRead, err := api.service.GetNoteLatestReadById(c.Request.Context(), req.Id)
	if err != nil {
		api.logger.Error("msg", "获取笔记最近阅读失败", "error", err.Error())
		response.ErrorNoData(c, "获取笔记最近阅读失败")
		return
	}

	resp := &pb.GetNoteLatestReadResponse{
		NoteLatestRead: &pb.NoteLatestRead{
			Id:     noteLatestRead.Id,
			NoteId: noteLatestRead.NoteId,
		},
	}

	response.Success(c, "success", resp)
}

// UpdateNoteLatestRead 更新笔记最近阅读
func (api *NoteLatestReadAPI) UpdateNoteLatestRead(c *gin.Context) {
	req := &pb.UpdateNoteLatestReadRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析更新笔记最近阅读请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析更新笔记最近阅读请求失败")
		return
	}

	// 先获取原始数据
	noteLatestRead, err := api.service.GetNoteLatestReadById(c.Request.Context(), req.Id)
	if err != nil {
		api.logger.Error("msg", "获取笔记最近阅读失败", "error", err.Error())
		response.ErrorNoData(c, "获取笔记最近阅读失败")
		return
	}

	// 更新数据
	success, err := api.service.UpdateNoteLatestRead(c.Request.Context(), noteLatestRead)
	if err != nil {
		api.logger.Error("msg", "更新笔记最近阅读失败", "error", err.Error())
		response.ErrorNoData(c, "更新笔记最近阅读失败")
		return
	}

	resp := &pb.UpdateNoteLatestReadResponse{
		Success: success,
	}

	response.Success(c, "success", resp)
}

// DeleteNoteLatestRead 删除笔记最近阅读
func (api *NoteLatestReadAPI) DeleteNoteLatestRead(c *gin.Context) {
	req := &pb.DeleteNoteLatestReadRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析删除笔记最近阅读请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析删除笔记最近阅读请求失败")
		return
	}

	// 直接使用请求的上下文，它已经在中间件中被更新
	_, err := api.service.DeleteNoteLatestReadById(c.Request.Context(), req.Id)
	if err != nil {
		api.logger.Error("msg", "删除笔记最近阅读失败", "error", err.Error())
		response.ErrorNoData(c, "删除笔记最近阅读失败")
		return
	}

	resp := &pb.DeleteNoteLatestReadResponse{
		Success: true,
	}

	response.Success(c, "success", resp)
}

// GetNoteLatestReadByNoteId 根据笔记ID获取笔记最近阅读
func (api *NoteLatestReadAPI) GetNoteLatestReadByNoteId(c *gin.Context) {
	req := &pb.GetNoteLatestReadByNoteIdRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析获取笔记最近阅读请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析获取笔记最近阅读请求失败")
		return
	}

	// 直接使用请求的上下文，它已经在中间件中被更新
	noteLatestRead, err := api.service.GetNoteLatestReadByNoteId(c.Request.Context(), req.NoteId)
	if err != nil {
		api.logger.Error("msg", "获取笔记最近阅读失败", "error", err.Error())
		response.ErrorNoData(c, "获取笔记最近阅读失败")
		return
	}

	resp := &pb.GetNoteLatestReadByNoteIdResponse{
		NoteLatestRead: &pb.NoteLatestRead{
			Id:     noteLatestRead.Id,
			NoteId: noteLatestRead.NoteId,
		},
	}

	response.Success(c, "success", resp)
}
