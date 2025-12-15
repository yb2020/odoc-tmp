package doc

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/cache"
	"github.com/yb2020/odoc/pkg/scheduler"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"github.com/yb2020/odoc/config"
	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/middleware"
	"github.com/yb2020/odoc/pkg/registry"
	"github.com/yb2020/odoc/services/doc/api"
	"github.com/yb2020/odoc/services/doc/dao"
	"github.com/yb2020/odoc/services/doc/factory"
	"github.com/yb2020/odoc/services/doc/service"
	membershipService "github.com/yb2020/odoc/services/membership/interfaces"
	noteInterface "github.com/yb2020/odoc/services/note/interfaces"
	ossService "github.com/yb2020/odoc/services/oss/service"
	paperService "github.com/yb2020/odoc/services/paper/service"
	pdfService "github.com/yb2020/odoc/services/pdf/service"
	userService "github.com/yb2020/odoc/services/user/service"
)

// 编译时类型检查：确保 DocModule 实现了 registry.Module 接口
var _ registry.Module = (*DocModule)(nil)

// Module 导出模块实例，用于自动发现和注册
var Module = &DocModule{}

// DocModule 文档模块
type DocModule struct {
	db                               *gorm.DB
	logger                           logging.Logger
	tracer                           opentracing.Tracer
	authMiddleware                   *middleware.AuthMiddleware
	cfg                              *config.Config
	userService                      *userService.UserService
	ossService                       ossService.OssServiceInterface
	paperJcrService                  *paperService.PaperJcrService
	paperService                     *paperService.PaperService
	paperPdfService                  *pdfService.PaperPdfService
	paperPdfParsedService            *paperService.PaperPdfParsedService
	membershipService                membershipService.IMembershipService
	noteService                      noteInterface.IPaperNoteService // 笔记服务，使用interface{}避免循环依赖
	userDocDAO                       *dao.UserDocDAO
	docClassifyRelationDAO           *dao.DocClassifyRelationDAO
	cslDAO                           *dao.CslDAO
	userCslRelationDAO               *dao.UserCslRelationDAO
	userDocAttachmentDAO             *dao.UserDocAttachmentDAO
	userDocClassifyDAO               *dao.UserDocClassifyDAO
	userDocFolderDAO                 *dao.UserDocFolderDAO
	userDocFolderRelationDAO         *dao.UserDocFolderRelationDAO
	doiMetaInfoDAO                   *dao.DoiMetaInfoDAO
	userDocService                   *service.UserDocService
	docClassifyRelationService       *service.DocClassifyRelationService
	cslService                       *service.CslService
	userCslRelationService           *service.UserCslRelationService
	userDocAttachmentService         *service.UserDocAttachmentService
	userDocClassifyService           *service.UserDocClassifyService
	userDocFolderService             *service.UserDocFolderService
	userDocFolderRelationService     *service.UserDocFolderRelationService
	docCiteSearchService             *service.DocCiteSearchService
	userDocAPI                       *api.UserDocAPI
	userDocFolderAPI                 *api.UserDocFolderAPI
	cslAPI                           *api.CslAPI
	docClassifyAPI                   *api.DocClassifyAPI
	transactionManager               *baseDao.TransactionManager
	docMetaInfoHandlerServiceFactory *factory.DocMetaInfoHandlerServiceFactory
	userDocUploadService             *service.UserDocUploadService
	userDocUploadLocalService        *service.UserDocUploadLocalService
	doiMetaInfoService               *service.DoiMetaInfoService
}

// NewDocModule 创建文档模块
func NewDocModule(db *gorm.DB,
	cfg *config.Config,
	logger logging.Logger,
	tracer opentracing.Tracer,
	authMiddleware *middleware.AuthMiddleware,
	userService *userService.UserService,
	ossService ossService.OssServiceInterface,
	paperJcrService *paperService.PaperJcrService,
	paperService *paperService.PaperService,
	paperPdfParsedService *paperService.PaperPdfParsedService,
	membershipService membershipService.IMembershipService,
	transactionManager *baseDao.TransactionManager,
) *DocModule {
	return &DocModule{
		db:                    db,
		logger:                logger,
		tracer:                tracer,
		cfg:                   cfg,
		authMiddleware:        authMiddleware,
		userService:           userService,
		ossService:            ossService,
		paperJcrService:       paperJcrService,
		paperService:          paperService,
		paperPdfParsedService: paperPdfParsedService,
		membershipService:     membershipService,
		transactionManager:    transactionManager,
	}
}

// Name 返回模块名称
func (m *DocModule) Name() string {
	return "doc"
}

// Shutdown 停止模块
func (m *DocModule) Shutdown() error {
	m.logger.Info("msg", "关闭文档模块")
	return nil
}

// RegisterGRPC 注册gRPC服务
func (m *DocModule) RegisterGRPC(server *grpc.Server) {
	// TODO: 实现gRPC服务注册
	m.logger.Debug("msg", "文档模块没有gRPC服务，跳过注册")
}

// RegisterJobSchedulers 注册Job定时任务
func (m *DocModule) RegisterJobSchedulers(scheduler *scheduler.Scheduler) {
	// TODO: 实现Job定时任务注册
	m.logger.Debug("msg", "文档模块没有Job定时任务，跳过注册")
}

// RegisterProviders 注册Provider
func (m *DocModule) RegisterProviders() {
	// TODO: 实现Provider注册
	m.logger.Debug("msg", "文档模块没有Provider，跳过注册")
}

// Initialize 初始化模块
func (m *DocModule) Initialize() error {
	m.RegisterProviders()
	m.logger.Info("msg", "初始化文档模块")
	cacheClient := cache.NewCache(m.logger, 30*time.Minute, m.Name())
	// 初始化DAO
	m.userDocDAO = dao.NewUserDocDAO(m.db, m.logger)
	m.userDocAttachmentDAO = dao.NewUserDocAttachmentDAO(m.db, m.logger)
	m.userDocClassifyDAO = dao.NewUserDocClassifyDAO(m.db, m.logger)
	m.userDocFolderDAO = dao.NewUserDocFolderDAO(m.db, m.logger)
	m.userDocFolderRelationDAO = dao.NewUserDocFolderRelationDAO(m.db, m.logger)
	m.docClassifyRelationDAO = dao.NewDocClassifyRelationDAO(m.db, m.logger)
	m.cslDAO = dao.NewCslDAO(m.db, m.logger)
	m.userCslRelationDAO = dao.NewUserCslRelationDAO(m.db, m.logger, m.tracer)
	m.doiMetaInfoDAO = dao.NewDoiMetaInfoDAO(m.db, m.logger)

	// 初始化DOI元信息服务
	m.doiMetaInfoService = service.NewDoiMetaInfoService(m.logger, m.tracer, m.doiMetaInfoDAO)

	// 初始化服务
	m.userDocAttachmentService = service.NewUserDocAttachmentService(m.logger, m.tracer, m.userDocAttachmentDAO)
	m.userDocFolderService = service.NewUserDocFolderService(m.logger, m.tracer, m.userDocFolderDAO, m.transactionManager)
	m.docClassifyRelationService = service.NewDocClassifyRelationService(m.logger, m.tracer, m.docClassifyRelationDAO)
	m.userDocClassifyService = service.NewUserDocClassifyService(m.logger, m.tracer, m.userDocClassifyDAO, m.docClassifyRelationService)
	m.userDocFolderRelationService = service.NewUserDocFolderRelationService(m.logger, m.tracer, m.userDocFolderRelationDAO, m.userDocFolderDAO)

	// 先创建 userDocService（不传 squidProxyService 和 difyDocStatusService，避免循环依赖，后续通过 setter 注入）
	m.userDocService = service.NewUserDocService(m.logger, m.tracer, cacheClient, m.userDocDAO, m.userDocAttachmentService,
		m.userDocClassifyService, m.userDocFolderService, m.docClassifyRelationService,
		m.userDocFolderRelationService, m.paperJcrService, m.paperService, m.ossService, m.cfg, m.membershipService,
		m.paperPdfParsedService)
	m.userCslRelationService = service.NewUserCslRelationService(m.userCslRelationDAO, m.logger, m.tracer)
	m.cslService = service.NewCslService(m.logger, m.tracer, m.cslDAO, m.userCslRelationService, m.userDocService)

	// 初始化文档元信息处理器工厂
	m.docMetaInfoHandlerServiceFactory = factory.NewDocMetaInfoHandlerServiceFactory(m.logger)
	doiHandler := service.NewDoiMetaInfoHandlerService(m.logger, m.tracer, cacheClient, m.doiMetaInfoService, m.cfg)
	m.docMetaInfoHandlerServiceFactory.Register(doiHandler)

	m.docCiteSearchService = service.NewDocCiteSearchService(m.logger, m.tracer, m.docMetaInfoHandlerServiceFactory)

	// 初始化上传服务
	m.userDocUploadService = service.NewUserDocUploadService(
		m.logger, m.tracer, cacheClient, m.cfg, m.userDocService,
		m.paperPdfService, m.paperPdfParsedService, m.ossService,
		m.membershipService,
	)

	// 初始化本地上传服务
	m.userDocUploadLocalService = service.NewUserDocUploadLocalService(
		m.logger, m.tracer, cacheClient, m.cfg, m.userDocService, m.ossService,
	)

	// 初始化API
	m.userDocAPI = api.NewUserDocAPI(m.userDocService, m.userDocUploadService, m.userDocUploadLocalService, m.ossService, m.docCiteSearchService, cacheClient,
		m.paperService, m.paperPdfParsedService, m.membershipService, m.logger, m.tracer, m.cfg)
	m.userDocFolderAPI = api.NewUserDocFolderAPI(m.userDocFolderService, m.userDocFolderRelationService, m.logger, m.tracer)
	m.cslAPI = api.NewCslAPI(m.cslService, m.userDocService, m.logger, m.tracer)
	m.docClassifyAPI = api.NewDocClassifyAPI(m.logger, m.tracer, m.userDocClassifyService, m.docClassifyRelationService)

	return nil
}

// RegisterRoutes 注册路由
func (m *DocModule) RegisterRoutes(r *gin.Engine) {
	docGroup := r.Group("/api")
	docGroup.Use(m.authMiddleware.AuthRequired())
	{
		// 用户文档API路由
		//获取上传状态
		docGroup.POST("/userDoc/GetUserDocCreateStatus", m.userDocAPI.GetUserUploadParseStatus)
		//获取上传信息
		docGroup.POST("/userDoc/GetUploadToken", m.userDocAPI.GetPdfUploadToken)

		//获取统一的上传token状态
		docGroup.POST("/userDoc/GetParseToken", m.userDocAPI.GetParseToken)

		//获取文档左侧列表
		docGroup.POST("/userDoc/getDocIndex", m.userDocAPI.GetDocIndex)
		//获取文档右侧列表
		docGroup.POST("/userDoc/getDocList", m.userDocAPI.GetDocList)
		//获取文档相关分类列表
		docGroup.POST("/userDoc/getDocRelatedClassifyList", m.userDocAPI.GetDocRelatedClassifyList)
		//获取文档相关作者列表
		docGroup.POST("/userDoc/getDocRelatedAuthorList", m.userDocAPI.GetDocRelatedAuthorList)
		//获取文档相关期刊列表
		docGroup.POST("/userDoc/getDocRelatedVenueList", m.userDocAPI.GetDocRelatedVenueList)
		//重命名文档
		docGroup.POST("/userDoc/renameUserDoc", m.userDocAPI.RenameDoc)
		//处理文件秒传逻辑
		docGroup.POST("/userDoc/HandleFileFastUpload", m.userDocAPI.HandleFileFastUpload)
		//复制文献或文件夹
		docGroup.POST("/client/doc/copyDocOrFolderToAnotherFolder", m.userDocFolderAPI.CopyDocOrFolderToAnotherFolder)
		//从文件夹中删除文献或文件夹
		docGroup.POST("/client/doc/removeDocFromFolder", m.userDocFolderAPI.RemoveDocFromFolder)
		//========================== 文献信息修改相关 =========================
		// 手动更新文档引用信息
		docGroup.POST("/userDoc/manualUpdateDocCiteInfo", m.userDocAPI.ManualUpdateDocCiteInfo)
		//修改作者信息
		docGroup.POST("/userDoc/updateAuthors", m.userDocAPI.UpdateUserDocAuthors)
		//查询论文作者信息
		docGroup.POST("/userDoc/getAuthors", m.userDocAPI.GetAuthors)
		//修改收录情况
		docGroup.POST("/userDoc/updateVenue", m.userDocAPI.UpdateUserDocVenue)
		//修改发布时间
		docGroup.POST("/userDoc/updatePublishDate", m.userDocAPI.UpdateUserDocPublishDate)
		//修改jcr分区
		docGroup.POST("/userDoc/update/jcr/partion", m.userDocAPI.UpdateUserDocJcrPartion)
		//修改备注
		docGroup.POST("/userDoc/updateDocRemark", m.userDocAPI.UpdateUserDocRemark)
		//修改影响因子
		docGroup.POST("/userDoc/update/impact/factor", m.userDocAPI.UpdateUserDocImpactFactor)
		//修改文献标签
		docGroup.POST("/userDoc/attachDocToClassify", m.userDocAPI.AttachDocToClassify)
		//删除文献标签
		docGroup.POST("/userDoc/removeDocFromClassify", m.userDocAPI.RemoveDocFromClassify)
		// 更新重要性评分
		docGroup.POST("/userDoc/importance/score", m.userDocAPI.UpdateUserDocImportanceScore)
		//========================== 文献信息搜索相关 =========================
		// 文档引用搜索API
		docGroup.POST("/userDoc/citeSearch/en", m.userDocAPI.EnDocCiteSearch)
		docGroup.POST("/userDoc/citeSearch/zh", m.userDocAPI.ZhDocCiteSearch)

		//========================== 文献信息操作相关 =========================
		// 获取文献信息
		docGroup.GET("/doc/userDoc", m.userDocAPI.GetUserDocById)
		// 获取文献状态信息
		docGroup.POST("/doc/userDocStatusByIds", m.userDocAPI.GetUserDocStatusByIds)

		// 更新文献阅读状态
		docGroup.POST("/doc/userDoc/updateReadStatus", m.userDocAPI.UpdateReadStatus)

		// JCR分区API路由
		docGroup.GET("/userDoc/jcr/partions", m.userDocAPI.GetJcrPartions)
		// 用户文档文件夹API路由
		docGroup.POST("/client/doc/addFolder", m.userDocFolderAPI.CreateUserDocFolder)
		//删除文件夹
		docGroup.POST("/client/doc/deleteFolder", m.userDocFolderAPI.DeleteUserDocFolder)
		//更新文件夹
		docGroup.POST("/client/doc/updateFolder", m.userDocFolderAPI.UpdateUserDocFolder)
		//移动文件夹或文献
		docGroup.POST("/client/doc/moveDocOrFolderToAnotherFolder", m.userDocFolderAPI.MoveDocOrFolderToAnotherFolder)
		//移动文献或文件夹
		docGroup.POST("/client/doc/moveFolderOrDoc", m.userDocFolderAPI.MoveFolderOrDoc)
		//删除文献
		docGroup.POST("/client/doc/deleteDoc", m.userDocAPI.DeleteDoc)

		//引用相关
		docGroup.POST("/docPublic/csl/getDefaultCslList", m.cslAPI.GetDefaultCslList)
		//获取文献信息
		docGroup.POST("/docPublic/docMetaInfo/getDocMetaInfo", m.cslAPI.GetDocMetaInfo)
		//获取用户自定义CSL列表
		docGroup.POST("/csl/myCslList", m.cslAPI.GetMyCslList)
		//获取文献类型列表
		docGroup.POST("/docPublic/docMetaInfo/getDocTypeList", m.cslAPI.GetDocTypeList)
		// BibTeX导出
		docGroup.POST("/doc/bibtex/exportByIds", m.cslAPI.ExportBibTexByIds)
		docGroup.POST("/doc/bibtex/exportByFolderId", m.cslAPI.ExportBibTexByFolderId)
		// 分类标签api
		//获取用户所有分类列表
		docGroup.POST("/userDoc/getUserAllClassifyList", m.docClassifyAPI.GetUserAllClassifyList)
		// 添加标签
		docGroup.POST("/userDoc/addClassify", m.docClassifyAPI.AddUserDocClassify)
		// 删除标签
		docGroup.POST("/userDoc/deleteClassify", m.docClassifyAPI.DeleteUserDocClassify)

		// 最近阅读文献
		docGroup.POST("/userDoc/getLatestReadDocList", m.userDocAPI.GetLatestReadDocList)
	}
}

// GetUserDocService 获取用户文档服务
func (m *DocModule) GetUserDocService() *service.UserDocService {
	return m.userDocService
}

// GetUserDocAttachmentService 获取用户文档附件服务
func (m *DocModule) GetUserDocAttachmentService() *service.UserDocAttachmentService {
	return m.userDocAttachmentService
}

// GetUserDocClassifyService 获取用户文档分类服务
func (m *DocModule) GetUserDocClassifyService() *service.UserDocClassifyService {
	return m.userDocClassifyService
}

// GetUserDocFolderService 获取用户文档文件夹服务
func (m *DocModule) GetUserDocFolderService() *service.UserDocFolderService {
	return m.userDocFolderService
}

// GetUserDocFolderRelationService 获取用户文档文件夹关系服务
func (m *DocModule) GetUserDocFolderRelationService() *service.UserDocFolderRelationService {
	return m.userDocFolderRelationService
}

// SetNoteService 设置笔记服务，用于解决循环依赖问题
func (m *DocModule) SetNoteService(noteService noteInterface.IPaperNoteService) error {
	if noteService == nil {
		return errors.New("noteService cannot be nil")
	}
	m.noteService = noteService

	// 将笔记服务注入到UserDocService
	if m.userDocService != nil {
		if err := m.userDocService.SetNoteService(noteService); err != nil {
			return err
		}
	}

	return nil
}

// SetUserDocFolderRelationToUserDocFolder 设置文件夹关系服务，用于解决循环依赖问题
func (m *DocModule) SetUserDocFolderRelationToUserDocFolder() error {
	if err := m.userDocFolderService.SetUserDocFolderRelationService(m.userDocFolderRelationService); err != nil {
		return err
	}
	return nil
}

// SetPaperPdfService 设置pdf服务，用于解决循环依赖问题
func (m *DocModule) SetPaperPdfService(paperPdfService *pdfService.PaperPdfService) error {
	if paperPdfService == nil {
		return errors.New("noteService cannot be nil")
	}
	m.paperPdfService = paperPdfService
	// 将笔记服务注入到UserDocService
	if m.userDocService != nil {
		if err := m.userDocService.SetPaperPdfService(paperPdfService); err != nil {
			return err
		}
	}
	if m.userDocAPI != nil {
		if err := m.userDocAPI.SetPaperPdfService(paperPdfService); err != nil {
			return err
		}
	}

	return nil
}
