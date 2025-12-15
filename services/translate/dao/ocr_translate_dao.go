package dao

import (
	"context"
	"time"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/translate/model"
	"gorm.io/gorm"
)

// OCRTranslateDAO GORM实现的OCR翻译DAO
type OCRTranslateDAO struct {
	*baseDao.GormBaseDAO[model.OCRTranslateLog]
	logger logging.Logger
}

// NewOCRTranslateDAO 创建一个新的 OCR翻译DAO
func NewOCRTranslateDAO(db *gorm.DB, logger logging.Logger) OCRTranslateDAO {
	return OCRTranslateDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.OCRTranslateLog](db, logger),
		logger:      logger,
	}
}

// GetUserOCRHistoryBetween 获取用户本周OCR翻译记录
func (d *OCRTranslateDAO) GetUserOCRHistoryBetween(ctx context.Context, userID string, start time.Time, end time.Time) ([]model.OCRTranslateLog, error) {
	var logs []model.OCRTranslateLog
	result := d.GetDB(ctx).
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, start, end).
		Find(&logs)

	if result.Error != nil {
		d.logger.Error("msg", "获取用户OCR翻译历史失败", "userID", userID, "error", result.Error.Error())
		return nil, result.Error
	}
	return logs, nil
}
