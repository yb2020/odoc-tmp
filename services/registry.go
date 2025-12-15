package services

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/internal/database"
	"github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/distlock"
	"github.com/yb2020/odoc/pkg/eventbus"
	"github.com/yb2020/odoc/pkg/http_client"
	"github.com/yb2020/odoc/pkg/i18n"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/middleware"
	"github.com/yb2020/odoc/pkg/oss"
	"github.com/yb2020/odoc/pkg/registry"
	"github.com/yb2020/odoc/pkg/scheduler"
	"github.com/yb2020/odoc/pkg/utils"
	"github.com/yb2020/odoc/services/doc"
	"github.com/yb2020/odoc/services/event_tracker"
	"github.com/yb2020/odoc/services/membership"
	"github.com/yb2020/odoc/services/nav"
	"github.com/yb2020/odoc/services/note"
	"github.com/yb2020/odoc/services/oauth2"
	"github.com/yb2020/odoc/services/oauth2/service"
	ossService "github.com/yb2020/odoc/services/oss"
	"github.com/yb2020/odoc/services/paper"
	"github.com/yb2020/odoc/services/parse"
	"github.com/yb2020/odoc/services/pay"
	"github.com/yb2020/odoc/services/pdf"
	"github.com/yb2020/odoc/services/translate"
	"github.com/yb2020/odoc/services/user"
	"gorm.io/gorm"
)

var (
	// 全局变量，用于存储已初始化的模块
	initializedModules     []registry.Module
	initializedGRPCModules []registry.GRPCModule
)

// InitializeModules 初始化所有模块
// 这个函数应该在应用程序启动时调用一次
func InitializeModules(db *gorm.DB, redis database.RedisClient, localizer i18n.Localizer, config *config.Config, logger logging.Logger, tracer opentracing.Tracer) error {
	// 清空已初始化的模块列表
	initializedModules = nil
	initializedGRPCModules = nil

	// 创建事件总线
	eventBus := eventbus.NewEventBus()

	// 初始化transactionManager
	transactionManager := dao.NewTransactionManager(db)

	// 初始化httpClient
	httpClient := http_client.NewHttpClient(logger)

	// 创建Redis分布式锁管理器
	redisLocker := distlock.NewRedisLocker(redis)
	lockTemplate := distlock.NewLockTemplate(redisLocker)

	// 初始化时参数排列约定：db, cache, config, logger, tracer, localizer, authMiddleware, service
	// 初始化用户模块
	userModule := user.NewUserModule(db, config, logger, tracer, localizer, eventBus, transactionManager)
	if err := userModule.Initialize(); err != nil {
		return err
	}
	initializedModules = append(initializedModules, userModule)

	// 初始化RSAUtil模块
	rsaUtil := utils.NewRSAUtil(config.OAuth2.RSA.PublicKey, config.OAuth2.RSA.PrivateKey)

	// 初始化OAuth2模块（依赖用户模块）
	oauth2Module := oauth2.NewOAuth2Module(db, redis, config, logger, tracer, localizer, userModule.GetUserService(), rsaUtil, eventBus)
	if err := oauth2Module.Initialize(); err != nil {
		return err
	}
	initializedModules = append(initializedModules, oauth2Module)

	// 创建认证适配器
	authAdapter := service.NewOAuth2AuthAdapter(oauth2Module.GetOAuth2Service())

	// 创建认证中间件
	authMiddleware := middleware.NewAuthMiddleware(*config, logger, localizer, authAdapter)

	// 更新已初始化模块的认证中间件
	userModule.SetAuthMiddleware(authMiddleware)
	oauth2Module.SetAuthMiddleware(authMiddleware)

	// 初始化OSS存储
	storage, err := oss.NewStorageInterface(config)
	if err != nil {
		return err
	}

	// 初始化OSS模块
	ossModule := ossService.NewOssModule(db, config, logger, tracer, storage)

	if err := ossModule.Initialize(); err != nil {
		return err
	}
	ossModule.SetAuthMiddleware(authMiddleware)
	initializedModules = append(initializedModules, ossModule)

	// 初始化会员模块
	membershipModule := membership.NewMembershipModule(db, config, logger, tracer, authMiddleware, transactionManager, eventBus, userModule.GetUserService())
	if err := membershipModule.Initialize(); err != nil {
		return err
	}
	initializedModules = append(initializedModules, membershipModule)

	// 初始化论文模块
	paperModule := paper.NewPaperModule(db, config, logger, tracer, authMiddleware, userModule.GetUserService())
	if err := paperModule.Initialize(); err != nil {
		return err
	}
	initializedModules = append(initializedModules, paperModule)

	// 初始化文档模块 - 不包含对笔记模块的依赖
	docModule := doc.NewDocModule(db, config, logger, tracer, authMiddleware, userModule.GetUserService(),
		ossModule.GetOssService(), paperModule.GetPaperJcrService(), paperModule.GetPaperService(),
		paperModule.GetPaperPdfParsedService(), membershipModule.GetMembershipService(),
		transactionManager)
	if err := docModule.Initialize(); err != nil {
		return err
	}
	initializedModules = append(initializedModules, docModule)

	// 初始化笔记模块
	noteModule := note.NewNoteModule(db, config, logger, tracer, authMiddleware, userModule.GetUserService(), docModule.GetUserDocService(), paperModule.GetPaperService())
	// 必须移到这里才能去初始化，不然设置的PaperPdfService将不会生效
	if err := noteModule.Initialize(); err != nil {
		return err
	}
	initializedModules = append(initializedModules, noteModule)

	// 初始化PDF模块
	pdfModule := pdf.NewPdfModule(db, config, logger, tracer, authMiddleware,
		userModule.GetUserService(), paperModule.GetPaperService(),
		docModule.GetUserDocService(), noteModule.GetPaperNoteService(),
		noteModule.GetPaperNoteAccessService(), ossModule.GetOssService(),
		paperModule.GetPaperAccessService(), paperModule.GetPaperPdfParsedService(),
	)
	if err := pdfModule.Initialize(); err != nil {
		return err
	}

	initializedModules = append(initializedModules, pdfModule)
	// 初始化翻译模块
	translateModule := translate.NewTranslateModule(db, config, logger, tracer, authMiddleware,
		httpClient,
		lockTemplate,
		noteModule.GetPaperNoteService(),
		pdfModule.GetPaperPdfService(),
		docModule.GetUserDocService(),
		ossModule.GetOssService(),
		membershipModule.GetMembershipService(),
	)
	if err := translateModule.Initialize(); err != nil {
		return err
	}
	initializedModules = append(initializedModules, translateModule)

	// 初始化支付模块
	payModule := pay.NewPayModule(db, config, logger, tracer, authMiddleware, transactionManager, eventBus, userModule.GetUserService(), membershipModule.GetOrderService())
	if err := payModule.Initialize(); err != nil {
		return err
	}
	initializedModules = append(initializedModules, payModule)

	// 初始化网站模块
	navModule := nav.NewNavModule(db, config, logger, tracer, authMiddleware, transactionManager, eventBus, userModule.GetUserService())
	if err := navModule.Initialize(); err != nil {
		return err
	}
	initializedModules = append(initializedModules, navModule)

	// 初始化解析模块
	parseModule := parse.NewParseModule(config, logger, tracer, httpClient, ossModule.GetOssService(), paperModule.GetPaperPdfParsedService())
	if err := parseModule.Initialize(); err != nil {
		return err
	}
	parseModule.SetAuthMiddleware(authMiddleware)
	initializedModules = append(initializedModules, parseModule)

	// 初始化事件追踪模块
	eventTrackerModule := event_tracker.NewEventTrackerModule(db, config, logger, tracer, authMiddleware)
	if err := eventTrackerModule.Initialize(); err != nil {
		return err
	}
	initializedModules = append(initializedModules, eventTrackerModule)

	//============================ 依赖注入 ============================
	// 解决循环依赖：为docModule注入noteModule的服务
	if err := docModule.SetNoteService(noteModule.GetPaperNoteService()); err != nil {
		logger.Error("msg", "设置笔记服务失败", "error", err.Error())
		return err
	}
	logger.Info("msg", "成功为文档模块设置笔记服务")
	if err := docModule.SetUserDocFolderRelationToUserDocFolder(); err != nil {
		logger.Error("msg", "设置文献文件夹和文献文件夹依赖注入失败", "error", err.Error())
		return err
	}
	if err := docModule.SetPaperPdfService(pdfModule.GetPaperPdfService()); err != nil {
		logger.Error("msg", "设置文献文件夹和文献文件夹依赖注入失败", "error", err.Error())
		return err
	}

	if err := noteModule.SetPaperPdfService(pdfModule.GetPaperPdfService()); err != nil {
		logger.Error("msg", "设置笔记服务失败", "error", err.Error())
		return err
	}
	logger.Info("msg", "成功为Note模块设置论文PDF服务")
	if err := ossModule.SetPaperPdfParsedService(paperModule.GetPaperPdfParsedService()); err != nil {
		logger.Error("msg", "设置笔记服务失败", "error", err.Error())
		return err
	}
	logger.Info("msg", "成功为Oss模块设置PaperPdfParsed服务")

	// 收集所有支持gRPC的模块
	for _, module := range initializedModules {
		if grpcModule, ok := module.(registry.GRPCModule); ok {
			initializedGRPCModules = append(initializedGRPCModules, grpcModule)
		}
	}

	return nil
}

// GetAllModules 返回所有已初始化的模块
// 这是一个集中管理所有业务模块的地方
func GetAllModules() []registry.Module {
	if len(initializedModules) == 0 {
		// 如果模块尚未初始化，则返回模块的声明
		// 注意：这些模块将需要在使用前进行初始化
		return []registry.Module{
			event_tracker.Module,
			user.Module,
			oauth2.Module,
			translate.Module,
			ossService.Module,
			note.Module,
			pdf.Module,
			paper.Module,
			doc.Module,
		}
	}
	return initializedModules
}

// GetAllGRPCModules 返回所有支持gRPC的已初始化模块
func GetAllGRPCModules() []registry.GRPCModule {
	if len(initializedGRPCModules) > 0 {
		return initializedGRPCModules
	}

	var modules []registry.GRPCModule

	// 遍历所有模块，检查是否实现了GRPCModule接口
	for _, module := range GetAllModules() {
		if grpcModule, ok := module.(registry.GRPCModule); ok {
			modules = append(modules, grpcModule)
		}
	}

	return modules
}

// RegisterAllModuleRoutes 注册所有模块的路由
func RegisterAllModuleRoutes(r *gin.Engine) {
	for _, module := range GetAllModules() {
		module.RegisterRoutes(r)
	}
}

// RegisterAllModuleJobSchedulers 注册所有JobScheduler的模块
func RegisterAllModuleJobSchedulers(scheduler *scheduler.Scheduler) {
	for _, module := range GetAllModules() {
		module.RegisterJobSchedulers(scheduler)
	}
}
