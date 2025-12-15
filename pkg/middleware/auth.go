package middleware

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/i18n"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
)

// Claims 令牌声明接口
// 这个接口允许不同的令牌实现（如OAuth2、JWT、基本认证等）
type Claims interface {
	// GetUserID 获取用户ID
	GetUserID() string
	// GetUsername 获取用户名
	GetUsername() string
	// GetRoles 获取用户角色
	GetRoles() []string
	// GetDevice 获取设备信息
	GetDevice() string
	// GetServiceId 获取服务ID
	GetServiceId() string
	// GetServiceName 获取服务名称
	GetServiceName() string
}

// AuthService 认证服务接口
// 这个接口允许不同的认证实现（如OAuth2、JWT、基本认证等）
type AuthService interface {
	// ValidateToken 验证令牌
	ValidateToken(ctx *gin.Context, token string) (Claims, error)
	// ValidateServiceToken 验证服务令牌
	ValidateServiceToken(ctx *gin.Context, token string) (Claims, error)
}

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	authService AuthService
	logger      logging.Logger
	config      config.Config
	localizer   i18n.Localizer
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(config config.Config, logger logging.Logger, localizer i18n.Localizer, authService AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		logger:      logger,
		config:      config,
		localizer:   localizer,
	}
}

// AuthRequired 需要认证的中间件
func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否是公开路径
		path := c.Request.URL.Path

		// 防御性检查：确保 PublicPaths 不为 nil
		if m.config.OAuth2.ResourceProtection.PublicPaths != nil {
			for _, publicPath := range m.config.OAuth2.ResourceProtection.PublicPaths {
				if strings.HasPrefix(path, publicPath) {
					c.Next()
					return
				}
			}
		}

		// 获取访问令牌
		accessToken := m.extractToken(c)
		if accessToken == "" {
			// 使用本地化器翻译错误消息
			msg := "认证令牌缺失" // 默认消息
			if m.localizer != nil {
				// 添加调试日志
				msg = m.localizer.Localize("auth.token.missing", c)
				m.logger.Debug("msg", "本地化结果", "result", msg)
			}
			response.SystemErrorNoData(c, response.Code_Unauthorized, response.Status_AuthInvalidToken, msg)
			c.Abort()
			return
		}

		// 验证令牌
		claims, err := m.authService.ValidateToken(c, accessToken)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		// 检查是否是管理员路径
		// 防御性检查：确保 AdminPaths 不为 nil
		if m.config.OAuth2.ResourceProtection.AdminPaths != nil {
			for _, adminPath := range m.config.OAuth2.ResourceProtection.AdminPaths {
				if strings.HasPrefix(path, adminPath) {
					// 检查是否有管理员角色
					hasAdminRole := false
					// 防御性检查：确保 AdminRoles 不为 nil
					if m.config.OAuth2.ResourceProtection.AdminRoles != nil && len(claims.GetRoles()) > 0 {
						for _, role := range claims.GetRoles() {
							for _, adminRole := range m.config.OAuth2.ResourceProtection.AdminRoles {
								if role == adminRole {
									hasAdminRole = true
									break
								}
							}
							if hasAdminRole {
								break
							}
						}
					}

					if !hasAdminRole {
						// 使用本地化器翻译错误消息
						msg := "需要管理员权限" // 默认消息
						if m.localizer != nil {
							msg = m.localizer.Localize("auth.admin.required", c)
						}
						response.SystemErrorNoData(c, response.Code_Unauthorized, response.Status_AuthInvalidRole, msg)
						c.Abort()
						return
					}
					break
				}
			}
		}

		m.setUserContext(c, claims, accessToken)
		c.Next()
	}
}

// OptionalAuth 可选认证的中间件, 有一些不需要认证的接口, 但是需要获取用户信息
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
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
}

func (m *AuthMiddleware) setUserContext(c *gin.Context, claims Claims, accessToken string) {
	// 将用户信息存储到上下文中
	c.Set(string(context.UserIDKey), claims.GetUserID())
	c.Set(string(context.UsernameKey), claims.GetUsername())
	c.Set(string(context.RolesKey), claims.GetRoles())
	c.Set(string(context.DeviceKey), claims.GetDevice())
	c.Set(string(context.ClaimsKey), claims)
	c.Set(string(context.AccessTokenKey), accessToken)
	c.Set(string(context.AuthenticatedKey), true)

	// 创建用户上下文并存储到标准库的context.Context中
	uc := context.NewUserContext().
		SetUserID(claims.GetUserID()).
		SetUsername(claims.GetUsername()).
		SetRoles(claims.GetRoles()).
		SetDevice(claims.GetDevice()).
		SetClaims(claims).
		SetAccessToken(accessToken).
		SetAuthenticated(true)

	// 将用户上下文转换为标准库的context.Context
	ctx := uc.ToContext(c.Request.Context())

	// 更新请求的上下文
	c.Request = c.Request.WithContext(ctx)
}

// extractToken 从请求中提取令牌
// 优先从 cookie 中获取，取不到再从请求头获取
func (m *AuthMiddleware) extractToken(c *gin.Context) string {
	// 1. 尝试从 cookie 中获取令牌
	appID := m.config.OAuth2.AppID

	// 构造 cookie 名称
	cookieName := fmt.Sprintf(m.config.OAuth2.TokenStorage.CookieTokenLabel, appID)
	cookie, err := c.Cookie(cookieName)
	if err == nil && cookie != "" {
		// 解码 base64 编码的令牌
		tokenBytes, err := base64.StdEncoding.DecodeString(cookie)
		if err == nil {
			return string(tokenBytes)
		}
		m.logger.Debug("msg", "从cookie解码令牌失败", "error", err.Error())
	}

	// 2. 尝试从请求头中获取令牌
	// 确定请求头名称
	headerName := m.config.OAuth2.TokenStorage.TokenHeaderName
	authHeader := c.GetHeader(headerName)

	// 处理自定义头的特殊格式（Bearer token）
	if headerName == m.config.OAuth2.TokenStorage.TokenHeaderName {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
		return ""
	}

	// 对于其他自定义头，直接返回值
	return authHeader
}

// ServiceAuthRequired middleware for service token authentication
func (m *AuthMiddleware) ServiceAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否是公开路径
		path := c.Request.URL.Path
		// 防御性检查：确保 PublicPaths 不为 nil
		if m.config.OAuth2.ServiceAccount.ResourceProtection.PublicPaths != nil {
			for _, publicPath := range m.config.OAuth2.ServiceAccount.ResourceProtection.PublicPaths {
				if strings.HasPrefix(path, publicPath) {
					c.Next()
					return
				}
			}
		}

		// Get token from header
		tokenHeader := m.config.OAuth2.ServiceAccount.TokenHeaderName
		serviceAccesstoken := c.GetHeader(tokenHeader)
		if serviceAccesstoken == "" {
			m.logger.Error("msg", "服务令牌未找到", "component", "service_token_middleware")
			response.SystemErrorNoData(c, response.Code_Unauthorized, response.Status_ServiceTokenNot, "oauth2.error.service_token_not_found")
			c.Abort()
			return
		}

		// Validate token
		claims, err := m.authService.ValidateServiceToken(c, serviceAccesstoken)
		if err != nil {
			m.logger.Warn("msg", "验证服务令牌失败", "error", err.Error())
			c.Error(err)
			c.Abort()
			return
		}

		if claims == nil || claims.GetServiceId() == "" || claims.GetServiceName() == "" {
			m.logger.Error("msg", "服务令牌无效", "component", "service_token_middleware")
			c.Error(errors.Biz("oauth2.error.invalid_service_token"))
			c.Abort()
			return
		}

		// 设置服务上下文
		m.setServiceContext(c, claims, serviceAccesstoken)
		c.Next()
	}

}

// set service context
func (m *AuthMiddleware) setServiceContext(c *gin.Context, claims Claims, serviceAccesstoken string) {
	// Create service context
	// 设置服务上下文
	c.Set(string(context.ServiceIdKey), claims.GetServiceId())
	c.Set(string(context.ServiceNameKey), claims.GetServiceName())
	// 设置服务令牌
	c.Set(string(context.AccessTokenKey), serviceAccesstoken)

	// 创建用户上下文并存储到标准库的context.Context中
	uc := context.NewUserContext().
		SetServiceId(claims.GetServiceId()).
		SetServiceName(claims.GetServiceName()).
		SetRoles(claims.GetRoles()).
		SetAccessToken(serviceAccesstoken).
		SetAuthenticated(true)

	// 将用户上下文转换为标准库的context.Context
	ctx := uc.ToContext(c.Request.Context())

	// 更新请求上下文
	c.Request = c.Request.WithContext(ctx)
}
