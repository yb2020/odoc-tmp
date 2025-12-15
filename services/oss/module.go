package oss

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/middleware"
	pkgoss "github.com/yb2020/odoc/pkg/oss"
	"github.com/yb2020/odoc/pkg/registry"
	"github.com/yb2020/odoc/pkg/scheduler"
	"github.com/yb2020/odoc/services/oss/api"
	"github.com/yb2020/odoc/services/oss/dao"
	"github.com/yb2020/odoc/services/oss/service"
	paperService "github.com/yb2020/odoc/services/paper/service"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var _ registry.Module = (*OssModule)(nil)

// Module 导出模块实例
var Module = &OssModule{}

// OssModule OSS模块
type OssModule struct {
	db                    *gorm.DB
	config                *config.Config
	logger                logging.Logger
	tracer                opentracing.Tracer
	ossService            service.OssServiceInterface
	callbackService       *service.CallbackService
	ossDao                *dao.OssDAO
	storage               pkgoss.StorageInterface
	authMiddleware        *middleware.AuthMiddleware
	callbackAPI           *api.CallbackAPI
	paperPdfParsedService *paperService.PaperPdfParsedService
	ossServiceAPI         *api.OssServiceAPI
	ossAPI                *api.OssAPI
}

// NewOssModule 创建新的OSS模块实例
func NewOssModule(
	db *gorm.DB,
	config *config.Config,
	logger logging.Logger,
	tracer opentracing.Tracer,
	storage pkgoss.StorageInterface,
) *OssModule {
	return &OssModule{
		db:      db,
		config:  config,
		logger:  logger,
		tracer:  tracer,
		storage: storage,
	}
}

func (m *OssModule) RegisterProviders() {
	// 注册本地 OSS 服务实现（开源版本默认）
	registry.Register[service.OssServiceInterface](m.Name(), registry.ProviderLocal, func(deps any) service.OssServiceInterface {
		d := deps.(*service.OssDeps)
		return service.NewLocalOssService(d.OssDao, d.Config, d.Logger, d.Tracer, d.Storage)
	})
}

// Initialize 初始化模块
func (m *OssModule) Initialize() error {
	m.RegisterProviders()
	m.logger.Info("msg", "初始化OSS模块")

	// 创建DAO实例
	m.ossDao = dao.NewOssDAO(m.db, m.logger)
	// 创建Service实例（根据配置自动选择实现）
	m.ossService = registry.Create[service.OssServiceInterface](m.Name(), m.config.Service.Type, &service.OssDeps{
		OssDao:  m.ossDao,
		Config:  m.config,
		Logger:  m.logger,
		Tracer:  m.tracer,
		Storage: m.storage,
	})
	// 创建API实例
	m.ossServiceAPI = api.NewOSSServiceAPI(m.logger, m.tracer, m.config, m.ossService)
	m.callbackAPI = api.NewCallbackAPI(m.config, m.logger, m.tracer, m.callbackService)
	m.ossAPI = api.NewOSSAPI(m.logger, m.tracer, m.ossService)
	return nil
}

// SetAuthMiddleware 设置认证中间件
func (m *OssModule) SetAuthMiddleware(middleware *middleware.AuthMiddleware) {
	m.authMiddleware = middleware
}

// GetOssService 获取文件服务实例
func (m *OssModule) GetOssService() service.OssServiceInterface {
	return m.ossService
}

// Name 返回模块名称
func (m *OssModule) Name() string {
	return "oss"
}

// RegisterRoutes 注册路由
func (m *OssModule) RegisterRoutes(r *gin.Engine) {
	ossServiceGroup := r.Group("/services/oss/s3")
	ossServiceGroup.Use(m.authMiddleware.ServiceAuthRequired())
	// 服务之间内部翻译接口
	ossServiceGroup.GET("/getS3UploadToken", m.ossServiceAPI.GetS3UploadToken)
	ossServiceGroup.POST("/get/parsed/downloadTempUrl", m.ossServiceAPI.GetParseDownloadTempUrl)
	ossServiceGroup.POST("/get/pdf/downloadTempUrl", m.ossServiceAPI.GetPdfDownloadTempUrl)
	// ossServiceGroup.POST("/test/upload", m.ossServiceAPI.UploadToS3)

	// oss回调接口
	r.POST("/services/oss/minio/upload/callback", m.callbackAPI.HandleCallback)
	r.POST("/services/oss/s3/upload/callback", m.callbackAPI.HandleS3Callback)

	ossGroup := r.Group("/api/oss/s3")
	ossGroup.POST("/getDownloadTempUrl", m.ossAPI.GetDownloadTempUrl)
	ossGroup.POST("/getS3UploadToken", m.ossAPI.GetS3UploadToken)
}

// RegisterGRPC 注册gRPC服务
func (m *OssModule) RegisterGRPC(server *grpc.Server) {
	// TODO: 实现gRPC服务注册
}

// RegisterJobSchedulers 注册Job定时任务
func (m *OssModule) RegisterJobSchedulers(scheduler *scheduler.Scheduler) {
	// TODO: 实现Job定时任务注册
	m.logger.Debug("msg", "OSS模块没有Job定时任务，跳过注册")
}

// Shutdown 关闭模块
func (m *OssModule) Shutdown() error {
	m.logger.Info("msg", "关闭OSS业务模块")
	if m.storage != nil {
		return m.storage.Close()
	}
	return nil
}

// SetPaperPdfParsedService 设置PaperPdfParsed服务
func (m *OssModule) SetPaperPdfParsedService(paperPdfParsedService *paperService.PaperPdfParsedService) error {
	if paperPdfParsedService == nil {
		return errors.New("paperPdfParsedService cannot be nil")
	}
	m.paperPdfParsedService = paperPdfParsedService

	// 将paperPdfParsedService服务注入到ossService
	if m.ossService != nil {
		if err := m.ossService.SetPaperPdfParsedService(paperPdfParsedService); err != nil {
			return err
		}
	}
	return nil
}
