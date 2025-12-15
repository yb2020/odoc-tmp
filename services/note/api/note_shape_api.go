package api

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	commonPb "github.com/yb2020/odoc-proto/gen/go/common"
	pb "github.com/yb2020/odoc-proto/gen/go/note"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	noteInterface "github.com/yb2020/odoc/services/note/interfaces"
	"github.com/yb2020/odoc/services/note/service"
)

// NoteShapeAPI 笔记形状API处理器
type NoteShapeAPI struct {
	service                *service.NoteShapeService
	paperNoteSrv           noteInterface.IPaperNoteService
	paperNoteAccessService *service.PaperNoteAccessService
	logger                 logging.Logger
	tracer                 opentracing.Tracer
}

// NewNoteShapeAPI 创建笔记形状API处理器
func NewNoteShapeAPI(service *service.NoteShapeService, paperNoteSrv noteInterface.IPaperNoteService, paperNoteAccessService *service.PaperNoteAccessService, logger logging.Logger, tracer opentracing.Tracer) *NoteShapeAPI {
	return &NoteShapeAPI{
		service:                service,
		paperNoteSrv:           paperNoteSrv,
		paperNoteAccessService: paperNoteAccessService,
		logger:                 logger,
		tracer:                 tracer,
	}
}

// @api /api/note/noteShape/getList
// @method POST
// @apiDescription 获取笔记形状列表
// GetNoteShapesByNoteId 根据笔记ID获取笔记形状列表
func (api *NoteShapeAPI) GetNoteShapesByNoteId(c *gin.Context) {
	// 直接获取用户ID，因为中间件已经确保了用户已认证
	userId, _ := userContext.GetUserID(c.Request.Context())
	api.logger.Info("msg", "获取论文笔记基础信息By ID", "userID", userId)

	req := &pb.GetNoteShapesByNoteIdRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析获取笔记形状列表请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析获取笔记形状列表请求失败")
		return
	}
	api.logger.Info("msg", "GetNoteShapesByNoteId 请求参数", "request", req)

	//权限检查逻辑，用于控制用户是否可以访问笔记的形状列表

	// 获取笔记访问权限记录
	paperNoteAccess, err := api.paperNoteAccessService.GetPaperNoteAccessByNoteId(c.Request.Context(), req.NoteId)
	if err != nil {
		api.logger.Error("msg", "获取笔记访问权限记录失败", "error", err.Error())
	}

	// 获取笔记信息
	paperNote, err := api.paperNoteSrv.GetPaperNoteById(c.Request.Context(), req.NoteId)
	if err != nil {
		api.logger.Error("msg", "获取笔记信息失败", "error", err.Error())
		response.Success(c, "success", &pb.GetNoteShapesByNoteIdResponse{List: []*commonPb.ShapeAnnotation{}})
		return
	}
	if paperNote == nil {
		response.Success(c, "success", &pb.GetNoteShapesByNoteIdResponse{List: []*commonPb.ShapeAnnotation{}})
		return
	}

	// 权限检查
	if paperNoteAccess == nil {
		// 如果没有访问权限记录，检查用户是否是笔记创建者
		if userId == "0" {
			response.Success(c, "success", &pb.GetNoteShapesByNoteIdResponse{List: []*commonPb.ShapeAnnotation{}})
			return
		}
		// TODO: 这里需要修改，因为PaperNote模型中没有CreatorId字段
		// 假设笔记创建者ID存储在某个字段中，这里需要根据实际情况调整
		if paperNote.CreatorId != userId {
			response.Success(c, "success", &pb.GetNoteShapesByNoteIdResponse{List: []*commonPb.ShapeAnnotation{}})
			return
		}
	} else {
		// 如果有访问权限记录，检查是否开放状态
		if !paperNoteAccess.OpenStatus {
			if userId == "" {
				response.Success(c, "success", &pb.GetNoteShapesByNoteIdResponse{List: []*commonPb.ShapeAnnotation{}})
				return
			} else {
				// TODO: 这里需要修改，因为PaperNote模型中没有CreatorId字段
				// 假设笔记创建者ID存储在某个字段中，这里需要根据实际情况调整
				// isAdmin := false // 这里需要实现管理员检查逻辑
				// if paperNote.CreatorId != userCtx.UserId && !isAdmin {
				// 	response.Success(c, "success", &pb.GetNoteShapesByNoteIdResponse{List: []*pb.NoteShape{}})
				// 	return
				// }
				if paperNote.CreatorId != userId {
					response.Success(c, "success", &pb.GetNoteShapesByNoteIdResponse{List: []*commonPb.ShapeAnnotation{}})
					return
				}
			}
		}
	}

	// 通过权限检查后，继续处理请求
	// 直接使用请求的上下文，它已经在中间件中被更新
	noteShapes, err := api.service.GetNoteShapesByNoteId(c.Request.Context(), req.NoteId)
	if err != nil {
		api.logger.Error("msg", "获取笔记形状列表失败", "error", err.Error())
		response.ErrorNoData(c, "获取笔记形状列表失败")
		return
	}

	api.logger.Info("msg", "GetNoteShapesByNoteId", "noteShapes", noteShapes)

	// 转换为响应格式
	pbNoteShapes := make([]*commonPb.ShapeAnnotation, 0, len(noteShapes))
	for _, noteShape := range noteShapes {
		pbNoteShapes = append(pbNoteShapes, &commonPb.ShapeAnnotation{
			Uuid:        noteShape.UUID,
			Type:        commonPb.ShapeType(commonPb.ShapeType_value[string(noteShape.Type)]),
			X:           noteShape.X,
			Y:           noteShape.Y,
			StrokeColor: commonPb.AnnotationColor(commonPb.AnnotationColor_value[string(noteShape.StrokeColor)]),
			Width:       noteShape.Width,
			Height:      noteShape.Height,
			RadiusX:     noteShape.RadiusX,
			RadiusY:     noteShape.RadiusY,
			EndX:        noteShape.EndX,
			EndY:        noteShape.EndY,
			PageNumber:  uint32(noteShape.PageNumber),
			ShapeId:     &noteShape.Id,
		})
	}

	resp := &pb.GetNoteShapesByNoteIdResponse{
		List: pbNoteShapes,
	}

	response.Success(c, "success", resp)
}

/**
 * @api_path: /api/note/noteShape/save
 * @method: POST
 * @content-type: application/json
 * @summary: 保存笔记形状
 */
func (api *NoteShapeAPI) SaveShapeRequest(c *gin.Context) {

	req := &pb.SaveShapeRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析保存笔记形状请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析保存笔记形状请求失败")
		return
	}
	api.logger.Info("msg", "SaveShapeRequest 请求参数", "request", req)

	noteShapeId, err := api.service.SaveNoteShape(c.Request.Context(), req.NoteId, req.ShapeAnnotation)
	if err != nil {
		api.logger.Error("msg", "保存笔记形状失败", "error", err.Error())
		response.ErrorNoData(c, "保存笔记形状失败")
		return
	}

	saveShapeResponse := &pb.SaveShapeResponse{
		ShapeId: noteShapeId,
	}

	response.Success(c, "success", saveShapeResponse)
}

/**
 * @api_path: /api/note/noteShape/delete
 * @method: POST
 * @content-type: application/json
 * @summary: 删除笔记形状
 */
func (api *NoteShapeAPI) DeleteShapeRequest(c *gin.Context) {

	req := &pb.DeleteShapeRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析删除笔记形状请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析删除笔记形状请求失败")
		return
	}

	// 实现删除笔记形状的逻辑
	if len(req.ShapeIds) == 0 {
		response.ErrorNoData(c, "shapeIds can not be empty")
		return
	}

	// 将[]uint64转换为[]int64
	ids := make([]string, len(req.ShapeIds))
	for i, id := range req.ShapeIds {
		ids[i] = id
	}

	// 调用服务层方法
	_, err := api.service.DeleteNoteShapeByIds(c.Request.Context(), ids)
	if err != nil {
		api.logger.Error("msg", "删除笔记形状失败", "error", err.Error())
		response.ErrorNoData(c, "删除笔记形状失败")
		return
	}

	response.SuccessNoData(c, "success")
}

/**
 * @api_path: /api/note/noteShape/update
 * @method: POST
 * @content-type: application/json
 * @summary: 更新笔记形状
 */
func (api *NoteShapeAPI) UpdateShapeRequest(c *gin.Context) {

	req := &pb.UpdateShapeRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析更新笔记形状请求失败", "error", err.Error())
		response.ErrorNoData(c, "解析更新笔记形状请求失败")
		return
	}

	if len(req.Annotations) == 0 {
		response.ErrorNoData(c, "annotations can not be empty")
		return
	}

	// TODO: 实现更新笔记形状的逻辑
	_, err := api.service.UpdateNoteShapeAnnotations(c.Request.Context(), req.Annotations)
	if err != nil {
		api.logger.Error("msg", "更新笔记形状失败", "error", err.Error())
		response.ErrorNoData(c, "更新笔记形状失败")
		return
	}

	response.SuccessNoData(c, "success")
}
