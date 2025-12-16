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

// PaperNoteAccessAPI 论文笔记访问记录API处理器
type PaperNoteAccessAPI struct {
	service *service.PaperNoteAccessService
	logger  logging.Logger
	tracer  opentracing.Tracer
}

// NewPaperNoteAccessAPI 创建论文笔记访问记录API处理器
func NewPaperNoteAccessAPI(service *service.PaperNoteAccessService, logger logging.Logger, tracer opentracing.Tracer) *PaperNoteAccessAPI {
	return &PaperNoteAccessAPI{
		service: service,
		logger:  logger,
		tracer:  tracer,
	}
}

// CreatePaperNoteAccess 创建论文笔记访问记录
func (api *PaperNoteAccessAPI) CreatePaperNoteAccess(c *gin.Context) {
	req := &pb.CreatePaperNoteAccessRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析创建论文笔记访问记录请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析创建论文笔记访问记录请求失败")
		return
	}

	// 转换请求参数
	paperNoteAccess := &model.PaperNoteAccess{
		NoteId:     req.NoteId,
		OpenStatus: req.OpenStatus,
	}

	// 直接使用请求的上下文，它已经在中间件中被更新
	id, err := api.service.CreatePaperNoteAccess(c.Request.Context(), paperNoteAccess)
	if err != nil {
		api.logger.Error("msg", "创建论文笔记访问记录失败", "error", err.Error())
		response.ErrorNoData(c, "创建论文笔记访问记录失败")
		return
	}

	resp := &pb.CreatePaperNoteAccessResponse{
		Id: id,
	}

	response.Success(c, "success", resp)
}

// GetPaperNoteAccess 获取论文笔记访问记录
func (api *PaperNoteAccessAPI) GetPaperNoteAccess(c *gin.Context) {
	req := &pb.GetPaperNoteAccessRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析获取论文笔记访问记录请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析获取论文笔记访问记录请求失败")
		return
	}

	// 直接使用请求的上下文，它已经在中间件中被更新
	paperNoteAccess, err := api.service.GetPaperNoteAccessById(c.Request.Context(), req.Id)
	if err != nil {
		api.logger.Error("msg", "获取论文笔记访问记录失败", "error", err.Error())
		response.ErrorNoData(c, "获取论文笔记访问记录失败")
		return
	}

	resp := &pb.GetPaperNoteAccessResponse{
		Access: &pb.PaperNoteAccess{
			Id:         paperNoteAccess.Id,
			NoteId:     paperNoteAccess.NoteId,
			OpenStatus: paperNoteAccess.OpenStatus,
		},
	}

	response.Success(c, "success", resp)
}

// GetPaperNoteAccessByNoteId 根据笔记ID获取论文笔记访问记录
func (api *PaperNoteAccessAPI) GetPaperNoteAccessByNoteId(c *gin.Context) {
	req := &pb.GetPaperNoteAccessByNoteIdRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析获取论文笔记访问记录请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析获取论文笔记访问记录请求失败")
		return
	}

	// 直接使用请求的上下文，它已经在中间件中被更新
	paperNoteAccess, err := api.service.GetPaperNoteAccessByNoteId(c.Request.Context(), req.NoteId)
	if err != nil {
		api.logger.Error("msg", "获取论文笔记访问记录失败", "error", err.Error())
		response.ErrorNoData(c, "获取论文笔记访问记录失败")
		return
	}

	resp := &pb.GetPaperNoteAccessByNoteIdResponse{
		Access: &pb.PaperNoteAccess{
			Id:         paperNoteAccess.Id,
			NoteId:     paperNoteAccess.NoteId,
			OpenStatus: paperNoteAccess.OpenStatus,
		},
	}

	response.Success(c, "success", resp)
}

// UpdatePaperNoteAccess 更新论文笔记访问记录
func (api *PaperNoteAccessAPI) UpdatePaperNoteAccess(c *gin.Context) {
	req := &pb.UpdatePaperNoteAccessRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析更新论文笔记访问记录请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析更新论文笔记访问记录请求失败")
		return
	}

	// 先获取原始数据
	paperNoteAccess, err := api.service.GetPaperNoteAccessById(c.Request.Context(), req.Id)
	if err != nil {
		api.logger.Error("msg", "获取论文笔记访问记录失败", "error", err.Error())
		response.ErrorNoData(c, "获取论文笔记访问记录失败")
		return
	}

	// 更新字段
	paperNoteAccess.OpenStatus = req.OpenStatus

	// 更新数据
	success, err := api.service.UpdatePaperNoteAccess(c.Request.Context(), paperNoteAccess)
	if err != nil {
		api.logger.Error("msg", "更新论文笔记访问记录失败", "error", err.Error())
		response.ErrorNoData(c, "更新论文笔记访问记录失败")
		return
	}

	resp := &pb.UpdatePaperNoteAccessResponse{
		Success: success,
	}

	response.Success(c, "success", resp)
}

// DeletePaperNoteAccess 删除论文笔记访问记录
func (api *PaperNoteAccessAPI) DeletePaperNoteAccess(c *gin.Context) {
	req := &pb.DeletePaperNoteAccessRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析删除论文笔记访问记录请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析删除论文笔记访问记录请求失败")
		return
	}

	// 直接使用请求的上下文，它已经在中间件中被更新
	err := api.service.DeletePaperNoteAccessById(c.Request.Context(), req.Id)
	if err != nil {
		api.logger.Error("msg", "删除论文笔记访问记录失败", "error", err.Error())
		response.ErrorNoData(c, "删除论文笔记访问记录失败")
		return
	}

	resp := &pb.DeletePaperNoteAccessResponse{
		Success: true,
	}

	response.Success(c, "success", resp)
}
