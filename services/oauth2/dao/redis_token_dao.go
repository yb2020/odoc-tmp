package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/internal/database"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/oauth2/model"
)

// RedisTokenDAO Redis实现的令牌DAO
type RedisTokenDAO struct {
	redis  database.RedisClient
	logger logging.Logger
	config *config.Config
}

// newRedisTokenDAO 创建Redis令牌DAO（包内私有）
func newRedisTokenDAO(redis database.RedisClient, logger logging.Logger, config *config.Config) *RedisTokenDAO {
	return &RedisTokenDAO{
		redis:  redis,
		logger: logger,
		config: config,
	}
}

// 生成Redis键名
func (dao *RedisTokenDAO) tokenKey(tokenID string) string {
	return fmt.Sprintf(dao.config.OAuth2.TokenStorage.RedisKeyPrefix.Token, tokenID)
}

func (dao *RedisTokenDAO) accessTokenKey(accessToken string) string {
	return fmt.Sprintf(dao.config.OAuth2.TokenStorage.RedisKeyPrefix.AccessToken, accessToken)
}

func (dao *RedisTokenDAO) refreshTokenKey(refreshToken string) string {
	return fmt.Sprintf(dao.config.OAuth2.TokenStorage.RedisKeyPrefix.RefreshToken, refreshToken)
}

func (dao *RedisTokenDAO) userTokensKey(userID string) string {
	return fmt.Sprintf(dao.config.OAuth2.TokenStorage.RedisKeyPrefix.UserTokens, userID)
}

// 使用自定义前缀生成Redis键名
func (dao *RedisTokenDAO) tokenKeyWithPrefix(tokenID string, keyPrefix string) string {
	return fmt.Sprintf("%s%s", keyPrefix, tokenID)
}

// SaveToken 保存令牌信息
func (dao *RedisTokenDAO) SaveToken(ctx context.Context, tokenInfo *model.OAuth2Token) error {
	// 检查用户令牌数量
	count, err := dao.CountTokensByUserID(ctx, tokenInfo.UserId)
	if err != nil {
		return err
	}

	// 如果超过最大数量，删除最旧的令牌
	maxTokensPerUser := dao.config.OAuth2.TokenStorage.MaxTokensPerUser
	if count >= int64(maxTokensPerUser) {
		// 获取用户的所有令牌
		tokens, err := dao.GetTokensByUserID(ctx, tokenInfo.UserId)
		if err == nil && len(tokens) > 0 {
			// 找到最旧的令牌（按创建时间排序）
			var oldestToken *model.OAuth2Token
			oldestTime := time.Now()

			for _, token := range tokens {
				if token.CreatedAt.Before(oldestTime) {
					oldestToken = token
					oldestTime = token.CreatedAt
				}
			}

			if oldestToken != nil {
				// 撤销最旧的令牌
				dao.RevokeToken(ctx, oldestToken.TokenId)
			}
		}
	}

	// 序列化令牌信息
	tokenJSON, err := json.Marshal(tokenInfo)
	if err != nil {
		dao.logger.Error("msg", "序列化令牌失败", "error", err.Error())
		return err
	}

	// 计算过期时间
	expiration := tokenInfo.ExpiresAt.Sub(time.Now()) + time.Hour

	// 存储令牌信息
	tokenKey := dao.tokenKey(tokenInfo.TokenId)
	if err := dao.redis.Set(ctx, tokenKey, string(tokenJSON), expiration).Err(); err != nil {
		dao.logger.Error("msg", "保存令牌失败", "error", err.Error())
		return err
	}

	// 创建访问令牌索引
	accessTokenKey := dao.accessTokenKey(tokenInfo.AccessToken)
	if err := dao.redis.Set(ctx, accessTokenKey, tokenInfo.TokenId, expiration).Err(); err != nil {
		dao.logger.Warn("msg", "创建访问令牌索引失败", "error", err.Error())
	}

	// 创建刷新令牌索引
	refreshTokenKey := dao.refreshTokenKey(tokenInfo.RefreshToken)
	if err := dao.redis.Set(ctx, refreshTokenKey, tokenInfo.TokenId, expiration).Err(); err != nil {
		dao.logger.Warn("msg", "创建刷新令牌索引失败", "error", err.Error())
	}

	// 将令牌ID添加到用户的令牌集合中
	userTokensKey := dao.userTokensKey(tokenInfo.UserId)
	if err := dao.redis.SAdd(ctx, userTokensKey, tokenInfo.TokenId).Err(); err != nil {
		dao.logger.Warn("msg", "添加令牌到用户集合失败", "error", err.Error())
	}
	// 设置用户令牌集合的过期时间
	dao.redis.Expire(ctx, userTokensKey, 30*24*time.Hour)

	return nil
}

// SaveTokenWithPrefix 使用自定义前缀保存令牌信息（用于服务令牌）
func (dao *RedisTokenDAO) SaveTokenWithPrefix(ctx context.Context, tokenInfo *model.OAuth2Token, keyPrefix string) error {
	// 序列化令牌信息
	tokenJSON, err := json.Marshal(tokenInfo)
	if err != nil {
		dao.logger.Error("msg", "序列化令牌失败", "error", err.Error())
		return err
	}

	// 计算过期时间
	expiration := tokenInfo.ExpiresAt.Sub(time.Now()) + time.Hour

	// 存储令牌信息
	tokenKey := dao.tokenKeyWithPrefix(tokenInfo.TokenId, keyPrefix)
	if err := dao.redis.Set(ctx, tokenKey, string(tokenJSON), expiration).Err(); err != nil {
		dao.logger.Error("msg", "保存服务令牌失败", "error", err.Error())
		return err
	}

	dao.logger.Info("msg", "成功保存服务令牌", "token_id", tokenInfo.TokenId)
	return nil
}

// GetTokenByID 根据令牌ID获取令牌信息
func (dao *RedisTokenDAO) GetTokenByID(ctx context.Context, tokenID string) (*model.OAuth2Token, error) {
	// 从Redis获取令牌信息
	tokenKey := dao.tokenKey(tokenID)
	tokenJSON, err := dao.redis.Get(ctx, tokenKey).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, nil
		}
		dao.logger.Error("msg", "获取令牌失败", "error", err.Error())
		return nil, err
	}

	// 反序列化令牌信息
	var tokenInfo model.OAuth2Token
	if err := json.Unmarshal([]byte(tokenJSON), &tokenInfo); err != nil {
		dao.logger.Error("msg", "反序列化令牌失败", "error", err.Error())
		return nil, err
	}

	return &tokenInfo, nil
}

// GetTokenByIDWithPrefix 使用自定义前缀根据令牌ID获取令牌信息（用于服务令牌）
func (dao *RedisTokenDAO) GetTokenByIDWithPrefix(ctx context.Context, tokenID string, keyPrefix string) (*model.OAuth2Token, error) {
	// 从Redis获取令牌信息
	tokenKey := dao.tokenKeyWithPrefix(tokenID, keyPrefix)
	tokenJSON, err := dao.redis.Get(ctx, tokenKey).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, nil
		}
		dao.logger.Error("msg", "获取服务令牌失败", "error", err.Error())
		return nil, err
	}

	// 反序列化令牌信息
	var token model.OAuth2Token
	if err := json.Unmarshal([]byte(tokenJSON), &token); err != nil {
		dao.logger.Error("msg", "反序列化服务令牌失败", "error", err.Error())
		return nil, err
	}

	return &token, nil
}

// GetTokenByAccessToken 根据访问令牌获取令牌信息
func (dao *RedisTokenDAO) GetTokenByAccessToken(ctx context.Context, accessToken string) (*model.OAuth2Token, error) {
	// 从Redis获取令牌ID
	accessTokenKey := dao.accessTokenKey(accessToken)
	tokenID, err := dao.redis.Get(ctx, accessTokenKey).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, nil
		}
		dao.logger.Error("msg", "获取访问令牌索引失败", "error", err.Error())
		return nil, err
	}

	// 使用令牌ID获取令牌信息
	return dao.GetTokenByID(ctx, tokenID)
}

// GetTokenByRefreshToken 根据刷新令牌获取令牌信息
func (dao *RedisTokenDAO) GetTokenByRefreshToken(ctx context.Context, refreshToken string) (*model.OAuth2Token, error) {
	// 从Redis获取令牌ID
	refreshTokenKey := dao.refreshTokenKey(refreshToken)
	tokenID, err := dao.redis.Get(ctx, refreshTokenKey).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, nil
		}
		dao.logger.Error("msg", "获取刷新令牌索引失败", "error", err.Error())
		return nil, err
	}

	// 使用令牌ID获取令牌信息
	return dao.GetTokenByID(ctx, tokenID)
}

// GetTokensByUserID 获取用户的所有令牌
func (dao *RedisTokenDAO) GetTokensByUserID(ctx context.Context, userID string) ([]*model.OAuth2Token, error) {
	// 获取用户的所有令牌ID
	userTokensKey := dao.userTokensKey(userID)
	tokenIDs, err := dao.redis.SMembers(ctx, userTokensKey).Result()
	if err != nil {
		dao.logger.Error("msg", "获取用户令牌集合失败", "error", err.Error())
		return nil, err
	}

	// 获取所有令牌信息
	var tokens []*model.OAuth2Token
	for _, tokenID := range tokenIDs {
		token, err := dao.GetTokenByID(ctx, tokenID)
		if err != nil {
			continue
		}
		if token != nil && !token.Revoked && token.ExpiresAt.After(time.Now()) {
			tokens = append(tokens, token)
		}
	}

	return tokens, nil
}

// RevokeToken 撤销令牌
func (dao *RedisTokenDAO) RevokeToken(ctx context.Context, tokenID string) error {
	// 获取令牌信息
	token, err := dao.GetTokenByID(ctx, tokenID)
	if err != nil {
		return err
	}
	if token == nil {
		return nil
	}

	// 标记令牌为已撤销
	token.Revoked = true

	// 序列化令牌信息
	tokenJSON, err := json.Marshal(token)
	if err != nil {
		dao.logger.Error("msg", "序列化令牌失败", "error", err.Error())
		return err
	}

	// 更新令牌信息
	tokenKey := dao.tokenKey(tokenID)
	if err := dao.redis.Set(ctx, tokenKey, string(tokenJSON), time.Hour).Err(); err != nil {
		dao.logger.Error("msg", "更新令牌失败", "error", err.Error())
		return err
	}

	// 删除访问令牌索引
	accessTokenKey := dao.accessTokenKey(token.AccessToken)
	dao.redis.Del(ctx, accessTokenKey)

	// 删除刷新令牌索引
	refreshTokenKey := dao.refreshTokenKey(token.RefreshToken)
	dao.redis.Del(ctx, refreshTokenKey)

	// 从用户令牌集合中删除
	userTokensKey := dao.userTokensKey(token.UserId)
	dao.redis.SRem(ctx, userTokensKey, tokenID)

	return nil
}

// RevokeTokenWithPrefix 使用自定义前缀撤销令牌（用于服务令牌）
func (dao *RedisTokenDAO) RevokeTokenWithPrefix(ctx context.Context, tokenID string, keyPrefix string) error {
	// 获取令牌信息
	token, err := dao.GetTokenByIDWithPrefix(ctx, tokenID, keyPrefix)
	if err != nil {
		return err
	}
	if token == nil {
		return nil
	}

	// 标记令牌为已撤销
	token.Revoked = true

	// 序列化令牌信息
	tokenJSON, err := json.Marshal(token)
	if err != nil {
		dao.logger.Error("msg", "序列化服务令牌失败", "error", err.Error())
		return err
	}

	// 更新令牌信息
	tokenKey := dao.tokenKeyWithPrefix(tokenID, keyPrefix)
	if err := dao.redis.Set(ctx, tokenKey, string(tokenJSON), time.Hour).Err(); err != nil {
		dao.logger.Error("msg", "更新服务令牌失败", "error", err.Error())
		return err
	}

	dao.logger.Info("msg", "成功撤销服务令牌", "token_id", tokenID)
	return nil
}

// RevokeAllTokensByUserID 撤销用户的所有令牌
func (dao *RedisTokenDAO) RevokeAllTokensByUserID(ctx context.Context, userID string) error {
	// 获取用户的所有令牌
	tokens, err := dao.GetTokensByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// 撤销所有令牌
	for _, token := range tokens {
		dao.RevokeToken(ctx, token.TokenId)
	}

	return nil
}

// CleanupExpiredTokens 清理过期令牌
func (dao *RedisTokenDAO) CleanupExpiredTokens(ctx context.Context) error {
	// Redis会自动清理过期的键，所以这里不需要额外的实现
	dao.logger.Info("msg", "Redis会自动清理过期的键，不需要额外实现")
	return nil
}

// CountTokensByUserID 获取用户的令牌数量
func (dao *RedisTokenDAO) CountTokensByUserID(ctx context.Context, userID string) (int64, error) {
	// 获取用户的所有令牌ID
	userTokensKey := dao.userTokensKey(userID)
	count, err := dao.redis.SCard(ctx, userTokensKey).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return 0, nil
		}
		dao.logger.Error("msg", "获取用户令牌数量失败", "error", err.Error())
		return 0, err
	}

	return count, nil
}
