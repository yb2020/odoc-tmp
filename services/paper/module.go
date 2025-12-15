package paper

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/scheduler"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/middleware"
	"github.com/yb2020/odoc/pkg/registry"
	"github.com/yb2020/odoc/services/paper/api"
	"github.com/yb2020/odoc/services/paper/dao"
	"github.com/yb2020/odoc/services/paper/service"
	userService "github.com/yb2020/odoc/services/user/service"
)

// 编译时类型检查：确保 PaperModule 实现了 registry.Module 接口
var _ registry.Module = (*PaperModule)(nil)

// Module 导出模块实例，用于自动发现和注册
var Module = &PaperModule{}

// PaperModule 实现论文服务模块
type PaperModule struct {
	db                      *gorm.DB
	logger                  logging.Logger
	tracer                  opentracing.Tracer
	authMiddleware          *middleware.AuthMiddleware
	cfg                     *config.Config
	paperDao                *dao.PaperDAO
	paperAccessDao          *dao.PaperAccessDAO
	paperAnswerDao          *dao.PaperAnswerDAO
	paperAttachmentDao      *dao.PaperAttachmentDAO
	paperCommentDao         *dao.PaperCommentDAO
	paperCommentApprovalDao *dao.PaperCommentApprovalDAO
	paperQuestionDao        *dao.PaperQuestionDAO
	paperResourcesDao       *dao.PaperResourcesDAO
	paperJcrDao             *dao.PaperJcrDAO
	paperPdfParsedDao       *dao.PaperPdfParsedDAO

	userService *userService.UserService
	// 服务实例
	paperService                *service.PaperService
	paperAccessService          *service.PaperAccessService
	paperAnswerService          *service.PaperAnswerService
	paperAttachmentService      *service.PaperAttachmentService
	paperCommentService         *service.PaperCommentService
	paperCommentApprovalService *service.PaperCommentApprovalService
	paperQuestionService        *service.PaperQuestionService
	paperResourcesService       *service.PaperResourcesService
	paperJcrService             *service.PaperJcrService
	paperPdfParsedService       *service.PaperPdfParsedService

	// API实例
	paperAPI *api.PaperAPI
}

// NewPaperModule 创建论文模块
func NewPaperModule(db *gorm.DB, cfg *config.Config, logger logging.Logger,
	tracer opentracing.Tracer, authMiddleware *middleware.AuthMiddleware,
	userService *userService.UserService,
) *PaperModule {
	return &PaperModule{
		db:             db,
		logger:         logger,
		tracer:         tracer,
		cfg:            cfg,
		authMiddleware: authMiddleware,
		userService:    userService,
	}
}

// Name 返回模块名称
func (m *PaperModule) Name() string {
	return "paper"
}

// Shutdown 停止模块
func (m *PaperModule) Shutdown() error {
	m.logger.Info("msg", "关闭论文模块")
	return nil
}

// RegisterGRPC 注册gRPC服务
func (m *PaperModule) RegisterGRPC(server *grpc.Server) {
	// TODO: 实现gRPC服务注册
	m.logger.Debug("msg", "论文模块没有gRPC服务，跳过注册")
}

// RegisterJobSchedulers 注册Job定时任务
func (m *PaperModule) RegisterJobSchedulers(scheduler *scheduler.Scheduler) {
	// TODO: 实现Job定时任务注册
	m.logger.Debug("msg", "论文模块没有Job定时任务，跳过注册")
}

func (m *PaperModule) RegisterProviders() {
	m.logger.Debug("msg", "论文模块没有Provider，跳过注册")
}

// Initialize 初始化模块
func (m *PaperModule) Initialize() error {
	m.RegisterProviders()
	m.logger.Info("msg", "初始化论文模块")

	// 初始化DAO
	m.paperDao = dao.NewPaperDAO(m.db, m.logger)
	m.paperAccessDao = dao.NewPaperAccessDAO(m.db, m.logger)
	m.paperAnswerDao = dao.NewPaperAnswerDAO(m.db, m.logger)
	m.paperAttachmentDao = dao.NewPaperAttachmentDAO(m.db, m.logger)
	m.paperCommentDao = dao.NewPaperCommentDAO(m.db, m.logger)
	m.paperCommentApprovalDao = dao.NewPaperCommentApprovalDAO(m.db, m.logger)
	m.paperQuestionDao = dao.NewPaperQuestionDAO(m.db, m.logger)
	m.paperResourcesDao = dao.NewPaperResourcesDAO(m.db, m.logger)
	// 初始化paper_jcr DAO
	m.paperJcrDao = dao.NewPaperJcrDAO(m.db, m.logger)
	m.paperPdfParsedDao = dao.NewPaperPdfParsedDAO(m.db, m.logger)

	// 初始化服务
	m.paperService = service.NewPaperService(m.logger, m.tracer, m.paperDao)
	m.paperAccessService = service.NewPaperAccessService(m.logger, m.tracer, m.paperAccessDao)
	m.paperAnswerService = service.NewPaperAnswerService(m.logger, m.tracer, m.paperAnswerDao)
	m.paperAttachmentService = service.NewPaperAttachmentService(m.logger, m.tracer, m.paperAttachmentDao)
	m.paperCommentService = service.NewPaperCommentService(m.logger, m.tracer, m.paperCommentDao)
	m.paperCommentApprovalService = service.NewPaperCommentApprovalService(m.logger, m.tracer, m.paperCommentApprovalDao)
	m.paperQuestionService = service.NewPaperQuestionService(m.logger, m.tracer, m.paperQuestionDao)
	m.paperResourcesService = service.NewPaperResourcesService(m.logger, m.tracer, m.paperResourcesDao)
	// 初始化paper_jcr Service
	m.paperJcrService = service.NewPaperJcrService(m.logger, m.tracer, m.paperJcrDao)
	// 初始化paper_pdf_parsed Service
	m.paperPdfParsedService = service.NewPaperPdfParsedService(m.logger, m.tracer, m.paperPdfParsedDao)

	// 初始化API
	m.paperAPI = api.NewPaperAPI(m.paperService, m.logger, m.tracer)

	return nil
}

// RegisterRoutes 注册路由
func (m *PaperModule) RegisterRoutes(r *gin.Engine) {
	paperGroup := r.Group("/api/paper")
	paperGroup.Use(m.authMiddleware.AuthRequired())
	{
		// 论文相关API路由
		// paperGroup.POST("/create", m.paperAPI.CreatePaper)
		// paperGroup.POST("/get", m.paperAPI.GetPaperById)
		// paperGroup.POST("/getBaseInfo", m.paperAPI.GetPaperBaseInfoById)
		// paperGroup.POST("/update", m.paperAPI.UpdatePaper)
		// paperGroup.POST("/delete", m.paperAPI.DeletePaper)
		// paperGroup.POST("/list", m.paperAPI.ListPapers)

		paperGroup.GET("/versions", m.paperAPI.GetVersions)
		paperGroup.POST("/getPaperDetailInfo", m.paperAPI.GetPaperDetailInfo)
	}
}

func (m *PaperModule) GetPaperService() *service.PaperService {
	return m.paperService
}

// GetPaperJcrService 返回paperJcrService实例
func (m *PaperModule) GetPaperJcrService() *service.PaperJcrService {
	return m.paperJcrService
}

func (m *PaperModule) GetPaperAccessService() *service.PaperAccessService {
	return m.paperAccessService
}

func (m *PaperModule) GetPaperPdfParsedService() *service.PaperPdfParsedService {
	return m.paperPdfParsedService
}
