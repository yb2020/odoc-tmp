package event_tracker

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/middleware"
	"github.com/yb2020/odoc/pkg/registry"
	"github.com/yb2020/odoc/pkg/scheduler"
	"github.com/yb2020/odoc/services/event_tracker/api"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// 编译时类型检查：确保 EventTrackerModule 实现了 registry.Module 接口
var _ registry.Module = (*EventTrackerModule)(nil)

// Module 导出模块实例，用于自动发现和注册
var Module = &EventTrackerModule{}

// EventTrackerModule 事件追踪模块
type EventTrackerModule struct {
	db              *gorm.DB
	logger          logging.Logger
	tracer          opentracing.Tracer
	config          *config.Config
	authMiddleware  *middleware.AuthMiddleware
	eventTrackerAPI *api.EventTrackerAPI
}

// NewModule 创建事件追踪模块
func NewEventTrackerModule(db *gorm.DB,
	config *config.Config,
	logger logging.Logger,
	tracer opentracing.Tracer,
	authMiddleware *middleware.AuthMiddleware,
) *EventTrackerModule {
	return &EventTrackerModule{
		db:             db,
		logger:         logger,
		tracer:         tracer,
		config:         config,
		authMiddleware: authMiddleware,
	}
}

// Name 返回模块名称
func (m *EventTrackerModule) Name() string {
	return "event_tracker"
}

// Initialize 初始化模块
func (m *EventTrackerModule) Initialize() error {
	m.logger.Info("msg", "初始化事件追踪模块")
	m.eventTrackerAPI = api.NewEventTrackerAPI(m.logger,
		m.tracer,
	)
	return nil
}

// Shutdown 关闭模块
func (m *EventTrackerModule) Shutdown() error {
	m.logger.Info("msg", "关闭事件追踪模块")
	return nil
}

// RegisterGRPC 注册gRPC服务
func (m *EventTrackerModule) RegisterGRPC(server *grpc.Server) {
	// 事件追踪模块没有gRPC服务，不需要注册
	m.logger.Debug("msg", "事件追踪模块没有gRPC服务，跳过注册")
}

// RegisterJobSchedulers 注册Job定时任务
func (m *EventTrackerModule) RegisterJobSchedulers(scheduler *scheduler.Scheduler) {
	// TODO: 实现Job定时任务注册
	m.logger.Debug("msg", "事件追踪模块没有Job定时任务，跳过注册")
}

// RegisterProviders 注册Provider
func (m *EventTrackerModule) RegisterProviders() {
	m.RegisterProviders()
	// TODO: 实现Provider注册
	m.logger.Debug("msg", "事件追踪模块没有Provider，跳过注册")
}

// RegisterRoutes 注册路由
func (m *EventTrackerModule) RegisterRoutes(r *gin.Engine) {
	r.Use(
		m.authMiddleware.OptionalAuth(),
	).POST(
		"/report/collection_tracking0", m.eventTrackerAPI.TrackEvent,
	)
}
