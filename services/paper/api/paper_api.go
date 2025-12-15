package api

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	pb "github.com/yb2020/odoc-proto/gen/go/paper"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	"github.com/yb2020/odoc/services/paper/service"
)

// PaperAPI 论文API处理器
type PaperAPI struct {
	paperService *service.PaperService
	logger       logging.Logger
	tracer       opentracing.Tracer
}

// NewPaperAPI 创建论文API处理器
func NewPaperAPI(
	paperService *service.PaperService,
	logger logging.Logger,
	tracer opentracing.Tracer,
) *PaperAPI {
	return &PaperAPI{
		paperService: paperService,
		logger:       logger,
		tracer:       tracer,
	}
}

// @api_path: /api/paper/versions
// @method: GET
// @content-type: application/json
// @summary: 获取论文版本列表
func (api *PaperAPI) GetVersions(c *gin.Context) {
	span, _ := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "PaperAPI.GetVersions")
	defer span.Finish()

	req := &pb.GetPaperVersionsRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析请求参数失败", "error", err.Error())
		response.ErrorNoData(c, "解析请求参数失败")
		return
	}

	versionResponse, err := api.paperService.GetPaperVersions(c.Request.Context(), req.PaperId)
	if err != nil {
		api.logger.Error("msg", "获取论文版本列表失败", "error", err.Error())
		response.ErrorNoData(c, "获取论文版本列表失败")
		return
	}

	response.Success(c, "success", versionResponse)
	// 生成静态模拟数据 MOCK:TODO
	// versionResponse := pb.GetPaperVersionsResponse{
	// 	// PublicVersions: []*pb.PaperVersionInfo{
	// 	// 	{
	// 	// 		Type:        pb.PaperVersionType_PAPER_PDF,
	// 	// 		Name:        "v1.0",
	// 	// 		PdfId:       req.PaperId,
	// 	// 		JumpUrl:     "https://example.com/papers/paper1_v1.pdf",
	// 	// 		CurVersion:  false,
	// 	// 		LastVersion: false,
	// 	// 		DatePrefix:  "2025-01-15",
	// 	// 	},
	// 	// 	{
	// 	// 		Type:        pb.PaperVersionType_PAPER_PDF,
	// 	// 		Name:        "v2.0",
	// 	// 		PdfId:       req.PaperId,
	// 	// 		JumpUrl:     "https://example.com/papers/paper1_v2.pdf",
	// 	// 		CurVersion:  true,
	// 	// 		LastVersion: true,
	// 	// 		DatePrefix:  "2025-04-10",
	// 	// 	},
	// 	// },
	// 	PrivateVersions: []*pb.PaperVersionInfo{
	// 		{
	// 			Type:        pb.PaperVersionType_PAPER_LINK,
	// 			Name:        "预发布版",
	// 			PdfId:       req.PaperId,
	// 			JumpUrl:     "https://example.com/papers/paper1_prerelease.html",
	// 			CurVersion:  false,
	// 			LastVersion: false,
	// 			DatePrefix:  "2024-12-05",
	// 		},
	// 	},
	// }
}

/*
@api_path: /api/paper/getPaperDetailInfo
@method: POST
@content-type: application/json
@summary: 获取论文详情信息
*/
func (api *PaperAPI) GetPaperDetailInfo(c *gin.Context) {

	span, _ := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "PaperAPI.GetPaperDetailInfo")
	defer span.Finish()

	req := &pb.GetPaperDetailInfoRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析请求参数失败", "error", err.Error())
		response.ErrorNoData(c, "解析请求参数失败")
		return
	}
	paperDetailInfo, err := api.paperService.GetPaperDetailInfo(c.Request.Context(), req.PaperId)
	if err != nil {
		api.logger.Error("msg", "获取论文详情信息失败", "error", err.Error())
		response.ErrorNoData(c, "获取论文详情信息失败")
		return
	}
	response.Success(c, "success", paperDetailInfo)

}
