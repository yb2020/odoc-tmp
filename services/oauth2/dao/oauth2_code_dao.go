package dao

import (
	"context"
	"time"

	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/oauth2/model"
	"gorm.io/gorm"
)

// OAuth2CodeDAO PostgreSQL实现的授权码DAO
type OAuth2CodeDAO struct {
	db     *gorm.DB
	logger logging.Logger
	config *config.Config
}

// NewAuthCodeDAO 创建授权码DAO
func NewAuthCodeDAO(db *gorm.DB, logger logging.Logger, config *config.Config) OAuth2CodeDAO {
	return OAuth2CodeDAO{
		db:     db,
		logger: logger,
		config: config,
	}
}

// SaveAuthCode 保存授权码
func (dao *OAuth2CodeDAO) SaveAuthCode(ctx context.Context, authCode *model.OAuth2AuthCode) error {
	return dao.db.WithContext(ctx).Create(authCode).Error
}

// GetAuthCode 根据授权码获取授权码信息
func (dao *OAuth2CodeDAO) GetAuthCode(ctx context.Context, code string) (*model.OAuth2AuthCode, error) {
	var authCode model.OAuth2AuthCode
	err := dao.db.WithContext(ctx).Where("code = ?", code).First(&authCode).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &authCode, nil
}

// MarkAuthCodeAsUsed 标记授权码为已使用
func (dao *OAuth2CodeDAO) MarkAuthCodeAsUsed(ctx context.Context, code string) error {
	return dao.db.WithContext(ctx).Model(&model.OAuth2AuthCode{}).Where("code = ?", code).Update("used", true).Error
}

// CleanupExpiredAuthCodes 清理过期的授权码
func (dao *OAuth2CodeDAO) CleanupExpiredAuthCodes(ctx context.Context) error {
	return dao.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&model.OAuth2AuthCode{}).Error
}
