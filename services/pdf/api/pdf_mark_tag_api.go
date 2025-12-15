package api

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	notePb "github.com/yb2020/odoc-proto/gen/go/note"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	"github.com/yb2020/odoc/services/pdf/service"
)

// PdfMarkTagAPI 论文PDF API处理器
type PdfMarkTagAPI struct {
	pdfMarkTagService *service.PdfMarkTagService
	logger            logging.Logger
	tracer            opentracing.Tracer
}

// NewPdfMarkTagAPI 创建论文PDF API处理器
func NewPdfMarkTagAPI(
	pdfMarkTagService *service.PdfMarkTagService,
	logger logging.Logger,
	tracer opentracing.Tracer,
) *PdfMarkTagAPI {
	return &PdfMarkTagAPI{
		pdfMarkTagService: pdfMarkTagService,
		logger:            logger,
		tracer:            tracer,
	}
}

// @api_path: /api/pdf/marktag/tags
// @method: GET
// @content-type: application/json
// @summary: 获取PDF标记标签列表
func (api *PdfMarkTagAPI) GetPdfMarkTagsRequest(c *gin.Context) {
	req := &notePb.GetAnnotateTagsRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "parse request parameters failed", "error", err.Error())
		response.ErrorNoData(c, "parse request parameters failed")
		return
	}

	userId, _ := userContext.GetUserID(c.Request.Context())
	api.logger.Info("msg", "获取PDF状态信息", "userId", userId)
	// 获取PDF标记标签
	tags, err := api.pdfMarkTagService.GetTagsByUserId(c.Request.Context(), userId, req.OnlyUsed)
	if err != nil {
		api.logger.Error("get pdf mark tag failed", "error", err)
		response.ErrorNoData(c, "get pdf mark tag failed")
		return
	}

	resp := &notePb.GetAnnotateTagsResponse{
		Tags: tags,
	}

	response.Success(c, "success", resp)
}

/**
 * @api_path: /api/pdf/marktag/save
 * @method: POST
 * @content-type: application/json
 * @summary: 创建标签信息
 */
func (api *PdfMarkTagAPI) CreatePdfMarkTagRequest(c *gin.Context) {
	req := &notePb.CreateAnnotateTagRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "parse request parameters failed", "error", err.Error())
		response.ErrorNoData(c, "parse request parameters failed")
		return
	}

	// 创建PDF标记标签
	id, err := api.pdfMarkTagService.SavePdfMarkTag(c.Request.Context(), req.MarkId, req.TagName)
	if err != nil {
		api.logger.Error("create pdf mark tag failed", "error", err)
		response.ErrorNoData(c, "create pdf mark tag failed")
		return
	}

	resp := &notePb.CreateAnnotateTagResponse{
		TagId: id,
	}

	response.Success(c, "success", resp)
}

/**
 * @api_path: /api/pdf/marktag/update
 * @method: POST
 * @content-type: application/json
 * @summary: 修改标签信息
 */
func (api *PdfMarkTagAPI) RenameAnnotateTagRequest(c *gin.Context) {
	req := &notePb.RenameAnnotateTagRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "parse request parameters failed", "error", err.Error())
		response.ErrorNoData(c, "parse request parameters failed")
		return
	}

	// 修改PDF标记标签
	ok, err := api.pdfMarkTagService.RenameTagById(c.Request.Context(), req.TagId, req.TagName)
	if err != nil {
		api.logger.Error("update pdf mark tag failed", "error", err)
		response.ErrorNoData(c, "update pdf mark tag failed")
		return
	}

	if !ok {
		response.ErrorNoData(c, "update pdf mark tag failed")
		return
	}

	resp := &notePb.RenameAnnotateTagResponse{}

	response.Success(c, "success", resp)
}

/**
 * @api_path: /api/pdf/marktag/delete
 * @method: POST
 * @content-type: application/json
 * @summary: 删除标签信息
 */
func (api *PdfMarkTagAPI) DeleteAnnotateTagRequest(c *gin.Context) {
	req := &notePb.DeleteAnnotateTagRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "parse request parameters failed", "error", err.Error())
		response.ErrorNoData(c, "parse request parameters failed")
		return
	}

	// 删除PDF标记标签
	ok, err := api.pdfMarkTagService.DeleteTagById(c.Request.Context(), req.TagId)
	if err != nil {
		api.logger.Error("delete pdf mark tag failed", "error", err)
		response.ErrorNoData(c, "delete pdf mark tag failed")
		return
	}

	if !ok {
		response.ErrorNoData(c, "delete pdf mark tag failed")
		return
	}

	resp := &notePb.DeleteAnnotateTagResponse{}

	response.Success(c, "success", resp)
}

/**
 * @api_path: /api/pdf/marktag/relation/mark/save
 * @method: POST
 * @content-type: application/json
 * @summary: 标注批量添加标签
 */
func (api *PdfMarkTagAPI) AddTagToAnnotateRequest(c *gin.Context) {
	req := &notePb.AddTagToAnnotateRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "parse request parameters failed", "error", err.Error())
		response.ErrorNoData(c, "parse request parameters failed")
		return
	}
	if req.MarkId == "" {
		api.logger.Error("msg", "parse request parameters failed", "error", "markId can not be null")
		response.ErrorNoData(c, "markId can not be null")
		return
	}
	if len(req.TagIds) == 0 {
		api.logger.Error("msg", "parse request parameters failed", "error", "tagIds can not be empty")
		response.ErrorNoData(c, "tagIds can not be empty")
		return
	}

	// 将 req.TagIds 转换为 []int64 类型的 ids
	ids := make([]string, len(req.TagIds))
	for i, tagId := range req.TagIds {
		ids[i] = tagId
	}

	// 给批注添加标签信息
	ok, err := api.pdfMarkTagService.AddTagIdsToAnnotation(c.Request.Context(), req.MarkId, ids)
	if err != nil {
		api.logger.Error("add tag to annotate failed", "error", err)
		response.ErrorNoData(c, "add tag to annotate failed")
		return
	}

	if !ok {
		response.ErrorNoData(c, "add tag to annotate failed")
		return
	}

	response.Success(c, "success", nil)
}

/**
 * @api_path: /api/pdf/marktag/relation/mark/delete
 * @method: POST
 * @content-type: application/json
 * @summary: 标注批量删除标签
 */
func (api *PdfMarkTagAPI) DeleteTagToAnnotateRequest(c *gin.Context) {
	req := &notePb.DeleteTagToAnnotateRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "parse request parameters failed", "error", err.Error())
		response.ErrorNoData(c, "parse request parameters failed")
		return
	}

	if req.MarkId == "" {
		api.logger.Error("msg", "parse request parameters failed", "error", "markId can not be null")
		response.ErrorNoData(c, "markId can not be null")
		return
	}
	if len(req.TagIds) == 0 {
		api.logger.Error("msg", "parse request parameters failed", "error", "tagIds can not be empty")
		response.ErrorNoData(c, "tagIds can not be empty")
		return
	}

	// 将 req.TagIds 转换为 []int64 类型的 ids
	ids := make([]string, len(req.TagIds))
	for i, tagId := range req.TagIds {
		ids[i] = tagId
	}

	// 删除批注的标签信息
	ok, err := api.pdfMarkTagService.DeleteTagsToAnnotate(c.Request.Context(), req.MarkId, ids)
	if err != nil {
		api.logger.Error("delete tag to annotate failed", "error", err)
		response.ErrorNoData(c, "delete tag to annotate failed")
		return
	}

	if !ok {
		response.ErrorNoData(c, "delete tag to annotate failed")
		return
	}

	response.Success(c, "success", nil)
}
