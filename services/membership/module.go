package membership

import (
	"context"

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
	"github.com/yb2020/odoc/services/membership/api"
	"github.com/yb2020/odoc/services/membership/dao"
	"github.com/yb2020/odoc/services/membership/interfaces"
	"github.com/yb2020/odoc/services/membership/service"
	payevent "github.com/yb2020/odoc/services/pay/event"
	userEvent "github.com/yb2020/odoc/services/user/event"
	userService "github.com/yb2020/odoc/services/user/service"
)

// 编译时类型检查：确保 MembershipModule 实现了 registry.Module 接口
var _ registry.Module = (*MembershipModule)(nil)

// Module 导出模块实例，用于自动发现和注册
var Module = &MembershipModule{}

// MembershipModule 实现会员服务模块
type MembershipModule struct {
	db                 *gorm.DB
	logger             logging.Logger
	tracer             opentracing.Tracer
	authMiddleware     *middleware.AuthMiddleware
	cfg                *config.Config
	transactionManager *baseDao.TransactionManager
	eventBus           *eventbus.EventBus
	userService        *userService.UserService

	msConfigService       *service.ConfigService
	membershipService     *service.MembershipService
	userMembershipService interfaces.IUserMembershipService
	creditService         interfaces.ICreditService
	orderService          interfaces.IOrderService
	creditBillService     *service.CreditBillService
	creditPaymentService  interfaces.ICreditPaymentService

	membershipAPI *api.MembershipApi
	orderAPI      *api.OrderApi
}

// NewModule 创建会员模块
func NewMembershipModule(db *gorm.DB, cfg *config.Config, logger logging.Logger,
	tracer opentracing.Tracer, authMiddleware *middleware.AuthMiddleware,
	transactionManager *baseDao.TransactionManager,
	eventBus *eventbus.EventBus,
	userService *userService.UserService,
) *MembershipModule {
	return &MembershipModule{
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
func (m *MembershipModule) Name() string {
	return "membership"
}

// Shutdown 停止模块
func (m *MembershipModule) Shutdown() error {
	m.logger.Info("msg", "关闭会员模块")
	return nil
}

// RegisterGRPC 注册gRPC服务
func (m *MembershipModule) RegisterGRPC(server *grpc.Server) {
	// TODO: 实现gRPC服务注册
	m.logger.Debug("msg", "会员模块没有gRPC服务，跳过注册")
}

// RegisterJobSchedulers 注册Job定时任务
func (m *MembershipModule) RegisterJobSchedulers(scheduler *scheduler.Scheduler) {
	// // TODO: 实现Job定时任务注册
	// m.logger.Debug("msg", "会员模块注册Job定时任务")
	// msJob := job.NewMembershipExpiredJob(m.logger, m.cfg, m.membershipService, m.userService)
	// ccpayExpiredJob := job.NewCreditPayConfirmExpiredJob(m.logger, m.cfg, m.membershipService, m.creditPaymentService, m.userService)

	// scheduler.RegisterJobs(msJob, ccpayExpiredJob)
}

// RegisterProviders 注册Provider
func (m *MembershipModule) RegisterProviders() {
	// TODO: 实现Provider注册
	m.logger.Debug("msg", "会员模块没有Provider，跳过注册")
}

// Initialize 初始化模块
func (m *MembershipModule) Initialize() error {
	m.RegisterProviders()
	m.logger.Info("msg", "初始化会员模块")

	// 初始化会员配置服务
	m.msConfigService = service.NewConfigService(m.logger, m.tracer, m.cfg)

	// 初始化积分账单服务
	membershipCreditBillDAO := dao.NewCreditBillDAO(m.db, m.logger)
	m.creditBillService = service.NewCreditBillService(m.logger, m.tracer, membershipCreditBillDAO)

	// 初始化会员支付服务
	creditPaymentRecordDAO := dao.NewCreditPaymentRecordDAO(m.db, m.logger)
	m.creditPaymentService = service.NewCreditPaymentService(m.logger, m.tracer, creditPaymentRecordDAO)

	// 初始化会员积分服务
	membershipCreditDAO := dao.NewCreditDAO(m.db, m.logger)
	m.creditService = service.NewCreditService(m.logger, m.tracer, membershipCreditDAO, m.creditBillService, m.creditPaymentService)

	// 初始化用户会员服务
	userMembershipDAO := dao.NewUserMembershipDAO(m.db, m.logger)
	m.userMembershipService = service.NewUserMembershipService(m.logger, m.tracer, m.userService, userMembershipDAO, m.msConfigService, m.creditService, m.transactionManager)

	// 初始化会员订单服务
	membershipSubOrderDAO := dao.NewOrderDAO(m.db, m.logger)
	m.orderService = service.NewOrderService(m.logger, m.tracer, membershipSubOrderDAO, m.msConfigService, m.creditService, m.userMembershipService)

	// 初始化会员服务
	m.membershipService = service.NewMembershipService(m.logger, m.tracer, m.msConfigService, m.userMembershipService, m.orderService, m.creditService, m.creditPaymentService)
	m.membershipAPI = api.NewMembershipApi(m.logger, m.tracer, m.cfg, m.membershipService, m.msConfigService, m.userService)

	m.orderAPI = api.NewOrderApi(m.logger, m.tracer, m.orderService, m.membershipService)

	// 订阅pay模块支付成功事件
	m.eventBus.Subscribe(payevent.PayNotifyEvent_PaySuccess, func(ctx context.Context, event eventbus.Event) {
		m.logger.Info("msg", "收到支付成功事件", "event", event)
		m.orderService.HandlePayNotifyEventHandler(ctx, event)
	})
	m.eventBus.Subscribe(payevent.PayNotifyEvent_PayFailed, func(ctx context.Context, event eventbus.Event) {
		m.logger.Info("msg", "收到支付失败事件", "event", event)
		m.orderService.HandlePayNotifyEventHandler(ctx, event)
	})
	m.eventBus.Subscribe(payevent.PayNotifyEvent_PayExpire, func(ctx context.Context, event eventbus.Event) {
		m.logger.Info("msg", "收到支付过期事件", "event", event)
		m.orderService.HandlePayNotifyEventHandler(ctx, event)
	})
	m.eventBus.Subscribe(payevent.PayNotifyEvent_CustomerSubscriptionsCreated, func(ctx context.Context, event eventbus.Event) {
		m.logger.Info("msg", "收到订阅创建事件", "event", event)
		m.orderService.HandlePayNotifyEventHandler(ctx, event)
	})
	m.eventBus.Subscribe(payevent.PayNotifyEvent_CustomerSubscriptionsUpdated, func(ctx context.Context, event eventbus.Event) {
		m.logger.Info("msg", "收到订阅更新事件", "event", event)
		m.orderService.HandlePayNotifyEventHandler(ctx, event)
	})
	m.eventBus.Subscribe(payevent.PayNotifyEvent_CustomerSubscriptionsDeleted, func(ctx context.Context, event eventbus.Event) {
		m.logger.Info("msg", "收到订阅删除事件", "event", event)
		m.orderService.HandlePayNotifyEventHandler(ctx, event)
	})
	m.eventBus.Subscribe(payevent.PayNotifyEvent_InvoicePaymentSucceeded, func(ctx context.Context, event eventbus.Event) {
		m.logger.Info("msg", "收到发票支付成功事件", "event", event)
		m.orderService.HandlePayNotifyEventHandler(ctx, event)
	})
	m.eventBus.Subscribe(payevent.PayNotifyEvent_InvoicePaymentFailed, func(ctx context.Context, event eventbus.Event) {
		m.logger.Info("msg", "收到发票支付失败事件", "event", event)
		m.orderService.HandlePayNotifyEventHandler(ctx, event)
	})

	// 订阅用户注册事件
	m.eventBus.Subscribe(userEvent.UserRegisterEvent, func(ctx context.Context, event eventbus.Event) {
		m.logger.Info("msg", "收到用户注册事件", "event", event)
		m.membershipService.HandleUserNotifyEventHandler(ctx, userEvent.UserRegisterEvent, event.Data.(string))
	})
	// 订阅用户删除事件
	m.eventBus.Subscribe(userEvent.UserDeletedEvent, func(ctx context.Context, event eventbus.Event) {
		m.logger.Info("msg", "收到用户删除事件", "event", event)
		m.membershipService.HandleUserNotifyEventHandler(ctx, userEvent.UserDeletedEvent, event.Data.(string))
	})

	return nil
}

// RegisterRoutes 注册路由
func (m *MembershipModule) RegisterRoutes(r *gin.Engine) {
	MembershipGroup := r.Group("/api/membership")
	MembershipGroup.Use(m.authMiddleware.AuthRequired())
	{
		MembershipGroup.POST("/test", m.membershipAPI.Test)
		MembershipGroup.POST("/testSubFree", m.membershipAPI.TestSubFree)
		MembershipGroup.POST("/testCallCreditFunsDocsUpload", m.membershipAPI.TestCallCreditFunsDocsUpload)
		MembershipGroup.POST("/testCallCreditFunsAi", m.membershipAPI.TestCallCreditFunsAi)
		MembershipGroup.POST("/testCallCreditFunsTranslate", m.membershipAPI.TestCallCreditFunsTranslate)
		MembershipGroup.POST("/testCallCreditFunsNote", m.membershipAPI.TestCallCreditFunsNote)

		MembershipGroup.GET("/get-base-info", m.membershipAPI.GetBaseInfo)
		MembershipGroup.GET("/get-info", m.membershipAPI.GetInfo)
		MembershipGroup.GET("/user/profile", m.membershipAPI.GetMembershipAndUserInfo)
		MembershipGroup.GET("/get-user-credit", m.membershipAPI.GetUserCredit)
	}

	MembershipPublicGroup := r.Group("/api/public/membership")
	{
		MembershipPublicGroup.GET("/get-subplan-infos", m.membershipAPI.GetSubPlanInfos)
		MembershipPublicGroup.GET("/get-login-page-config", m.membershipAPI.GetLoginPageConfig)
	}

	// 订单服务组
	OrderGroup := r.Group("/api/order")
	OrderGroup.Use(m.authMiddleware.AuthRequired())
	{
		OrderGroup.GET("/get-order-info", m.orderAPI.GetOrderInfo)
		OrderGroup.POST("/create-order", m.orderAPI.CreateOrder)
	}
}

// GetOrderService 获取订单服务
func (m *MembershipModule) GetOrderService() interfaces.IOrderService {
	return m.orderService
}

// GetMembershipService 获取会员服务
func (m *MembershipModule) GetMembershipService() interfaces.IMembershipService {
	return m.membershipService
}
