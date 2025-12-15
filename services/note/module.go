package note

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/scheduler"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/middleware"
	"github.com/yb2020/odoc/pkg/registry"
	userDocService "github.com/yb2020/odoc/services/doc/service"
	"github.com/yb2020/odoc/services/note/api"
	"github.com/yb2020/odoc/services/note/dao"
	noteInterface "github.com/yb2020/odoc/services/note/interfaces"
	"github.com/yb2020/odoc/services/note/service"
	paperService "github.com/yb2020/odoc/services/paper/service"
	pdfInterface "github.com/yb2020/odoc/services/pdf/interfaces"
	userService "github.com/yb2020/odoc/services/user/service"
)

// 编译时类型检查：确保 NoteModule 实现了 registry.Module 接口
var _ registry.Module = (*NoteModule)(nil)

// Module 导出模块实例，用于自动发现和注册
var Module = &NoteModule{}

// NoteModule 实现论文笔记服务模块
type NoteModule struct {
	db                     *gorm.DB
	logger                 logging.Logger
	tracer                 opentracing.Tracer
	paperNoteAPI           *api.PaperNoteAPI
	noteShapeAPI           *api.NoteShapeAPI
	noteSummaryAPI         *api.NoteSummaryAPI
	noteWordAPI            *api.NoteWordAPI
	noteReadLocationAPI    *api.NoteReadLocationAPI
	noteManageAPI          *api.NoteManageAPI
	noteSummaryService     *service.NoteSummaryService
	authMiddleware         *middleware.AuthMiddleware
	cfg                    *config.Config
	userService            *userService.UserService
	paperNoteService       noteInterface.IPaperNoteService
	paperNoteAccessService *service.PaperNoteAccessService
	userDocService         *userDocService.UserDocService
	pdfService             pdfInterface.IPaperPdfService
	paperService           *paperService.PaperService
	noteWordService        noteInterface.INoteWordService
}

// NewModule 创建论文笔记模块
func NewNoteModule(db *gorm.DB, cfg *config.Config, logger logging.Logger,
	tracer opentracing.Tracer, authMiddleware *middleware.AuthMiddleware,
	userService *userService.UserService,
	userDocService *userDocService.UserDocService,
	paperService *paperService.PaperService,
) *NoteModule {
	return &NoteModule{
		db:             db,
		logger:         logger,
		tracer:         tracer,
		cfg:            cfg,
		authMiddleware: authMiddleware,
		userService:    userService,
		userDocService: userDocService,
		paperService:   paperService,
	}
}

// Name 返回模块名称
func (m *NoteModule) Name() string {
	return "note"
}

// Shutdown 停止模块
func (m *NoteModule) Shutdown() error {
	m.logger.Info("msg", "关闭论文笔记模块")
	return nil
}

// RegisterGRPC 注册gRPC服务
func (m *NoteModule) RegisterGRPC(server *grpc.Server) {
	// TODO: 实现gRPC服务注册
	m.logger.Debug("msg", "翻译模块没有gRPC服务，跳过注册")
}

// RegisterJobSchedulers 注册Job定时任务
func (m *NoteModule) RegisterJobSchedulers(scheduler *scheduler.Scheduler) {
	// TODO: 实现Job定时任务注册
	m.logger.Debug("msg", "翻译模块没有Job定时任务，跳过注册")
}

// RegisterProviders 注册Provider
func (m *NoteModule) RegisterProviders() {
	// TODO: 实现Provider注册
	m.logger.Debug("msg", "翻译模块没有Provider，跳过注册")
}

// Initialize 初始化模块
func (m *NoteModule) Initialize() error {
	m.RegisterProviders()
	m.logger.Info("msg", "初始化论文笔记模块")

	//外部依赖模块====

	// paperService := paperService.NewPaperService(m.logger, m.tracer, paperDao.NewPaperDAO(m.db, m.logger))

	// 创建PaperNoteService
	paperNoteDAO := dao.NewPaperNoteDAO(m.db, m.logger)
	m.paperNoteService = service.NewPaperNoteService(m.logger, m.tracer, paperNoteDAO, m.pdfService, m.userDocService, m.userService, m.paperService, m.noteWordService)

	// 创建PaperNoteAPI
	m.paperNoteAPI = api.NewPaperNoteAPI(m.paperNoteService, m.pdfService, m.userDocService, m.userService, m.paperService, m.logger, m.tracer)

	// 创建PaperNoteAccessService
	paperNoteAccessDAO := dao.NewPaperNoteAccessDAO(m.db, m.logger)
	m.paperNoteAccessService = service.NewPaperNoteAccessService(m.logger, m.tracer, paperNoteAccessDAO)

	// 创建NoteShapeAPI
	noteShapeDAO := dao.NewNoteShapeDAO(m.db, m.logger)
	noteShapeService := service.NewNoteShapeService(m.logger, m.tracer, noteShapeDAO)
	m.noteShapeAPI = api.NewNoteShapeAPI(noteShapeService, m.paperNoteService, m.paperNoteAccessService, m.logger, m.tracer)

	// 创建NoteSummaryAPI
	noteSummaryDAO := dao.NewNoteSummaryDAO(m.db, m.logger)
	m.noteSummaryService = service.NewNoteSummaryService(m.logger, m.tracer, noteSummaryDAO, m.userDocService, m.paperNoteService, m.pdfService)
	m.noteSummaryAPI = api.NewNoteSummaryAPI(m.noteSummaryService, m.userDocService, m.logger, m.tracer, m.paperNoteService)

	// 创建NoteWordAPI
	noteWordDAO := dao.NewNoteWordDAO(m.db, m.logger)
	noteWordConfigDAO := dao.NewNoteWordConfigDAO(m.db, m.logger)
	noteWordConfigService := service.NewNoteWordConfigService(m.logger, m.tracer, noteWordConfigDAO)
	m.noteWordService = service.NewNoteWordService(m.logger, m.tracer, noteWordDAO, noteWordConfigService, m.userDocService, m.paperNoteService)
	m.noteWordAPI = api.NewNoteWordAPI(m.noteWordService, noteWordConfigService, m.logger, m.tracer)
	m.paperNoteService.SetNoteWordService(m.noteWordService)

	// 创建NoteReadLocationService
	noteReadLocationDAO := dao.NewNoteReadLocationDAO(m.db, m.logger)
	noteReadLocationService := service.NewNoteReadLocationService(m.logger, m.tracer, noteReadLocationDAO)

	// 创建NoteReadLocationAPI
	m.noteReadLocationAPI = api.NewNoteReadLocationAPI(noteReadLocationService, m.logger, m.tracer)

	// 创建NoteReadLocationAPI
	m.noteManageAPI = api.NewNoteManageAPI(m.logger, m.tracer, m.noteWordService, m.noteSummaryService, m.pdfService)

	return nil
}

// RegisterRoutes 注册路由
func (m *NoteModule) RegisterRoutes(r *gin.Engine) {
	noteGroup := r.Group("/api/note")
	noteGroup.Use(m.authMiddleware.AuthRequired())
	{
		// 论文笔记相关API旧接口扩展
		noteGroup.POST("/paperNote/getPaperNoteBaseInfoById", m.paperNoteAPI.GetPaperNoteBaseInfoById)
		noteGroup.POST("/paperNote/getOwnerPaperNoteBaseInfo", m.paperNoteAPI.GetOwnerPaperNoteBaseInfo)

		// 笔记形状相关API
		noteGroup.POST("/noteShape/getList", m.noteShapeAPI.GetNoteShapesByNoteId)
		noteGroup.POST("/noteShape/save", m.noteShapeAPI.SaveShapeRequest)
		noteGroup.POST("/noteShape/delete", m.noteShapeAPI.DeleteShapeRequest)
		noteGroup.POST("/noteShape/update", m.noteShapeAPI.UpdateShapeRequest)

		// 笔记摘要相关API
		noteGroup.POST("/paperNote/summary/getByNoteId", m.noteSummaryAPI.GetNoteSummaryByNoteId)
		noteGroup.POST("/paperNote/summary/saveOrUpdate", m.noteSummaryAPI.SaveOrUpdateSummary)

		// 笔记单词相关API
		noteGroup.POST("/paperNote/word/getByNoteId", m.noteWordAPI.GetNoteWordsByNoteId)
		noteGroup.POST("/paperNote/word/config", m.noteWordAPI.SaveOrUpdateNoteWordConfig)
		noteGroup.POST("/paperNote/word/delete", m.noteWordAPI.DeleteNoteWord)
		noteGroup.POST("/paperNote/word/update", m.noteWordAPI.UpdateNoteWord)
		noteGroup.POST("/paperNote/word/save", m.noteWordAPI.SaveNoteWord)

		// 笔记PDF相关API
		noteGroup.GET("/paperNote/downloadNotePdf", m.paperNoteAPI.DownloadNotePdf)
		noteGroup.GET("/paperNote/downloadNoteMarkdown", m.paperNoteAPI.DownloadNoteMarkdown)

		// 笔记阅读位置相关API
		noteGroup.POST("/noteReadLocation/getLocation", m.noteReadLocationAPI.GetNoteReadLocation)
		// 记录笔记阅读位置
		noteGroup.POST("/noteReadLocation/record", m.noteReadLocationAPI.RecordNoteReadLocation)

		//---管理端接口---//
		// 笔记摘要相关API
		noteGroup.POST("/noteManage/summary/getList", m.noteManageAPI.GetUserSummaryList)
		noteGroup.POST("/noteManage/summary/getListByFolderId", m.noteManageAPI.GetUserSummaryListByFolderId)

		noteGroup.POST("/noteManage/word/getList", m.noteManageAPI.GetUserWordList)
		noteGroup.POST("/noteManage/word/getListByFolderId", m.noteManageAPI.GetWordListByFolderIdReq)
		noteGroup.POST("/noteManage/extract/getList", m.noteManageAPI.GetExtractListReq)
		noteGroup.POST("/noteManage/extract/getMarkTagListByFolderId", m.noteManageAPI.GetMarkTagListByFolderIdReq)
	}
}

// GetPaperNoteService 获取论文笔记服务
func (m *NoteModule) GetPaperNoteService() noteInterface.IPaperNoteService {
	return m.paperNoteService
}

// GetPaperNoteAccessService 获取论文笔记访问服务
func (m *NoteModule) GetPaperNoteAccessService() *service.PaperNoteAccessService {
	return m.paperNoteAccessService
}

// SetPaperPdfService 设置论文PDF服务，用于解决循环依赖问题
func (m *NoteModule) SetPaperPdfService(pdfService pdfInterface.IPaperPdfService) error {
	if pdfService == nil {
		return errors.New("pdfService cannot be nil")
	}
	m.pdfService = pdfService

	// 将笔记服务注入到PaperNoteAPI
	if m.paperNoteAPI != nil {
		if err := m.paperNoteAPI.SetPaperPdfService(pdfService); err != nil {
			return err
		}
	}

	// 将笔记服务注入到PaperNoteService
	if m.paperNoteService != nil {
		if err := m.paperNoteService.SetPaperPdfService(pdfService); err != nil {
			return err
		}
	}

	// 将笔记服务注入到noteSummaryService
	if m.noteSummaryService != nil {
		if err := m.noteSummaryService.SetPaperPdfService(pdfService); err != nil {
			return err
		}
	}

	if m.noteManageAPI != nil {
		if err := m.noteManageAPI.SetPaperPdfService(pdfService); err != nil {
			return err
		}
	}

	return nil
}
