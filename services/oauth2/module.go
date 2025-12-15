package oauth2

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/internal/database"
	"github.com/yb2020/odoc/pkg/eventbus"
	"github.com/yb2020/odoc/pkg/i18n"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/middleware"
	"github.com/yb2020/odoc/pkg/registry"
	"github.com/yb2020/odoc/pkg/scheduler"
	"github.com/yb2020/odoc/pkg/utils"
	"github.com/yb2020/odoc/services/oauth2/api"
	"github.com/yb2020/odoc/services/oauth2/dao"
	"github.com/yb2020/odoc/services/oauth2/service"
	userEvent "github.com/yb2020/odoc/services/user/event"
	userService "github.com/yb2020/odoc/services/user/service"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// 编译时类型检查：确保 OAuth2Module 实现了 registry.Module 接口
var _ registry.Module = (*OAuth2Module)(nil)

// Module 导出模块实例，用于自动发现和注册
var Module = &OAuth2Module{}

// OAuth2Module 实现OAuth2模块，实现registry.Module接口
type OAuth2Module struct {
	API                *api.OAuth2API
	OAuth2Service      service.OAuth2Service
	GoogleOAuthAPI     *api.GoogleOAuthAPI
	GoogleOAuthService *service.GoogleOAuthService
	db                 *gorm.DB
	config             *config.Config
	logger             logging.Logger
	tracer             opentracing.Tracer
	localizer          i18n.Localizer
	userService        *userService.UserService
	authMiddleware     *middleware.AuthMiddleware
	redis              database.RedisClient
	RSAUtil            *utils.RSAUtil
	eventBus           *eventbus.EventBus
}

// NewOAuth2Module 创建一个新的OAuth2模块实例
func NewOAuth2Module(db *gorm.DB, redis database.RedisClient, config *config.Config, logger logging.Logger,
	tracer opentracing.Tracer, localizer i18n.Localizer, userService *userService.UserService, RSAUtil *utils.RSAUtil, eventBus *eventbus.EventBus) *OAuth2Module {
	return &OAuth2Module{
		db:          db,
		config:      config,
		logger:      logger,
		tracer:      tracer,
		localizer:   localizer,
		userService: userService,
		redis:       redis,
		RSAUtil:     RSAUtil,
		eventBus:    eventBus,
	}
}

// Name 返回模块名称
func (m *OAuth2Module) Name() string {
	return "oauth2"
}

// RegisterProviders 注册Provider
func (m *OAuth2Module) RegisterProviders() {
	// TODO: 实现Provider注册
	m.logger.Debug("msg", "OAuth2模块没有Provider，跳过注册")
}

// Initialize 初始化模块
func (m *OAuth2Module) Initialize() error {
	m.RegisterProviders()
	m.logger.Info("msg", "初始化OAuth2模块")

	// 创建DAO
	tokenDAO := dao.NewTokenDAO(m.db, m.redis, m.logger, m.config)
	clientDAO := dao.NewOAuth2ClientsDAO(m.db, m.logger, m.config)
	authCodeDAO := dao.NewAuthCodeDAO(m.db, m.logger, m.config)

	// 创建OAuth2服务
	m.OAuth2Service = service.NewOAuth2Service(tokenDAO, m.userService, m.logger,
		m.tracer, m.config, m.localizer, clientDAO, authCodeDAO)

	// 订阅用户删除事件
	m.eventBus.Subscribe(userEvent.UserDeletedEvent, func(ctx context.Context, event eventbus.Event) {
		if userId, ok := event.Data.(string); ok {
			m.OAuth2Service.RevokeUserTokens(ctx, userId)
		}
	})

	// 创建Google OAuth2服务
	googleOAuthSvc, err := service.NewGoogleOAuthService(m.config, m.logger, *m.userService)
	if err != nil {
		m.logger.Error("Failed to initialize Google OAuth service", "error", err)
		return err
	}
	m.GoogleOAuthService = googleOAuthSvc

	// 创建API层
	m.API = api.NewOAuth2API(m.logger, m.tracer, m.localizer, m.config, m.OAuth2Service, m.RSAUtil)
	m.GoogleOAuthAPI = api.NewGoogleOAuthAPI(m.logger, m.GoogleOAuthService, m.OAuth2Service, m.config)

	return nil
}

func (m *OAuth2Module) Shutdown() error {
	m.logger.Info("msg", "关闭OAuth2模块")
	return nil
}

// RegisterGRPC 注册gRPC服务
func (m *OAuth2Module) RegisterGRPC(server *grpc.Server) {
	// OAuth2模块没有gRPC服务，不需要注册
	m.logger.Debug("msg", "OAuth2模块没有gRPC服务，跳过注册")
}

// RegisterJobSchedulers 注册Job定时任务
func (m *OAuth2Module) RegisterJobSchedulers(scheduler *scheduler.Scheduler) {
	// TODO: 实现Job定时任务注册
	m.logger.Debug("msg", "OAuth2模块没有Job定时任务，跳过注册")
}

// SetAuthMiddleware 设置认证中间件
func (m *OAuth2Module) SetAuthMiddleware(authMiddleware *middleware.AuthMiddleware) {
	m.authMiddleware = authMiddleware
}

// RegisterRoutes 注册路由
func (m *OAuth2Module) RegisterRoutes(r *gin.Engine) {
	// 注册路由
	apiGroup := r.Group("/api/oauth2")

	// 公开路由，不需要认证
	apiGroup.POST("/rsa_key", m.API.GetRSAKeyHandler)
	apiGroup.GET("/code", m.API.GetAuthCodeHandler)
	apiGroup.POST("/sign_in", m.API.SignInAuthCodeHandler)
	apiGroup.POST("/validate", m.API.ValidateHandler)

	// 公开路由，不需要认证 Google OAuth2认证
	apiGroup.GET("/google/login", m.GoogleOAuthAPI.LoginHandler)
	apiGroup.GET("/google/callback", m.GoogleOAuthAPI.CallbackHandler)

	// 需要认证的路由
	authRouter := apiGroup.Group("")
	authRouter.Use(m.authMiddleware.AuthRequired())
	{
		authRouter.POST("/refresh", m.API.RefreshHandler)
		authRouter.POST("/sign_out", m.API.SignOutHandler)
	}

	// 服务令牌路由
	serviceGroup := r.Group("/services/oauth2")
	serviceGroup.Use(m.authMiddleware.ServiceAuthRequired())
	serviceGroup.POST("/authorize/token", m.API.GenerateServiceTokenHandler)
	serviceGroup.GET("/revoke", m.API.RevokeServiceTokenHandler)
	serviceGroup.POST("/validate", m.API.ValidateServiceTokenHandler)

	m.logger.Info("msg", "OAuth2服务路由注册成功", "prefix", "/api/oauth2")
}

// GetOAuth2Service 返回OAuth2服务实例
func (m *OAuth2Module) GetOAuth2Service() service.OAuth2Service {
	return m.OAuth2Service
}
