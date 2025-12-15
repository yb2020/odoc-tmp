package api

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	pb "github.com/yb2020/odoc-proto/gen/go/note"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	"github.com/yb2020/odoc/services/note/model"
	"github.com/yb2020/odoc/services/note/service"
)

// NoteReadLocationAPI 笔记阅读位置API处理器
type NoteReadLocationAPI struct {
	service *service.NoteReadLocationService
	logger  logging.Logger
	tracer  opentracing.Tracer
}

// NewNoteReadLocationAPI 创建笔记阅读位置API处理器
func NewNoteReadLocationAPI(service *service.NoteReadLocationService, logger logging.Logger, tracer opentracing.Tracer) *NoteReadLocationAPI {
	return &NoteReadLocationAPI{
		service: service,
		logger:  logger,
		tracer:  tracer,
	}
}

// @old_api_path /noteReadLocation/getLocation
// @api_path /api/note/noteReadLocation/getLocation
// @api_method POST
// @api_summary 获取笔记阅读位置
// GetNoteReadLocation 获取笔记阅读位置
func (api *NoteReadLocationAPI) GetNoteReadLocation(c *gin.Context) {

	// 直接获取用户ID，因为中间件已经确保了用户已认证
	userId, _ := userContext.GetUserID(c.Request.Context())

	req := &pb.GetNoteReadLocationRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析获取笔记阅读位置请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析获取笔记阅读位置请求失败")
		return
	}

	// 检查请求参数
	// 从REDIS获取和设置笔记阅读位置MOCK:TODO

	// 直接使用请求的上下文，它已经在中间件中被更新
	noteReadLocation, err := api.service.GetNoteReadLocationByUserIdAndNoteId(c.Request.Context(), userId, req.NoteId)
	if err != nil {
		api.logger.Error("msg", "获取笔记阅读位置失败", "error", err.Error())
		response.ErrorNoData(c, "获取笔记阅读位置失败")
		return
	}

	var location string
	if noteReadLocation != nil {
		location = noteReadLocation.Location
	}

	resp := &pb.GetNoteReadLocationResponse{
		Location: location,
	}

	response.Success(c, "success", resp)
}

// @api {post} /api/note/noteReadLocation/record 记录笔记阅读位置
// RecordNoteReadLocation 记录笔记阅读位置
func (api *NoteReadLocationAPI) RecordNoteReadLocation(c *gin.Context) {

	// 直接获取用户ID，因为中间件已经确保了用户已认证
	userId, _ := userContext.GetUserID(c.Request.Context())

	req := &pb.RecordNoteReadLocationRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析记录笔记阅读位置请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析记录笔记阅读位置请求失败")
		return
	}

	// TODO: 需要实现REDIS缓存

	// 直接使用请求的上下文，它已经在中间件中被更新
	record, err := api.service.GetNoteReadLocationByUserIdAndNoteId(c.Request.Context(), userId, req.NoteId)
	if err != nil {
		api.logger.Error("msg", "记录笔记阅读位置失败", "error", err.Error())
		response.ErrorNoData(c, "记录笔记阅读位置失败")
		return
	}
	if record != nil {
		record.Location = req.Location
		api.service.UpdateNoteReadLocation(c.Request.Context(), record)
	} else {
		record = &model.NoteReadLocation{
			NoteId:   req.NoteId,
			Location: req.Location,
		}
		api.service.CreateNoteReadLocation(c.Request.Context(), record)
	}

	response.SuccessNoData(c, "success")
}
