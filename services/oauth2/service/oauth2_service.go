package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/config"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/i18n"
	"github.com/yb2020/odoc/pkg/logging"
	pkgmodel "github.com/yb2020/odoc/pkg/model"
	"github.com/yb2020/odoc/pkg/utils"
	pb "github.com/yb2020/odoc/proto/gen/go/oauth2"
	"github.com/yb2020/odoc/services/oauth2/dao"
	"github.com/yb2020/odoc/services/oauth2/model"
	userservice "github.com/yb2020/odoc/services/user/service"
)

// OAuth2Service OAuth2服务实现
type OAuth2Service struct {
	tokenDAO      dao.TokenDAO
	userService   userservice.UserService
	logger        logging.Logger
	tracer        opentracing.Tracer
	localizer     i18n.Localizer
	jwtSecret     string
	jwtIssuer     string
	jwtExpiry     time.Duration
	refreshExpiry time.Duration
	config        *config.Config
	clientDAO     dao.OAuth2ClientsDAO
	authCodeDAO   dao.OAuth2CodeDAO
}

// NewOAuth2Service 创建OAuth2服务
func NewOAuth2Service(
	tokenDAO dao.TokenDAO,
	userService *userservice.UserService,
	logger logging.Logger,
	tracer opentracing.Tracer,
	config *config.Config,
	localizer i18n.Localizer,
	clientDAO dao.OAuth2ClientsDAO,
	authCodeDAO dao.OAuth2CodeDAO,
) OAuth2Service {
	// 从配置中获取JWT密钥
	jwtSecret := config.OAuth2.JWT.Secret
	if jwtSecret == "" {
		logger.Warn("未配置JWT密钥，使用默认密钥，这在生产环境中是不安全的", "component", "oauth2_service", "jwt_secret", "default-jwt-secret")
		jwtSecret = "default-jwt-secret" // 默认密钥，实际项目中应该使用更安全的密钥
	}

	// 从配置中获取JWT发行者
	jwtIssuer := config.OAuth2.JWT.Issuer
	if jwtIssuer == "" {
		logger.Warn("未配置JWT发行者，使用默认值", "component", "oauth2_service", "jwt_issuer", "go-sea")
		jwtIssuer = "go-sea" // 默认发行者
	}

	// 使用GetInt获取秒数，然后转换为time.Duration
	jwtExpirySeconds := config.OAuth2.JWT.Expiry
	jwtExpiry := time.Duration(jwtExpirySeconds) * time.Second
	if jwtExpiry <= 0 {
		logger.Warn("未配置JWT过期时间或配置无效，使用默认值24小时", "component", "oauth2_service", "default_expiry", "24h")
		jwtExpiry = 24 * time.Hour // 默认24小时
	}

	// 使用GetInt获取秒数，然后转换为time.Duration
	refreshExpirySeconds := config.OAuth2.JWT.RefreshExpiry
	refreshExpiry := time.Duration(refreshExpirySeconds) * time.Second
	if refreshExpiry <= 0 {
		logger.Warn("未配置刷新令牌过期时间或配置无效，使用默认值7天", "component", "oauth2_service", "default_refresh_expiry", "168h")
		refreshExpiry = 7 * 24 * time.Hour // 默认7天
	}

	return OAuth2Service{
		tokenDAO:      tokenDAO,
		userService:   *userService,
		logger:        logger,
		tracer:        tracer,
		localizer:     localizer,
		jwtSecret:     jwtSecret,
		jwtIssuer:     jwtIssuer,
		jwtExpiry:     jwtExpiry,
		refreshExpiry: refreshExpiry,
		config:        config,
		clientDAO:     clientDAO,
		authCodeDAO:   authCodeDAO,
	}
}

// CreateTokenForGoogleUser creates tokens for a user authenticated via Google OAuth2.
func (s *OAuth2Service) CreateTokenForGoogleUser(ctx context.Context, userID string, username string, roles []string) (*model.TokenResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OAuth2Service.CreateTokenForGoogleUser")
	defer span.Finish()

	// For OAuth2 logins, the device can be identified as such.
	device := "google_oauth2_login"

	tokenInfo, err := s.generateAndSaveToken(ctx, userID, username, roles, device)
	if err != nil {
		return nil, errors.BizWrap("oauth2.error.token_generation_failed", err)
	}

	return &model.TokenResponse{
		UserId:       tokenInfo.UserId,
		AccessToken:  tokenInfo.AccessToken,
		RefreshToken: tokenInfo.RefreshToken,
		ExpiresAt:    uint64(tokenInfo.ExpiresAt.Unix()),
		ExpiresIn:    uint64(s.jwtExpiry.Seconds()),
	}, nil
}

// generateAndSaveToken 生成并保存令牌的通用方法
func (s *OAuth2Service) generateAndSaveToken(
	ctx context.Context,
	userID string,
	username string,
	roles []string,
	device string,
) (*model.OAuth2Token, error) {
	// 生成新的令牌ID
	tokenID := uuid.New().String()

	// 创建JWT声明
	now := time.Now()
	expiresAt := now.Add(s.jwtExpiry)
	refreshExpiresAt := now.Add(s.refreshExpiry)

	claims := model.Claims{
		UserId:   userID,
		Roles:    roles,
		Device:   device,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			Issuer:    s.jwtIssuer,
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: expiresAt.Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			Id:        tokenID,
		},
	}

	// 生成JWT令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		s.logger.Error("生成JWT令牌失败", "component", "oauth2_service", "error", err)
		return nil, errors.BizWrap("oauth2.error.token_generation_failed", err)
	}

	// 生成刷新令牌
	refreshTokenID := uuid.New().String()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    s.jwtIssuer,
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: refreshExpiresAt.Unix(),
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		Id:        refreshTokenID,
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		s.logger.Error("生成刷新令牌失败", "component", "oauth2_service", "error", err)
		return nil, errors.BizWrap("oauth2.error.token_generation_failed", err)
	}

	// 保存令牌信息到数据库
	tokenInfo := &model.OAuth2Token{
		TokenId:      tokenID,
		UserId:       userID,
		AccessToken:  tokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    expiresAt,
		Device:       device,
	}

	err = s.tokenDAO.SaveToken(ctx, tokenInfo)
	if err != nil {
		s.logger.Error("保存令牌失败", "component", "oauth2_service", "error", err)
		return nil, errors.BizWrap("oauth2.error.token_storage_failed", err)
	}

	// 返回令牌响应
	return tokenInfo, nil
}

// GenerateToken 生成令牌
func (s *OAuth2Service) GenerateToken(ctx context.Context, request *model.TokenRequest) (*model.TokenResponse, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OAuth2Service.GenerateToken")
	defer span.Finish()

	// 实现用户认证逻辑
	// 使用电子邮件地址作为用户名进行认证
	username := request.GetUsername()
	s.logger.Info("尝试认证用户", "component", "oauth2_service", "username", username)

	// 使用新添加的GetUserByEmail方法获取用户
	user, err := s.userService.GetUserByEmail(ctx, username)
	if err != nil {
		s.logger.Error("用户认证失败", "component", "oauth2_service", "username", username, "error", err)
		return nil, errors.Biz("invalid_credentials")
	}

	// 验证密码
	password := utils.StrengthenPassword(request.GetPassword(), user.Id)
	if password != user.Password {
		s.logger.Error("用户认证失败", "component", "oauth2_service", "username", username, "error", errors.Biz("invalid_credentials"))
		return nil, errors.Biz("invalid_credentials")
	}
	s.logger.Info("用户认证成功", "component", "oauth2_service", "user_id", user.Id)

	// GenerateToken 方法中
	roles := []string{}
	for _, role := range user.Roles {
		roles = append(roles, role.String())
	}

	tokenInfo, err := s.generateAndSaveToken(ctx, user.Id, user.Username, roles, request.Device)
	if err != nil {
		return nil, errors.BizWrap("oauth2.error.token_generation_failed", err)
	}

	// 返回令牌响应
	return &model.TokenResponse{
		UserId:       tokenInfo.UserId,
		AccessToken:  tokenInfo.AccessToken,
		RefreshToken: tokenInfo.RefreshToken,
		ExpiresAt:    uint64(tokenInfo.ExpiresAt.Unix()),
		ExpiresIn:    uint64(s.jwtExpiry.Seconds()),
	}, nil
}

// RefreshToken 刷新令牌
func (s *OAuth2Service) RefreshToken(ctx context.Context, request *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OAuth2Service.RefreshToken")
	defer span.Finish()

	// 解析刷新令牌
	token, err := jwt.ParseWithClaims(request.RefreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			msg := fmt.Sprintf("unexpected signing method: %v", token.Header["alg"])
			return nil, errors.Biz(msg)
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		s.logger.Error("解析刷新令牌失败", "component", "oauth2_service", "error", err)
		return nil, errors.BizWrap("oauth2.error.invalid_refresh_token", err)
	}

	// 验证令牌
	if !token.Valid {
		return nil, errors.Biz("oauth2.error.invalid_refresh_token")
	}

	// 获取声明
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, errors.BizWrap("oauth2.error.invalid_refresh_token", nil)
	}

	// 获取令牌信息
	tokenInfo, err := s.tokenDAO.GetTokenByRefreshToken(ctx, request.RefreshToken)
	if err != nil {
		s.logger.Error("获取令牌信息失败", "component", "oauth2_service", "error", err)
		return nil, errors.BizWrap("oauth2.error.token_retrieval_failed", err)
	}

	if tokenInfo == nil || tokenInfo.Revoked {
		return nil, errors.Biz("oauth2.error.token_revoked_or_not_found")
	}

	// 验证用户ID
	userID, ok := ctx.Value(userContext.UserIDKey).(string)
	if !ok {
		s.logger.Error("获取用户ID失败", "component", "oauth2_service", "error", errors.Biz("oauth2.error.invalid_refresh_token"))
		return nil, errors.BizWrap("oauth2.error.invalid_refresh_token", nil)
	}
	_, err = fmt.Sscanf(claims.Subject, "%s", &userID)
	if err != nil {
		s.logger.Error("用户ID格式错误", "component", "oauth2_service", "subject", claims.Subject, "error", err)
		return nil, errors.BizWrap("oauth2.error.invalid_refresh_token", err)
	}

	if userID != tokenInfo.UserId {
		s.logger.Error("用户ID不匹配", "component", "oauth2_service", "tokenUserID", tokenInfo.UserId, "claimUserID", userID)
		return nil, errors.BizWrap("oauth2.error.invalid_refresh_token", nil)
	}

	// 获取用户信息
	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Error("获取用户信息失败", "component", "oauth2_service", "userID", userID, "error", err)
		return nil, errors.BizWrap("oauth2.error.user_not_found", err)
	}

	// 获取用户角色
	roles := []string{}
	for _, role := range user.Roles {
		roles = append(roles, role.String())
	}

	// 实际项目中应该根据用户实际角色来判断
	// RefreshToken 方法中
	// 在验证刷新令牌并获取用户信息后
	response, err := s.generateAndSaveToken(ctx, user.Id, user.Username, roles, tokenInfo.Device)
	if err != nil {
		return nil, errors.BizWrap("oauth2.error.token_generation_failed", err)
	}

	// 撤销旧令牌
	err = s.tokenDAO.RevokeToken(ctx, tokenInfo.TokenId)
	if err != nil {
		s.logger.Error("撤销旧令牌失败", "component", "oauth2_service", "tokenID", tokenInfo.TokenId, "error", err)
		// 继续处理，不返回错误
	}

	// 返回令牌响应
	return &pb.RefreshTokenResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		ExpiresAt:    uint64(response.ExpiresAt.Unix()),
		ExpiresIn:    uint64(s.jwtExpiry.Seconds()),
	}, nil
}

// ValidateToken 验证令牌
func (s *OAuth2Service) ValidateToken(ctx context.Context, request *pb.ValidateRequest) (*model.Claims, error) {
	// 创建跟踪span
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OAuth2Service.ValidateToken")
	defer span.Finish()

	// 解析令牌
	token, err := jwt.ParseWithClaims(request.AccessToken, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		s.logger.Error("解析令牌失败", "component", "oauth2_service", "error", err)
		return nil, errors.BizWrap("oauth2.error.invalid_token", err)
	}

	// 验证令牌
	if !token.Valid {
		return nil, errors.BizWrap("oauth2.error.invalid_token", nil)
	}

	// 获取声明
	claims, ok := token.Claims.(*model.Claims)
	if !ok {
		return nil, errors.BizWrap("oauth2.error.invalid_token", nil)
	}

	// 根据配置决定验证方式
	s.logger.Debug("验证令牌", "component", "oauth2_service", "token_id", claims.Id)

	// 验证令牌状态（无论使用哪种存储方式，接口调用是相同的）
	tokenInfo, err := s.tokenDAO.GetTokenByAccessToken(ctx, request.AccessToken)

	if err != nil {
		s.logger.Error("获取令牌信息失败", "component", "oauth2_service", "error", err)
		return nil, errors.BizWrap("oauth2.error.token_retrieval_failed", err)
	}

	if tokenInfo == nil || tokenInfo.Revoked {
		return nil, errors.Biz("oauth2.error.token_revoked_or_not_found")
	}

	// 验证令牌是否过期
	if time.Now().After(tokenInfo.ExpiresAt) {
		return nil, errors.Biz("oauth2.error.token_expired")
	}

	// 返回验证结果
	return claims, nil
}

// RevokeUserToken 撤销令牌
func (s *OAuth2Service) RevokeUserToken(ctx context.Context, accessToken string) error {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OAuth2Service.RevokeToken")
	defer span.Finish()

	// 先通过访问令牌获取令牌信息
	tokenInfo, err := s.tokenDAO.GetTokenByAccessToken(ctx, accessToken)
	if err != nil {
		s.logger.Error("获取令牌信息失败", "component", "oauth2_service", "error", err)
		return errors.BizWrap("oauth2.error.token_retrieval_failed", err)
	}

	if tokenInfo == nil {
		s.logger.Warn("令牌不存在", "component", "oauth2_service", "access_token", accessToken)
		return errors.Biz("oauth2.error.token_not_found")
	}

	// 撤销令牌
	if err := s.tokenDAO.RevokeToken(ctx, tokenInfo.TokenId); err != nil {
		s.logger.Error("撤销令牌失败", "component", "oauth2_service", "error", err)
		return errors.BizWrap("oauth2.error.token_revocation_failed", err)
	}

	return nil
}

// 获取用户的所有令牌
func (s *OAuth2Service) GetUserTokens(ctx context.Context, userID string) ([]*model.OAuth2Token, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OAuth2Service.GetUserTokens")
	defer span.Finish()

	// 从数据库获取用户的所有令牌
	tokens, err := s.tokenDAO.GetTokensByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("获取用户令牌失败", "component", "oauth2_service", "userID", userID, "error", err)
		return nil, errors.BizWrap("oauth2.error.token_retrieval_failed", err)
	}

	return tokens, nil
}

// 撤销用户的所有令牌
func (s *OAuth2Service) RevokeUserTokens(ctx context.Context, userID string) error {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OAuth2Service.RevokeUserTokens")
	defer span.Finish()

	// 撤销用户的所有令牌
	if err := s.tokenDAO.RevokeAllTokensByUserID(ctx, userID); err != nil {
		s.logger.Error("撤销用户所有令牌失败", "component", "oauth2_service", "userID", userID, "error", err)
		return errors.BizWrap("oauth2.error.token_revocation_failed", err)
	}

	return nil
}

// 清理过期令牌
func (s *OAuth2Service) CleanupExpiredTokens(ctx context.Context) error {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OAuth2Service.CleanupExpiredTokens")
	defer span.Finish()

	// 清理过期令牌
	if err := s.tokenDAO.CleanupExpiredTokens(ctx); err != nil {
		s.logger.Error("清理过期令牌失败", "component", "oauth2_service", "error", err)
		return errors.Biz("oauth2.error.token_cleanup_failed")
	}

	return nil
}

// GetRSAKey 获取RSA公钥
func (s *OAuth2Service) GetRSAKey(ctx context.Context, request *pb.GetRSAKeyRequest) (*pb.GetRSAKeyResponse, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OAuth2Service.GetRSAKey")
	defer span.Finish()

	// 从配置中获取RSA公钥
	publicKey := s.config.OAuth2.RSA.PublicKey
	if publicKey == "" {
		s.logger.Error("msg", "RSA公钥未配置")
		return nil, errors.Biz("oauth2.error.rsa_public_key_not_configured")
	}

	// 创建响应
	response := &pb.GetRSAKeyResponse{
		PublicKey: publicKey,
	}

	return response, nil
}

// GetAuthCode 获取授权码
func (s *OAuth2Service) GetAuthCode(ctx context.Context, request *pb.GetAuthCodeRequest) (*pb.GetAuthCodeResponse, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OAuth2Service.GetAuthCode")
	defer span.Finish()

	// 验证客户端是否存在且有效
	client, err := s.clientDAO.GetClientByID(ctx, request.GetClientId())
	if err != nil {
		s.logger.Error("获取客户端信息失败", "component", "oauth2_service", "error", err)
		return nil, errors.BizWrap("oauth2.error.client_retrieval_failed", err)
	}

	if client == nil {
		s.logger.Error("客户端不存在", "component", "oauth2_service", "client_id", request.GetClientId())
		return nil, errors.Biz("oauth2.error.client_not_found")
	}

	// 验证重定向URI是否合法
	validRedirectURI := false
	for _, uri := range client.RedirectUris {
		if uri == request.GetRedirectUri() {
			validRedirectURI = true
			break
		}
	}

	if !validRedirectURI {
		s.logger.Error("重定向URI不合法", "component", "oauth2_service", "redirect_uri", request.GetRedirectUri())
		return nil, errors.Biz("oauth2.error.invalid_redirect_uri")
	}

	// 验证作用域
	requestedScopes := strings.Split(request.GetScope(), " ")
	for _, scope := range requestedScopes {
		if scope != "" && !client.IsValidScope(scope) {
			s.logger.Error("请求的作用域不合法", "component", "oauth2_service", "scope", scope)
			return nil, errors.BizWrap("oauth2.error.invalid_scope", nil)
		}
	}

	// 生成授权码
	authCode := uuid.New().String()

	// 创建授权码记录
	// 将 int 转换为 time.Duration（假设配置中的值是秒）
	authCodeLifetime := time.Duration(s.config.OAuth2.TokenStorage.AuthCode.Lifetime) * time.Second
	expiresAt := time.Now().Add(authCodeLifetime)
	authCodeInfo := &model.OAuth2AuthCode{
		Code:        authCode,
		ClientId:    request.GetClientId(),
		UserId:      request.GetUserId(),
		RedirectUri: request.GetRedirectUri(),
		Scope:       request.GetScope(),
		ExpiresAt:   expiresAt,
		CreatedAt:   time.Now(),
	}

	// 存储授权码
	err = s.authCodeDAO.SaveAuthCode(ctx, authCodeInfo)
	if err != nil {
		s.logger.Error("保存授权码失败", "component", "oauth2_service", "error", err)
		return nil, errors.BizWrap("oauth2.error.auth_code_save_failed", err)
	}

	s.logger.Info("成功生成授权码", "component", "oauth2_service", "client_id", request.GetClientId(), "user_id", request.GetUserId())

	// 创建响应
	response := &pb.GetAuthCodeResponse{
		Code: authCode,
	}

	return response, nil
}

// GenerateServiceToken 生成服务间通信的长期令牌
func (s *OAuth2Service) GenerateServiceToken(ctx context.Context, serviceId, serviceName string) (*model.TokenResponse, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OAuth2Service.GenerateServiceToken")
	defer span.Finish()

	// 验证服务是否在配置中
	var serviceFound bool

	for _, service := range s.config.OAuth2.ServiceAccount.Services {
		if service.Id == serviceId && service.Name == serviceName {
			serviceFound = true
			break
		}
	}

	if !serviceFound {
		s.logger.Error("服务未在配置中注册", "component", "oauth2_service", "service_id", serviceId, "service_name", serviceName)
		return nil, errors.Biz("oauth2.error.service_not_registered")
	}

	// 生成新的令牌ID
	tokenID := uuid.New().String()

	// 创建JWT声明
	now := time.Now()
	serviceTokenExpiry := time.Duration(s.config.OAuth2.ServiceAccount.TokenExpiry) * time.Second
	expiresAt := now.Add(serviceTokenExpiry)

	// 使用服务角色
	roles := s.config.OAuth2.ServiceAccount.Roles

	claims := model.Claims{
		UserId:   "0", // 服务令牌不关联用户
		Roles:    roles,
		Device:   "service",
		Username: serviceName,
		ServiceInfo: &model.ServiceInfo{
			ServiceId:   serviceId,
			ServiceName: serviceName,
		},
		StandardClaims: jwt.StandardClaims{
			Issuer:    s.jwtIssuer,
			Subject:   serviceId,
			ExpiresAt: expiresAt.Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			Id:        tokenID,
		},
	}

	// 生成JWT令牌
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		s.logger.Error("生成JWT令牌失败", "component", "oauth2_service", "error", err)
		return nil, errors.BizWrap("oauth2.error.token_generation_failed", err)
	}

	// 生成刷新令牌
	refreshTokenID := uuid.New().String()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    s.jwtIssuer,
		Subject:   fmt.Sprintf("%d", serviceId),
		ExpiresAt: expiresAt.Unix(),
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		Id:        refreshTokenID,
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		s.logger.Error("生成刷新令牌失败", "component", "oauth2_service", "error", err)
		return nil, errors.BizWrap("oauth2.error.token_generation_failed", err)
	}

	// 创建令牌记录
	oauth2Token := &model.OAuth2Token{
		TokenId:      tokenID,
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		UserId:       "0", // 服务令牌不关联用户
		Device:       "service",
		ExpiresAt:    expiresAt,
		Roles:        pkgmodel.StringSlice(roles),
	}

	// 使用服务令牌的Redis键前缀
	redisKeyPrefix := s.config.OAuth2.ServiceAccount.TokenStorage.RedisKeyPrefix

	// 存储令牌
	err = s.tokenDAO.SaveTokenWithPrefix(ctx, oauth2Token, redisKeyPrefix)
	if err != nil {
		s.logger.Error("保存服务令牌失败", "component", "oauth2_service", "error", err)
		return nil, errors.BizWrap("oauth2.error.token_save_failed", err)
	}

	s.logger.Info("成功生成服务令牌", "component", "oauth2_service", "service_id", serviceId, "service_name", serviceName)

	// 创建响应
	response := &model.TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString, // 服务令牌不需要刷新令牌
		ExpiresIn:    uint64(serviceTokenExpiry.Seconds()),
		TokenType:    "Bearer",
	}

	return response, nil
}

// ValidateServiceToken 验证服务令牌
func (s *OAuth2Service) ValidateServiceToken(ctx context.Context, request *pb.ValidateRequest) (*model.Claims, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OAuth2Service.ValidateServiceToken")
	defer span.Finish()

	// 解析JWT令牌
	token, err := jwt.ParseWithClaims(request.GetAccessToken(), &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		s.logger.Error("解析JWT令牌失败", "component", "oauth2_service", "error", err)
		return nil, errors.BizWrap("oauth2.error.token_parse_failed", err)
	}

	// 验证令牌是否有效
	if !token.Valid {
		s.logger.Error("JWT令牌无效", "component", "oauth2_service")
		return nil, errors.Biz("oauth2.error.invalid_token")
	}

	// 获取声明
	claims, ok := token.Claims.(*model.Claims)
	if !ok {
		s.logger.Error("获取JWT声明失败", "component", "oauth2_service")
		return nil, errors.Biz("oauth2.error.invalid_claims")
	}

	// 验证是否为服务令牌
	if claims.ServiceInfo == nil {
		s.logger.Error("不是服务令牌", "component", "oauth2_service")
		return nil, errors.Biz("oauth2.error.not_service_token")
	}

	// 验证服务是否在配置中
	serviceFound := false
	for _, service := range s.config.OAuth2.ServiceAccount.Services {
		if service.Id == claims.ServiceInfo.ServiceId && service.Name == claims.ServiceInfo.ServiceName {
			serviceFound = true
			break
		}
	}

	if !serviceFound {
		s.logger.Error("服务未在配置中注册", "component", "oauth2_service", "service_id", claims.ServiceInfo.ServiceId, "service_name", claims.ServiceInfo.ServiceName)
		return nil, errors.Biz("oauth2.error.service_not_registered")
	}

	// 使用服务令牌的Redis键前缀
	redisKeyPrefix := s.config.OAuth2.ServiceAccount.TokenStorage.RedisKeyPrefix

	// 从存储中验证令牌
	tokenID := claims.Id
	storedToken, err := s.tokenDAO.GetTokenByIDWithPrefix(ctx, tokenID, redisKeyPrefix)
	if err != nil {
		s.logger.Error("从存储中获取令牌失败", "component", "oauth2_service", "error", err)
		return nil, errors.BizWrap("oauth2.error.token_retrieval_failed", err)
	}

	if storedToken == nil {
		s.logger.Error("令牌不存在或已过期", "component", "oauth2_service")
		return nil, errors.Biz("oauth2.error.token_not_found")
	}

	// 验证令牌是否过期
	if time.Now().After(storedToken.ExpiresAt) {
		s.logger.Error("令牌已过期", "component", "oauth2_service")
		return nil, errors.Biz("oauth2.error.token_expired")
	}

	return claims, nil
}

// RevokeServiceToken 撤销服务令牌
func (s *OAuth2Service) RevokeServiceToken(ctx context.Context, tokenString string) error {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OAuth2Service.RevokeServiceToken")
	defer span.Finish()

	// 验证服务令牌
	request := &pb.ValidateRequest{}
	request.AccessToken = tokenString
	claims, err := s.ValidateServiceToken(ctx, request)
	if err != nil {
		s.logger.Error("验证服务令牌失败", "component", "oauth2_service", "error", err)
		return errors.BizWrap("oauth2.error.invalid_service_token", err)
	}

	// 撤销令牌
	if err := s.tokenDAO.RevokeTokenWithPrefix(ctx, claims.Id, s.config.OAuth2.ServiceAccount.TokenStorage.RedisKeyPrefix); err != nil {
		s.logger.Error("撤销服务令牌失败", "component", "oauth2_service", "error", err)
		return errors.BizWrap("oauth2.error.token_revocation_failed", err)
	}

	return nil
}
