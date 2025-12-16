package api

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	"github.com/yb2020/odoc/proto/gen/go/translate"
	"github.com/yb2020/odoc/services/translate/service"
)

// GlossaryAPI 术语库API处理器
type GlossaryAPI struct {
	glossaryService service.GlossaryService
	logger          logging.Logger
	tracer          opentracing.Tracer
}

// NewGlossaryAPI 创建术语库API处理器
func NewGlossaryAPI(glossaryService *service.GlossaryService, logger logging.Logger, tracer opentracing.Tracer) *GlossaryAPI {
	return &GlossaryAPI{
		glossaryService: *glossaryService,
		logger:          logger,
		tracer:          tracer,
	}
}

// AddGlossary 添加术语条目
func (api *GlossaryAPI) AddGlossary(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "GlossaryAPI.AddGlossary")
	defer span.Finish()

	// 直接获取用户ID，因为中间件已经确保了用户已认证
	userId, _ := userContext.GetUserID(c.Request.Context())

	// 使用 Proto 绑定器解析请求体
	req := &translate.AddGlossaryReq{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("msg", "解析添加术语条目请求失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 调用服务添加术语条目
	id, err := api.glossaryService.AddGlossary(ctx, userId, req)
	if err != nil {
		api.logger.Warn("msg", "添加术语条目失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 创建响应
	resp := &translate.AddGlossaryResponse{
		Id: id,
	}

	// 返回响应
	response.Success(c, "success", resp)
}

// UpdateGlossary 更新术语条目
func (api *GlossaryAPI) UpdateGlossary(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "GlossaryAPI.UpdateGlossary")
	defer span.Finish()

	// 直接获取用户ID，因为中间件已经确保了用户已认证
	userId, _ := userContext.GetUserID(c.Request.Context())

	// 使用 Proto 绑定器解析请求体
	req := &translate.UpdateGlossaryReq{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("msg", "解析更新术语条目请求失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 调用服务更新术语条目
	err := api.glossaryService.UpdateGlossary(ctx, userId, req)
	if err != nil {
		api.logger.Warn("msg", "更新术语条目失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 返回成功响应
	response.SuccessNoData(c, "success")
}

// DeleteGlossary 删除术语条目
func (api *GlossaryAPI) DeleteGlossary(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "GlossaryAPI.DeleteGlossary")
	defer span.Finish()

	// 直接获取用户ID，因为中间件已经确保了用户已认证
	userId, _ := userContext.GetUserID(c.Request.Context())

	// 使用 Proto 绑定器解析请求体
	req := &translate.DeleteGlossaryReq{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("msg", "解析删除术语条目请求失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 遍历删除每个术语条目
	for _, id := range req.Ids {
		err := api.glossaryService.DeleteGlossary(ctx, userId, id)
		if err != nil {
			api.logger.Warn("msg", "删除术语条目失败", "id", id, "error", err.Error())
			c.Error(err)
			return
		}
	}

	// 返回成功响应
	response.SuccessNoData(c, "success")
}

// GetGlossaryList 获取术语条目列表
func (api *GlossaryAPI) GetGlossaryList(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "GlossaryAPI.GetGlossaryList")
	defer span.Finish()

	// 直接获取用户ID，因为中间件已经确保了用户已认证
	userId, _ := userContext.GetUserID(c.Request.Context())

	// 使用 Proto 绑定器解析请求体
	req := &translate.GetGlossaryListReq{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("msg", "解析获取术语条目列表请求失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 将请求参数转换为服务层需要的格式
	// 注意：proto 中的字段名与服务层实现不一致，需要转换
	if req.CurrentPage <= 0 {
		req.CurrentPage = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	// 调用服务获取术语条目列表
	resp, err := api.glossaryService.GetGlossaryList(ctx, userId, req)
	if err != nil {
		api.logger.Warn("msg", "获取术语条目列表失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 返回响应
	response.Success(c, "success", resp)
}
