package api

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	commonPb "github.com/yb2020/odoc/proto/gen/go/common"
	notePb "github.com/yb2020/odoc/proto/gen/go/note"
	pb "github.com/yb2020/odoc/proto/gen/go/pdf"
	"github.com/yb2020/odoc/services/pdf/dto"
	"github.com/yb2020/odoc/services/pdf/interfaces"
)

// PdfMarkAPI 论文PDF标记API处理器
type PdfMarkAPI struct {
	pdfMarkService interfaces.IPdfMarkService
	logger         logging.Logger
	tracer         opentracing.Tracer
}

// NewPdfMarkAPI 创建PDF标记API处理器
func NewPdfMarkAPI(
	pdfMarkService interfaces.IPdfMarkService,
	logger logging.Logger,
	tracer opentracing.Tracer,

) *PdfMarkAPI {
	return &PdfMarkAPI{
		pdfMarkService: pdfMarkService,
		logger:         logger,
		tracer:         tracer,
	}
}

/*
* @api_path: /api/pdf/pdfMark/v3/web/getByNote
* @method: GET
* @content-type: application/json
* @summary: 根据笔记ID查询笔记-摘要
 */
func (api *PdfMarkAPI) GetNoteAnnotationListByNoteId(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "PaperPdfAPI.GetNoteAnnotationListByNoteId")
	defer span.Finish()

	var req pb.GetNoteAnnotationListByNoteIdRequest
	if err := transport.BindProto(c, &req); err != nil {
		api.logger.Error("msg", "解析请求参数失败", "error", err.Error())
		response.ErrorNoData(c, "解析请求参数失败")
		return
	}

	// // 查询
	// annotateionRawModels, err := api.pdfMarkService.GetAnnotationRawModelsByNoteId(ctx, int64(req.NoteId))
	// if err != nil {

	// }
	// api.logger.Debug("list", annotateionRawModels)

	// webAnnotationModels := make([]*notePb.WebNoteAnnotationModel, 0, len(annotateionRawModels))
	// for _, annotation := range annotateionRawModels {
	// 	webNoteProtoTransformer := &proto.WebNoteProtoTransformer{}
	// 	webAnnotation, _ := webNoteProtoTransformer.WebAnnotationModel(annotation)

	// 	//TODO mock
	// 	// 1.设置DeleteAuthority属性判断和设置
	// 	webAnnotation.DeleteAuthority = true
	// 	// 2.设置关联的tags
	// 	//var markId int64 = 2840323191133811456
	// 	tags, _ := api.pdfMarkService.GetAnnotateTagsByMarkId(ctx, int64(webAnnotation.Id))
	// 	webAnnotation.Tags = tags
	// 	// 3.设置CommentatorInfoView信息

	// 	webAnnotationModels = append(webAnnotationModels, webAnnotation)
	// }

	// var resp = &pb.GetNoteAnnotationListByNoteIdResponse{
	// 	Annotations: webAnnotationModels,
	// }

	// // 按照实际数据结构返回标注列表
	// response.Success(c, "success", resp)

	webAnnotationModels, err := api.pdfMarkService.GetWebNoteAnnotationModelsByNoteId(ctx, req.NoteId)
	if err != nil {
		api.logger.Error("获取笔记标注列表失败", "error", err)
		response.ErrorNoData(c, "获取笔记标注列表失败")
		return
	}
	var resp = &pb.GetNoteAnnotationListByNoteIdResponse{
		Annotations: webAnnotationModels,
	}
	// 按照实际数据结构返回标注列表
	response.Success(c, "success", resp)

}

/*
* 旧接口在microservice-note微服务
* @api_path: /api/pdf/pdfMark/v3/web/draw/getByNote
* @method: GET
* @content-type: application/json
* @summary: 获取IOS笔记标注列表
 */
func (api *PdfMarkAPI) GetDrawNoteAnnotationListByNoteId(c *gin.Context) {
	span, _ := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "PaperPdfAPI.GetDrawNoteAnnotationListByNoteId")
	defer span.Finish()

	userId, _ := userContext.GetUserID(c.Request.Context())
	api.logger.Info("msg", "获取PDF状态信息", "userId", userId)

	var req pb.GetDrawNoteAnnotationListByNoteIdRequest
	if err := transport.BindProto(c, &req); err != nil {
		api.logger.Error("msg", "解析请求参数失败", "error", err.Error())
		response.ErrorNoData(c, "解析请求参数失败")
		return
	}

	// 1.判断PaperNote访问权限

	// 2.获取webDraws列表

	// TODO: 生成静态模拟数据
	list := []*notePb.WebNoteAnnotationModel{}

	var resp = &pb.GetDrawNoteAnnotationListByNoteIdResponse{
		Annotations: list,
	}

	// 按照实际数据结构返回标注列表
	response.Success(c, "success", resp)
}

/*
* @api_path: /api/pdf/pdfMark/v2/web/hotSelect
* @method: GET
* @content-type: application/json
* @summary: 获取热选笔记标注列表
 */
func (api *PdfMarkAPI) HotSelect(c *gin.Context) {
	span, _ := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "PaperPdfAPI.HotSelect")
	defer span.Finish()

	userId, _ := userContext.GetUserID(c.Request.Context())
	api.logger.Info("msg", "获取PDF状态信息", "userId", userId)

	var req pb.HotSelectRequest
	if err := transport.BindProto(c, &req); err != nil {
		api.logger.Error("msg", "解析请求参数失败", "error", err.Error())
		response.ErrorNoData(c, "解析请求参数失败")
		return
	}

	// TODO: 新版本暂无此功能，使用模拟返回空数据
	list := []*notePb.WebNoteAnnotationModel{}

	var resp = &pb.HotSelectResponse{
		Annotations: list,
	}

	response.Success(c, "success", resp)
}

// @api_path: /api/pdf/pdfMark/v2/web/save
// @method post
// @summary 保存标注笔记
func (api *PdfMarkAPI) SavePdfMark(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "PaperPdfAPI.SavePdfMark")
	defer span.Finish()

	var req notePb.WebNoteAnnotationModel
	if err := transport.BindProto(c, &req); err != nil {
		api.logger.Error("msg", "解析请求参数失败", "error", err.Error())
		response.ErrorNoData(c, "解析请求参数失败")
		return
	}

	id, err := api.pdfMarkService.SavePdfMarkByAnnotation(ctx, &req)
	if err != nil {
		api.logger.Error("msg", "保存PDF标记失败", "error", err.Error())
		response.ErrorNoData(c, "保存PDF标记失败")
		return
	}
	api.logger.Info("msg", "保存PDF标记成功", "id", id)

	response.Success(c, "success", &pb.SavePdfMarkResponse{Uuid: id})
}

// @api_path: /api/pdf/pdfMark/v2/web/update
// @method post
// @summary 更新标注笔记
func (api *PdfMarkAPI) UpdatePdfMark(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "PaperPdfAPI.UpdatePdfMark")
	defer span.Finish()

	var req notePb.WebNoteAnnotationModel
	if err := transport.BindProto(c, &req); err != nil {
		api.logger.Error("msg", "解析请求参数失败", "error", err.Error())
		response.ErrorNoData(c, "解析请求参数失败")
		return
	}

	id, err := api.pdfMarkService.UpdatePdfMarkByAnnotation(ctx, &req)
	if err != nil {
		api.logger.Error("msg", "更新PDF标记失败", "error", err.Error())
		response.ErrorNoData(c, "更新PDF标记失败")
		return
	}

	response.Success(c, "success", &pb.UpdatePdfMarkResponse{Uuid: id})
}

// @api_path: /api/pdf/pdfMark/v2/web/delete
// @method post
// @summary 删除标注笔记
func (api *PdfMarkAPI) DeletePdfMark(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "PaperPdfAPI.DeletePdfMark")
	defer span.Finish()

	var req commonPb.AnnotationPointer
	if err := transport.BindProto(c, &req); err != nil {
		api.logger.Error("msg", "解析请求参数失败", "error", err.Error())
		response.ErrorNoData(c, "解析请求参数失败")
		return
	}

	// 参数校验
	if req.Id == "0" {
		response.ErrorNoData(c, "参数id不能为空")
		return
	}

	_, err := api.pdfMarkService.DeleteByAnnotationPointer(ctx, &req)
	if err != nil {
		api.logger.Error("msg", "删除标注笔记", "error", err.Error())
		response.ErrorNoData(c, "删除标注笔记")
		return
	}

	response.Success(c, "success", nil)
}

/**
 * @api_path: /api/pdf/pdfMark/v2/web/getMyNoteMarkList
 * @method: POST
 * @content-type: application/json
 * @summary: 查询我的笔记列表
 */
func (api *PdfMarkAPI) GetMyNoteMarkListReq(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "PaperPdfAPI.GetMyNoteMarkListReq")
	defer span.Finish()

	req := &pb.GetMyNoteMarkListReq{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析请求参数失败", "error", err.Error())
		response.ErrorNoData(c, "解析请求参数失败")
		return
	}
	api.logger.Info("msg", "获取笔记标记列表", "req", req)

	// req 转换成 dto
	searchDto := &dto.PdfMarkSearchPageDto{}
	// 处理基本类型转换
	if req.FolderId != nil && *req.FolderId != "0" {
		searchDto.FolderId = *req.FolderId
	}
	searchDto.SearchContent = req.SearchContent
	if req.SortType > 0 {
		searchDto.SortType = int32(req.SortType)
	}
	if req.DocId != nil && *req.DocId != "0" {
		searchDto.DocId = *req.DocId
	}

	// 处理分页参数
	if req.CurrentPage > 0 {
		searchDto.CurrentPage = int32(req.CurrentPage)
	} else {
		searchDto.CurrentPage = 1 // 默认第一页
	}
	if req.PageSize > 0 {
		searchDto.PageSize = int32(req.PageSize)
	} else {
		searchDto.PageSize = 10 // 默认每页10条
	}

	// 处理数组类型转换
	if len(req.TagIdList) > 0 {
		searchDto.TagIdList = make([]string, len(req.TagIdList))
		for i, id := range req.TagIdList {
			searchDto.TagIdList[i] = id
		}
	}
	if len(req.StyleIdList) > 0 {
		searchDto.StyleIdList = make([]int64, len(req.StyleIdList))
		for i, id := range req.StyleIdList {
			searchDto.StyleIdList[i] = int64(id)
		}
	}

	userId, _ := userContext.GetUserID(ctx)
	markList, total, err := api.pdfMarkService.GetUserPdfMarkPage(ctx, userId, searchDto)
	if err != nil {
		api.logger.Error("msg", "获取笔记标记列表失败", "error", err.Error())
		response.ErrorNoData(c, "获取笔记标记列表失败")
		return
	}
	var resp = &pb.GetMyNoteMarkListResponse{
		AnnotationModelList: markList,
		Total:               total,
	}
	response.Success(c, "success", resp)
}
