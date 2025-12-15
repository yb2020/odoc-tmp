package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	google_oauth "golang.org/x/oauth2/google"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	usermodel "github.com/yb2020/odoc/services/user/model"
	userservice "github.com/yb2020/odoc/services/user/service"
)

// GoogleOAuthService 负责处理 Google 特定的 OAuth2 逻辑
type GoogleOAuthService struct {
	config       *oauth2.Config
	redirectPath string
	logger       logging.Logger
	userService  userservice.UserService
}

// NewGoogleOAuthService 创建一个新的 GoogleOAuthService
func NewGoogleOAuthService(cfg *config.Config, logger logging.Logger, userService userservice.UserService) (*GoogleOAuthService, error) {
	googleCfg := cfg.OAuth2.GoogleOAuth2
	if googleCfg.ClientID == "" || googleCfg.ClientSecret == "" {
		logger.Warn("Google OAuth provider is not configured.")
		// 返回nil而不是错误，允许应用在没有配置的情况下启动
		return nil, nil
	}

	oauthConfig := &oauth2.Config{
		ClientID:     googleCfg.ClientID,
		ClientSecret: googleCfg.ClientSecret,
		RedirectURL:  "", // 动态设置，不在这里固定
		Scopes:       googleCfg.Scopes,
		Endpoint:     google_oauth.Endpoint,
	}
	logger.Info("Google OAuth provider configured successfully.")

	return &GoogleOAuthService{
		config:       oauthConfig,
		redirectPath: googleCfg.RedirectPath,
		logger:       logger,
		userService:  userService,
	}, nil
}

// GetGoogleLoginURL 生成认证URL
func (s *GoogleOAuthService) GetGoogleLoginURL(state string, host string, isHTTPS bool) (string, error) {
	if s.config == nil {
		return "", errors.Biz("oauth2.error.google_not_configured")
	}
	
	// 创建配置副本，避免并发修改原始配置
	configCopy := *s.config
	
	// 动态生成redirect_uri，根据当前请求的host和协议
	protocol := "http"
	if isHTTPS {
		protocol = "https"
	}
	
	if host != "" && s.redirectPath != "" {
		configCopy.RedirectURL = fmt.Sprintf("%s://%s%s", protocol, host, s.redirectPath)
	}
	
	authURL := configCopy.AuthCodeURL(state)
	s.logger.Info("生成Google OAuth2认证URL", "redirect_uri", configCopy.RedirectURL, "host", host, "protocol", protocol)
	
	return authURL, nil
}

// HandleGoogleCallback 处理来自 Google 的回调
func (s *GoogleOAuthService) HandleGoogleCallback(ctx context.Context, code string, host string, isHTTPS bool) (*usermodel.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GoogleOAuthService.HandleGoogleCallback")
	defer span.Finish()

	if s.config == nil {
		return nil, errors.Biz("oauth2.error.google_not_configured")
	}

	// 创建配置副本，设置正确的RedirectURL用于token exchange
	configCopy := *s.config
	protocol := "http"
	if isHTTPS {
		protocol = "https"
	}
	
	if host != "" && s.redirectPath != "" {
		configCopy.RedirectURL = fmt.Sprintf("%s://%s%s", protocol, host, s.redirectPath)
	}

	// 步骤 1: 用 code 换取 Google Token
	token, err := configCopy.Exchange(ctx, code)
	if err != nil {
		s.logger.Error("用code换取Google token失败", "error", err, "redirect_uri", configCopy.RedirectURL)
		return nil, errors.BizWrap("oauth2.error.google_exchange_failed", err)
	}

	// 步骤 2: 获取 Google 用户信息
	client := s.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		s.logger.Error("获取Google用户信息失败", "error", err)
		return nil, errors.BizWrap("oauth2.error.google_userinfo_failed", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("获取Google用户信息请求非200", "status", resp.Status, "body", string(body))
		return nil, errors.Biz("oauth2.error.google_userinfo_failed")
	}

	var userInfo usermodel.GoogleUserInfo // Assuming this is defined in usermodel
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		s.logger.Error("解析Google用户信息失败", "error", err)
		return nil, errors.BizWrap("oauth2.error.google_userinfo_decode_failed", err)
	}

	if !userInfo.VerifiedEmail {
		return nil, errors.Biz("oauth2.error.google_email_not_verified")
	}

	// 步骤 3: 调用 UserService 的统一方法来处理用户查找、绑定或创建
	user, err := s.userService.FindOrCreateByGoogleID(ctx, &userInfo)
	if err != nil {
		s.logger.Error("处理Google用户登录失败", "error", err, "email", userInfo.Email)
		// 将底层错误包装成一个对上层更友好的业务错误
		return nil, errors.BizWrap("oauth2.error.user_processing_failed", err)
	}

	return user, nil
}
