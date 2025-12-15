package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/oauth2/helper"
	"github.com/yb2020/odoc/services/oauth2/service"
)

const googleStateCookie = "oauth_google_state"

// GoogleOAuthAPI 封装了 Google OAuth2 的路由处理器
type GoogleOAuthAPI struct {
	logger             logging.Logger
	googleOAuthService *service.GoogleOAuthService
	oauth2Service      service.OAuth2Service
	config             *config.Config
}

// NewGoogleOAuthAPI 创建一个新的 GoogleOAuthAPI
func NewGoogleOAuthAPI(logger logging.Logger, googleOAuthService *service.GoogleOAuthService, oauth2Service service.OAuth2Service, cfg *config.Config) *GoogleOAuthAPI {
	return &GoogleOAuthAPI{
		logger:             logger,
		googleOAuthService: googleOAuthService,
		oauth2Service:      oauth2Service,
		config:             cfg,
	}
}

// RegisterRoutes 注册 Google OAuth2 相关的路由
// func (api *GoogleOAuthAPI) RegisterRoutes(router *gin.RouterGroup) {
// 	if api.googleOAuthService == nil {
// 		api.logger.Info("Google OAuth service not configured, skipping route registration.")
// 		return
// 	}
// 	googleGroup := router.Group("/google")
// 	{
// 		googleGroup.GET("/login", api.LoginHandler)
// 		googleGroup.GET("/callback", api.CallbackHandler)
// 	}
// }

// LoginHandler 处理登录请求，重定向到 Google
func (api *GoogleOAuthAPI) LoginHandler(c *gin.Context) {
	state := uuid.New().String()
	// 将 state 保存到安全、HttpOnly 的 cookie 中
	// 根据协议动态设置secure标志
	isHTTPS := c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https"
	
	// 获取域名，处理端口号和子域名
	host := c.Request.Host
	domain := host
	if colonIndex := strings.Index(host, ":"); colonIndex != -1 {
		domain = host[:colonIndex]
	}
	
	// 处理子域名跨域问题：如果域名包含www前缀，设置为顶级域名以支持子域名共享
	// 例如：www.example.com -> .example.com，支持example.com和www.example.com互通
	if strings.HasPrefix(domain, "www.") {
		domain = domain[3:] // 去掉www.前缀
		domain = "." + domain // 添加.前缀表示顶级域名
	}
	
	// 设置cookie，使用正确的domain
	c.SetCookie(googleStateCookie, state, 3600, "/", domain, isHTTPS, true)

	// 使用之前获取的host变量
	url, err := api.googleOAuthService.GetGoogleLoginURL(state, host, isHTTPS)
	if err != nil {
		c.Error(err)
		return
	}

	api.logger.Info("msg", "Google OAuth2 login", "url", url, "host", host, "https", isHTTPS)

	c.Redirect(http.StatusTemporaryRedirect, url)
}

// CallbackHandler 处理 Google 的回调请求
func (api *GoogleOAuthAPI) CallbackHandler(c *gin.Context) {
	// 验证 state
	cookieState, err := c.Cookie(googleStateCookie)
	if err != nil {
		c.Error(errors.Biz("oauth2.error.missing_state_cookie"))
		return
	}

	if c.Query("state") != cookieState {
		c.Error(errors.Biz("oauth2.error.invalid_state"))
		return
	}

	// 清除 state cookie，使用相同的domain逻辑
	host := c.Request.Host
	domain := host
	if colonIndex := strings.Index(host, ":"); colonIndex != -1 {
		domain = host[:colonIndex]
	}
	
	// 处理子域名跨域问题：如果域名包含www前缀，设置为顶级域名以支持子域名共享
	if strings.HasPrefix(domain, "www.") {
		domain = domain[3:] // 去掉www.前缀
		domain = "." + domain // 添加.前缀表示顶级域名
	}
	
	isHTTPS := c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https"
	c.SetCookie(googleStateCookie, "", -1, "/", domain, isHTTPS, true)

	code := c.Query("code")
	if code == "" {
		c.Error(errors.Biz("oauth2.error.missing_code"))
		return
	}
	api.logger.Info("msg", "Google OAuth2 callback", "code", code)

	// 使用之前获取的host和isHTTPS变量进行token exchange
	
	// 使用从 Google 获取的授权码(code)来处理回调逻辑，
	// 这通常包括用 code 换取 access_token，然后用 access_token 获取用户信息，
	// 最后根据用户信息在本地数据库查找或创建一个对应的用户。
	user, err := api.googleOAuthService.HandleGoogleCallback(c.Request.Context(), code, host, isHTTPS)
	// 检查上述过程是否发生错误（例如，网络问题、Google API 返回错误等）。
	if err != nil {
		// 如果有错误，将错误传递给 Gin 的错误处理中间件，并立即停止处理该请求。
		c.Error(err)
		return
	}

	// 在成功获取或创建用户后，调用通用的认证服务为该用户生成一个系��内部的认证令牌（通常是 JWT）。
	// 这个令牌将用于后续所有需要用户登录才能访问的 API 请求。
	tokenResponse, err := api.oauth2Service.CreateTokenForGoogleUser(c.Request.Context(), user.Id, user.Username, user.Roles.ToStringSlice())
	// 检查令牌生成过程是否出错。
	if err != nil {
		// 如果生成令牌失败，同样将错误传递给错误处理中间件并停止处理。
		c.Error(err)
		return
	}

	// 步骤 7: 登录成功，写入 Cookie，与内建登录流程保持一致
	helper.WriteSuccessLoginCookie(c, tokenResponse.UserId, tokenResponse.AccessToken, tokenResponse.ExpiresAt, api.config.OAuth2.AppID)

	// 步骤 8: 重定向用户到前端主页
	// 注意：这里的重定向URL应该在配置文件中定义
	redirectURL := api.config.OAuth2.GoogleOAuth2.LoginSuccessURL
	if redirectURL == "" {
		// 提供一个默认的回退URL
		redirectURL = "/"
		api.logger.Warn("msg", "未配置前端重定向URL (App.FrontendURL)，将重定向到根目录'/'")
	}
	c.Redirect(http.StatusFound, redirectURL)
}
