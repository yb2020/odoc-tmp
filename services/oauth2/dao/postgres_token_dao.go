package dao

import (
	"context"
	"errors"
	"time"

	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/internal/database"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/oauth2/model"
	"gorm.io/gorm"
)

// PostgresTokenDAO PostgreSQL实现的令牌DAO
type PostgresTokenDAO struct {
	db     *gorm.DB
	redis  database.RedisClient
	logger logging.Logger
	config *config.Config
}

// newPostgresTokenDAO 创建PostgreSQL令牌DAO（包内私有）
func newPostgresTokenDAO(db *gorm.DB, redis database.RedisClient, logger logging.Logger, config *config.Config) *PostgresTokenDAO {
	return &PostgresTokenDAO{
		db:     db,
		redis:  redis,
		logger: logger,
		config: config,
	}
}

// SaveToken 保存令牌信息
func (dao *PostgresTokenDAO) SaveToken(ctx context.Context, tokenInfo *model.OAuth2Token) error {
	// 检查用户令牌数量
	count, err := dao.CountTokensByUserID(ctx, tokenInfo.UserId)
	if err != nil {
		return err
	}

	// 如果超过最大数量，删除最旧的令牌
	maxTokensPerUser := dao.config.OAuth2.TokenStorage.MaxTokensPerUser
	if count >= int64(maxTokensPerUser) {
		var oldestToken model.OAuth2Token
		if err := dao.db.WithContext(ctx).Where("user_id = ? AND revoked = ?", tokenInfo.UserId, false).
			Order("created_at ASC").First(&oldestToken).Error; err == nil {
			// 撤销最旧的令牌
			dao.RevokeToken(ctx, oldestToken.TokenId)
		}
	}

	// 保存到数据库
	if err := dao.db.WithContext(ctx).Create(tokenInfo).Error; err != nil {
		dao.logger.Error("msg", "保存令牌失败", "error", err.Error())
		return err
	}

	// 缓存到Redis
	if dao.redis != nil {
		// 使用令牌ID作为键
		key := "token:" + tokenInfo.TokenId
		// 设置过期时间比令牌稍长一些
		expiration := tokenInfo.ExpiresAt.Sub(time.Now()) + time.Hour

		// 将令牌信息序列化为JSON并存储
		if err := dao.redis.Set(ctx, key, tokenInfo.AccessToken, expiration).Err(); err != nil {
			dao.logger.Warn("msg", "缓存令牌失败", "error", err.Error())
			// 不返回错误，因为数据库已保存成功
		}
	}

	return nil
}

// GetTokenByID 根据令牌ID获取令牌信息
func (dao *PostgresTokenDAO) GetTokenByID(ctx context.Context, tokenID string) (*model.OAuth2Token, error) {
	var tokenInfo model.OAuth2Token

	// 从数据库查询
	if err := dao.db.WithContext(ctx).Where("token_id = ?", tokenID).First(&tokenInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		dao.logger.Error("msg", "查询令牌失败", "error", err.Error())
		return nil, err
	}

	return &tokenInfo, nil
}

// GetTokenByAccessToken 根据访问令牌获取令牌信息
func (dao *PostgresTokenDAO) GetTokenByAccessToken(ctx context.Context, accessToken string) (*model.OAuth2Token, error) {
	var tokenInfo model.OAuth2Token

	// 从数据库查询
	if err := dao.db.WithContext(ctx).Where("access_token = ? AND revoked = ?", accessToken, false).First(&tokenInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		dao.logger.Error("msg", "查询令牌失败", "error", err.Error())
		return nil, err
	}

	return &tokenInfo, nil
}

// GetTokenByRefreshToken 根据刷新令牌获取令牌信息
func (dao *PostgresTokenDAO) GetTokenByRefreshToken(ctx context.Context, refreshToken string) (*model.OAuth2Token, error) {
	var tokenInfo model.OAuth2Token

	// 从数据库查询
	if err := dao.db.WithContext(ctx).Where("refresh_token = ? AND revoked = ?", refreshToken, false).First(&tokenInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		dao.logger.Error("msg", "查询令牌失败", "error", err.Error())
		return nil, err
	}

	return &tokenInfo, nil
}

// GetTokensByUserID 获取用户的所有令牌
func (dao *PostgresTokenDAO) GetTokensByUserID(ctx context.Context, userID string) ([]*model.OAuth2Token, error) {
	var tokens []*model.OAuth2Token

	// 从数据库查询
	if err := dao.db.WithContext(ctx).Where("user_id = ? AND revoked = ?", userID, false).Find(&tokens).Error; err != nil {
		dao.logger.Error("msg", "查询用户令牌失败", "error", err.Error())
		return nil, err
	}

	return tokens, nil
}

// RevokeToken 撤销令牌
func (dao *PostgresTokenDAO) RevokeToken(ctx context.Context, tokenID string) error {
	// 更新数据库
	result := dao.db.WithContext(ctx).Model(&model.OAuth2Token{}).
		Where("token_id = ?", tokenID).
		Update("revoked", true)

	if result.Error != nil {
		dao.logger.Error("msg", "撤销令牌失败", "error", result.Error.Error())
		return result.Error
	}

	// 从Redis删除
	if dao.redis != nil {
		key := "token:" + tokenID
		if err := dao.redis.Del(ctx, key).Err(); err != nil {
			dao.logger.Warn("msg", "从缓存删除令牌失败", "error", err.Error())
			// 不返回错误，因为数据库已更新成功
		}
	}

	return nil
}

// RevokeAllTokensByUserID 撤销用户的所有令牌
func (dao *PostgresTokenDAO) RevokeAllTokensByUserID(ctx context.Context, userID string) error {
	// 更新数据库
	result := dao.db.WithContext(ctx).Model(&model.OAuth2Token{}).
		Where("user_id = ? AND revoked = ?", userID, false).
		Update("revoked", true)

	if result.Error != nil {
		dao.logger.Error("msg", "撤销用户所有令牌失败", "error", result.Error.Error())
		return result.Error
	}

	// 从Redis删除
	if dao.redis != nil {
		// 获取所有需要删除的令牌
		var tokens []*model.OAuth2Token
		if err := dao.db.WithContext(ctx).Where("user_id = ?", userID).Find(&tokens).Error; err == nil {
			for _, token := range tokens {
				key := "token:" + token.TokenId
				dao.redis.Del(ctx, key)
			}
		}
	}

	return nil
}

// CleanupExpiredTokens 清理过期令牌
func (dao *PostgresTokenDAO) CleanupExpiredTokens(ctx context.Context) error {
	// 更新数据库，将过期的令牌标记为已撤销
	result := dao.db.WithContext(ctx).Model(&model.OAuth2Token{}).
		Where("expires_at < ? AND revoked = ?", time.Now(), false).
		Update("revoked", true)

	if result.Error != nil {
		dao.logger.Error("msg", "清理过期令牌失败", "error", result.Error.Error())
		return result.Error
	}

	return nil
}

// CountTokensByUserID 获取用户的令牌数量
func (dao *PostgresTokenDAO) CountTokensByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64

	// 从数据库查询
	if err := dao.db.WithContext(ctx).Model(&model.OAuth2Token{}).
		Where("user_id = ? AND revoked = ?", userID, false).
		Count(&count).Error; err != nil {
		dao.logger.Error("msg", "查询用户令牌数量失败", "error", err.Error())
		return 0, err
	}

	return count, nil
}

// SaveTokenWithPrefix 使用自定义前缀保存令牌信息（用于服务令牌）
func (dao *PostgresTokenDAO) SaveTokenWithPrefix(ctx context.Context, tokenInfo *model.OAuth2Token, keyPrefix string) error {
	// 保存到数据库
	if err := dao.db.WithContext(ctx).Create(tokenInfo).Error; err != nil {
		dao.logger.Error("msg", "保存服务令牌失败", "error", err.Error())
		return err
	}

	// 缓存到Redis
	if dao.redis != nil {
		// 使用自定义前缀和令牌ID作为键
		key := keyPrefix + tokenInfo.TokenId
		// 设置过期时间比令牌稍长一些
		expiration := tokenInfo.ExpiresAt.Sub(time.Now()) + time.Hour

		// 将令牌信息序列化为JSON并存储
		if err := dao.redis.Set(ctx, key, tokenInfo.AccessToken, expiration).Err(); err != nil {
			dao.logger.Warn("msg", "缓存服务令牌失败", "error", err.Error())
			// 不返回错误，因为数据库已保存成功
		}
	}

	return nil
}

// GetTokenByIDWithPrefix 使用自定义前缀根据令牌ID获取令牌信息（用于服务令牌）
func (dao *PostgresTokenDAO) GetTokenByIDWithPrefix(ctx context.Context, tokenID string, keyPrefix string) (*model.OAuth2Token, error) {
	var tokenInfo model.OAuth2Token

	// 首先尝试从Redis获取
	if dao.redis != nil {
		key := keyPrefix + tokenID
		accessToken, err := dao.redis.Get(ctx, key).Result()
		if err == nil && accessToken != "" {
			// 从Redis获取成功，再从数据库获取完整信息
			if err := dao.db.WithContext(ctx).Where("token_id = ? AND access_token = ?", tokenID, accessToken).First(&tokenInfo).Error; err == nil {
				return &tokenInfo, nil
			}
		}
	}

	// 从数据库查询
	if err := dao.db.WithContext(ctx).Where("token_id = ?", tokenID).First(&tokenInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		dao.logger.Error("msg", "查询服务令牌失败", "error", err.Error())
		return nil, err
	}

	return &tokenInfo, nil
}

// RevokeTokenWithPrefix 使用自定义前缀撤销令牌（用于服务令牌）
func (dao *PostgresTokenDAO) RevokeTokenWithPrefix(ctx context.Context, tokenID string, keyPrefix string) error {
	// 更新数据库
	result := dao.db.WithContext(ctx).Model(&model.OAuth2Token{}).
		Where("token_id = ?", tokenID).
		Update("revoked", true)

	if result.Error != nil {
		dao.logger.Error("msg", "撤销服务令牌失败", "error", result.Error.Error())
		return result.Error
	}

	// 从Redis删除
	if dao.redis != nil {
		key := keyPrefix + tokenID
		if err := dao.redis.Del(ctx, key).Err(); err != nil {
			dao.logger.Warn("msg", "从缓存删除服务令牌失败", "error", err.Error())
			// 不返回错误，因为数据库已更新成功
		}
	}

	return nil
}
