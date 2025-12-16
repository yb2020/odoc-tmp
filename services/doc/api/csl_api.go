package api

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	pb "github.com/yb2020/odoc/proto/gen/go/doc"
	"github.com/yb2020/odoc/services/doc/service"
)

// CslAPI 引用样式API处理器
type CslAPI struct {
	service        *service.CslService
	userDocService *service.UserDocService
	logger         logging.Logger
	tracer         opentracing.Tracer
}

// NewCslAPI 创建新的引用样式API处理器
func NewCslAPI(
	service *service.CslService,
	userDocService *service.UserDocService,
	logger logging.Logger,
	tracer opentracing.Tracer,
) *CslAPI {
	return &CslAPI{
		service:        service,
		userDocService: userDocService,
		logger:         logger,
		tracer:         tracer,
	}
}

// GetDefaultCslList 获取默认的引用样式列表
func (api *CslAPI) GetDefaultCslList(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "CslAPI.GetDefaultCslList")
	defer span.Finish()

	// 判断是否是国际版   todo:  这里暂时缺少方法   需要引入  localizer   i18n.Localizer
	isI18n := true
	// 调用服务层方法
	resp, err := api.service.GetDefaultCslList(ctx, isI18n)
	if err != nil {
		response.ErrorNoData(c, "get default citation styles failed")
		return
	}
	// 返回成功响应
	response.Success(c, "get default citation styles successfully", resp)
}

// GetDocMetaInfo 获取文档元数据信息
func (api *CslAPI) GetDocMetaInfo(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "CslAPI.GetDocMetaInfo")
	defer span.Finish()

	// 解析请求参数
	req := &pb.DocMetaInfoHandlerReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "invalid request parameters")
		return
	}
	// 调用服务层方法
	resp, err := api.service.GetDocMetaInfo(ctx, req)
	if err != nil {
		response.ErrorNoData(c, "get document metadata failed")
		return
	}
	// 返回成功响应
	response.Success(c, "get document metadata successfully", resp)
}

// GetMyCslList 获取我的引用样式列表
func (api *CslAPI) GetMyCslList(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "CslAPI.GetMyCslList")
	defer span.Finish()

	// 解析请求参数
	req := &pb.GetMyCslListReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "invalid request parameters")
		return
	}
	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())
	// 调用服务层方法
	resp, err := api.service.GetMyCslList(ctx, req, userId)
	if err != nil {
		response.ErrorNoData(c, "get my citation styles failed")
		return
	}
	// 返回成功响应
	response.Success(c, "get my citation styles successfully", resp)
}

// GetDocTypeList 获取文档类型列表
func (api *CslAPI) GetDocTypeList(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "CslAPI.GetDocTypeList")
	defer span.Finish()

	// 调用服务层方法
	resp, err := api.service.GetDocTypeList(ctx)
	if err != nil {
		response.ErrorNoData(c, "get document type list failed")
		return
	}
	// 返回成功响应
	response.Success(c, "get document type list successfully", resp)
}

// ExportBibTexByIds 根据文档ID列表导出BibTeX格式的引用
func (api *CslAPI) ExportBibTexByIds(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "CslAPI.ExportBibTexByIds")
	defer span.Finish()

	// 解析请求参数
	var req struct {
		DocIds []string `json:"docIds" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorNoData(c, "invalid request parameters")
		return
	}

	// 检查文档ID列表是否为空
	if len(req.DocIds) == 0 {
		response.ErrorNoData(c, "document id list cannot be empty")
		return
	}

	// 调用服务层方法导出BibTeX
	err := api.service.DownloadBibTex(ctx, req.DocIds, c.Writer)
	if err != nil {
		api.logger.Error("export bibtex failed", "error", err.Error())
		response.ErrorNoData(c, "export bibtex failed")
		return
	}
}

// ExportBibTexByFolderId 根据文件夹ID导出BibTeX格式的引用
func (api *CslAPI) ExportBibTexByFolderId(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "CslAPI.ExportBibTexByFolderId")
	defer span.Finish()

	// 解析请求参数
	var req struct {
		FolderId string `json:"folderId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorNoData(c, "invalid request parameters")
		return
	}

	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())

	// 构建获取文档列表的请求
	folderId := req.FolderId
	// 创建获取文档列表的请求
	sortType := pb.UserDocListSortType_LAST_ADD
	currentPage := uint32(1)
	pageSize := uint32(200)
	ascSort := false
	getDocListReq := &pb.GetDocListReq{
		SortType:    &sortType,
		CurrentPage: &currentPage,
		PageSize:    &pageSize,
		AscSort:     &ascSort,
		UserId:      &userId,
		FolderId:    &folderId,
	}
	// 调用服务获取文档列表
	getDocListResp, err := api.userDocService.GetDocList(ctx, getDocListReq)
	if err != nil {
		api.logger.Error("get doc list failed", "error", err.Error())
		response.ErrorNoData(c, "get document list failed")
		return
	}

	// 检查文档列表是否为空
	if getDocListResp == nil || len(getDocListResp.GetDocList()) == 0 {
		response.ErrorNoData(c, "document list is empty")
		return
	}

	// 提取文档ID列表
	docIds := make([]string, 0, len(getDocListResp.GetDocList()))
	for _, doc := range getDocListResp.GetDocList() {
		docIds = append(docIds, doc.GetDocId())
	}
	// 调用服务层方法导出BibTeX
	err = api.service.DownloadBibTex(ctx, docIds, c.Writer)
	if err != nil {
		api.logger.Error("export bibtex failed", "error", err.Error())
		response.ErrorNoData(c, "export bibtex failed")
		return
	}
}
