package api

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	pb "github.com/yb2020/odoc-proto/gen/go/note"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	"github.com/yb2020/odoc/services/note/model"
	"github.com/yb2020/odoc/services/note/service"
)

// NoteDrawEntityAPI 笔记绘制实体API处理器
type NoteDrawEntityAPI struct {
	service *service.NoteDrawEntityService
	logger  logging.Logger
	tracer  opentracing.Tracer
}

// NewNoteDrawEntityAPI 创建笔记绘制实体API处理器
func NewNoteDrawEntityAPI(service *service.NoteDrawEntityService, logger logging.Logger, tracer opentracing.Tracer) *NoteDrawEntityAPI {
	return &NoteDrawEntityAPI{
		service: service,
		logger:  logger,
		tracer:  tracer,
	}
}

// CreateNoteDrawEntity 创建笔记绘制实体
func (api *NoteDrawEntityAPI) CreateNoteDrawEntity(c *gin.Context) {
	req := &pb.CreateNoteDrawEntityRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析创建笔记绘制实体请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析创建笔记绘制实体请求失败")
		return
	}

	// 转换请求参数
	noteDrawEntity := &model.NoteDrawEntity{
		LineHexColor:  req.LineHexColor,
		LineAlpha:     float32(req.LineAlpha),
		Points:        req.Points,
		LineWidth:     float32(req.LineWidth),
		ToolType:      model.ToolType(req.ToolType),
		DrawShapeType: model.DrawShapeType(req.DrawShapeType),
		NoteId:        req.NoteId,
		PageNumber:    int(req.PageNumber),
	}

	// 直接使用请求的上下文，它已经在中间件中被更新
	id, err := api.service.CreateNoteDrawEntity(c.Request.Context(), noteDrawEntity)
	if err != nil {
		api.logger.Error("msg", "创建笔记绘制实体失败", "error", err.Error())
		response.ErrorNoData(c, "创建笔记绘制实体失败")
		return
	}

	resp := &pb.CreateNoteDrawEntityResponse{
		Id: id,
	}

	response.Success(c, "success", resp)
}

// GetNoteDrawEntity 获取笔记绘制实体
func (api *NoteDrawEntityAPI) GetNoteDrawEntity(c *gin.Context) {
	req := &pb.GetNoteDrawEntityRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析获取笔记绘制实体请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析获取笔记绘制实体请求失败")
		return
	}

	// 直接使用请求的上下文，它已经在中间件中被更新
	noteDrawEntity, err := api.service.GetNoteDrawEntityById(c.Request.Context(), req.Id)
	if err != nil {
		api.logger.Error("msg", "获取笔记绘制实体失败", "error", err.Error())
		response.ErrorNoData(c, "获取笔记绘制实体失败")
		return
	}

	resp := &pb.GetNoteDrawEntityResponse{
		NoteDrawEntity: &pb.NoteDrawEntity{
			Id:         noteDrawEntity.Id,
			NoteId:     noteDrawEntity.NoteId,
			PageNumber: uint32(noteDrawEntity.PageNumber),
		},
	}

	response.Success(c, "success", resp)
}

// UpdateNoteDrawEntity 更新笔记绘制实体
func (api *NoteDrawEntityAPI) UpdateNoteDrawEntity(c *gin.Context) {
	req := &pb.UpdateNoteDrawEntityRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析更新笔记绘制实体请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析更新笔记绘制实体请求失败")
		return
	}

	// 先获取原始数据
	noteDrawEntity, err := api.service.GetNoteDrawEntityById(c.Request.Context(), req.Id)
	if err != nil {
		api.logger.Error("msg", "获取笔记绘制实体失败", "error", err.Error())
		response.ErrorNoData(c, "获取笔记绘制实体失败")
		return
	}

	// 更新数据
	success, err := api.service.UpdateNoteDrawEntity(c.Request.Context(), noteDrawEntity)
	if err != nil {
		api.logger.Error("msg", "更新笔记绘制实体失败", "error", err.Error())
		response.ErrorNoData(c, "更新笔记绘制实体失败")
		return
	}

	resp := &pb.UpdateNoteDrawEntityResponse{
		Success: success,
	}

	response.Success(c, "success", resp)
}

// DeleteNoteDrawEntity 删除笔记绘制实体
func (api *NoteDrawEntityAPI) DeleteNoteDrawEntity(c *gin.Context) {
	req := &pb.DeleteNoteDrawEntityRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析删除笔记绘制实体请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析删除笔记绘制实体请求失败")
		return
	}

	// 直接使用请求的上下文，它已经在中间件中被更新
	_, err := api.service.DeleteNoteDrawEntityById(c.Request.Context(), req.Id)
	if err != nil {
		api.logger.Error("msg", "删除笔记绘制实体失败", "error", err.Error())
		response.ErrorNoData(c, "删除笔记绘制实体失败")
		return
	}

	resp := &pb.DeleteNoteDrawEntityResponse{
		Success: true,
	}

	response.Success(c, "success", resp)
}

// GetNoteDrawEntitiesByNoteId 根据笔记ID获取笔记绘制实体列表
func (api *NoteDrawEntityAPI) GetNoteDrawEntitiesByNoteId(c *gin.Context) {
	req := &pb.GetNoteDrawEntitiesByNoteIdRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析获取笔记绘制实体列表请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析获取笔记绘制实体列表请求失败")
		return
	}

	// 直接使用请求的上下文，它已经在中间件中被更新
	noteDrawEntities, err := api.service.GetNoteDrawEntitiesByNoteId(c.Request.Context(), req.NoteId)
	if err != nil {
		api.logger.Error("msg", "获取笔记绘制实体列表失败", "error", err.Error())
		response.ErrorNoData(c, "获取笔记绘制实体列表失败")
		return
	}

	// 转换为响应格式
	pbNoteDrawEntities := make([]*pb.NoteDrawEntity, 0, len(noteDrawEntities))
	for _, noteDrawEntity := range noteDrawEntities {
		pbNoteDrawEntities = append(pbNoteDrawEntities, &pb.NoteDrawEntity{
			Id:         noteDrawEntity.Id,
			NoteId:     noteDrawEntity.NoteId,
			PageNumber: uint32(noteDrawEntity.PageNumber),
		})
	}

	resp := &pb.GetNoteDrawEntitiesByNoteIdResponse{
		NoteDrawEntities: pbNoteDrawEntities,
	}

	response.Success(c, "success", resp)
}

// GetNoteDrawEntitiesByNoteIdAndPage 根据笔记ID和页码获取笔记绘制实体列表
func (api *NoteDrawEntityAPI) GetNoteDrawEntitiesByNoteIdAndPage(c *gin.Context) {
	req := &pb.GetNoteDrawEntitiesByNoteIdAndPageRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析获取笔记绘制实体列表请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析获取笔记绘制实体列表请求失败")
		return
	}

	// 直接使用请求的上下文，它已经在中间件中被更新
	noteDrawEntities, err := api.service.GetNoteDrawEntitiesByNoteIdAndPage(c.Request.Context(), req.NoteId, int(req.PageNumber))
	if err != nil {
		api.logger.Error("msg", "获取笔记绘制实体列表失败", "error", err.Error())
		response.ErrorNoData(c, "获取笔记绘制实体列表失败")
		return
	}

	// 转换为响应格式
	pbNoteDrawEntities := make([]*pb.NoteDrawEntity, 0, len(noteDrawEntities))
	for _, noteDrawEntity := range noteDrawEntities {
		pbNoteDrawEntities = append(pbNoteDrawEntities, &pb.NoteDrawEntity{
			Id:         noteDrawEntity.Id,
			NoteId:     noteDrawEntity.NoteId,
			PageNumber: uint32(noteDrawEntity.PageNumber),
		})
	}

	resp := &pb.GetNoteDrawEntitiesByNoteIdAndPageResponse{
		NoteDrawEntities: pbNoteDrawEntities,
	}

	response.Success(c, "success", resp)
}
