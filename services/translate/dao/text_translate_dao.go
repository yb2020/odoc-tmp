package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/proto/gen/go/translate"
	"github.com/yb2020/odoc/services/translate/model"
	"gorm.io/gorm"
)

// TextTranslateDAO 实现的翻译DAO
type TextTranslateDAO struct {
	*baseDao.GormBaseDAO[model.TextTranslateLog]
	logger logging.Logger
}

// NewTextTranslateDAO 创建一个新的GORM翻译DAO
func NewTextTranslateDAO(db *gorm.DB, logger logging.Logger) TextTranslateDAO {
	return TextTranslateDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.TextTranslateLog](db, logger),
		logger:      logger,
	}
}

// FindByMD5 根据MD5查找翻译日志
func (d *TextTranslateDAO) FindByMD5(ctx context.Context, md5Hash string) (*model.TextTranslateLog, error) {
	var log model.TextTranslateLog
	result := d.GetDB(ctx).
		Where("md5_hash = ? AND status = ?", md5Hash, 1).
		Order("created_at DESC").
		First(&log)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 日志不存在
		}
		d.logger.Error("msg", "根据MD5查找翻译日志失败", "md5", md5Hash, "error", result.Error.Error())
		return nil, result.Error
	}
	return &log, nil
}

// FindByRequestIDAndChannel 根据请求ID和渠道查找翻译日志
func (d *TextTranslateDAO) FindByRequestIDAndChannel(ctx context.Context, requestID string, channel translate.TranslateChannel) (*model.TextTranslateLog, error) {
	var log model.TextTranslateLog
	result := d.GetDB(ctx).
		Where("request_id = ? AND channel = ?", requestID, channel.String()).
		First(&log)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 日志不存在
		}
		d.logger.Error("msg", "根据请求ID查找翻译日志失败", "requestID", requestID, "error", result.Error.Error())
		return nil, result.Error
	}
	return &log, nil
}

// GetUserTranslateHistory 获取用户的翻译历史
func (d *TextTranslateDAO) GetUserTranslateHistory(ctx context.Context, userID string, limit, offset int) ([]model.TextTranslateLog, error) {
	var logs []model.TextTranslateLog
	result := d.GetDB(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs)

	if result.Error != nil {
		d.logger.Error("msg", "获取用户翻译历史失败", "userID", userID, "error", result.Error.Error())
		return nil, result.Error
	}
	return logs, nil
}

// GetUserTranslateHistoryCount 获取用户的翻译历史总数
func (d *TextTranslateDAO) GetUserTranslateHistoryCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	result := d.GetDB(ctx).
		Model(&model.TextTranslateLog{}).
		Where("user_id = ?", userID).
		Count(&count)

	if result.Error != nil {
		d.logger.Error("msg", "获取用户翻译历史总数失败", "userID", userID, "error", result.Error.Error())
		return 0, result.Error
	}
	return count, nil
}
