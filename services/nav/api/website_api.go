package api

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	pb "github.com/yb2020/odoc/proto/gen/go/nav"
	"github.com/yb2020/odoc/services/nav/service"
)

// WebsiteAPI 网站API实现
type WebsiteAPI struct {
	websiteService *service.WebsiteService
	logger         logging.Logger
	tracer         opentracing.Tracer
}

// NewWebsiteAPI 创建新的网站API
func NewWebsiteAPI(
	logger logging.Logger,
	tracer opentracing.Tracer,
	websiteService *service.WebsiteService,
) *WebsiteAPI {
	return &WebsiteAPI{
		websiteService: websiteService,
		logger:         logger,
		tracer:         tracer,
	}
}

// @api_path: /api/nav/website/create
// @method: POST
// @summary: 创建用户学术网站
func (api *WebsiteAPI) CreateWebsite(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "WebsiteAPI.CreateWebsiteRequest")
	defer span.Finish()

	req := pb.CreateWebsiteRequest{}
	if err := transport.BindProto(c, &req); err != nil {
		api.logger.Error("CreateWebsite failed parse param", "error", err)
		response.ErrorNoData(c, err.Error())
		return
	}

	req.OpenType = pb.WebsiteOpenType_WebsiteOpenType_NewTab

	userId, _ := userContext.GetUserID(ctx)
	id, err := api.websiteService.CreateWebsite(ctx, userId, int32(pb.WebsiteSource_WebsiteSource_User), &req)
	if err != nil {
		api.logger.Error("CreateWebsite failed", "error", err)
		response.ErrorNoData(c, "CreateWebsite failed")
		return
	}

	rsp := &pb.CreateWebsiteResponse{
		Id: id,
	}

	response.Success(c, "Success", rsp)
}

// @api_path: /api/nav/website/delete
// @method: POST
// @summary: 删除用户学术网站
func (api *WebsiteAPI) DeleteWebsite(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "WebsiteAPI.DeleteWebsiteRequest")
	defer span.Finish()

	req := pb.DeleteWebsiteRequest{}
	if err := transport.BindProto(c, &req); err != nil {
		api.logger.Error("DeleteWebsite failed", "error", err)
		response.ErrorNoData(c, "parameter error")
		return
	}

	if req.Id == "0" {
		api.logger.Error("DeleteWebsite failed", "error", "id is empty")
		response.ErrorNoData(c, "id is empty")
		return
	}

	err := api.websiteService.DeleteWebsite(ctx, req.Id)
	if err != nil {
		api.logger.Error("DeleteWebsite failed", "error", err)
		response.ErrorNoData(c, err.Error())
		return
	}

	response.SuccessNoData(c, "Success")
}

// @api_path: /api/nav/website/getById
// @method: GET
// @summary: 获取用户学术网站
func (api *WebsiteAPI) GetWebsiteById(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "WebsiteAPI.GetWebsiteByIdRequest")
	defer span.Finish()

	req := pb.GetWebsiteByIdRequest{}
	if err := transport.BindProto(c, &req); err != nil {
		api.logger.Error("parameter error", "error", err)
		response.ErrorNoData(c, "parameter error")
		return
	}

	website, err := api.websiteService.GetWebsitePbById(ctx, req.Id)
	if err != nil {
		api.logger.Error("GetWebsite failed", "error", err)
		response.ErrorNoData(c, err.Error())
		return
	}
	rep := &pb.GetWebsiteByIdResponse{
		Website: website,
	}

	response.Success(c, "Success", rep)
}

// @api_path: /api/nav/website/update
// @method: POST
// @summary: 更新用户学术网站
func (api *WebsiteAPI) UpdateWebsite(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "WebsiteAPI.UpdateWebsiteRequest")
	defer span.Finish()

	req := pb.UpdateWebsiteRequest{}
	if err := transport.BindProto(c, &req); err != nil {
		api.logger.Error("UpdateWebsite parameter error", "error", err)
		response.ErrorNoData(c, err.Error())
		return
	}

	if req.Id == "0" {
		api.logger.Error("UpdateWebsite parameter error", "error", "id is empty")
		response.ErrorNoData(c, "id is empty")
		return
	}

	id, err := api.websiteService.UpdateWebsite(ctx, &req)
	if err != nil {
		api.logger.Error("UpdateWebsite failed", "error", err)
		response.ErrorNoData(c, err.Error())
		return
	}

	rep := &pb.UpdateWebsiteResponse{
		Id: id,
	}

	response.Success(c, "Success", rep)
}

// @api_path: /api/nav/website/getList
// @method: GET
// @summary: 获取用户学术网站列表
func (api *WebsiteAPI) GetWebsiteList(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "WebsiteAPI.GetWebsiteListRequest")
	defer span.Finish()

	req := pb.GetWebsiteListRequest{}
	if err := transport.BindProto(c, &req); err != nil {
		api.logger.Error("parameter error", "error", err)
		response.ErrorNoData(c, "parameter error")
		return
	}

	userId, _ := userContext.GetUserID(ctx)

	// 检查用户是否初始化了系统网站
	isInit, err := api.websiteService.CheckUserInitSystemWebsite(ctx, userId)
	if err != nil {
		api.logger.Error("CheckUserInitSystemWebsite failed", "error", err)
		response.ErrorNoData(c, err.Error())
		return
	}
	if !isInit {
		// 如果用户没有初始化系统网站，初始化系统网站
		err := api.websiteService.InitUserSystemWebsite(ctx, userId)
		if err != nil {
			api.logger.Error("InitUserSystemWebsite failed", "error", err)
			response.ErrorNoData(c, err.Error())
			return
		}
	}

	websiteList, err := api.websiteService.GetUserWebsiteListBySortOrder(ctx, userId)
	if err != nil {
		api.logger.Error("GetUserWebsiteListBySortOrder failed", "error", err)
		response.ErrorNoData(c, err.Error())
		return
	}

	websiteListPb := make([]*pb.Website, len(websiteList))
	for i, website := range websiteList {
		websiteListPb[i] = &pb.Website{
			Id:        website.Id,
			UserId:    website.UserId,
			Source:    pb.WebsiteSource(website.Source),
			IconUrl:   website.IconUrl,
			Name:      website.Name,
			Url:       website.Url,
			OpenType:  pb.WebsiteOpenType(website.OpenType),
			SortOrder: website.SortOrder,
		}
	}
	rep := &pb.GetWebsiteListResponse{
		Websites: websiteListPb,
	}

	response.Success(c, "Success", rep)
}

// @api_path: /api/nav/website/reorder
// @method: POST
// @summary: 重新排序用户学术网站
func (api *WebsiteAPI) ReorderWebsites(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "WebsiteAPI.ReorderWebsitesRequest")
	defer span.Finish()

	req := pb.ReorderWebsitesRequest{}
	if err := transport.BindProto(c, &req); err != nil {
		api.logger.Error("parameter error", "error", err)
		response.ErrorNoData(c, "parameter error")
		return
	}
	if req.Id == "0" {
		api.logger.Error("parameter error", "error", "id is empty")
		response.ErrorNoData(c, "id is empty")
		return
	}

	userId, _ := userContext.GetUserID(ctx)

	rebalanced, updateWebsites, err := api.websiteService.ReorderWebsites(ctx, userId, &req)
	if err != nil {
		api.logger.Error("ReorderWebsites failed", "error", err)
		response.ErrorNoData(c, err.Error())
		return
	}

	rsp := &pb.ReorderWebsitesResponse{
		Rebalanced: rebalanced,
		Updates:    updateWebsites,
	}

	response.Success(c, "Success", rsp)
}
