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

type DocClassifyAPI struct {
	service            *service.UserDocClassifyService
	clsRelationService *service.DocClassifyRelationService
	logger             logging.Logger
	tracer             opentracing.Tracer
}

func NewDocClassifyAPI(
	logger logging.Logger,
	tracer opentracing.Tracer,
	userDocClassifyService *service.UserDocClassifyService,
	clsRelationService *service.DocClassifyRelationService,
) *DocClassifyAPI {
	return &DocClassifyAPI{
		logger:             logger,
		tracer:             tracer,
		service:            userDocClassifyService,
		clsRelationService: clsRelationService,
	}
}

// GetUserAllClassifyList 获取用户所有的文档分类列表
func (api *DocClassifyAPI) GetUserAllClassifyList(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "DocClassifyAPI.GetUserAllClassifyList")
	defer span.Finish()
	userId, _ := userContext.GetUserID(c.Request.Context())
	resp, err := api.service.GetUserAllClassifyList(ctx, userId)
	if err != nil {
		response.ErrorNoData(c, "获取用户文档分类列表失败")
		return
	}
	response.Success(c, "获取用户所有文档分类列表成功", resp)
}

func (api *DocClassifyAPI) AddUserDocClassify(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "DocClassifyAPI.AddUserDocClassify")
	defer span.Finish()

	// 解析请求参数
	req := &pb.AddClassifyReq{}
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
	// 调用服务层添加标签
	classify, err := api.service.AddUserDocClassify(ctx, userId, req.GetClassifyName(), req.GetRemark())
	if err != nil {
		response.ErrorNoData(c, "add user doc classify failed")
		return
	}
	resp := &pb.AddClassifyResponse{
		ClassifyId:   classify.Id,
		ClassifyName: classify.Name,
	}
	response.Success(c, "add user doc classify success", resp)
}

func (api *DocClassifyAPI) DeleteUserDocClassify(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "DocClassifyAPI.DeleteUserDocClassify")
	defer span.Finish()

	// 解析请求参数
	req := &pb.DeleteClassifyReq{}
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
	// 调用服务层删除标签
	err := api.service.DeleteUserDocClassify(ctx, userId, req.GetClassifyId())
	if err != nil {
		response.ErrorNoData(c, "delete user doc classify failed")
		return
	}
	response.SuccessNoData(c, "delete user doc classify success")
}
