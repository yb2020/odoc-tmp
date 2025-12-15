package dao

import (
	"context"

	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/internal/database"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/oauth2/model"
	"gorm.io/gorm"
)

// TokenDAO 令牌数据访问接口
type TokenDAO interface {
	// 保存令牌信息
	SaveToken(ctx context.Context, tokenInfo *model.OAuth2Token) error

	// 使用自定义前缀保存令牌信息（用于服务令牌）
	SaveTokenWithPrefix(ctx context.Context, tokenInfo *model.OAuth2Token, keyPrefix string) error

	// 根据令牌ID获取令牌信息
	GetTokenByID(ctx context.Context, tokenID string) (*model.OAuth2Token, error)

	// 使用自定义前缀根据令牌ID获取令牌信息（用于服务令牌）
	GetTokenByIDWithPrefix(ctx context.Context, tokenID string, keyPrefix string) (*model.OAuth2Token, error)

	// 根据访问令牌获取令牌信息
	GetTokenByAccessToken(ctx context.Context, accessToken string) (*model.OAuth2Token, error)

	// 根据刷新令牌获取令牌信息
	GetTokenByRefreshToken(ctx context.Context, refreshToken string) (*model.OAuth2Token, error)

	// 获取用户的所有令牌
	GetTokensByUserID(ctx context.Context, userID string) ([]*model.OAuth2Token, error)

	// 撤销令牌
	RevokeToken(ctx context.Context, tokenID string) error

	// 使用自定义前缀撤销令牌（用于服务令牌）
	RevokeTokenWithPrefix(ctx context.Context, tokenID string, keyPrefix string) error

	// 撤销用户的所有令牌
	RevokeAllTokensByUserID(ctx context.Context, userID string) error

	// 清理过期令牌
	CleanupExpiredTokens(ctx context.Context) error

	// 获取用户的令牌数量
	CountTokensByUserID(ctx context.Context, userID string) (int64, error)
}

// NewTokenDAO 创建令牌DAO，根据配置选择不同的实现
func NewTokenDAO(db *gorm.DB, redis database.RedisClient, logger logging.Logger, config *config.Config) TokenDAO {
	// 获取存储类型
	storageType := config.OAuth2.TokenStorage.Type

	// 根据存储类型选择不同的实现
	switch storageType {
	case "redis":
		// Redis实现
		logger.Info("msg", "使用Redis存储令牌")
		return newRedisTokenDAO(redis, logger, config)
	default:
		// 默认使用PostgreSQL存储
		logger.Info("msg", "使用PostgreSQL存储令牌")
		return newPostgresTokenDAO(db, redis, logger, config)
	}
}
