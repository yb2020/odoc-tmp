package api

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	pb "github.com/yb2020/odoc-proto/gen/go/oauth2"
	"github.com/yb2020/odoc/config"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/i18n"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	"github.com/yb2020/odoc/pkg/utils"
	"github.com/yb2020/odoc/services/oauth2/helper"
	"github.com/yb2020/odoc/services/oauth2/model"
	"github.com/yb2020/odoc/services/oauth2/service"
)

// OAuth2APIHandler OAuth2处理程序
type OAuth2API struct {
	oauth2Service service.OAuth2Service
	logger        logging.Logger
	tracer        opentracing.Tracer
	localizer     i18n.Localizer
	config        *config.Config
	rsaUtil       *utils.RSAUtil
}

// NewOAuth2API 创建OAuth2 API
func NewOAuth2API(
	logger logging.Logger,
	tracer opentracing.Tracer,
	localizer i18n.Localizer,
	config *config.Config,
	oauth2Service service.OAuth2Service,
	rsaUtil *utils.RSAUtil,
) *OAuth2API {
	return &OAuth2API{
		oauth2Service: oauth2Service,
		logger:        logger,
		tracer:        tracer,
		localizer:     localizer,
		config:        config,
		rsaUtil:       rsaUtil,
	}
}

// RefreshHandler 刷新令牌处理程序
func (api *OAuth2API) RefreshHandler(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "/api/oauth2/refresh")
	defer span.Finish()

	// 使用 Proto 绑定器解析请求
	req := &pb.RefreshTokenRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("msg", "绑定刷新令牌请求失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 刷新令牌
	tokenResponse, err := api.oauth2Service.RefreshToken(ctx, req)
	if err != nil {
		api.logger.Warn("msg", "刷新令牌失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 返回响应
	response.Success(c, "token_refreshed", tokenResponse)
}

// ValidateHandler 验证令牌处理程序
func (api *OAuth2API) ValidateHandler(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "OAuth2API.ValidateHandler")
	defer span.Finish()

	// 使用 Proto 绑定器解析请求
	req := &pb.ValidateRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("msg", "绑定验证令牌请求失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 验证令牌
	claims, err := api.oauth2Service.ValidateToken(ctx, req)
	if err != nil {
		api.logger.Warn("msg", "验证令牌失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 返回响应
	response.Success(c, "token_valid", &pb.ValidateResponse{
		Valid: claims != nil,
	})
}

// GetRSAKeyHandler 获取RSA公钥处理程序
func (api *OAuth2API) GetRSAKeyHandler(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "OAuth2API.GetRSAKeyHandler")
	defer span.Finish()

	// 使用 Proto 绑定器解析请求
	req := &pb.GetRSAKeyRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("msg", "绑定获取RSA公钥请求失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 获取RSA公钥
	resp, err := api.oauth2Service.GetRSAKey(ctx, req)
	if err != nil {
		api.logger.Warn("msg", "获取RSA公钥失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 返回响应
	response.Success(c, "rsa_key_retrieved", resp)
}

// GetAuthCodeHandler 获取授权码处理程序
func (api *OAuth2API) GetAuthCodeHandler(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "OAuth2API.GetAuthCodeHandler")
	defer span.Finish()

	// 使用 Proto 绑定器解析请求
	req := &pb.GetAuthCodeRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("msg", "绑定授权码请求失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 获取授权码
	resp, err := api.oauth2Service.GetAuthCode(ctx, req)
	if err != nil {
		api.logger.Warn("msg", "获取授权码失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 返回响应
	response.Success(c, "auth_code_generated", resp)
}

// SignInAuthCodeHandler 登录/获取访问令牌处理程序
func (api *OAuth2API) SignInAuthCodeHandler(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "OAuth2API.SignInAuthCodeHandler")
	defer span.Finish()

	// 使用 Proto 绑定器解析请求
	req := &pb.SignInAuthCodeRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("msg", "绑定登录请求失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 转换请求
	tokenReq := &model.TokenRequest{
		Username: req.Username,
		Password: req.Password,
	}

	//RSA解密
	tokenReq.Password, _ = api.rsaUtil.DecryptBase64(tokenReq.Password)
	tokenReq.Username, _ = api.rsaUtil.DecryptBase64(tokenReq.Username)

	// 验证授权码并获取访问令牌
	resp, err := api.oauth2Service.GenerateToken(ctx, tokenReq)
	if err != nil {
		api.logger.Warn("msg", "授权码换取令牌失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 设置 cookie
	helper.WriteSuccessLoginCookie(c, resp.UserId, resp.AccessToken, resp.ExpiresAt, api.config.OAuth2.AppID)

	// 转换响应
	respProto := &pb.SignInAuthCodeResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    resp.ExpiresAt,
		ExpiresIn:    resp.ExpiresIn,
	}

	// 返回响应
	response.Success(c, "signin_successful", respProto)
}

// SignOutHandler 登出/撤销访问令牌处理程序
func (api *OAuth2API) SignOutHandler(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "OAuth2API.SignOutHandler")
	defer span.Finish()

	// 使用 Proto 绑定器解析请求
	req := &pb.SignOutRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("msg", "绑定登出请求失败", "error", err.Error())
		c.Error(err)
		return
	}
	accessToken := ctx.Value(userContext.AccessTokenKey)
	if accessToken == nil {
		api.logger.Error("msg", "访问令牌未找到", "component", "oauth2_api")
		c.Error(errors.Biz("oauth2.error.token_not_found"))
		return
	}
	// 撤销令牌
	err := api.oauth2Service.RevokeUserToken(ctx, accessToken.(string))
	if err != nil {
		api.logger.Warn("msg", "撤销用户令牌失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 从配置中获取 appId
	appId := api.config.OAuth2.AppID

	// 清除登录Cookie
	helper.DestroyCookieByAppId(c, appId)

	// 返回成功响应
	response.SuccessNoData(c, "signout_successful")
}

// GenerateServiceTokenHandler 生成服务令牌处理程序
func (api *OAuth2API) GenerateServiceTokenHandler(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "OAuth2API.GenerateServiceTokenHandler")
	defer span.Finish()

	// 解析请求
	var req pb.ServiceTokenRequest
	if err := transport.BindProto(c, &req); err != nil {
		api.logger.Warn("msg", "绑定服务令牌请求失败", "error", err.Error())
		c.Error(errors.BizWrap("service_token.error.invalid_request", err))
		return
	}

	// 验证服务密钥
	var validSecret bool
	for _, service := range api.config.OAuth2.ServiceAccount.Services {
		if service.Id == req.ServiceId && service.Name == req.ServiceName && service.Secret == req.Secret {
			validSecret = true
			break
		}
	}

	if !validSecret {
		api.logger.Error("msg", "服务密钥验证失败", "service_id", req.ServiceId, "service_name", req.ServiceName)
		c.Error(errors.Biz("service_token.error.invalid_service_credentials"))
		return
	}

	// 生成服务令牌
	tokenResp, err := api.oauth2Service.GenerateServiceToken(ctx, req.ServiceId, req.ServiceName)
	if err != nil {
		api.logger.Warn("msg", "生成服务令牌失败", "error", err.Error(), "service_id", req.ServiceId, "service_name", req.ServiceName)
		c.Error(err)
		return
	}

	// 转换为 protobuf 响应
	respProto := &pb.ServiceTokenResponse{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    tokenResp.ExpiresAt,
		ExpiresIn:    tokenResp.ExpiresIn,
	}

	// 返回响应
	response.Success(c, "service_token_generated", respProto)
}

// ValidateServiceTokenHandler 验证服务令牌处理程序
func (api *OAuth2API) ValidateServiceTokenHandler(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "OAuth2API.ValidateServiceTokenHandler")
	defer span.Finish()

	// 从请求头获取服务令牌
	tokenHeader := api.config.OAuth2.ServiceAccount.TokenHeaderName
	tokenString := c.GetHeader(tokenHeader)
	if tokenString == "" {
		api.logger.Error("msg", "服务令牌未找到", "component", "oauth2_api")
		c.Error(errors.Biz("service_token.error.service_token_not_found"))
		return
	}

	// 验证服务令牌
	request := &pb.ValidateRequest{}
	request.AccessToken = tokenString
	claims, err := api.oauth2Service.ValidateServiceToken(ctx, request)
	if err != nil {
		api.logger.Warn("msg", "验证服务令牌失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 返回响应
	respProto := &pb.ValidateResponse{
		Valid: claims != nil,
	}
	response.Success(c, "service_token_valid", respProto)
}

// RevokeServiceTokenHandler 撤销服务令牌处理程序
func (api *OAuth2API) RevokeServiceTokenHandler(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "OAuth2API.RevokeServiceTokenHandler")
	defer span.Finish()

	accessToken := ctx.Value(userContext.AccessTokenKey)
	if accessToken == nil {
		api.logger.Error("msg", "访问令牌未找到", "component", "oauth2_api")
		c.Error(errors.Biz("service_token.error.token_not_found"))
		return
	}

	// 撤销服务令牌
	err := api.oauth2Service.RevokeServiceToken(ctx, accessToken.(string))
	if err != nil {
		api.logger.Warn("msg", "撤销服务令牌失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 返回成功响应
	response.SuccessNoData(c, "service_token_revoked")
}
