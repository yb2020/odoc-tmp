package app

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"

	"github.com/yb2020/odoc/config"
	configAdapter "github.com/yb2020/odoc/config/adapter"

	"github.com/yb2020/odoc/internal/database"
	"github.com/yb2020/odoc/internal/i18n"
	"github.com/yb2020/odoc/pkg/errors"
	pkgi18n "github.com/yb2020/odoc/pkg/i18n"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/metrics"
	"github.com/yb2020/odoc/pkg/middleware"
	"github.com/yb2020/odoc/pkg/registry"
	"github.com/yb2020/odoc/pkg/scheduler"
	"github.com/yb2020/odoc/pkg/service"
	"github.com/yb2020/odoc/pkg/tracing"
	"github.com/yb2020/odoc/pkg/transport"
	"github.com/yb2020/odoc/pkg/utils"
	"github.com/yb2020/odoc/services"
)

// App 表示一个完整的应用程序实例
type App struct {
	Config          *config.Config
	Logger          logging.Logger
	Tracer          opentracing.Tracer
	TracerCloser    io.Closer
	MetricsHandler  *metrics.Metrics
	ErrorHandler    *errors.ErrorHandler
	Service         service.Service
	GinEngine       *gin.Engine
	Server          *http.Server
	SeaHandlers     []func(*gin.Engine, opentracing.Tracer, logging.Logger)
	ShutdownTimeout time.Duration
	DB              *gorm.DB             // 数据库连接
	RedisClient     database.RedisClient // Redis客户端
	Localizer       pkgi18n.Localizer
	Scheduler       *scheduler.Scheduler // 定时调度器
	ErrorReporter   middleware.ErrorReporter

	// gRPC服务器相关
	GRPCServer *grpc.Server // gRPC服务器
	GRPCAddr   string       // gRPC服务器地址
}

// Option 定义 App 的选项函数类型
type Option func(*App)

// WithConfig 设置自定义配置
func WithConfig(cfg *config.Config) Option {
	return func(a *App) {
		a.Config = cfg
	}
}

// WithConfigPath 从指定路径加载配置
func WithConfigPath(configPath string) Option {
	return func(a *App) {
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			fmt.Printf("Configuration file %s not found, using default configuration\n", configPath)
			a.Config = config.GetConfig()
		} else {
			// 使用Viper加载配置，支持环境变量覆盖
			cfg, err := configAdapter.LoadConfigWithViper(configPath)
			if err != nil {
				fmt.Printf("Error loading configuration: %v, using default configuration\n", err)
				a.Config = config.GetConfig()
			} else {
				a.Config = cfg
				// 更新全局配置，确保其他地方通过 config.GetConfig() 获取的是最新配置
				config.SetConfig(cfg)
				fmt.Printf("Configuration loaded from %s with environment variable support\n", configPath)
			}
		}
	}
}

// WithLogger 设置自定义日志记录器
func WithLogger(logger logging.Logger) Option {
	return func(a *App) {
		a.Logger = logger
	}
}

// WithTracer 设置自定义跟踪器
func WithTracer(tracer opentracing.Tracer, closer io.Closer) Option {
	return func(a *App) {
		a.Tracer = tracer
		a.TracerCloser = closer
	}
}

// WithMetrics 设置自定义指标处理器
func WithMetrics(metrics *metrics.Metrics) Option {
	return func(a *App) {
		a.MetricsHandler = metrics
	}
}

// WithService 设置自定义服务实现
func WithService(svc service.Service) Option {
	return func(a *App) {
		a.Service = svc
	}
}

// WithShutdownTimeout 设置关闭超时时间
func WithShutdownTimeout(timeout time.Duration) Option {
	return func(a *App) {
		a.ShutdownTimeout = timeout
	}
}

// WithGRPCAddr 设置gRPC服务器地址
func WithGRPCAddr(addr string) Option {
	return func(a *App) {
		a.GRPCAddr = addr
	}
}

// WithHTTPPort 设置 HTTP 服务器端口
func WithHTTPPort(port int) Option {
	return func(a *App) {
		// 修改配置中的 HTTP 端口
		a.Config.Server.Port = port
	}
}

// WithGRPCPort 设置 gRPC 服务器端口
func WithGRPCPort(port int) Option {
	return func(a *App) {
		// 修改配置中的 gRPC 端口
		a.Config.Server.GRPC.Port = port
		// 同时更新 GRPCAddr
		a.GRPCAddr = fmt.Sprintf("%s:%d", a.Config.Server.GRPC.Host, port)
	}
}

// NewApp 创建一个新的应用程序实例
func NewApp(options ...Option) (*App, error) {
	// 创建默认应用程序
	app := &App{
		Config:          config.GetConfig(),
		ShutdownTimeout: 5 * time.Second,
		GRPCAddr:        fmt.Sprintf("%s:%d", config.GetConfig().Server.GRPC.Host, config.GetConfig().Server.GRPC.Port), // 从配置文件读取GRPC地址
	}

	// 应用选项
	for _, option := range options {
		option(app)
	}

	// 初始化日志
	fmt.Println("Initializing logging reporter")

	// 检查是否在调试模式下
	_, isDebugMode := os.LookupEnv("GO_DEBUG")

	// 准备日志选项
	var logOptions []logging.LogOption

	// 如果配置了日志文件路径，则添加文件日志选项
	if app.Config.Logging.Path != "" {
		// 构建日志文件完整路径
		var logFilePath string

		projectRoot := utils.FindProjectRoot()
		logFilePath = filepath.Join(projectRoot, app.Config.Logging.Path, "app.log")

		// 确保日志目录存在
		logDir := filepath.Dir(logFilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Printf("创建日志目录失败: %v\n", err)
		}

		logOptions = append(logOptions, logging.WithLogFile(
			logFilePath,
			app.Config.Logging.MaxSize,
			app.Config.Logging.MaxAge,
			app.Config.Logging.MaxBackups,
		))
	}

	// 在调试模式下也使用配置文件中指定的日志格式，但强制使用 debug 级别
	if isDebugMode {
		app.Logger = logging.NewLogger("debug", app.Config.Logging.Format, logOptions...)
		app.Logger.Info("msg", "Running in debug mode, all logs will be shown")
	} else {
		// 使用配置文件中指定的日志级别和格式
		app.Logger = logging.NewLogger(app.Config.Logging.Level, app.Config.Logging.Format, logOptions...)
	}

	// 初始化指标（如果未提供）
	if app.MetricsHandler == nil {
		app.MetricsHandler = metrics.NewMetrics("go_sea")
	}

	// 初始化错误通知器
	var notifiers []errors.ErrorNotifier
	if app.Config.ErrorNotification.Enabled {
		switch app.Config.ErrorNotification.Type {
		case "slack":
			notifiers = append(notifiers, errors.NewSlackNotifier(
				app.Config.ErrorNotification.URL,
				"go-sea-service",
				app.Logger,
			))
		default:
			app.Logger.Warn("msg", "Unknown error notification type", "type", app.Config.ErrorNotification.Type)
		}
	}

	// 初始化错误处理器
	app.ErrorHandler = errors.NewErrorHandler(app.Logger, notifiers...)

	// 初始化跟踪器（如果未提供）
	if app.Tracer == nil {
		if app.Config.Tracing.Enabled {
			tracer, closer, err := tracing.InitTracer(
				app.Config.Tracing.ServiceName,
				app.Config.Tracing.JaegerURL,
				app.Config.Tracing.SampleRate,
				app.Logger,
			)
			if err != nil {
				app.Logger.Error("msg", "Failed to initialize tracer", "error", err.Error())
				app.Tracer = opentracing.NoopTracer{}
			} else {
				app.Tracer = tracer
				app.TracerCloser = closer
			}
		} else {
			app.Tracer = opentracing.NoopTracer{}
		}
	}

	// 初始化服务（如果未提供）
	if app.Service == nil {
		app.Service = service.NewService(app.Logger)
	}

	// // 初始化错误报告器
	// app.ErrorReporter = middleware.NewDefaultErrorReporter(app.Logger, app.Tracer)

	// // 创建中间件
	// var endpointMiddleware []kitendpoint.Middleware
	// endpointMiddleware = append(endpointMiddleware, middleware.LoggingMiddleware(app.Logger))
	// endpointMiddleware = append(endpointMiddleware, middleware.ErrorHandlingMiddleware(app.ErrorHandler, app.Logger))
	// endpointMiddleware = append(endpointMiddleware, middleware.MetricsMiddleware(app.MetricsHandler))
	// endpointMiddleware = append(endpointMiddleware, tracing.NewEndpointTracingMiddleware(app.Tracer, "endpoint", app.Logger))

	return app, nil
}

// Setup 设置应用程序
func (a *App) Setup() error {
	a.Logger.Info("msg", "Setting up application", "port", a.Config.Server.Port, "host", a.Config.Server.Host)

	// 初始化数据库连接（如果配置了数据库且尚未提供连接）
	if a.DB == nil && a.Config.Database.Enabled {
		a.Logger.Info("msg", "初始化数据库连接", "type", a.Config.Database.Type)

		db, err := database.NewDB(a.Config, a.Logger)
		if err != nil {
			a.Logger.Error("msg", "数据库连接失败", "error", err.Error())
			return fmt.Errorf("failed to connect to database: %w", err)
		}

		a.DB = db
		a.Logger.Info("msg", "数据库连接成功", "type", a.Config.Database.Type)
	}

	// 初始化 Redis 客户端（如果配置了 Redis 且尚未提供连接）
	if a.RedisClient == nil && a.Config.Redis.Enabled {
		a.Logger.Info("msg", "初始化 Redis 客户端")

		redisConfig := database.RedisConfig{
			Enabled:         a.Config.Redis.Enabled,
			Host:            a.Config.Redis.Host,
			Port:            a.Config.Redis.Port,
			Password:        a.Config.Redis.Password,
			DB:              a.Config.Redis.DB,
			PoolSize:        a.Config.Redis.PoolSize,
			MinIdleConns:    a.Config.Redis.MinIdleConns,
			DialTimeout:     time.Duration(a.Config.Redis.DialTimeout),
			ReadTimeout:     time.Duration(a.Config.Redis.ReadTimeout),
			WriteTimeout:    time.Duration(a.Config.Redis.WriteTimeout),
			MaxConnAge:      time.Duration(a.Config.Redis.MaxConnAge),
			MaxRetries:      a.Config.Redis.MaxRetries,
			MinRetryBackoff: time.Duration(a.Config.Redis.MinRetryBackoff),
			MaxRetryBackoff: time.Duration(a.Config.Redis.MaxRetryBackoff),
		}

		// 初始化全局Redis客户端
		if err := database.InitRedisClient(redisConfig, a.Logger); err != nil {
			a.Logger.Error("msg", "Redis 客户端初始化失败", "error", err.Error())
			// Redis 连接失败不应该导致整个应用程序启动失败
			a.Logger.Warn("msg", "Redis 连接失败，应用程序将继续运行但不使用 Redis")
		} else {
			// 获取初始化后的全局Redis客户端
			a.RedisClient = database.GetRedisClient()
			a.Logger.Info("msg", "Redis 客户端初始化成功", "host", a.Config.Redis.Host, "port", a.Config.Redis.Port)
		}
	}

	// 初始化 transport 包中的 logger
	transport.InitLogger(a.Logger)
	a.Logger.Info("msg", "初始化 transport 包中的 logger 成功")

	// 初始化国际化资源
	a.Logger.Info("msg", "初始化国际化资源")
	workDir, err := os.Getwd()
	if err != nil {
		a.Logger.Error("msg", "获取工作目录失败", "error", err.Error())
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// 初始化i18n系统
	defaultLang := pkgi18n.GetDefaultLanguage()
	supportedLangs := pkgi18n.GetSupportedLanguages()
	fallbackLang := pkgi18n.GetFallbackLanguage()

	_, err = i18n.InitI18n(workDir, defaultLang, supportedLangs, fallbackLang, a.Logger)
	if err != nil {
		a.Logger.Warn("msg", "初始化国际化资源失败", "error", err.Error())
		// 国际化初始化失败不应该导致整个应用程序启动失败
		a.Logger.Warn("msg", "国际化资源初始化失败，将使用消息ID作为默认消息")
	} else {
		a.Logger.Info("msg", "国际化资源初始化成功", "defaultLang", defaultLang)
	}

	a.Localizer = i18n.GetLocalizer()

	// 初始化定时调度器
	if a.Scheduler == nil {
		a.Logger.Info("msg", "初始化定时调度器")
		// Assuming a.RedisClient is a non-nil, valid *database.RedisClientWrapper
		a.Scheduler = scheduler.NewScheduler(a.Logger, a.RedisClient.(*database.RedisClientWrapper).Client)
		if err := a.Scheduler.Init(); err != nil {
			a.Logger.Error("msg", "定时调度器初始化失败", "error", err.Error())
			// 定时调度器初始化失败不应该导致整个应用程序启动失败
			a.Logger.Warn("msg", "定时调度器初始化失败，应用程序将继续运行但不使用定时调度器")
		} else {
			//a.Logger.Info("msg", "Adding scheduled jobs...")
			// Example: Add a simple job that runs every 1 minute
			// a.Scheduler.RegisterJobs(
			// 	scheduler.NewSimpleJob(
			// 		a.Logger,
			// 		"0 * * * * *",
			// 		"my-unique-job-key",
			// 		1*time.Minute,
			// 	),
			// )
			a.Scheduler.Start()
			a.Logger.Info("msg", "定时调度器初始化成功")
		}
	}

	// 初始化和注册所有模块
	a.Logger.Info("msg", "初始化和注册所有业务模块")

	// 初始化所有模块
	if err := services.InitializeModules(a.DB, a.RedisClient, a.Localizer, a.Config, a.Logger, a.Tracer); err != nil {
		a.Logger.Error("msg", "初始化模块失败", "error", err.Error())
		return fmt.Errorf("failed to initialize modules: %w", err)
	}

	// 注册所有模块到注册表
	for _, module := range services.GetAllModules() {
		registry.RegisterModule(module.Name(), module)
		a.Logger.Info("msg", "注册模块", "name", module.Name())
	}

	httpHandler := transport.NewHTTPHandler(a.Logger, a.Tracer, a.MetricsHandler, a.SeaHandlers...)

	// 获取底层的 gin.Engine
	if ginEngine, ok := httpHandler.(*gin.Engine); ok {
		a.GinEngine = ginEngine
	} else {
		return fmt.Errorf("failed to get gin.Engine instance")
	}

	// 注册所有模块的路由
	services.RegisterAllModuleRoutes(a.GinEngine)

	// 打印所有已注册的端点
	a.Logger.Info("msg", "已注册的端点列表")
	a.Logger.Info("endpoint", "HealthEndpoint")
	a.Logger.Info("endpoint", "EchoEndpoint")

	// 打印所有已注册的路由
	a.Logger.Info("msg", "路由注册完成")

	// 打印所有路由信息
	routes := a.GinEngine.Routes()
	a.Logger.Info("msg", fmt.Sprintf("已注册的路由数量: %d", len(routes)))

	// 按模块分组打印路由
	routesByPrefix := make(map[string][]gin.RouteInfo)

	for _, route := range routes {
		if route.Path == "" {
			continue
		}

		// 提取路由前缀（模块名称）
		prefix := "/"
		parts := strings.Split(route.Path, "/")
		if len(parts) > 1 && parts[1] != "" {
			prefix = "/" + parts[1]
		}

		// 按前缀分组
		routesByPrefix[prefix] = append(routesByPrefix[prefix], route)
	}

	// 按模块打印路由信息
	for prefix, routes := range routesByPrefix {
		a.Logger.Info("msg", fmt.Sprintf("模块 '%s' 的路由信息：", prefix))
		for _, route := range routes {
			a.Logger.Info("msg", "路由详情",
				"method", route.Method,
				"path", route.Path,
				"handler", fmt.Sprintf("%s", route.Handler))
		}
	}

	// 创建 HTTP 服务器
	httpAddr := fmt.Sprintf("%s:%d", a.Config.Server.Host, a.Config.Server.Port)
	a.Server = &http.Server{
		Addr:         httpAddr,
		Handler:      a.GinEngine,
		ReadTimeout:  time.Duration(a.Config.Server.Timeout) * time.Second,
		WriteTimeout: time.Duration(a.Config.Server.Timeout) * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 创建gRPC服务器
	// 从配置文件读取GRPC地址
	a.GRPCAddr = fmt.Sprintf("%s:%d", a.Config.Server.GRPC.Host, a.Config.Server.GRPC.Port)
	a.Logger.Info("msg", "配置GRPC服务器地址", "addr", a.GRPCAddr)

	// 创建额外的 gRPC 拦截器
	extraUnaryInterceptors := []grpc.UnaryServerInterceptor{
		grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(a.Tracer)),
	}

	// 使用我们的统一错误处理机制创建 gRPC 服务器选项
	serverOptions := middleware.CreateGRPCServerOptions(a.ErrorReporter, extraUnaryInterceptors, nil)

	// 创建 gRPC 服务器
	a.GRPCServer = grpc.NewServer(serverOptions...)

	// 注册所有支持gRPC的模块
	for _, grpcModule := range services.GetAllGRPCModules() {
		a.Logger.Info("msg", "注册gRPC模块", "module", reflect.TypeOf(grpcModule).String())
		grpcModule.RegisterGRPC(a.GRPCServer)
	}

	// 启用gRPC反射服务，便于客户端发现服务
	reflection.Register(a.GRPCServer)

	//注册所有JobScheduler
	services.RegisterAllModuleJobSchedulers(a.Scheduler)

	return nil
}

// Run 运行应用程序
func (a *App) Run() error {
	// 确保应用程序已设置
	if a.Server == nil {
		if err := a.Setup(); err != nil {
			return err
		}
	}

	// 创建一个通道用于接收信号
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// 创建一个通道用于接收服务器错误
	errorCh := make(chan error, 1)

	// 在后台启动 HTTP 服务器
	go func() {
		a.Logger.Info("msg", "Starting HTTP server", "addr", a.Server.Addr)
		if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Logger.Error("msg", "HTTP server error", "err", err)
			errorCh <- err
		}
	}()

	// 如果配置了gRPC服务器，则在后台启动
	if a.GRPCServer != nil {
		go func() {
			// 监听gRPC端口
			lis, err := net.Listen("tcp", a.GRPCAddr)
			if err != nil {
				a.Logger.Error("msg", "Failed to listen for gRPC", "addr", a.GRPCAddr, "err", err)
				errorCh <- err
				return
			}

			a.Logger.Info("msg", "Starting gRPC server", "addr", a.GRPCAddr)
			if err := a.GRPCServer.Serve(lis); err != nil {
				a.Logger.Error("msg", "gRPC server error", "err", err)
				errorCh <- err
			}
		}()
	}

	a.Logger.Info("msg", "Service started")

	// 等待信号或错误
	var err error
	select {
	case s := <-stop:
		a.Logger.Info("msg", "Received signal", "signal", s)
		err = fmt.Errorf("received signal: %s", s)
	case err = <-errorCh:
		a.Logger.Error("msg", "Service error", "err", err)
	}

	// 优雅关闭 HTTP 服务器
	ctx, cancel := context.WithTimeout(context.Background(), a.ShutdownTimeout)
	defer cancel()

	a.Logger.Info("msg", "Shutting down HTTP server")
	if shutdownErr := a.Server.Shutdown(ctx); shutdownErr != nil {
		a.Logger.Error("msg", "HTTP server shutdown error", "err", shutdownErr)
		if err == nil {
			err = shutdownErr
		}
	}

	// 如果有gRPC服务器，优雅关闭
	if a.GRPCServer != nil {
		a.Logger.Info("msg", "Shutting down gRPC server")
		// 优雅停止gRPC服务器
		a.GRPCServer.GracefulStop()
	}

	a.Logger.Info("msg", "Service stopped", "err", err)
	return err
}

// Close 关闭应用程序资源
func (a *App) Close() error {
	// 关闭所有已初始化的模块
	for _, module := range services.GetAllModules() {
		// 检查模块是否实现了Shutdown方法
		moduleName := module.Name()
		a.Logger.Info("msg", "关闭模块", "module", moduleName)
		if err := module.Shutdown(); err != nil {
			a.Logger.Error("msg", "关闭模块失败", "module", moduleName, "error", err.Error())
			// 继续关闭其他模块
		}
	}

	// 关闭调度器
	if a.Scheduler != nil {
		if err := a.Scheduler.Stop(); err != nil {
			a.Logger.Error("msg", "关闭调度器失败", "error", err.Error())
		}
	}

	// 关闭跟踪器
	if a.TracerCloser != nil {
		if err := a.TracerCloser.Close(); err != nil {
			a.Logger.Error("msg", "关闭跟踪器失败", "error", err.Error())
		}
	}

	// 关闭 Redis 客户端
	if a.RedisClient != nil {
		if err := a.RedisClient.Close(); err != nil {
			a.Logger.Error("msg", "Redis 客户端关闭失败", "error", err.Error())
		}
	}

	return nil
}
