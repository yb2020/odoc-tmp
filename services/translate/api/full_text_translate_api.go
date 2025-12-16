package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/config"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	pb "github.com/yb2020/odoc/proto/gen/go/membership"
	"github.com/yb2020/odoc/proto/gen/go/translate"
	"github.com/yb2020/odoc/services/membership/interfaces"
	"github.com/yb2020/odoc/services/translate/service"
)

// FullTextTranslateAPI 全文翻译API处理器
type FullTextTranslateAPI struct {
	fullTextService   *service.FullTextTranslateService
	membershipService interfaces.IMembershipService
	config            *config.Config
	logger            logging.Logger
	tracer            opentracing.Tracer
}

// NewFullTextTranslateAPI 创建全文翻译API处理器
func NewFullTextTranslateAPI(
	config *config.Config,
	logger logging.Logger,
	tracer opentracing.Tracer,
	fullTextService *service.FullTextTranslateService,
	membershipService interfaces.IMembershipService,
) *FullTextTranslateAPI {
	return &FullTextTranslateAPI{
		fullTextService:   fullTextService,
		membershipService: membershipService,
		config:            config,
		logger:            logger,
		tracer:            tracer,
	}
}

// GetHistoryList 获取用户的全文翻译历史
func (api *FullTextTranslateAPI) GetHistoryList(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "FullTextTranslateAPI.GetHistoryList")
	defer span.Finish()

	userId, _ := userContext.GetUserID(c.Request.Context())
	resp, err := api.fullTextService.GetHistoryList(ctx, userId)
	if err != nil {
		api.logger.Error("获取全文翻译历史失败", "error", err.Error())
		c.Error(err)
		return
	}
	response.Success(c, "success", resp)
}

// GetRightInfo 获取用户全文翻译权限信息
func (api *FullTextTranslateAPI) GetRightInfo(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "FullTextTranslateAPI.GetRightInfo")
	defer span.Finish()

	userId, _ := userContext.GetUserID(c.Request.Context())
	resp, err := api.fullTextService.GetRightInfo(ctx, userId)
	if err != nil {
		api.logger.Error("获取全文翻译权限失败", "error", err.Error())
		c.Error(err)
		return
	}
	response.Success(c, "success", resp)
}

// Translate 发起全文翻译
func (api *FullTextTranslateAPI) Translate(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "FullTextTranslateAPI.Translate")
	defer span.Finish()

	// 解析请求体
	req := &translate.FullTextTranslateRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("解析全文翻译请求失败", "error", err.Error())
		c.Error(err)
		return
	}
	// 参数校验（可按需补充）
	if req.NeedTranslateFileUrl == "" {
		c.Error(errors.System(errors.ErrorTypeInvalidArgument, "needTranslateFileUrl不能为空", nil))
		return
	}

	var finalResp *translate.FullTextTranslateResponse
	// 使用CreditFunTranslate包装全文翻译逻辑 //TODO filePageSize设置为真实的页数   这里不需要自动提交，拿到sessionId后，进行二次确认
	err := api.membershipService.CreditFunTranslate(ctx, pb.CreditServiceType_CREDIT_SERVICE_TYPE_TRANSLATE_FULLTEXT, 0, func(xctx context.Context, sessionId string) error {
		// 调用服务
		resp, err := api.fullTextService.Translate(xctx, req, sessionId)
		if err != nil {
			api.logger.Error("全文翻译失败", "error", err.Error())
			return err
		}
		finalResp = resp
		return nil
	}, false)
	if err != nil {
		c.Error(err)
		return
	}

	// 在积分扣除成功后再返回响应
	if finalResp != nil {
		response.Success(c, "success", finalResp)
	}
}

// GetTranslateStatus 查询全文翻译状态
func (api *FullTextTranslateAPI) GetTranslateStatus(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "FullTextTranslateAPI.GetTranslateStatus")
	defer span.Finish()

	req := &translate.GetTranslateStatusReq{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("解析翻译状态请求失败", "error", err.Error())
		c.Error(err)
		return
	}
	resp, err := api.fullTextService.GetTranslateStatus(ctx, req)
	if err != nil {
		api.logger.Error("查询翻译状态失败", "error", err.Error())
		c.Error(err)
		return
	}
	response.Success(c, "success", resp)
}

// // ReTranslate 管理员重试全文翻译
// func (api *FullTextTranslateAPI) ReTranslate(c *gin.Context) {
// 	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "FullTextTranslateAPI.ReTranslate")
// 	defer span.Finish()

// 	req := &translate.FullTextReTranslateReq{}
// 	if err := transport.BindProto(c, req); err != nil {
// 		api.logger.Warn("解析重试翻译请求失败", "error", err.Error())
// 		c.Error(err)
// 		return
// 	}
// 	// 权限校验（如有需要）
// 	// ...
// 	resp, err := api.fullTextService.ReTranslate(ctx, req)
// 	if err != nil {
// 		api.logger.Error("重试全文翻译失败", "error", err.Error())
// 		c.Error(err)
// 		return
// 	}
// 	response.Success(c, "success", resp)
// }

// // GetReTranslateResult 管理员获取重试翻译结果
// func (api *FullTextTranslateAPI) GetReTranslateResult(c *gin.Context) {
// 	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "FullTextTranslateAPI.GetReTranslateResult")
// 	defer span.Finish()

// 	req := &translate.FullTextReTranslateResultReq{}
// 	if err := transport.BindProto(c, req); err != nil {
// 		api.logger.Warn("解析获取重试结果请求失败", "error", err.Error())
// 		c.Error(err)
// 		return
// 	}
// 	resp, err := api.fullTextService.GetReTranslateResult(ctx, req)
// 	if err != nil {
// 		api.logger.Error("获取重试翻译结果失败", "error", err.Error())
// 		c.Error(err)
// 		return
// 	}
// 	response.Success(c, "success", resp)
// }
