package pdf

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/middleware"
	"github.com/yb2020/odoc/pkg/registry"
	"github.com/yb2020/odoc/pkg/scheduler"
	userDocService "github.com/yb2020/odoc/services/doc/service"
	noteInterfaces "github.com/yb2020/odoc/services/note/interfaces"
	paperNotesService "github.com/yb2020/odoc/services/note/service"
	ossService "github.com/yb2020/odoc/services/oss/service"
	paperService "github.com/yb2020/odoc/services/paper/service"
	"github.com/yb2020/odoc/services/pdf/api"
	"github.com/yb2020/odoc/services/pdf/dao"
	"github.com/yb2020/odoc/services/pdf/interfaces"
	"github.com/yb2020/odoc/services/pdf/service"
	userService "github.com/yb2020/odoc/services/user/service"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// 编译时类型检查：确保 PdfModule 实现了 registry.Module 接口
var _ registry.Module = (*PdfModule)(nil)

// Module 导出模块实例，用于自动发现和注册
var Module = &PdfModule{}

// PdfModule 实现PDF服务模块
type PdfModule struct {
	db                     *gorm.DB
	logger                 logging.Logger
	tracer                 opentracing.Tracer
	authMiddleware         *middleware.AuthMiddleware
	cfg                    *config.Config
	userService            *userService.UserService
	userDocService         *userDocService.UserDocService
	paperNoteService       noteInterfaces.IPaperNoteService
	paperNoteAccessService *paperNotesService.PaperNoteAccessService
	ossService             ossService.OssServiceInterface
	// DAO实例
	paperPdfDAO             *dao.PaperPDFDAO
	paperPdfSelectRecordDAO *dao.PaperPdfSelectRecordDAO
	pdfAnnotationDAO        *dao.PdfAnnotationDAO
	pdfCommentDAO           *dao.PdfCommentDAO
	pdfMarkDAO              *dao.PdfMarkDAO
	pdfMarkBackupDAO        *dao.PdfMarkBackupDAO
	pdfMarkTagDAO           *dao.PdfMarkTagDAO
	pdfMarkTagRelationDAO   *dao.PdfMarkTagRelationDAO
	pdfReaderSettingDAO     *dao.PdfReaderSettingDAO
	pdfThumbDAO             *dao.PdfThumbDAO
	pdfSummaryDAO           *dao.PdfSummaryDAO

	// 服务实例
	paperPdfService             *service.PaperPdfService
	paperPdfSelectRecordService *service.PaperPdfSelectRecordService
	pdfAnnotationService        *service.PdfAnnotationService
	pdfCommentService           *service.PdfCommentService
	pdfMarkService              interfaces.IPdfMarkService
	pdfMarkBackupService        *service.PdfMarkBackupService
	pdfMarkTagService           *service.PdfMarkTagService
	pdfMarkTagRelationService   *service.PdfMarkTagRelationService
	pdfReaderSettingService     *service.PdfReaderSettingService
	pdfThumbService             *service.PdfThumbService
	paperService                *paperService.PaperService
	paperAccessService          *paperService.PaperAccessService
	paperPdfParsedService       *paperService.PaperPdfParsedService
	pdfParseService             *service.PdfParseService
	pdfSummaryService           *service.PdfSummaryService
	// API实例
	paperPdfAPI   *api.PaperPdfAPI
	pdfParseAPI   *api.PdfParseAPI
	pdfMarkAPI    *api.PdfMarkAPI
	pdfMarkTagAPI *api.PdfMarkTagAPI
}

// NewPdfModule 创建PDF模块
func NewPdfModule(db *gorm.DB, cfg *config.Config, logger logging.Logger,
	tracer opentracing.Tracer, authMiddleware *middleware.AuthMiddleware,
	userService *userService.UserService,
	paperService *paperService.PaperService,
	userDocService *userDocService.UserDocService,
	paperNoteService noteInterfaces.IPaperNoteService,
	paperNoteAccessService *paperNotesService.PaperNoteAccessService,
	ossService ossService.OssServiceInterface,
	paperAccessService *paperService.PaperAccessService,
	paperPdfParsedService *paperService.PaperPdfParsedService,
) *PdfModule {
	return &PdfModule{
		db:                     db,
		logger:                 logger,
		tracer:                 tracer,
		cfg:                    cfg,
		authMiddleware:         authMiddleware,
		userService:            userService,
		paperService:           paperService,
		userDocService:         userDocService,
		paperNoteService:       paperNoteService,
		paperNoteAccessService: paperNoteAccessService,
		ossService:             ossService,
		paperAccessService:     paperAccessService,
		paperPdfParsedService:  paperPdfParsedService,
	}
}

// Name 返回模块名称
func (m *PdfModule) Name() string {
	return "pdf"
}

// Shutdown 停止模块
func (m *PdfModule) Shutdown() error {
	m.logger.Info("msg", "关闭PDF模块")
	return nil
}

// RegisterGRPC 注册gRPC服务
func (m *PdfModule) RegisterGRPC(server *grpc.Server) {
	// TODO: 实现gRPC服务注册
	m.logger.Debug("msg", "PDF模块没有gRPC服务，跳过注册")
}

// RegisterJobSchedulers 注册Job定时任务
func (m *PdfModule) RegisterJobSchedulers(scheduler *scheduler.Scheduler) {
	// TODO: 实现Job定时任务注册
	m.logger.Debug("msg", "PDF模块没有Job定时任务，跳过注册")
}

// RegisterProviders 注册Provider
func (m *PdfModule) RegisterProviders() {
	// TODO: 实现Provider注册
	m.logger.Debug("msg", "PDF模块没有Provider，跳过注册")
}

// Initialize 初始化模块
func (m *PdfModule) Initialize() error {
	m.RegisterProviders()
	m.logger.Info("msg", "初始化PDF模块")
	// 初始化DAO
	m.paperPdfDAO = dao.NewPaperPDFDAO(m.db, m.logger)
	m.paperPdfSelectRecordDAO = dao.NewPaperPdfSelectRecordDAO(m.db, m.logger)
	m.pdfAnnotationDAO = dao.NewPdfAnnotationDAO(m.db, m.logger)
	m.pdfCommentDAO = dao.NewPdfCommentDAO(m.db, m.logger)
	m.pdfMarkDAO = dao.NewPdfMarkDAO(m.db, m.logger)
	m.pdfMarkBackupDAO = dao.NewPdfMarkBackupDAO(m.db, m.logger)
	m.pdfMarkTagDAO = dao.NewPdfMarkTagDAO(m.db, m.logger)
	m.pdfMarkTagRelationDAO = dao.NewPdfMarkTagRelationDAO(m.db, m.logger)
	m.pdfReaderSettingDAO = dao.NewPdfReaderSettingDAO(m.db, m.logger)
	m.pdfThumbDAO = dao.NewPdfThumbDAO(m.db, m.logger)
	m.pdfSummaryDAO = dao.NewPdfSummaryDAO(m.db, m.logger)

	// 初始化服务
	m.paperPdfSelectRecordService = service.NewPaperPdfSelectRecordService(m.logger, m.tracer, m.paperPdfSelectRecordDAO)
	m.pdfAnnotationService = service.NewPdfAnnotationService(m.logger, m.tracer, m.pdfAnnotationDAO)
	m.pdfCommentService = service.NewPdfCommentService(m.logger, m.tracer, m.pdfCommentDAO)
	m.pdfMarkBackupService = service.NewPdfMarkBackupService(m.logger, m.tracer, m.pdfMarkBackupDAO)
	m.pdfMarkTagRelationService = service.NewPdfMarkTagRelationService(m.logger, m.tracer, m.pdfMarkTagRelationDAO)

	// 解决服务间的循环依赖问题
	// 1. 先初始化 paperPdfService 和 pdfMarkTagService，此时它们的 pdfMarkService 依赖暂时传入 nil
	m.paperPdfService = service.NewPaperPdfService(m.cfg, m.logger, m.tracer, m.paperPdfDAO, m.ossService, m.paperAccessService, m.paperService, m.userService, m.userDocService, m.paperNoteService, m.paperNoteAccessService, nil)
	m.pdfMarkTagService = service.NewPdfMarkTagService(m.logger, m.tracer, m.pdfMarkTagDAO, nil, m.pdfMarkTagRelationService, m.paperNoteService)

	// 2. 接着初始化 pdfMarkService，此时它可以安全地注入已经实例化的 paperPdfService 和 pdfMarkTagService
	m.pdfMarkService = service.NewPdfMarkService(m.logger, m.tracer, m.pdfMarkDAO, m.pdfMarkTagRelationService, m.pdfMarkTagService, m.paperNoteService, m.paperPdfService)

	// 3. 最后，通过 Setter 方法将 pdfMarkService 实例注入回 paperPdfService 和 pdfMarkTagService，完成依赖闭环
	m.paperPdfService.SetPdfMarkService(m.pdfMarkService)
	m.pdfMarkTagService.SetPdfMarkService(m.pdfMarkService)

	m.pdfReaderSettingService = service.NewPdfReaderSettingService(m.logger, m.tracer, m.pdfReaderSettingDAO)
	m.pdfThumbService = service.NewPdfThumbService(m.logger, m.tracer, m.pdfThumbDAO)

	m.pdfSummaryService = service.NewPdfSummaryService(m.logger, m.tracer, m.pdfSummaryDAO)

	// 初始化API
	m.paperPdfAPI = api.NewPaperPdfAPI(m.paperPdfService, m.logger, m.tracer, m.pdfReaderSettingService)
	m.pdfParseAPI = api.NewPdfParseAPI(m.paperPdfService, m.logger, m.tracer, m.paperService, m.userService, m.userDocService, m.paperNoteService, m.paperNoteAccessService, m.pdfReaderSettingService, m.pdfParseService)
	m.pdfMarkAPI = api.NewPdfMarkAPI(m.pdfMarkService, m.logger, m.tracer)

	m.pdfMarkTagAPI = api.NewPdfMarkTagAPI(m.pdfMarkTagService, m.logger, m.tracer)

	return nil
}

// RegisterRoutes 注册路由
func (m *PdfModule) RegisterRoutes(r *gin.Engine) {
	pdfGroup := r.Group("/api/pdf")
	pdfGroup.Use(m.authMiddleware.AuthRequired())
	{
		// PDF相关API路由
		pdfGroup.POST("/getPdfStatusInfo/v2", m.paperPdfAPI.GetPdfStatusInfo)
		pdfGroup.POST("/pdfReader/getSetting", m.paperPdfAPI.GetPdfReaderSetting)
		pdfGroup.POST("/pdfReader/recordSetting", m.paperPdfAPI.RecordPdfReaderSetting)

		// PDF解析API路由
		pdfGroup.POST("/parser/getReferenceMarkers", m.pdfParseAPI.GetReferenceMarkers)
		pdfGroup.POST("/parser/getFiguresAndTables", m.pdfParseAPI.GetFiguresAndTables)
		pdfGroup.POST("/parser/getReference", m.pdfParseAPI.GetReference)
		pdfGroup.POST("/parser/getCatalogue", m.pdfParseAPI.GetCatalogue)
		//重新解析
		pdfGroup.POST("/parser/reParse", m.pdfParseAPI.ReParse)

		// PDF标记API路由
		pdfGroup.GET("/pdfMark/v3/web/getByNote", m.pdfMarkAPI.GetNoteAnnotationListByNoteId)
		pdfGroup.GET("/pdfMark/v3/web/draw/getByNote", m.pdfMarkAPI.GetDrawNoteAnnotationListByNoteId)
		pdfGroup.GET("/pdfMark/v2/web/hotSelect", m.pdfMarkAPI.HotSelect)
		pdfGroup.POST("/pdfMark/v2/web/save", m.pdfMarkAPI.SavePdfMark)
		pdfGroup.POST("/pdfMark/v2/web/update", m.pdfMarkAPI.UpdatePdfMark)
		pdfGroup.POST("/pdfMark/v2/web/delete", m.pdfMarkAPI.DeletePdfMark)
		pdfGroup.POST("/pdfMark/v2/web/getMyNoteMarkList", m.pdfMarkAPI.GetMyNoteMarkListReq)

		// PDF标记标签API路由
		pdfGroup.GET("/marktag/tags", m.pdfMarkTagAPI.GetPdfMarkTagsRequest)
		pdfGroup.POST("/marktag/save", m.pdfMarkTagAPI.CreatePdfMarkTagRequest)
		pdfGroup.POST("/marktag/update", m.pdfMarkTagAPI.RenameAnnotateTagRequest)
		pdfGroup.POST("/marktag/delete", m.pdfMarkTagAPI.DeleteAnnotateTagRequest)
		pdfGroup.POST("/marktag/relation/mark/save", m.pdfMarkTagAPI.AddTagToAnnotateRequest)
		pdfGroup.POST("/marktag/relation/mark/delete", m.pdfMarkTagAPI.DeleteTagToAnnotateRequest)

	}
}

// GetPaperPdfService 获取论文PDF服务
func (m *PdfModule) GetPaperPdfService() *service.PaperPdfService {
	return m.paperPdfService
}

// GetPaperPdfSelectRecordService 获取论文PDF选择记录服务
func (m *PdfModule) GetPaperPdfSelectRecordService() *service.PaperPdfSelectRecordService {
	return m.paperPdfSelectRecordService
}

// GetPdfAnnotationService 获取PDF注释服务
func (m *PdfModule) GetPdfAnnotationService() *service.PdfAnnotationService {
	return m.pdfAnnotationService
}

// GetPdfCommentService 获取PDF评论服务
func (m *PdfModule) GetPdfCommentService() *service.PdfCommentService {
	return m.pdfCommentService
}

// GetPdfMarkService 获取PDF标记服务
func (m *PdfModule) GetPdfMarkService() interfaces.IPdfMarkService {
	return m.pdfMarkService
}

// GetPdfMarkBackupService 获取PDF标记备份服务
func (m *PdfModule) GetPdfMarkBackupService() *service.PdfMarkBackupService {
	return m.pdfMarkBackupService
}

// GetPdfMarkTagService 获取PDF标记标签服务
func (m *PdfModule) GetPdfMarkTagService() *service.PdfMarkTagService {
	return m.pdfMarkTagService
}

// GetPdfMarkTagRelationService 获取PDF标记标签关系服务
func (m *PdfModule) GetPdfMarkTagRelationService() *service.PdfMarkTagRelationService {
	return m.pdfMarkTagRelationService
}

// GetPdfReaderSettingService 获取PDF阅读设置服务
func (m *PdfModule) GetPdfReaderSettingService() *service.PdfReaderSettingService {
	return m.pdfReaderSettingService
}

// GetPdfThumbService 获取PDF缩略图服务
func (m *PdfModule) GetPdfThumbService() *service.PdfThumbService {
	return m.pdfThumbService
}

// GetPdfParseService 获取PDF解析服务
func (m *PdfModule) GetPdfParseService() *service.PdfParseService {
	return m.pdfParseService
}

// GetPdfSummaryService 获取PDF摘要服务
func (m *PdfModule) GetPdfSummaryService() *service.PdfSummaryService {
	return m.pdfSummaryService
}
