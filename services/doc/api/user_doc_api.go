package api

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	config "github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/cache"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	"github.com/yb2020/odoc/pkg/utils"
	pb "github.com/yb2020/odoc/proto/gen/go/doc"
	osspb "github.com/yb2020/odoc/proto/gen/go/oss"
	docBean "github.com/yb2020/odoc/services/doc/bean"
	"github.com/yb2020/odoc/services/doc/service"
	membershipInterface "github.com/yb2020/odoc/services/membership/interfaces"
	ossService "github.com/yb2020/odoc/services/oss/service"
	paperService "github.com/yb2020/odoc/services/paper/service"
	pdfInterface "github.com/yb2020/odoc/services/pdf/interfaces"
)

// UserDocAPI 用户文档API处理器
type UserDocAPI struct {
	service               *service.UserDocService
	uploadService         *service.UserDocUploadService
	uploadLocalService    *service.UserDocUploadLocalService
	logger                logging.Logger
	tracer                opentracing.Tracer
	cache                 cache.Cache
	ossService            ossService.OssServiceInterface
	docCiteSearchService  *service.DocCiteSearchService
	paperPdfService       pdfInterface.IPaperPdfService // PDF服务，使用interface{}避免循环依赖
	paperService          *paperService.PaperService
	paperPdfParsedService *paperService.PaperPdfParsedService
	membershipService     membershipInterface.IMembershipService
	config                *config.Config // 添加配置字段
}

// UploadRequest 上传请求参数   todo  这里的定义可以写到proto文件中
type UploadRequest struct {
	BucketType string `json:"bucketType" binding:"required"` // 存储桶类型
	Filename   string `json:"filename" binding:"required"`   // 文件名
}

// NewUserDocAPI 创建用户文档API处理器
func NewUserDocAPI(
	service *service.UserDocService,
	uploadService *service.UserDocUploadService,
	uploadLocalService *service.UserDocUploadLocalService,
	ossService ossService.OssServiceInterface,
	docCiteSearchService *service.DocCiteSearchService,
	cache cache.Cache,
	paperService *paperService.PaperService,
	paperPdfParsedService *paperService.PaperPdfParsedService,
	membershipService membershipInterface.IMembershipService,
	logger logging.Logger,
	tracer opentracing.Tracer,
	config *config.Config,
) *UserDocAPI {
	return &UserDocAPI{
		service:               service,
		uploadService:         uploadService,
		uploadLocalService:    uploadLocalService,
		logger:                logger,
		tracer:                tracer,
		cache:                 cache,
		ossService:            ossService,
		docCiteSearchService:  docCiteSearchService,
		paperService:          paperService,
		paperPdfParsedService: paperPdfParsedService,
		membershipService:     membershipService,
		config:                config,
	}
}

// SetPaperPdfService 设置pdf服务，用于解决循环依赖问题
func (s *UserDocAPI) SetPaperPdfService(paperPdfService pdfInterface.IPaperPdfService) error {
	if paperPdfService == nil {
		return errors.Biz("noteService cannot be nil")
	}
	s.paperPdfService = paperPdfService
	return nil
}

// GetUserUploadParseStatus 获取用户上传解析状态
func (api *UserDocAPI) GetUserUploadParseStatus(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.GetUserUploadParseStatus")
	defer span.Finish()
	req := &pb.GetUserDocCreateStatusRequest{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "参数错误")
		return
	}
	userId, _ := userContext.GetUserID(c.Request.Context())
	resp, err := api.service.GetUserUploadParseStatus(ctx, userId, req.GetToken())
	if err != nil {
		response.ErrorNoData(c, "获取状态失败")
		return
	}
	response.Success(c, "获取用户上传解析状态成功", resp)
}

// GetParseToken 获取解析token状态
func (api *UserDocAPI) GetParseToken(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.GetParseToken")
	defer span.Finish()

	userId, _ := userContext.GetUserID(c.Request.Context())
	//生成token
	token, err := utils.GenerateServiceTokenDefaultWithDefaults(fmt.Sprintf("%d", userId))
	if err != nil {
		response.ErrorNoData(c, "get token failed")
		return
	}
	//生成token默认redis状态为准备中
	if err := api.service.ChangeUploadTokenStatus(ctx, token, pb.UserDocParsedStatusEnum_READY); err != nil {
		response.ErrorNoData(c, "get token failed")
		return
	}
	pbResp := &pb.GetParseTokenResponse{
		Token: token,
	}
	response.Success(c, "get token success", pbResp)
}

// GetPdfUploadToken 处理文件上传请求
func (api *UserDocAPI) GetPdfUploadToken(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.GetPdfUploadToken")
	defer span.Finish()

	req := &osspb.GetUploadTokenRequest{}
	if err := transport.BindProto(c, req); err != nil {
		c.Error(err)
		return
	}

	// 检查文件是否为PDF格式
	if req.FileName == "" || !strings.HasSuffix(strings.ToLower(strings.TrimSpace(req.FileName)), ".pdf") {
		c.Error(errors.Biz("please upload PDF file with .pdf extension"))
		return
	}
	//获取用户的ID
	userId, _ := userContext.GetUserID(ctx)
	// 检查文件的参数
	if api.config.Service.Type == "local" {
		if req.GetLocalFilePath() == "" {
			response.ErrorNoData(c, "local file path is empty")
			return
		}
		// 本地文件上传
		resp, err := api.uploadLocalService.UploadLocalFile(ctx, userId, req.GetFileName(), req.GetLocalFilePath(), req.GetFolderId())
		if err != nil {
			c.Error(err)
			return
		}
		response.Success(c, "success", resp)
		return
	}
	if req.GetFileSHA256() == "" || req.GetFileSize() == 0 || req.GetFilePage() == 0 || req.GetFileName() == "" {
		response.ErrorNoData(c, "file sha256 or file size or file page or file name is empty")
		return
	}
	if req.GetFilePage() == 0 || req.GetFileSize() == 0 {
		response.ErrorNoData(c, "file page or file size is empty")
		return
	}

	// 普通上传（S3上传）
	resp, err := api.uploadService.GetPdfUploadToken(ctx, userId, req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, "success", resp)
}

// GetDocIndex 获取用户文档索引
func (api *UserDocAPI) GetDocIndex(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.GetDocIndex")
	defer span.Finish()
	userId, _ := userContext.GetUserID(ctx)
	resp, err := api.service.GetDocIndex(ctx, userId)
	if err != nil {
		response.ErrorNoData(c, "get doc index failed")
		return
	}
	response.Success(c, "get doc index successfully", resp)
}

// GetDocList 获取用户文档列表
func (api *UserDocAPI) GetDocList(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.GetDocList")
	defer span.Finish()
	req := &pb.GetDocListReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "parameter error")
		return
	}
	userId, _ := userContext.GetUserID(c.Request.Context())
	req.UserId = &userId
	resp, err := api.service.GetDocList(ctx, req)
	if err != nil {
		response.ErrorNoData(c, "get doc list failed")
		return
	}
	response.Success(c, "get doc list successfully", resp)
}

// GetDocRelatedVenueList 获取文档相关的发表场所列表
func (api *UserDocAPI) GetDocRelatedVenueList(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.GetDocRelatedVenueList")
	defer span.Finish()
	req := &pb.GetDocRelatedVenueListReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "parameter error")
		return
	}
	userId, _ := userContext.GetUserID(c.Request.Context())
	req.UserId = &userId
	resp, err := api.service.GetDocRelatedVenueList(ctx, req)
	if err != nil {
		response.ErrorNoData(c, "get doc related venue list failed")
		return
	}
	response.Success(c, "get doc related venue list successfully", resp)
}

// GetDocRelatedAuthorList 获取文档相关的作者列表
func (api *UserDocAPI) GetDocRelatedAuthorList(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.GetDocRelatedAuthorList")
	defer span.Finish()
	req := &pb.GetDocRelatedAuthorListReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "parameter error")
		return
	}
	userId, _ := userContext.GetUserID(c.Request.Context())
	req.UserId = &userId
	resp, err := api.service.GetDocRelatedAuthorList(ctx, req)
	if err != nil {
		response.ErrorNoData(c, "get doc related author list failed")
		return
	}
	response.Success(c, "get doc related author list successfully", resp)
}

// GetDocRelatedClassifyList 获取文档相关的分类列表
func (api *UserDocAPI) GetDocRelatedClassifyList(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.GetDocRelatedClassifyList")
	defer span.Finish()
	req := &pb.GetDocRelatedClassifyListReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "parameter error")
		return
	}
	userId, _ := userContext.GetUserID(c.Request.Context())
	req.UserId = &userId
	resp, err := api.service.GetDocRelatedClassifyList(ctx, req)
	if err != nil {
		response.ErrorNoData(c, "get doc related classify list failed")
		return
	}
	response.Success(c, "get doc related classify list successfully", resp)
}

// GetJcrPartions 获取JCR分区列表
func (api *UserDocAPI) GetJcrPartions(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.GetJcrPartions")
	defer span.Finish()
	req := &pb.JcrPartionsReq{}
	if err := transport.BindProto(c, req); err != nil {
		c.Error(err)
		return
	}
	userId, _ := userContext.GetUserID(c.Request.Context())
	resp, err := api.service.GetAllJcrPartionsByFolderId(ctx, req.GetFolderId(), userId)
	if err != nil {
		response.ErrorNoData(c, "get jcr partions failed")
		return
	}
	response.Success(c, "get jcr partions successfully", resp)
}

// GetUserDocById 获取文献信息
func (api *UserDocAPI) GetUserDocById(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.GetUserDocById")
	defer span.Finish()
	req := &pb.GetUserDocRequest{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "parameter error")
		return
	}

	// 调用服务层方法获取文献详情
	docDetailInfo, err := api.service.GetUserDocDetailInfoById(ctx, req.NoteId)
	if err != nil {
		response.ErrorNoData(c, "get doc detail info failed")
		return
	}

	response.Success(c, "get doc detail info successfully", docDetailInfo)
}

func (api *UserDocAPI) GetUserDocStatusByIds(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.GetUserDocStatusByIds")
	defer span.Finish()
	req := &pb.UserDocStatusByIdsRequest{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "parameter error")
		return
	}
	// 调用服务层方法获取文献状态
	// 将 []string 转换为 []string
	docIds := make([]string, len(req.DocIds))
	for i, id := range req.DocIds {
		docIds[i] = id
	}
	docStatus, err := api.service.GetUserDocStatusByIds(ctx, docIds)
	if err != nil {
		response.ErrorNoData(c, "get doc status failed")
		return
	}
	response.Success(c, "get doc status successfully", docStatus)
}

// @api /api/doc/userDoc/updateReadStatus
// @method POST
// @apiDescription 更新文献阅读状态
func (api *UserDocAPI) UpdateReadStatus(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.UpdateReadStatus")
	defer span.Finish()

	req := &pb.UpdateReadStatusRequest{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "parameter error")
		return
	}
	userId, _ := userContext.GetUserID(c.Request.Context())
	// 调用服务层更新文献阅读状态
	err := api.service.UpdateReadStatus(ctx, userId, req)
	if err != nil {
		response.ErrorNoData(c, "update read status failed")
		return
	}
	response.SuccessNoData(c, "update read status successfully")
}

// DeleteDoc 删除用户文献
func (api *UserDocAPI) DeleteDoc(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.DeleteDoc")
	defer span.Finish()
	// 解析请求参数
	req := &pb.DeleteDocReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "parameter error")
		return
	}
	// 验证请求参数
	if len(req.DocIds) == 0 {
		response.ErrorNoData(c, "doc ids is empty")
		return
	}
	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())
	// 调用服务层删除文献
	err := api.service.DeleteUserDocs(ctx, req.DocIds, userId)
	if err != nil {
		response.ErrorNoData(c, "delete user docs failed")
		return
	}
	response.SuccessNoData(c, "delete user docs successfully")
}

// RenameDoc 重命名用户文献
func (api *UserDocAPI) RenameDoc(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.RenameDoc")
	defer span.Finish()
	// 解析请求参数
	req := &pb.RenameUserDocReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "rename user doc error , bind proto failed")
		return
	}
	// 验证请求参数
	if req.DocId == "0" || req.DocName == "" {
		response.ErrorNoData(c, "rename user doc error , doc id and doc name is empty")
		return
	}
	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())
	//验证是否存在同名文件
	// exists, err := api.service.CheckDocExistsByFileName(ctx, req.DocName, userID)
	// if err != nil {
	// 	response.ErrorNoData(c, "rename user doc error , check doc exists failed")
	// 	return
	// }
	// if exists {
	// 	response.ErrorNoData(c, "A file with the same name already exists")
	// 	return
	// }
	// 调用服务层重命名文献
	err := api.service.RenameUserDoc(ctx, req.DocId, req.DocName, userId)
	if err != nil {
		response.ErrorNoData(c, "rename user doc failed")
		return
	}
	response.SuccessNoData(c, "rename user doc success")
}

// HandleFileFastUpload 处理文件秒传的逻辑
func (api *UserDocAPI) HandleFileFastUpload(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.HandleFileFastUpload")
	defer span.Finish()
	// 解析请求参数
	req := &pb.HandleFileFastUploadReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "file fast upload error , bind proto failed")
		return
	}
	// 只在 needUpload 为 false 时执行逻辑
	if req.NeedUpload {
		response.ErrorNoData(c, "file fast upload error , need upload is false")
		return
	}
	// 检查 OssInfo 是否存在
	if req.OssInfo == nil {
		response.ErrorNoData(c, "file fast upload error , oss info is nil")
		return
	}
	resp := &pb.HandleFileFastUploadResp{}
	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())
	fileName := req.OssInfo.FileName
	// 根据文件名和用户ID查询是否存在重名文件
	// exists, err := api.service.CheckDocExistsByFileName(ctx, fileName, userID)
	// if err != nil {
	// 	response.ErrorNoData(c, "file fast upload error , check file exists failed")
	// 	return
	// }
	// // 如果存在重名文件，返回错误
	// if exists {
	// 	response.ErrorNoData(c, "A file with the same name already exists")
	// 	return
	// }
	// 根据FileSHA256值查询PaperPdf记录
	fileSHA256 := req.OssInfo.FileSHA256
	if fileSHA256 == "" {
		response.ErrorNoData(c, "file fast upload error , file sha256 is empty")
		return
	}
	// 调用PaperPdfService的GetByFileSHA256方法查询PaperPdf记录
	paperPdf, err := api.paperPdfService.GetByFileSHA256(ctx, fileSHA256)
	if err != nil {
		response.ErrorNoData(c, "file fast upload error , paper pdf not found")
		return
	}
	userDoc, err := api.service.GetByUserIdAndPdfId(ctx, paperPdf.CreatorId, paperPdf.Id)
	if err != nil {
		response.ErrorNoData(c, "file upload error , document not found ")
		return
	}
	paper, err := api.paperService.GetPaperById(ctx, paperPdf.PaperId)
	if err != nil {
		response.ErrorNoData(c, "file fast upload error , paper not found")
		return
	}
	resp, err = api.service.HandleFileFastUpload(ctx, userId, fileName, userDoc, paperPdf, paper)
	if err != nil {
		response.ErrorNoData(c, "file fast upload error")
		return
	}
	// 返回成功响应，包含PDF ID
	response.Success(c, "file fast upload success", resp)

}

// ManualUpdateDocCiteInfo 手动更新文档引用信息
func (api *UserDocAPI) ManualUpdateDocCiteInfo(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.ManualUpdateDocCiteInfo")
	defer span.Finish()

	// 解析请求参数
	req := &pb.ManualUpdateDocCiteInfoReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "parameter error")
		return
	}
	//验证请求参数
	if req.GetDocName() == "" {
		response.ErrorNoData(c, "document title is empty")
		return
	}
	// 验证请求参数
	if req.PaperId == nil || req.PdfId == nil {
		response.ErrorNoData(c, "paper id and pdf id is empty")
		return
	}
	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())
	// 调用服务层更新文档引用信息
	err := api.service.ManualUpdateDocCiteInfo(ctx, userId, req)
	if err != nil {
		response.ErrorNoData(c, "manual update doc cite info failed")
		return
	}
	response.SuccessNoData(c, "manual update doc cite info success")
}

// EnDocCiteSearch 英文文档引用搜索
func (api *UserDocAPI) EnDocCiteSearch(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.EnDocCiteSearch")
	defer span.Finish()

	// 解析请求参数
	req := &pb.EnDocCiteSearchReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "parameter error")
		return
	}

	// 验证请求参数
	if req.GetSearchContent() == "" {
		response.ErrorNoData(c, "search content is empty")
		return
	}

	// 调用服务层进行英文文档引用搜索
	resp, err := api.docCiteSearchService.EnDocCiteSearch(ctx, req)
	if err != nil {
		response.ErrorNoData(c, "english document citation search failed")
		return
	}

	response.Success(c, "english document citation search success", resp)
}

// ZhDocCiteSearch 中文文档引用搜索
func (api *UserDocAPI) ZhDocCiteSearch(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.ZhDocCiteSearch")
	defer span.Finish()

	// 解析请求参数
	req := &pb.ZhDocCiteSearchReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "parameter error")
		return
	}
	// 验证请求参数
	if req.GetSearchContent() == "" {
		response.ErrorNoData(c, "search content is empty")
		return
	}
	// 调用服务层进行中文文档引用搜索
	resp, err := api.docCiteSearchService.ZhDocCiteSearch(ctx, req)
	if err != nil {
		response.ErrorNoData(c, "chinese document citation search failed")
		return
	}
	response.Success(c, "chinese document citation search success", resp)
}

func (api *UserDocAPI) UpdateUserDocAuthors(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.UpdateUserDocAuthors")
	defer span.Finish()

	// 解析请求参数
	req := &pb.UpdateAuthorsRequest{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "parameter error")
		return
	}
	// 验证请求参数
	if req.GetDocId() == "0" {
		response.ErrorNoData(c, "doc id is empty")
		return
	}
	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())
	// 调用服务层更新作者
	baseResp, err := api.service.UpdateUserDocByCustomType(ctx, userId, req.GetDocId(), &docBean.UserDocUpdateRequest{
		UpdateType: docBean.UpdateTypeAuthors,
		Authors:    req.GetAuthors(),
	})
	if err != nil {
		response.ErrorNoData(c, "update user doc authors failed")
		return
	}
	resp := baseResp.AuthorsResponse
	response.Success(c, "update user doc authors success", resp)
}

func (api *UserDocAPI) GetAuthors(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.GetAuthors")
	defer span.Finish()

	// 解析请求参数
	req := &pb.GetAuthorsRequest{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "parameter error")
		return
	}
	// 验证请求参数
	if req.GetDocId() == "0" {
		response.ErrorNoData(c, "doc id is empty")
		return
	}
	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())
	if userId == "" {
		response.ErrorNoData(c, "user not login")
		return
	}
	// 调用服务层获取作者
	resp, err := api.service.GetAuthors(ctx, req.GetDocId())
	if err != nil {
		response.ErrorNoData(c, "get authors failed")
		return
	}
	response.Success(c, "get authors success", resp)
}

func (api *UserDocAPI) UpdateUserDocVenue(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.UpdateUserDocVenue")
	defer span.Finish()

	// 解析请求参数
	req := &pb.UpdateVenueRequest{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "parameter error")
		return
	}
	// 验证请求参数
	if req.GetDocId() == "0" {
		response.ErrorNoData(c, "doc id is empty")
		return
	}
	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())
	if userId == "" {
		response.ErrorNoData(c, "user not login")
		return
	}
	// 调用服务层更新收录情况
	baseResp, err := api.service.UpdateUserDocByCustomType(ctx, userId, req.GetDocId(), &docBean.UserDocUpdateRequest{
		UpdateType: docBean.UpdateTypeVenue,
		Venue:      req.GetVenue(),
	})
	if err != nil {
		response.ErrorNoData(c, "update user doc venue failed")
		return
	}
	resp := baseResp.VenueResponse
	response.Success(c, "update user doc venue success", resp)
}

func (api *UserDocAPI) UpdateUserDocPublishDate(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.UpdateUserDocPublishDate")
	defer span.Finish()

	// 解析请求参数
	req := &pb.UpdatePublishDateRequest{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "parameter error")
		return
	}
	// 验证请求参数
	if req.GetDocId() == "0" {
		response.ErrorNoData(c, "doc id is empty")
		return
	}
	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())
	if userId == "" {
		response.ErrorNoData(c, "user not login")
		return
	}
	// 调用服务层更新发布时间
	baseResp, err := api.service.UpdateUserDocByCustomType(ctx, userId, req.GetDocId(), &docBean.UserDocUpdateRequest{
		UpdateType:  docBean.UpdateTypePublishDate,
		PublishDate: req.GetPublishDate(),
	})
	if err != nil {
		response.ErrorNoData(c, "update user doc publish date failed")
		return
	}
	resp := baseResp.PublishDateResponse
	response.Success(c, "update user doc publish date success", resp)
}

func (api *UserDocAPI) UpdateUserDocJcrPartion(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.UpdateUserDocJcrPartion")
	defer span.Finish()

	// 解析请求参数
	req := &pb.JcrPartionUpdateReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "parameter error")
		return
	}
	// 验证请求参数
	if req.GetDocId() == "0" {
		response.ErrorNoData(c, "doc id is empty")
		return
	}
	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())
	if userId == "" {
		response.ErrorNoData(c, "user not login")
		return
	}
	// 调用服务层更新jcr分区
	baseResp, err := api.service.UpdateUserDocByCustomType(ctx, userId, req.GetDocId(), &docBean.UserDocUpdateRequest{
		UpdateType: docBean.UpdateTypeJcrPartion,
		JcrPartion: req.GetJcrPartion(),
	})
	if err != nil {
		response.ErrorNoData(c, "update user doc jcr partion failed")
		return
	}
	resp := baseResp.JcrPartionUpdateResponse
	response.Success(c, "update user doc jcr partion success", resp)
}

func (api *UserDocAPI) UpdateUserDocRemark(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.UpdateUserDocRemark")
	defer span.Finish()

	// 解析请求参数
	req := &pb.UpdateDocRemarkReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "parameter error")
		return
	}
	// 验证请求参数
	if req.GetDocId() == "0" {
		response.ErrorNoData(c, "doc id is empty")
		return
	}
	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())
	if userId == "" {
		response.ErrorNoData(c, "user not login")
		return
	}
	// 调用服务层更新备注
	_, err := api.service.UpdateUserDocByCustomType(ctx, userId, req.GetDocId(), &docBean.UserDocUpdateRequest{
		UpdateType: docBean.UpdateTypeRemark,
		Remark:     req.GetRemark(),
	})
	if err != nil {
		response.ErrorNoData(c, "update user doc remark failed")
		return
	}
	response.SuccessNoData(c, "update user doc remark success")
}

func (api *UserDocAPI) UpdateUserDocImpactFactor(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.UpdateUserDocImpactFactor")
	defer span.Finish()

	// 解析请求参数
	req := &pb.ImpactFactorReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "parameter error")
		return
	}
	// 验证请求参数
	if req.GetDocId() == "0" {
		response.ErrorNoData(c, "doc id is empty")
		return
	}
	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())
	if userId == "" {
		response.ErrorNoData(c, "user not login")
		return
	}
	// 调用服务层更新影响因子
	baseResp, err := api.service.UpdateUserDocByCustomType(ctx, userId, req.GetDocId(), &docBean.UserDocUpdateRequest{
		UpdateType:   docBean.UpdateTypeImpactFactor,
		ImpactFactor: req.ImpactOfFactor,
	})
	if err != nil {
		response.ErrorNoData(c, "update user doc impact factor failed")
		return
	}
	resp := baseResp.ImpactFactorResponse
	response.Success(c, "update user doc impact factor success", resp)
}

func (api *UserDocAPI) UpdateUserDocImportanceScore(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.UpdateUserDocImportanceScore")
	defer span.Finish()

	// 解析请求参数
	req := &pb.ImportanceScoreReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "parameter error")
		return
	}
	// 验证请求参数
	if req.GetDocId() == "0" {
		response.ErrorNoData(c, "doc id is empty")
		return
	}
	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())
	if userId == "" {
		response.ErrorNoData(c, "user not login")
		return
	}
	// 调用服务层更新重要性评分
	baseResp, err := api.service.UpdateUserDocByCustomType(ctx, userId, req.GetDocId(), &docBean.UserDocUpdateRequest{
		UpdateType:      docBean.UpdateTypeImportanceScore,
		ImportanceScore: req.GetScore(),
	})
	if err != nil {
		response.ErrorNoData(c, "update user doc importance score failed")
		return
	}
	resp := baseResp.ImpactFactorResponse
	response.Success(c, "update user doc importance score success", resp)
}

func (api *UserDocAPI) AttachDocToClassify(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.AttachDocToClassify")
	defer span.Finish()

	// 解析请求参数
	req := &pb.AttachDocToClassifyReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "parameter error")
		return
	}
	// 验证请求参数
	if req.GetDocId() == "0" {
		response.ErrorNoData(c, "doc id or classify id is empty")
		return
	}
	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())
	if userId == "" {
		response.ErrorNoData(c, "user not login")
		return
	}
	// 调用服务层添加标签
	err := api.service.AttachDocToClassify(ctx, userId, req.GetDocId(), req.GetClassifyId(), req.GetClassifyName())
	if err != nil {
		response.ErrorNoData(c, "attach doc to classify failed")
		return
	}
	response.SuccessNoData(c, "attach doc to classify success")
}

func (api *UserDocAPI) RemoveDocFromClassify(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.RemoveDocFromClassify")
	defer span.Finish()

	// 解析请求参数
	req := &pb.DeleteDocClassifyReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "parameter error")
		return
	}
	// 验证请求参数
	if req.GetDocId() == "0" || req.GetClassifyId() == "0" {
		response.ErrorNoData(c, "doc id or classify id is empty")
		return
	}
	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())
	if userId == "" {
		response.ErrorNoData(c, "user not login")
		return
	}
	// 调用服务层删除标签
	err := api.service.RemoveDocFromClassify(ctx, userId, req.GetDocId(), req.GetClassifyId())
	if err != nil {
		response.ErrorNoData(c, "remove doc from classify failed")
		return
	}
	response.SuccessNoData(c, "remove doc from classify success")
}

func (api *UserDocAPI) GetLatestReadDocList(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocAPI.GetLatestReadDocList")
	defer span.Finish()

	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())
	if userId == "" {
		response.ErrorNoData(c, "user not login")
		return
	}
	//
	defaultLatestReadSize := api.config.Personal.LatestReadSize
	// 调用服务层获取最近阅读文献
	resp, err := api.service.GetLatestReadDocList(ctx, userId, defaultLatestReadSize)
	if err != nil {
		response.ErrorNoData(c, "get latest read doc list failed")
		return
	}
	response.Success(c, "get latest read doc list successfully", resp)
}
