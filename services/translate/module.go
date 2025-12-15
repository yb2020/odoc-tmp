package translate

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/config"
	ocr "github.com/yb2020/odoc/external/ocr/api"
	"github.com/yb2020/odoc/internal/database"
	"github.com/yb2020/odoc/pkg/cache"
	"github.com/yb2020/odoc/pkg/distlock"
	"github.com/yb2020/odoc/pkg/http_client"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/metrics"
	"github.com/yb2020/odoc/pkg/middleware"
	"github.com/yb2020/odoc/pkg/ratelimit"
	"github.com/yb2020/odoc/pkg/registry"
	"github.com/yb2020/odoc/pkg/scheduler"
	docService "github.com/yb2020/odoc/services/doc/service"
	membershipService "github.com/yb2020/odoc/services/membership/interfaces"
	noteInterface "github.com/yb2020/odoc/services/note/interfaces"
	ossService "github.com/yb2020/odoc/services/oss/service"
	pdfService "github.com/yb2020/odoc/services/pdf/service"
	"github.com/yb2020/odoc/services/translate/api"
	"github.com/yb2020/odoc/services/translate/dao"
	"github.com/yb2020/odoc/services/translate/service"
	translateService "github.com/yb2020/odoc/services/translate/service"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// 编译时类型检查：确保 TranslateModule 实现了 registry.Module 接口
var _ registry.Module = (*TranslateModule)(nil)

// Module 导出模块实例，用于自动发现和注册
var Module = &TranslateModule{}

// TranslateModule 翻译模块
type TranslateModule struct {
	db                       *gorm.DB
	logger                   logging.Logger
	tracer                   opentracing.Tracer
	config                   *config.Config
	glossaryAPI              *api.GlossaryAPI
	textTranslateAPI         *api.TextTranslateAPI
	ocrTextTranslateAPI      *api.OCRTextTranslateAPI
	authMiddleware           *middleware.AuthMiddleware
	httpClient               http_client.HttpClient
	fullTextTranslateAPI     *api.FullTextTranslateAPI
	paperNoteService         noteInterface.IPaperNoteService
	paperPdfService          *pdfService.PaperPdfService
	fullTextTranslateService *translateService.FullTextTranslateService
	userDocService           *docService.UserDocService
	ossService               ossService.OssServiceInterface
	lockTemplate             *distlock.LockTemplate
	membershipService        membershipService.IMembershipService
}

// NewModule 创建翻译模块
func NewTranslateModule(db *gorm.DB,
	config *config.Config,
	logger logging.Logger,
	tracer opentracing.Tracer,
	authMiddleware *middleware.AuthMiddleware,
	httpClient http_client.HttpClient,
	lockTemplate *distlock.LockTemplate,
	paperNoteService noteInterface.IPaperNoteService,
	paperPdfService *pdfService.PaperPdfService,
	userDocService *docService.UserDocService,
	ossService ossService.OssServiceInterface,
	membershipService membershipService.IMembershipService,
) *TranslateModule {
	return &TranslateModule{
		db:                db,
		logger:            logger,
		tracer:            tracer,
		config:            config,
		authMiddleware:    authMiddleware,
		httpClient:        httpClient,
		lockTemplate:      lockTemplate,
		userDocService:    userDocService,
		paperNoteService:  paperNoteService,
		paperPdfService:   paperPdfService,
		ossService:        ossService,
		membershipService: membershipService,
	}
}

// Name 返回模块名称
func (m *TranslateModule) Name() string {
	return "translate"
}

// RegisterProviders 注册Provider
func (m *TranslateModule) RegisterProviders() {
	// TODO: 实现Provider注册
	m.logger.Debug("msg", "翻译模块没有Provider，跳过注册")
}

// Initialize 初始化模块
func (m *TranslateModule) Initialize() error {
	m.RegisterProviders()
	m.logger.Info("msg", "初始化翻译模块")
	// glossary
	glossaryDAO := dao.NewGlossaryDAO(m.db, m.logger)
	glossaryService := service.NewGlossaryService(m.config, m.logger, m.tracer, glossaryDAO)
	m.glossaryAPI = api.NewGlossaryAPI(glossaryService, m.logger, m.tracer)

	// external ocr api
	imageOCRApiService := ocr.NewImageOCRApiService(m.logger, m.config.Translate.OCR.ExtractTextURL, m.config, m.httpClient)

	// word pronunciation
	wordPronunciationDAO := dao.NewWordPronunciationDAO(m.db, m.logger)
	wordPronunciationService := service.NewWordPronunciationService(m.logger, m.config, m.tracer, wordPronunciationDAO)

	// 创建限流服务
	redisClient := database.GetRedisClient()
	rateLimiterService := ratelimit.NewRateLimiterService(redisClient, m.logger)

	// 创建metrics
	metrics := metrics.NewMetrics(m.Name())

	// text translate
	textTranslateDAO := dao.NewTextTranslateDAO(m.db, m.logger)
	textTranslateService := service.NewTextTranslateService(m.config, m.logger, m.tracer, m.httpClient, textTranslateDAO)
	m.textTranslateAPI = api.NewTextTranslateAPI(m.config, m.logger, m.tracer, textTranslateService, glossaryService, wordPronunciationService, rateLimiterService, m.membershipService)

	cacheClient := cache.NewCache(m.logger, 30*time.Minute, m.Name())

	// full text translate
	fullTextTranslateDAO := dao.NewFullTextTranslateHistoryDAO(m.db, m.logger)
	fullTextTranslateService := service.NewFullTextTranslateService(m.config, m.logger, m.tracer, cacheClient, m.httpClient, metrics, fullTextTranslateDAO, m.lockTemplate, m.paperNoteService, m.paperPdfService, m.userDocService, m.ossService, m.membershipService)
	m.fullTextTranslateAPI = api.NewFullTextTranslateAPI(m.config, m.logger, m.tracer, fullTextTranslateService, m.membershipService)
	m.fullTextTranslateService = fullTextTranslateService

	// ocr translate
	ocrTranslateDAO := dao.NewOCRTranslateDAO(m.db, m.logger)
	ocrTranslateService := service.NewOCRTranslateService(m.config, m.logger, m.tracer, ocrTranslateDAO, textTranslateService, glossaryService, rateLimiterService, imageOCRApiService)
	m.ocrTextTranslateAPI = api.NewOCRTextTranslateAPI(m.config, m.logger, m.tracer, ocrTranslateService, rateLimiterService, m.membershipService)

	return nil
}

// Shutdown 关闭模块
func (m *TranslateModule) Shutdown() error {
	m.logger.Info("msg", "关闭翻译模块")
	return nil
}

// RegisterGRPC 注册gRPC服务
func (m *TranslateModule) RegisterGRPC(server *grpc.Server) {
	// 翻译模块没有gRPC服务，不需要注册
	m.logger.Debug("msg", "翻译模块没有gRPC服务，跳过注册")
}

// RegisterJobSchedulers 注册Job定时任务
func (m *TranslateModule) RegisterJobSchedulers(scheduler *scheduler.Scheduler) {
	// TODO: 实现Job定时任务注册
	m.logger.Debug("msg", "翻译模块没有Job定时任务，跳过注册")
}

// RegisterRoutes 注册路由
func (m *TranslateModule) RegisterRoutes(r *gin.Engine) {
	glossaryGroup := r.Group("/api/glossary")
	glossaryGroup.Use(m.authMiddleware.AuthRequired())
	{
		// 术语条目相关API
		glossaryGroup.POST("/add", m.glossaryAPI.AddGlossary)
		glossaryGroup.POST("/update", m.glossaryAPI.UpdateGlossary)
		glossaryGroup.POST("/delete", m.glossaryAPI.DeleteGlossary)
		glossaryGroup.POST("/list", m.glossaryAPI.GetGlossaryList)
	}

	textTranslateGroup := r.Group("/api/text")
	textTranslateGroup.Use(m.authMiddleware.AuthRequired())
	{
		// 获取翻译标签
		textTranslateGroup.GET("/getTranslateTabs", m.textTranslateAPI.GetTranslateTabs)
		// 文本翻译相关API
		textTranslateGroup.POST("/translate", m.textTranslateAPI.Translate)
		// 翻译反馈
		textTranslateGroup.POST("/translate/correct", m.textTranslateAPI.TranslateCorrect)
		// ocr翻译
		textTranslateGroup.POST("/ocr/translate", m.ocrTextTranslateAPI.OCRTranslate)
		// ocr提取文本
		textTranslateGroup.POST("/ocr/extractText", m.ocrTextTranslateAPI.ExtractText)
		// 流式翻译
		textTranslateGroup.POST("/translate/completions", m.textTranslateAPI.AiTranslate)
	}

	fullTextTranslateGroup := r.Group("/api/fullTextTranslate")
	fullTextTranslateGroup.Use(m.authMiddleware.AuthRequired())
	{
		// 获取用户全文翻译历史记录
		fullTextTranslateGroup.GET("/getHistoryList", m.fullTextTranslateAPI.GetHistoryList)
		// 获取用户全文翻译权限信息
		fullTextTranslateGroup.GET("/getRightInfo", m.fullTextTranslateAPI.GetRightInfo)
		// 发起全文翻译
		fullTextTranslateGroup.POST("/translate", m.fullTextTranslateAPI.Translate)
		// 查询全文翻译状态
		fullTextTranslateGroup.GET("/getTranslateStatus", m.fullTextTranslateAPI.GetTranslateStatus)
	}

	fullTextTranslateServiceGroup := r.Group("/services")
	fullTextTranslateServiceGroup.Use(m.authMiddleware.ServiceAuthRequired())
	{
		// 服务之间内部翻译接口
		fullTextTranslateServiceGroup.POST("/text/translate/internal", m.textTranslateAPI.TranslateInternal)
	}

	fullTextTranslateAdminGroup := r.Group("/api/admin/fullTextTranslate")
	fullTextTranslateAdminGroup.Use(m.authMiddleware.AuthRequired())
	{
		// // 管理员重试全文翻译
		// fullTextTranslateAdminGroup.POST("/reTranslate", m.fullTextTranslateAPI.ReTranslate)
		// // 管理员获取重试翻译结果
		// fullTextTranslateAdminGroup.POST("/getReTranslateResult", m.fullTextTranslateAPI.GetReTranslateResult)
	}
}

// GetFullTextTranslateService 获取全文翻译服务
func (m *TranslateModule) GetFullTextTranslateService() *translateService.FullTextTranslateService {
	return m.fullTextTranslateService
}
