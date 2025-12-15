package parse

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/http_client"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/middleware"
	"github.com/yb2020/odoc/pkg/registry"
	"github.com/yb2020/odoc/pkg/scheduler"
	ossService "github.com/yb2020/odoc/services/oss/service"
	paperService "github.com/yb2020/odoc/services/paper/service"
	"github.com/yb2020/odoc/services/parse/api"
	"github.com/yb2020/odoc/services/parse/service"
	"google.golang.org/grpc"
)

var _ registry.Module = (*ParseModule)(nil)

// Module 导出模块实例
var Module = &ParseModule{}

// ParseModule PDF解析模块
type ParseModule struct {
	config                *config.Config
	logger                logging.Logger
	tracer                opentracing.Tracer
	httpClient            http_client.HttpClient
	parseOperateService   *service.ParseOperateService
	grobidPdfParseService *service.GrobidPDFParseService
	mineruPdfParseService *service.MineruPDFParseService
	parseApi              *api.TestParseAPI
	authMiddleware        *middleware.AuthMiddleware
	ossService            ossService.OssServiceInterface
	paperPdfParsedService *paperService.PaperPdfParsedService
}

// NewParseModule 创建新的Parse模块实例
func NewParseModule(
	config *config.Config,
	logger logging.Logger,
	tracer opentracing.Tracer,
	httpClient http_client.HttpClient,
	ossService ossService.OssServiceInterface,
	paperPdfParsedService *paperService.PaperPdfParsedService,
) *ParseModule {
	return &ParseModule{
		config:                config,
		logger:                logger,
		tracer:                tracer,
		httpClient:            httpClient,
		ossService:            ossService,
		paperPdfParsedService: paperPdfParsedService,
	}
}

// RegisterProviders 注册Provider
func (m *ParseModule) RegisterProviders() {
	m.logger.Debug("msg", "解析模块没有Provider，跳过注册")
}

// Initialize 初始化模块
func (m *ParseModule) Initialize() error {
	m.RegisterProviders()
	m.logger.Info("msg", "初始化Parse模块")

	// 创建Service实例
	m.parseOperateService = service.NewParseOperateService(
		m.config,
		m.tracer,
		m.logger,
		m.paperPdfParsedService,
		m.ossService,
		m.httpClient,
	)
	m.grobidPdfParseService = service.NewGrobidPDFParseService(
		m.config,
		m.logger,
		m.tracer,
		m.httpClient,
		m.parseOperateService,
	)
	m.mineruPdfParseService = service.NewMineruPDFParseService(
		m.config,
		m.logger,
		m.tracer,
		m.httpClient,
		m.parseOperateService,
		m.ossService,
	)

	m.parseApi = api.NewTestParseApi(m.tracer, m.grobidPdfParseService, m.mineruPdfParseService)
	return nil
}

// SetAuthMiddleware 设置认证中间件
func (m *ParseModule) SetAuthMiddleware(middleware *middleware.AuthMiddleware) {
	m.authMiddleware = middleware
}

// GetGrobidPDFParseService 获取GrobidPDF解析服务实例
func (m *ParseModule) GetGrobidPDFParseService() *service.GrobidPDFParseService {
	return m.grobidPdfParseService
}

// GetParseOperateService 获取解析操作服务实例
func (m *ParseModule) GetParseOperateService() *service.ParseOperateService {
	return m.parseOperateService
}

// GetMineruPDFParseService 获取MineruPDF解析服务实例
func (m *ParseModule) GetMineruPDFParseService() *service.MineruPDFParseService {
	return m.mineruPdfParseService
}

// Name 返回模块名称
func (m *ParseModule) Name() string {
	return "parse"
}

// RegisterRoutes 注册路由
func (m *ParseModule) RegisterRoutes(r *gin.Engine) {
	parseServiceGroup := r.Group("/api/grobid")
	// parseServiceGroup.Use(m.authMiddleware.ServiceAuthRequired())
	{
		// 服务之间内部解析接口
		parseServiceGroup.POST("/v8/parseFulltext", m.parseApi.ParseFulltextV8)
		parseServiceGroup.POST("/v8/parseHeader", m.parseApi.ParseHeaderV8)
		parseServiceGroup.POST("/v8/parseMetadata", m.parseApi.ParseMetadataV8)
		parseServiceGroup.POST("/v8/parseMineru", m.parseApi.ParseMineru)
	}
}

// RegisterGRPC 注册gRPC服务
func (m *ParseModule) RegisterGRPC(server *grpc.Server) {
	// TODO: 实现gRPC服务注册
}

// RegisterJobSchedulers 注册Job定时任务
func (m *ParseModule) RegisterJobSchedulers(scheduler *scheduler.Scheduler) {
	// TODO: 实现Job定时任务注册
	m.logger.Debug("msg", "解析模块没有Job定时任务，跳过注册")
}

// Shutdown 关闭模块
func (m *ParseModule) Shutdown() error {
	m.logger.Info("msg", "关闭Parse模块")
	return nil
}
