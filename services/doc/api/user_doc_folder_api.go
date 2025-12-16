package api

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	docpb "github.com/yb2020/odoc/proto/gen/go/doc"
	"github.com/yb2020/odoc/services/doc/service"
)

// UserDocFolderAPI 用户文档文件夹API处理器
type UserDocFolderAPI struct {
	service                      *service.UserDocFolderService
	userDocFolderRelationService *service.UserDocFolderRelationService
	logger                       logging.Logger
	tracer                       opentracing.Tracer
}

// NewUserDocFolderAPI 创建新的用户文档文件夹API处理器
func NewUserDocFolderAPI(
	service *service.UserDocFolderService,
	userDocFolderRelationService *service.UserDocFolderRelationService,
	logger logging.Logger,
	tracer opentracing.Tracer,
) *UserDocFolderAPI {
	return &UserDocFolderAPI{
		service:                      service,
		userDocFolderRelationService: userDocFolderRelationService,
		logger:                       logger,
		tracer:                       tracer,
	}
}

// CreateUserDocFolder 创建用户文档文件夹
func (api *UserDocFolderAPI) CreateUserDocFolder(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocFolderAPI.CreateUserDocFolder")
	defer span.Finish()
	req := &docpb.CreateUserDocFolderRequest{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "参数错误")
		return
	}
	userId, _ := userContext.GetUserID(c.Request.Context())
	resp, err := api.service.CreateUserDocFolder(ctx, userId, req)
	if err != nil {
		response.ErrorNoData(c, "创建文件夹失败")
		return
	}
	response.Success(c, "创建用户文档文件夹成功", resp)
}

// DeleteUserDocFolder 删除用户文档文件夹
func (api *UserDocFolderAPI) DeleteUserDocFolder(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocFolderAPI.DeleteUserDocFolder")
	defer span.Finish()

	// 解析请求参数
	req := &docpb.DeleteUserDocFolderReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "参数错误")
		return
	}

	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())
	// 调用服务层删除文件夹 TODO:jinzhi 这个方法需要进行优化，下个版本解决
	err := api.service.DeleteUserDocFolder(ctx, req, userId)
	if err != nil {
		response.ErrorNoData(c, "删除文件夹失败")
		return
	}

	response.SuccessNoData(c, "删除文件夹成功")
}

// UpdateUserDocFolder 更新用户文档文件夹
func (api *UserDocFolderAPI) UpdateUserDocFolder(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocFolderAPI.UpdateUserDocFolder")
	defer span.Finish()
	// 解析请求参数
	req := &docpb.UpdateUserDocFolderReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "参数错误")
		return
	}
	// 验证必填参数
	if req.FolderId == "0" {
		response.ErrorNoData(c, "文件夹ID不能为空")
		return
	}
	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())
	// 调用服务层更新文件夹
	err := api.service.UpdateUserDocFolder(ctx, req, userId)
	if err != nil {
		response.ErrorNoData(c, "更新文件夹失败")
		return
	}
	response.SuccessNoData(c, "更新成功")
}

// 移动文件夹或者文献
func (api *UserDocFolderAPI) MoveDocOrFolderToAnotherFolder(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocFolderAPI.MoveDocOrFolderToAnotherFolder")
	defer span.Finish()
	// 解析请求参数
	req := &docpb.MoveDocOrFolderToAnotherFolderReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "参数错误")
		return
	}
	// 验证请求参数
	if len(req.MovedFolderIds) == 0 && len(req.MovedDocItems) == 0 {
		response.ErrorNoData(c, "移动的文件夹或文献不能为空")
		return
	}
	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())

	// 调用服务层移动文件夹或文献
	err := api.service.MoveDocOrFolderToAnotherFolder(ctx, req, userId)
	if err != nil {
		response.ErrorNoData(c, "移动失败")
		return
	}

	response.SuccessNoData(c, "操作成功")
}

// 拖拽移动文件夹或者文献
func (api *UserDocFolderAPI) MoveFolderOrDoc(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocFolderAPI.MoveFolderOrDoc")
	defer span.Finish()
	// 解析请求参数
	req := &docpb.MoveFolderOrDocReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "参数错误")
		return
	}
	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())
	// 验证请求参数
	if len(req.TargetFolderItems) == 0 {
		response.ErrorNoData(c, "目标文件夹项不能为空")
		return
	}
	// 根据类型调用不同的服务方法
	var err error
	if req.Type == 0 { // 移动文档
		// 调用文档关系服务移动文档
		err = api.userDocFolderRelationService.MoveDoc(ctx, req, userId)
	} else { // 移动文件夹
		// 调用文件夹服务移动文件夹
		err = api.service.MoveFolder(ctx, req, userId)
	}
	if err != nil {
		response.ErrorNoData(c, "移动失败")
		return
	}
	response.SuccessNoData(c, "操作成功")
}

func (api *UserDocFolderAPI) RemoveDocFromFolder(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocFolderAPI.RemoveDocFromFolder")
	defer span.Finish()

	// 解析请求参数
	req := &docpb.RemoveDocFromFolderReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "parameter error")
		return
	}
	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())
	if userId == "" {
		response.ErrorNoData(c, "user not login")
		return
	}
	// 调用服务层从文件夹中删除文献或文件夹
	err := api.service.RemoveDocFromFolder(ctx, userId, req)
	if err != nil {
		response.ErrorNoData(c, "remove doc from folder failed")
		return
	}
	response.SuccessNoData(c, "remove doc from folder success")
}

func (api *UserDocFolderAPI) CopyDocOrFolderToAnotherFolder(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserDocFolderAPI.CopyDocOrFolderToAnotherFolder")
	defer span.Finish()

	// 解析请求参数
	req := &docpb.CopyDocOrFolderToAnotherFolderReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "parameter error")
		return
	}
	// 验证请求参数
	if len(req.GetDocIds()) == 0 && len(req.GetFolderIds()) == 0 {
		response.ErrorNoData(c, " please select doc or folder")
		return
	}
	if len(req.GetDocIds()) > 0 && len(req.GetFolderIds()) > 0 {
		response.ErrorNoData(c, "can not copy doc and folder at the same time")
		return
	}
	if req.GetTargetFolderId() == "0" {
		response.ErrorNoData(c, "target folder id is empty")
		return
	}
	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())
	if userId == "" {
		response.ErrorNoData(c, "user not login")
		return
	}
	// 调用服务层复制文献或文件夹
	err := api.service.CopyDocOrFolderToAnotherFolder(ctx, userId, req)
	if err != nil {
		response.ErrorNoData(c, "copy doc or folder to another folder failed")
		return
	}
	response.SuccessNoData(c, "copy doc or folder to another folder success")
}
