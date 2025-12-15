package user

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/cache"
	"github.com/yb2020/odoc/pkg/scheduler"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"github.com/yb2020/odoc/config"
	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/eventbus"
	"github.com/yb2020/odoc/pkg/i18n"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/middleware"
	"github.com/yb2020/odoc/pkg/registry"
	"github.com/yb2020/odoc/services/user/api"
	"github.com/yb2020/odoc/services/user/dao"
	usergrpc "github.com/yb2020/odoc/services/user/grpc"
	"github.com/yb2020/odoc/services/user/service"
)

// 编译时类型检查：确保 UserModule 实现了 registry.Module 接口
var _ registry.Module = (*UserModule)(nil)

// Module 导出模块实例，用于自动发现和注册
var Module = &UserModule{}

// UserModule 实现用户服务模块
type UserModule struct {
	API                *api.UserAPI
	GRPCService        *usergrpc.UserGRPCServer
	UserService        *service.UserService
	authMiddleware     *middleware.AuthMiddleware
	db                 *gorm.DB
	config             *config.Config
	logger             logging.Logger
	tracer             opentracing.Tracer
	localizer          i18n.Localizer
	eventBus           *eventbus.EventBus
	transactionManager *baseDao.TransactionManager
}

// NewUserModule 创建一个新的用户模块实例
func NewUserModule(db *gorm.DB, config *config.Config, logger logging.Logger,
	tracer opentracing.Tracer, localizer i18n.Localizer,
	eventBus *eventbus.EventBus,
	transactionManager *baseDao.TransactionManager,
) *UserModule {
	return &UserModule{
		db:                 db,
		config:             config,
		logger:             logger,
		tracer:             tracer,
		localizer:          localizer,
		eventBus:           eventBus,
		transactionManager: transactionManager,
	}
}

// RegisterProviders 注册Provider
func (m *UserModule) RegisterProviders() {
	m.logger.Debug("msg", "用户模块没有Provider，跳过注册")
}

// Initialize 初始化用户模块
func (m *UserModule) Initialize() error {
	m.RegisterProviders()
	m.logger.Info("msg", "初始化用户模块")

	userDAO := dao.NewUserDAO(m.db, m.logger)
	cacheClient := cache.NewCache(m.logger, 30*time.Minute, m.Name())
	m.UserService = service.NewUserService(cacheClient, m.logger, m.tracer, m.localizer, userDAO, m.eventBus, m.transactionManager)

	m.API = api.NewUserAPI(m.logger, m.tracer, m.localizer, m.UserService)

	m.GRPCService = usergrpc.NewUserGRPCServer(m.logger, m.tracer, m.UserService)

	return nil
}

// Shutdown 关闭模块
func (m *UserModule) Shutdown() error {
	// 目前没有需要清理的资源
	return nil
}

// Name 返回模块名称
func (m *UserModule) Name() string {
	return "user"
}

// RegisterRoutes 注册用户服务路由
func (m *UserModule) RegisterRoutes(r *gin.Engine) {
	userPublicGroup := r.Group("/api/public/user")
	{
		userPublicGroup.GET("/email/exists", m.API.CheckEmailExists)
		userPublicGroup.POST("/register", m.API.Register)
	}

	userGroup := r.Group("/api/user")
	userGroup.Use(m.authMiddleware.AuthRequired())
	{
		userGroup.GET("/profile", m.API.GetProfile)
		userGroup.POST("/profile/update", m.API.UpdateProfile)
	}

	adminGroup := r.Group("/api/admin/user")
	adminGroup.Use(m.authMiddleware.AuthRequired())
	{
		adminGroup.GET("/getById", m.API.GetById)
		adminGroup.GET("/getByIds", m.API.GetByIds)
		adminGroup.GET("/pagination", m.API.PaginationUsers)
		adminGroup.POST("/create", m.API.CreateUser)
		adminGroup.POST("/update", m.API.UpdateUser)
		adminGroup.DELETE("/deleteById", m.API.DeleteUserById)
		adminGroup.DELETE("/deleteByIds", m.API.DeleteUserByIds)
	}
	m.logger.Info("msg", "用户服务路由注册成功", "prefix", "/api/users")
}

// RegisterGRPC 注册gRPC服务
func (m *UserModule) RegisterGRPC(server *grpc.Server) {
	if m.GRPCService == nil {
		// 如果gRPC服务未初始化，记录警告
		if m.logger != nil {
			m.logger.Warn("msg", "gRPC服务未初始化", "module", m.Name())
		}
		return
	}

	m.GRPCService.RegisterServer(server)
}

// RegisterJobSchedulers 注册Job定时任务
func (m *UserModule) RegisterJobSchedulers(scheduler *scheduler.Scheduler) {
	// TODO: 实现Job定时任务注册
	m.logger.Debug("msg", "用户模块没有Job定时任务，跳过注册")
}

// SetAuthMiddleware 设置认证中间件
func (m *UserModule) SetAuthMiddleware(authMiddleware *middleware.AuthMiddleware) {
	m.authMiddleware = authMiddleware
}

// GetUserService 返回用户服务实例
func (m *UserModule) GetUserService() *service.UserService {
	return m.UserService
}
