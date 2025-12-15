package nav

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/scheduler"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"github.com/yb2020/odoc/config"
	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/eventbus"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/middleware"
	"github.com/yb2020/odoc/pkg/registry"
	"github.com/yb2020/odoc/services/nav/api"
	"github.com/yb2020/odoc/services/nav/dao"
	"github.com/yb2020/odoc/services/nav/service"
	userService "github.com/yb2020/odoc/services/user/service"
)

// 编译时类型检查：确保 WebsiteModule 实现了 registry.Module 接口
var _ registry.Module = (*NavModule)(nil)

// Module 导出模块实例，用于自动发现和注册
var Module = &NavModule{}

// NavModule 实现网站服务模块
type NavModule struct {
	db                 *gorm.DB
	logger             logging.Logger
	tracer             opentracing.Tracer
	authMiddleware     *middleware.AuthMiddleware
	cfg                *config.Config
	eventBus           *eventbus.EventBus
	transactionManager *baseDao.TransactionManager
	userService        *userService.UserService

	websiteService *service.WebsiteService
	websiteAPI     *api.WebsiteAPI
}

// NewNavModule 创建网站模块
func NewNavModule(db *gorm.DB, cfg *config.Config, logger logging.Logger,
	tracer opentracing.Tracer, authMiddleware *middleware.AuthMiddleware,
	transactionManager *baseDao.TransactionManager,
	eventBus *eventbus.EventBus,
	userService *userService.UserService,
) *NavModule {
	return &NavModule{
		db:                 db,
		logger:             logger,
		tracer:             tracer,
		cfg:                cfg,
		authMiddleware:     authMiddleware,
		transactionManager: transactionManager,
		eventBus:           eventBus,
		userService:        userService,
	}
}

// Name 返回模块名称
func (m *NavModule) Name() string {
	return "nav"
}

// Shutdown 停止模块
func (m *NavModule) Shutdown() error {
	m.logger.Info("msg", "关闭导航模块")
	return nil
}

// RegisterGRPC 注册gRPC服务
func (m *NavModule) RegisterGRPC(server *grpc.Server) {
	// TODO: 实现gRPC服务注册
	m.logger.Debug("msg", "导航模块没有gRPC服务，跳过注册")
}

// RegisterJobSchedulers 注册Job定时任务
func (m *NavModule) RegisterJobSchedulers(scheduler *scheduler.Scheduler) {
	// TODO: 实现Job定时任务注册
	m.logger.Debug("msg", "导航模块没有Job定时任务，跳过注册")
}

// RegisterProviders 注册Provider
func (m *NavModule) RegisterProviders() {
	// TODO: 实现Provider注册
	m.logger.Debug("msg", "导航模块没有Provider，跳过注册")
}

// Initialize 初始化模块
func (m *NavModule) Initialize() error {
	m.RegisterProviders()
	m.logger.Info("msg", "初始化导航模块")

	websiteDAO := dao.NewWebsiteDAO(m.db, m.logger)
	websiteService := service.NewWebsiteService(websiteDAO, m.logger, m.tracer, &m.cfg.Nav)
	m.websiteService = websiteService

	m.websiteAPI = api.NewWebsiteAPI(m.logger, m.tracer, m.websiteService)

	return nil
}

// RegisterRoutes 注册路由
func (m *NavModule) RegisterRoutes(r *gin.Engine) {
	websiteGroup := r.Group("/api/nav")
	websiteGroup.Use(m.authMiddleware.AuthRequired())
	{
		websiteGroup.POST("/website/create", m.websiteAPI.CreateWebsite)
		websiteGroup.POST("/website/delete", m.websiteAPI.DeleteWebsite)
		websiteGroup.GET("/website/getById", m.websiteAPI.GetWebsiteById)
		websiteGroup.POST("/website/update", m.websiteAPI.UpdateWebsite)
		websiteGroup.GET("/website/getList", m.websiteAPI.GetWebsiteList)
		websiteGroup.POST("/website/reorder", m.websiteAPI.ReorderWebsites)
	}
}
