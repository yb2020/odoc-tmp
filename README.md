# Go-Sea

基于 Go 1.24、Gin 和 go-kit 构建的微服务脚手架，提供完整的企业级功能支持。Go-Sea 可以作为独立服务运行，也可以作为依赖库集成到其他项目中，提供基础设施能力。

## 功能特性

- **上下文支持**：完整的 context 传递，支持请求跟踪和取消，安全的 goroutine 上下文传递
- **调试能力**：灵活的日志级别和格式配置
- **容错处理**：强大的错误处理机制和自定义错误类型
- **全链路跟踪**：基于 Jaeger 的分布式追踪系统集成
- **错误推送**：支持多种错误通知渠道（Slack, Email等）
- **指标监控**：集成 Prometheus 监控系统
- **业务错误处理**：统一的业务错误处理系统，支持状态码和国际化

## 项目结构     

```
.
├── cmd/
│   └── server/         # 服务器入口
├── config/            # 配置文件和配置加载
├── internal/          # 内部共享代码
│   ├── database/      # 数据库相关
│   ├── cache/         # 缓存相关
│   ├── i18n/          # 国际化资源
│   ├── biz/           # 业务状态码
│   └── response/      # 响应处理
├── pkg/
│   ├── endpoint/      # 服务端点定义
│   ├── errors/        # 错误处理和通知
│   ├── logging/       # 日志处理
│   ├── metrics/       # 指标收集
│   ├── middleware/    # 中间件（日志、跟踪、错误处理、认证）
│   ├── context/       # 上下文管理
│   ├── service/       # 业务逻辑
│   ├── tracing/       # 分布式追踪
│   ├── i18n/          # 国际化接口
│   └── transport/     # HTTP传输层
├── services/          # 各个服务模块
│   ├── user/          # 用户服务
│   │   ├── api/       # HTTP API
│   │   ├── grpc/      # gRPC服务
│   │   ├── service/   # 业务逻辑
│   │   ├── model/     # 数据模型
│   │   ├── dao/       # 数据访问
│   │   └── module.go  # 模块注册
│   └── oauth2/        # 认证服务
└── go.mod             # Go模块定义
```



原文解析文件目录，桶说明：
├── bucket/	  # 桶
│   └── {filesha256}/                    # 文件sha256
│	│	├── origin/                   # 上传的原始文件
│	│	├── parsed/                   # 解析后的结果
│	│	├── full_text_translated/ 	  # 全文翻译后的结果
│	│	├── user/                     # 用户
│	│	│   ├── {draw}/               # 操作场景
│	│	│	│   ├── {userid}/		  # 用户id

## 快速开始

### 前提条件

- Go 1.24 或更高版本
- PostgreSQL 13 或更高版本（生产环境）
- SQLite（本地开发环境）
- Redis 6.2 或更高版本
- 可选：Jaeger 服务器（用于分布式追踪）

### 本地开发环境（SQLite）

本地开发时可以使用 SQLite 作为数据库，无需安装 PostgreSQL。

#### 1. 生成 SQLite 数据库表

```bash
# 创建所有表（自动确认）
go run tools/db_gen.go -config config/config.local.develop.yaml -all -y

# 或者只创建特定模型的表
go run tools/db_gen.go -config config/config.local.develop.yaml -model User -y

# 列出所有可用模型
go run tools/db_gen.go -config config/config.local.develop.yaml -list
```

#### 2. 参数说明

| 参数 | 说明 |
|------|------|
| `-config` | 配置文件路径 |
| `-all` | 创建/更新所有表 |
| `-y` | 自动确认所有操作 |
| `-model` | 指定单个模型名称 |
| `-list` | 仅列出可用模型 |
| `-gen` | 同时生成模型代码 |

#### 3. 数据库文件位置

SQLite 数据库文件默认生成在 `./data/app.db`，可在 `config/config.local.develop.yaml` 中修改：

```yaml
database:
  type: sqlite
  sqlite:
    dbPath: ./data/app.db
```

### 作为独立服务运行

```bash
# 从项目根目录运行
go run cmd/server/main.go
```

默认情况下，服务将在 `http://0.0.0.0:8080` 上启动。

### 作为依赖库集成到其他项目

#### 1. 添加 Go-Sea 作为依赖

```bash
go get github.com/yb2020/go-sea
```

#### 2. 在你的项目中使用 Go-Sea

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	"github.com/yb2020/go-sea/config"
	"github.com/yb2020/go-sea/pkg/logging"
	"github.com/yb2020/go-sea/pkg/metrics"
	"github.com/yb2020/go-sea/pkg/middleware"
	"github.com/yb2020/go-sea/pkg/tracing"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	logger := logging.NewLogger(cfg.Logging.Level, cfg.Logging.Format)

	// 初始化指标收集
	metricsHandler := metrics.NewMetrics("my_service")

	// 初始化追踪
	tracer, tracerCloser, err := tracing.NewJaegerTracer(
		"my_service",
		cfg.Tracing.JaegerURL,
		cfg.Tracing.SampleRate,
		logger,
	)
	if err != nil {
		logger.Error("msg", "Failed to initialize tracer", "error", err.Error())
		tracer = opentracing.NoopTracer{}
	} else {
		defer tracerCloser.Close()
	}

	// 创建Gin引擎
	r := gin.Default()

	// 添加中间件
	r.Use(middleware.GinTracingMiddleware(tracer))
	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(middleware.ErrorHandler())

	// 添加自定义路由
	r.GET("/my-endpoint", func(c *gin.Context) {
		// 你的处理逻辑
	})

	// 启动服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: r,
	}

	// 优雅关闭
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("msg", "Starting server", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("msg", "Server error", "error", err.Error())
			os.Exit(1)
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Info("msg", "Shutting down server")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("msg", "Server shutdown error", "error", err.Error())
	}
}

### 配置

可以通过修改配置文件来配置服务，使用YAML格式：

```bash
# 使用YAML配置
go run cmd/server/main.go -config=/path/to/custom-config.yaml
```

项目提供了YAML格式的配置文件示例：
- `config/config.yaml` - YAML格式配置（易读，支持注释）

## API 端点

### 用户服务 (/api/user)
- `GET /api/user/profile` - 获取用户个人资料
- `GET /api/user/exists` - 检查邮箱是否存在
- `POST /api/user/register` - 用户注册
- `POST /api/user/profile/update` - 更新用户个人资料
- `GET /api/user/getById` - 根据ID获取用户信息

### 用户管理 (/api/admin/user)
- `GET /api/admin/user/getById` - 根据ID获取用户信息
- `GET /api/admin/user/getByIds` - 根据多个ID获取用户信息
- `GET /api/admin/user/pagination` - 分页获取用户列表
- `POST /api/admin/user/create` - 创建用户
- `POST /api/admin/user/update` - 更新用户信息
- `POST /api/admin/user/delete` - 删除用户
- `POST /api/admin/user/batchDelete` - 批量删除用户

### 认证服务 (/api/oauth2)
- `POST /api/oauth2/token` - 获取访问令牌
- `POST /api/oauth2/refresh` - 刷新访问令牌
- `POST /api/oauth2/revoke` - 撤销访问令牌
- `POST /api/oauth2/validate` - 验证访问令牌
- `POST /api/oauth2/sign_in` - 用户登录
- `POST /api/oauth2/sign_out` - 用户登出

### 系统服务
- `GET /health-check` - 健康检查
- `GET /metrics` - Prometheus 指标

## 错误处理系统

Go-Sea 提供了强大的错误处理系统，包括标准错误和业务错误两种类型。

### 标准错误 (pkg/errors/errors.go)

标准错误系统提供了基本的错误处理功能：

- 错误类型定义
- 错误堆栈跟踪
- 错误包装和扩展

### 业务错误 (pkg/errors/biz_errors.go)

业务错误系统是标准错误系统的扩展，专门用于处理业务逻辑相关的错误：

```go
// BizError 业务错误
type BizError struct {
    Status int32  // 状态码
    MsgID  string // 消息ID，用于国际化
    Err    error  // 原始错误
}
```

提供了以下功能：

- `Biz` - 创建新的业务错误
- `BizWrap` - 将现有错误包装为业务错误
- `BizWithStatus` - 创建带有状态码的业务错误

### 业务错误中间件 (pkg/middleware/biz_errors_handler.go)

业务错误中间件能够自动处理业务错误，并将其转换为适当的HTTP响应：

- 捕获API处理函数中的业务错误
- 根据错误类型设置适当的HTTP状态码
- 使用国际化系统翻译错误消息
- 生成统一的错误响应格式

### 使用指南

#### 1. 在服务层创建业务错误

```go
// 创建新的业务错误
if user == nil {
    return nil, errors.Biz("user.error.user_not_found")
}

// 包装现有错误
if err != nil {
    return nil, errors.BizWrap("user.error.create_failed", err)
}

// 使用自定义状态码
if err != nil {
    return nil, errors.BizWithStatus(biz.User_StatusInvalidUserData, "user.error.invalid_data", err)
}
```

#### 2. 在API层处理错误

```go
func (api *UserAPI) GetProfile(c *gin.Context) {
    // 处理请求...
    resp, err := api.userService.GetProfile(ctx, req)
    if err != nil {
        // 只需记录错误并添加到上下文
        api.logger.Error("msg", "获取用户个人资料失败", "error", err.Error())
        c.Error(err)
        return
    }
    
    // 成功响应
    response.Success(c, "success", resp)
}
```

#### 3. 配置错误处理中间件

在HTTP服务器初始化时添加错误处理中间件：

```go
// 在其他中间件之后添加
r.Use(middleware.ErrorHandler())
```

## 上下文系统

### Gin 和 Go 标准库上下文集成

Go-Sea 实现了 Gin 和 Go 标准库上下文的无缝集成，解决了在 goroutine 中安全使用上下文的问题。这个系统主要包括以下组件：

#### 用户上下文 (pkg/context/user_context.go)

```go
// UserContext 用户上下文，包含用户相关的信息
type UserContext struct {
	UserID        int64
	AccessToken   string
	Roles         []string
	Device        string
	Username      string
	Authenticated bool
	Claims        any
	// 额外数据存储
	extraData     map[string]interface{}
}
```

提供了以下功能：

- `NewUserContext` - 创建新的用户上下文
- `GetUserContext` - 从标准库上下文中获取用户上下文
- `FromGinContext` - 从Gin上下文中提取用户上下文
- `ToGinContext` - 将用户上下文存储到Gin上下文中
- `ToContext` - 将用户上下文转换为标准库的context.Context
- `GetUserID`, `GetRoles`, `HasRole` - 辅助函数获取用户信息
- `IsAuthenticated` - 检查用户是否已认证

#### 安全的 Goroutine (pkg/context/goroutine.go)

提供了在 goroutine 中安全使用上下文的功能：

- `SafeGoroutine` - 在 goroutine 中安全使用标准库上下文
- `GinSafeGoroutine` - 在 goroutine 中安全使用 Gin 上下文
- `WithTimeout` - 创建带超时的上下文
- `WithCancel` - 创建可取消的上下文
- `BackgroundContext` - 创建带用户上下文的后台上下文
- `RunWithUserContext` - 使用用户上下文运行函数
- `RunAsyncWithUserContext` - 使用用户上下文异步运行函数

```go
// 示例：在 goroutine 中安全使用 Gin 上下文
context.GinSafeGoroutine(c, func(ctx context.Context, uc *UserContext) {
    // 异步处理逻辑
    logger.Info("异步处理", "user_id", uc.UserID)
})
```

### 事务管理 (pkg/dao/transaction_manager.go)

提供了以下功能：

- `GetDB` - 获取数据库连接，如果上下文中有事务则返回事务对象
- `ExecuteInTransaction` - 在事务中执行函数

```go
// 示例：在事务中执行函数
tm := NewTransactionManager(db)

// 在事务中执行函数
tm.ExecuteInTransaction(ctx, func(ctx context.Context) error {
    // 事务逻辑
    return nil
})
```

### 事件总线 (pkg/eventbus/event_bus.go)

- `Publish` - 发布事件
- `Subscribe` - 订阅事件
- 支持异步和同步事件处理
- 支持安全上下文

```go
// 示例：发布异步消费事件
eventBus := eventbus.NewEventBus()
eventBus.Publish("user.created", &UserCreatedEvent{
    UserID: 1,
}, true)

// 示例：发布同步消费事件
eventBus.Publish("user.created", &UserCreatedEvent{
    UserID: 1,
}, false)

// 示例：订阅事件
sub := eventBus.Subscribe("user.created", func(event eventbus.Event) {
    fmt.Println("用户创建事件", event.Data)
})

```

### 分页组件 (pkg/paginate/paginate.go)

- 链式条件查询

```go
// 使用链式调用创建并配置分页查询选项
options := paginate.NewOptions().
	AddLike("username", req.Username).
	AddLike("email", req.Email).
	AddLike("nickname", req.Nickname).
	AddTimeRange("created_at", req.FromCreateTime, req.ToCreateTime).
	AddOrder("id", "DESC")

// 使用数据库分页查询
users, total, err := s.userDAO.Paginate(ctx, page, size, options)
```

### 认证中间件 (pkg/middleware/auth.go)

认证中间件提供了以下功能：

- `AuthRequired` - 强制要求认证的中间件
- `OptionalAuth` - 可选认证的中间件
- 支持从请求头或Cookie中提取令牌
- 支持公开路径配置，无需认证即可访问
- 支持管理员路径配置，需要特定角色才能访问
- 将用户信息存储在 Gin 上下文和标准库上下文中
- 支持多语言错误消息
- 处理认证失败的情况

### 使用指南

#### 1. 在请求处理中获取用户上下文

```go
func (api *UserAPI) GetUser(c *gin.Context) {
    // 从标准库上下文中获取用户上下文
    ctx := c.Request.Context()
    userCtx := context.GetUserContext(ctx)
    
    // 使用用户上下文
    logger.Info("获取用户信息", "operator_id", userCtx.UserID)
}
```

#### 2. 在异步操作中使用上下文

```go
func (api *UserAPI) CreateUser(c *gin.Context) {
    // 处理请求...
    
    // 启动异步任务
    context.GinSafeGoroutine(c, func(ctx context.Context, uc *UserContext) {
        // 延迟执行，模拟异步处理
        time.Sleep(100 * time.Millisecond)
        
        // 获取用户上下文
        userCtx := context.GetUserContext(ctx)
        
        // 使用上下文信息
        logger.Info("异步处理完成",
            "user_id", createdUser.ID,
            "operator_id", userCtx.UserID,
            "operator_roles", userCtx.Roles,
        )
    })
    
    // 立即返回响应
    c.JSON(http.StatusOK, response.NewResult(response.StatusSuccess, user))
}
```

#### 3. 在中间件中设置用户上下文

```go
func (m *AuthMiddleware) AuthRequired(c *gin.Context) {
    // 提取并验证令牌...
    
    // 设置用户上下文
    m.setUserContext(c, claims, accessToken)
    
    c.Next()
}

// OptionalAuth 可选认证的中间件, 有一些不需要认证的接口, 但是需要获取用户信息
func (m *AuthMiddleware) OptionalAuth(c *gin.Context) {
    // 获取访问令牌
    accessToken := m.extractToken(c)
    if accessToken == "" {
        // 没有访问令牌，继续处理
        c.Next()
        return
    }

    // 验证令牌
    claims, err := m.authService.ValidateToken(c, accessToken)
    if err != nil {
        // 令牌无效，但不阻止请求
        c.Next()
        return
    }

    // 设置用户上下文
    m.setUserContext(c, claims, accessToken)

    c.Next()
}
```

#### 使用示例:

```go
// 在路由中使用
router.GET("/user", m.OptionalAuth(), func(c *gin.Context) {
    // 获取用户上下文
    userCtx := context.GetUserContext(c.Request.Context())
    
    // 使用用户上下文
    logger.Info("获取用户信息", "operator_id", userCtx.UserID)
})

// 在路由组中使用
userGroup.Use(m.authMiddleware.AuthRequired())

// 在路由组中使用可选认证
userGroup.Use(m.authMiddleware.OptionalAuth())

```

### 在程序中使用go异步调用

```go
// 4. 异步执行下载与落盘、更新状态（使用用户上下文，避免请求结束导致取消）
	uc := userContext.GetUserContext(ctx)
	userContext.RunAsyncWithUserContext(uc, func(bgCtx context.Context) {
		// 这里执行你的异步方法
		
	})

```

参考代码如下： DownloadFileService中的downloadAndFinalize方法

#### 代码中获取上下文中的用户信息

```go
// gin context
userId := userContext.GetUserID(c.Request.Context())

// service/dao context
userID, ok := ctx.Value(userContext.UserIDKey).(int64)
if !ok {
    logger.Error("获取用户ID失败", "component", "oauth2_service", "error", errors.Biz("oauth2.error.invalid_context"))
    return
}
```

### 通过模型文件生成数据库表
``` shell
# 查看帮助
go run tools/db_gen.go -help

# 全量模型生成数据库表
go run tools/db_gen.go

# 指定模型文件生成数据库表
go run tools/db_gen.go -model string 

# -list 列出所有模型
go run tools/db_gen.go -list

# gen 通过数据库表生成模型文件
go run tools/db_gen.go -gen
```

## 扩展服务

### 作为独立服务扩展

要添加新的服务功能，请按照以下步骤操作：

1. 在 `pkg/service/service.go` 中扩展 Service 接口
2. 在 `pkg/endpoint/endpoint.go` 中添加新的端点
3. 在 `pkg/transport/http.go` 中添加新的 HTTP 处理程序

### 作为依赖库使用时添加自定义端点

在你的项目中，可以通过以下方式添加自定义端点：

```go
// 定义自定义处理函数
customHandlers := func(r *gin.RouterGroup, tracer opentracing.Tracer, logger logging.Logger) {
	// 添加你的端点
	r.GET("/custom-endpoint", func(c *gin.Context) {
		// 创建 span
		span := tracer.StartSpan("custom-endpoint")
		defer span.Finish()

		// 添加到上下文
		ctx := opentracing.ContextWithSpan(c.Request.Context(), span)

		// 记录日志
		logger.Info("msg", "处理自定义请求")

		// 处理逻辑...
	})
}

// 使用 NewHTTPHandler 时传入自定义处理函数
httpHandler := transport.NewHTTPHandler(endpoints, logger, tracer, metricsHandler, customHandlers)
```

## 错误通知配置

要启用错误通知，请在配置文件中设置：

```json
"errorNotification": {
  "enabled": true,
  "type": "slack",
  "url": "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
}
```

## 分布式追踪

服务默认配置为使用 Jaeger 进行分布式追踪。确保 Jaeger 服务器正在运行，或在配置中禁用追踪。

## 调试和测试

Go-Sea 提供了多种调试和测试功能，帮助开发者快速定位和解决问题。

### 测试端点

- `GET /hello` - 简单的 Hello World 测试
- `GET /panic-test` - 测试 panic 恢复机制

### 日志级别调整

可以通过配置文件动态调整日志级别：

```yaml
logging:
  level: debug  # 可选值: debug, info, warn, error
  format: json  # 可选值: json, text
```

### 开启请求日志

在config配置文件中设置，开启后日志会打印请求和响应的详细信息，方便调试。

```yaml
debug:
  enableRequestLogging: true # 开启请求日志
  logRequestBody: true # 记录请求体
  logResponseBody: true # 记录响应体
  maxRequestBodySize: 4096 # 最大请求体大小
```

## 贡献指南

欢迎贡献代码或提出建议！请遵循以下步骤：

1. Fork 项目
2. 创建你的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交你的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启一个 Pull Request


### 程序逻辑问题

1. 解析pdf文件的时候，没有质量体系检查，需要详细的记录，哪些数据解析出来了，哪些数据没有解析出来
