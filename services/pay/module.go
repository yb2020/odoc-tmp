package pay

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
	msInterfaces "github.com/yb2020/odoc/services/membership/interfaces"
	"github.com/yb2020/odoc/services/pay/api"
	"github.com/yb2020/odoc/services/pay/dao"
	"github.com/yb2020/odoc/services/pay/provider"
	"github.com/yb2020/odoc/services/pay/service"
	userService "github.com/yb2020/odoc/services/user/service"
)

// 编译时类型检查：确保 PayModule 实现了 registry.Module 接口
var _ registry.Module = (*PayModule)(nil)

// Module 导出模块实例，用于自动发现和注册
var Module = &PayModule{}

// PayModule 实现支付服务模块
type PayModule struct {
	db                 *gorm.DB
	logger             logging.Logger
	tracer             opentracing.Tracer
	authMiddleware     *middleware.AuthMiddleware
	cfg                *config.Config
	eventBus           *eventbus.EventBus
	transactionManager *baseDao.TransactionManager
	userService        *userService.UserService
	orderService       msInterfaces.IOrderService

	paymentService *service.PaymentService
	paymentAPI     *api.PaymentAPI

	paymentSubscriptionService *service.PaymentSubscriptionService
	checkoutService            *service.StripeCheckoutService
	sripeCheckoutAPI           *api.StripeCheckoutAPI
}

// NewModule 创建会员模块
func NewPayModule(db *gorm.DB, cfg *config.Config, logger logging.Logger,
	tracer opentracing.Tracer, authMiddleware *middleware.AuthMiddleware,
	transactionManager *baseDao.TransactionManager,
	eventBus *eventbus.EventBus,
	userService *userService.UserService,
	orderService msInterfaces.IOrderService,
) *PayModule {
	return &PayModule{
		db:                 db,
		logger:             logger,
		tracer:             tracer,
		cfg:                cfg,
		authMiddleware:     authMiddleware,
		transactionManager: transactionManager,
		eventBus:           eventBus,
		userService:        userService,
		orderService:       orderService,
	}
}

// Name 返回模块名称
func (m *PayModule) Name() string {
	return "pay"
}

// Shutdown 停止模块
func (m *PayModule) Shutdown() error {
	m.logger.Info("msg", "关闭会员模块")
	return nil
}

// RegisterGRPC 注册gRPC服务
func (m *PayModule) RegisterGRPC(server *grpc.Server) {
	// TODO: 实现gRPC服务注册
	m.logger.Debug("msg", "会员模块没有gRPC服务，跳过注册")
}

// RegisterJobSchedulers 注册Job定时任务
func (m *PayModule) RegisterJobSchedulers(scheduler *scheduler.Scheduler) {
	// TODO: 实现Job定时任务注册
	m.logger.Debug("msg", "会员模块没有Job定时任务，跳过注册")
}

// RegisterProviders 注册Provider
func (m *PayModule) RegisterProviders() {
	m.logger.Debug("msg", "会员模块没有Provider，跳过注册")
}

// Initialize 初始化模块
func (m *PayModule) Initialize() error {
	m.RegisterProviders()
	m.logger.Info("msg", "初始化会员模块")

	paymentSubscriptionDAO := dao.NewPaymentSubscriptionDAO(m.db, m.logger)
	m.paymentSubscriptionService = service.NewPaymentSubscriptionService(paymentSubscriptionDAO, m.logger)

	paymentRecordDAO := dao.NewPaymentRecordDAO(m.db, m.logger)
	paymentProviderFactory := provider.NewPaymentProviderFactory(m.logger)
	if m.cfg.Pay.Stripe.IsEnable {
		paymentProviderFactory.RegisterStripeProvider(m.cfg.Pay.Stripe.SecretKey, m.cfg.Pay.Stripe.WebhookSecret)
	}

	m.paymentService = service.NewPaymentService(paymentRecordDAO, paymentProviderFactory, m.logger)

	m.paymentAPI = api.NewPaymentAPI(m.paymentService, m.logger)

	stripeCheckoutConfig := service.StripeCheckoutConfig{
		PublishableKey:     m.cfg.Pay.Stripe.PublishableKey,
		SecretKey:          m.cfg.Pay.Stripe.SecretKey,
		WebhookSecret:      m.cfg.Pay.Stripe.WebhookSecret,
		CheckoutSuccessURL: m.cfg.Pay.Stripe.CheckoutSuccessURL,
		CheckoutCancelURL:  m.cfg.Pay.Stripe.CheckoutCancelURL,
	}
	m.checkoutService = service.NewStripeCheckoutService(stripeCheckoutConfig, *paymentRecordDAO, m.logger, m.eventBus, m.paymentSubscriptionService)
	m.sripeCheckoutAPI = api.NewStripeCheckoutAPI(m.checkoutService, m.orderService, m.logger, m.tracer)

	return nil
}

// RegisterRoutes 注册路由
func (m *PayModule) RegisterRoutes(r *gin.Engine) {
	payGroup := r.Group("/api/pay")
	payGroup.Use(m.authMiddleware.AuthRequired())
	{
		// payGroup.POST("/prepay", m.paymentAPI.PrePay)
		// payGroup.POST("/stripe/webhook", m.paymentAPI.HandleStripeWebhook)
		// payGroup.GET("/:payment_id/status", m.paymentAPI.GetPaymentStatus)
		// payGroup.POST("/:payment_id/refund", m.paymentAPI.CreateRefund)
		// payGroup.GET("/user", m.paymentAPI.GetUserPayments)

		// stripe
		payGroup.GET("/stripe/get-publishable-key", m.sripeCheckoutAPI.GetStripePublishableKey)
		payGroup.POST("/stripe-checkout/create-session", m.sripeCheckoutAPI.CreateCheckoutSession)
		payGroup.POST("/stripe-checkout/cancel-subscription-at-period-end", m.sripeCheckoutAPI.CancelSubscriptionAtPeriodEnd)
		payGroup.POST("/stripe-checkout/cancel-subscription-immediately", m.sripeCheckoutAPI.CancelSubscriptionImmediately)
	}

	// 免验证回调 stripe checkout webhook
	r.POST("/services/pay/stripe-checkout/webhook", m.sripeCheckoutAPI.HandleCheckoutWebhook)
}
