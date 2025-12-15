package service

import (
	"github.com/gin-gonic/gin"
	pb "github.com/yb2020/odoc-proto/gen/go/oauth2"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/middleware"
	"github.com/yb2020/odoc/services/oauth2/model"
)

// 确保Claims实现了middleware.Claims接口
var _ middleware.Claims = (*model.Claims)(nil)

// OAuth2AuthAdapter 适配OAuth2服务到通用认证接口
type OAuth2AuthAdapter struct {
	oauth2Service OAuth2Service
}

// NewOAuth2AuthAdapter 创建OAuth2认证适配器
func NewOAuth2AuthAdapter(oauth2Service OAuth2Service) *OAuth2AuthAdapter {
	return &OAuth2AuthAdapter{
		oauth2Service: oauth2Service,
	}
}

// ValidateToken 验证令牌
func (a *OAuth2AuthAdapter) ValidateToken(ctx *gin.Context, token string) (middleware.Claims, error) {
	request := &pb.ValidateRequest{}
	request.AccessToken = token
	claims, err := a.oauth2Service.ValidateToken(ctx, request)

	if err != nil {
		return nil, errors.BizWrap("oauth2.validate_token.errors.invalid_token", err)
	}

	return claims, nil
}

// RevokeToken 撤销令牌
func (a *OAuth2AuthAdapter) RevokeToken(ctx *gin.Context, access_token string) error {
	return a.oauth2Service.RevokeUserToken(ctx, access_token)
}

// RevokeAllTokens 撤销所有令牌
func (a *OAuth2AuthAdapter) RevokeUserTokens(ctx *gin.Context, userId string) error {
	return a.oauth2Service.RevokeUserTokens(ctx, userId)
}

// ValidateServiceToken 验证服务令牌
func (a *OAuth2AuthAdapter) ValidateServiceToken(ctx *gin.Context, token string) (middleware.Claims, error) {
	request := &pb.ValidateRequest{}
	request.AccessToken = token
	claims, err := a.oauth2Service.ValidateServiceToken(ctx, request)

	if err != nil {
		return nil, errors.BizWrap("oauth2.validate_service_token.errors.invalid_token", err)
	}

	return claims, nil
}
